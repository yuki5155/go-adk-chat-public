import * as cdk from 'aws-cdk-lib';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as ecr from 'aws-cdk-lib/aws-ecr';
import * as ecs from 'aws-cdk-lib/aws-ecs';
import * as elbv2 from 'aws-cdk-lib/aws-elasticloadbalancingv2';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as logs from 'aws-cdk-lib/aws-logs';
import * as acm from 'aws-cdk-lib/aws-certificatemanager';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as route53Targets from 'aws-cdk-lib/aws-route53-targets';
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager';
import { Construct } from 'constructs';
import { RdsRequests, ContainerConfigRequests, CustomSecretsRequests } from './requests';

export interface BackendStackProps extends cdk.StackProps {
  projectName: string;
  environment: string;
  domainName?: string;
  databaseStackName?: string;
  isDatabaseStackDeployed?: boolean;
  rdsRequests?: RdsRequests;
  containerConfigRequests?: ContainerConfigRequests;
  containerPort?: number;
  imageTag?: string;
  cpu?: number;
  memory?: number;
  desiredCount?: number;
  containerCommand?: string[] | null;
  migrationCommand?: string[] | null;
}

export class BackendStack extends cdk.Stack {
  public readonly ecrRepository: ecr.Repository;
  public readonly cluster: ecs.Cluster;
  public readonly service: ecs.FargateService;
  public readonly loadBalancer: elbv2.ApplicationLoadBalancer;

  constructor(scope: Construct, id: string, props: BackendStackProps) {
    super(scope, id, props);

    const {
      projectName,
      environment,
      domainName,
      isDatabaseStackDeployed = false,
      containerPort = 8000,
      imageTag = 'latest',
      cpu = 256,
      memory = 512,
      desiredCount = 1,
    } = props;

    // ── ECR Repository ────────────────────────────────────────────────────────
    this.ecrRepository = new ecr.Repository(this, 'BackendRepository', {
      repositoryName: `${projectName}-${environment}-backend`,
      removalPolicy: environment === 'prod' ? cdk.RemovalPolicy.RETAIN : cdk.RemovalPolicy.DESTROY,
      emptyOnDelete: environment !== 'prod',
    });

    // ── VPC (imported from NetworkStack) ─────────────────────────────────────
    const vpc = ec2.Vpc.fromLookup(this, 'ImportedVpc', {
      tags: { Project: projectName, Environment: environment, StackType: 'Network' },
    });

    // ── ECS Cluster ───────────────────────────────────────────────────────────
    this.cluster = new ecs.Cluster(this, 'BackendCluster', {
      vpc,
      clusterName: `${projectName}-${environment}-backend-cluster`,
      containerInsights: environment === 'prod',
    });

    // ── IAM Roles ─────────────────────────────────────────────────────────────
    const taskRole = new iam.Role(this, 'BackendTaskRole', {
      assumedBy: new iam.ServicePrincipal('ecs-tasks.amazonaws.com'),
      description: `Task role for ${projectName}-${environment} backend`,
    });
    taskRole.addToPolicy(new iam.PolicyStatement({
      actions: ['secretsmanager:GetSecretValue', 'secretsmanager:DescribeSecret'],
      resources: [`arn:aws:secretsmanager:${this.region}:${this.account}:secret:${projectName}/${environment}/*`],
    }));

    const executionRole = new iam.Role(this, 'BackendExecutionRole', {
      assumedBy: new iam.ServicePrincipal('ecs-tasks.amazonaws.com'),
      managedPolicies: [
        iam.ManagedPolicy.fromAwsManagedPolicyName('service-role/AmazonECSTaskExecutionRolePolicy'),
      ],
    });
    executionRole.addToPolicy(new iam.PolicyStatement({
      actions: ['secretsmanager:GetSecretValue', 'secretsmanager:DescribeSecret'],
      resources: [`arn:aws:secretsmanager:${this.region}:${this.account}:secret:*`],
    }));

    // ── Log Group ─────────────────────────────────────────────────────────────
    const logGroup = new logs.LogGroup(this, 'BackendLogGroup', {
      logGroupName: `/ecs/${projectName}-${environment}-backend`,
      retention: logs.RetentionDays.ONE_WEEK,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    // ── Check for images ──────────────────────────────────────────────────────
    const hasImages = this.node.tryGetContext('hasImages');
    if (hasImages === 'false' || hasImages === false || !hasImages) {
      console.log('No images found — deploying ECR repository only. Use --context hasImages=true after pushing images.');
      this.setupOutputs(projectName, environment, domainName);
      return;
    }

    // ── SSL / Domain ──────────────────────────────────────────────────────────
    let certificate: acm.ICertificate | undefined;
    let hostedZone: route53.IHostedZone | undefined;
    if (domainName) {
      const rootDomain = domainName.split('.').slice(-2).join('.');
      hostedZone = route53.HostedZone.fromLookup(this, 'HostedZone', { domainName: rootDomain });
      certificate = new acm.Certificate(this, 'BackendCertificate', {
        domainName,
        validation: acm.CertificateValidation.fromDns(hostedZone),
      });
    }

    // ── Task Definition ───────────────────────────────────────────────────────
    const taskDefinition = new ecs.FargateTaskDefinition(this, 'BackendTaskDefinition', {
      taskRole,
      executionRole,
      cpu,
      memoryLimitMiB: memory,
    });

    // Build environment variables and secrets from ContainerConfigRequests
    const envVars = this.buildEnvironmentVariables(props, environment, containerPort);
    const secrets = this.buildSecrets(props);

    const container = taskDefinition.addContainer('backend-container', {
      image: ecs.ContainerImage.fromEcrRepository(this.ecrRepository, imageTag),
      containerName: 'backend',
      logging: ecs.LogDrivers.awsLogs({ streamPrefix: 'backend', logGroup }),
      environment: envVars,
      secrets,
      healthCheck: {
        command: ['CMD-SHELL', `curl -f http://localhost:${containerPort}/health || exit 1`],
        interval: cdk.Duration.seconds(30),
        timeout: cdk.Duration.seconds(5),
        retries: 3,
        startPeriod: cdk.Duration.seconds(90),
      },
      ...(props.containerCommand && { command: props.containerCommand }),
    });
    container.addPortMappings({ containerPort, protocol: ecs.Protocol.TCP });

    // ── Security Group ────────────────────────────────────────────────────────
    const ecsTaskSg = new ec2.SecurityGroup(this, 'EcsTaskSecurityGroup', {
      vpc,
      description: 'Security group for ECS backend tasks',
      allowAllOutbound: true,
    });
    ecsTaskSg.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(containerPort), 'Allow ALB to ECS tasks');

    // ── Fargate Service ───────────────────────────────────────────────────────
    const hasPrivateSubnets = vpc.privateSubnets.length > 0;
    this.service = new ecs.FargateService(this, 'BackendService', {
      cluster: this.cluster,
      taskDefinition,
      serviceName: `${projectName}-${environment}-backend-service`,
      desiredCount,
      assignPublicIp: !hasPrivateSubnets,
      vpcSubnets: {
        subnetType: hasPrivateSubnets ? ec2.SubnetType.PRIVATE_WITH_EGRESS : ec2.SubnetType.PUBLIC,
      },
      securityGroups: [ecsTaskSg],
      enableExecuteCommand: true,
      circuitBreaker: { enable: true, rollback: true },
      healthCheckGracePeriod: cdk.Duration.seconds(60),
    });

    // ── ALB ───────────────────────────────────────────────────────────────────
    this.loadBalancer = new elbv2.ApplicationLoadBalancer(this, 'BackendLoadBalancer', {
      vpc,
      internetFacing: true,
      loadBalancerName: `${projectName}-${environment}-backend-alb`,
    });

    const targetGroup = new elbv2.ApplicationTargetGroup(this, 'BackendTargetGroup', {
      port: containerPort,
      protocol: elbv2.ApplicationProtocol.HTTP,
      targetType: elbv2.TargetType.IP,
      vpc,
      healthCheck: {
        path: '/health',
        healthyThresholdCount: 2,
        unhealthyThresholdCount: 3,
        interval: cdk.Duration.seconds(30),
        healthyHttpCodes: '200',
      },
    });

    this.service.attachToApplicationTargetGroup(targetGroup);

    if (certificate && domainName) {
      this.loadBalancer.addListener('HttpsListener', {
        port: 443,
        protocol: elbv2.ApplicationProtocol.HTTPS,
        certificates: [certificate],
        defaultTargetGroups: [targetGroup],
      });
      this.loadBalancer.addListener('HttpRedirectListener', {
        port: 80,
        protocol: elbv2.ApplicationProtocol.HTTP,
        defaultAction: elbv2.ListenerAction.redirect({ protocol: 'HTTPS', port: '443', permanent: true }),
      });
      if (hostedZone) {
        new route53.ARecord(this, 'BackendAliasRecord', {
          zone: hostedZone,
          recordName: domainName,
          target: route53.RecordTarget.fromAlias(new route53Targets.LoadBalancerTarget(this.loadBalancer)),
        });
      }
    } else {
      this.loadBalancer.addListener('HttpListener', {
        port: 80,
        protocol: elbv2.ApplicationProtocol.HTTP,
        defaultTargetGroups: [targetGroup],
      });
    }

    this.setupOutputs(projectName, environment, domainName);
  }

  private buildEnvironmentVariables(
    props: BackendStackProps,
    environment: string,
    containerPort: number,
  ): Record<string, string> {
    const base: Record<string, string> = {
      ENV: environment,
      PORT: containerPort.toString(),
    };
    if (props.containerConfigRequests?.hasCustomEnvironmentVariables()) {
      return { ...base, ...props.containerConfigRequests.customEnvironmentVariables };
    }
    return base;
  }

  private buildSecrets(props: BackendStackProps): Record<string, ecs.Secret> {
    const secrets: Record<string, ecs.Secret> = {};
    const customSecrets = props.containerConfigRequests?.customSecretsRequests;
    if (!customSecrets?.hasSecrets()) return secrets;

    for (const config of customSecrets.secretConfigurations) {
      const secret = config.secretName
        ? secretsmanager.Secret.fromSecretNameV2(this, `Secret-${config.secretName}`, config.secretName)
        : secretsmanager.Secret.fromSecretCompleteArn(this, `Secret-${config.secretArn}`, config.secretArn!);

      for (const mapping of config.keyMappings) {
        secrets[mapping.envVarName] = ecs.Secret.fromSecretsManager(secret, mapping.secretKey);
      }
    }
    return secrets;
  }

  private setupOutputs(projectName: string, environment: string, domainName?: string) {
    new cdk.CfnOutput(this, 'RepositoryUri', {
      value: this.ecrRepository.repositoryUri,
      description: 'ECR Repository URI',
    });

    if (this.cluster) {
      new cdk.CfnOutput(this, 'ClusterName', {
        value: this.cluster.clusterName,
        description: 'ECS Cluster Name',
      });
    }

    if (this.loadBalancer) {
      new cdk.CfnOutput(this, 'LoadBalancerDNS', {
        value: this.loadBalancer.loadBalancerDnsName,
        description: 'Load Balancer DNS Name',
      });
    }

    if (domainName) {
      new cdk.CfnOutput(this, 'BackendUrl', {
        value: `https://${domainName}`,
        description: 'Backend API URL',
      });
    }
  }
}
