import * as cdk from 'aws-cdk-lib';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import { Construct } from 'constructs';

export interface NetworkStackProps extends cdk.StackProps {
  projectName: string;
  environment?: string;
  costLevel?: 'minimal' | 'standard' | 'high-availability';
}

export class NetworkStack extends cdk.Stack {
  public readonly vpc: ec2.Vpc;

  constructor(scope: Construct, id: string, props: NetworkStackProps) {
    super(scope, id, props);

    const costLevel = props.costLevel || 'standard';
    const environment = props.environment || 'dev';
    const { projectName } = props;

    const vpcConfig = this.getVpcConfiguration(costLevel);

    this.vpc = new ec2.Vpc(this, 'MainVpc', {
      ipAddresses: ec2.IpAddresses.cidr('10.0.0.0/16'),
      maxAzs: vpcConfig.maxAzs,
      subnetConfiguration: vpcConfig.subnetConfiguration,
      enableDnsHostnames: true,
      enableDnsSupport: true,
      natGateways: vpcConfig.natGateways,
    });

    cdk.Tags.of(this.vpc).add('Name', `${projectName}-${environment}-Vpc`);
    cdk.Tags.of(this.vpc).add('CostLevel', costLevel);
    cdk.Tags.of(this.vpc).add('Environment', environment);
    cdk.Tags.of(this.vpc).add('Project', projectName);

    this.createVpcEndpoints();

    new cdk.CfnOutput(this, 'VpcId', {
      value: this.vpc.vpcId,
      description: 'VPC ID',
    });

    new cdk.CfnOutput(this, 'CostLevel', {
      value: costLevel,
      description: 'Cost optimization level',
    });

    new cdk.CfnOutput(this, 'Environment', {
      value: environment,
      description: 'Environment name',
    });

    new cdk.CfnOutput(this, 'PublicSubnetIds', {
      value: this.vpc.publicSubnets.map(s => s.subnetId).join(','),
      description: 'Public subnet IDs',
    });

    new cdk.CfnOutput(this, 'PrivateSubnetIds', {
      value: this.vpc.privateSubnets.map(s => s.subnetId).join(','),
      description: 'Private subnet IDs',
    });
  }

  private getVpcConfiguration(costLevel: string) {
    switch (costLevel) {
      case 'minimal':
        return {
          maxAzs: 2,
          natGateways: 0,
          subnetConfiguration: [
            {
              cidrMask: 24,
              name: 'public-subnet',
              subnetType: ec2.SubnetType.PUBLIC,
            },
          ],
        };

      case 'standard':
        return {
          maxAzs: 2,
          natGateways: 1,
          subnetConfiguration: [
            {
              cidrMask: 24,
              name: 'public-subnet',
              subnetType: ec2.SubnetType.PUBLIC,
            },
            {
              cidrMask: 24,
              name: 'private-subnet',
              subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
            },
          ],
        };

      case 'high-availability':
        return {
          maxAzs: 2,
          natGateways: 2,
          subnetConfiguration: [
            {
              cidrMask: 24,
              name: 'public-subnet',
              subnetType: ec2.SubnetType.PUBLIC,
            },
            {
              cidrMask: 24,
              name: 'private-subnet',
              subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
            },
          ],
        };

      default:
        throw new Error(`Unknown costLevel: ${costLevel}`);
    }
  }

  private createVpcEndpoints() {
    const vpcEndpointSg = new ec2.SecurityGroup(this, 'VpcEndpointSecurityGroup', {
      vpc: this.vpc,
      description: 'Security group for VPC endpoints',
      allowAllOutbound: false,
    });

    vpcEndpointSg.addIngressRule(
      ec2.Peer.ipv4(this.vpc.vpcCidrBlock),
      ec2.Port.tcp(443),
      'Allow HTTPS from VPC'
    );

    const hasPrivateSubnets = this.vpc.privateSubnets.length > 0;
    const subnetSelection = hasPrivateSubnets
      ? { subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS }
      : { subnetType: ec2.SubnetType.PUBLIC };

    this.vpc.addInterfaceEndpoint('SecretsManagerEndpoint', {
      service: ec2.InterfaceVpcEndpointAwsService.SECRETS_MANAGER,
      subnets: subnetSelection,
      securityGroups: [vpcEndpointSg],
      privateDnsEnabled: true,
    });

    this.vpc.addInterfaceEndpoint('EcrDockerEndpoint', {
      service: ec2.InterfaceVpcEndpointAwsService.ECR_DOCKER,
      subnets: subnetSelection,
      securityGroups: [vpcEndpointSg],
      privateDnsEnabled: true,
    });

    this.vpc.addInterfaceEndpoint('EcrApiEndpoint', {
      service: ec2.InterfaceVpcEndpointAwsService.ECR,
      subnets: subnetSelection,
      securityGroups: [vpcEndpointSg],
      privateDnsEnabled: true,
    });

    this.vpc.addInterfaceEndpoint('CloudWatchLogsEndpoint', {
      service: ec2.InterfaceVpcEndpointAwsService.CLOUDWATCH_LOGS,
      subnets: subnetSelection,
      securityGroups: [vpcEndpointSg],
      privateDnsEnabled: true,
    });

    this.vpc.addGatewayEndpoint('S3Endpoint', {
      service: ec2.GatewayVpcEndpointAwsService.S3,
      subnets: hasPrivateSubnets
        ? [{ subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS }]
        : [{ subnetType: ec2.SubnetType.PUBLIC }],
    });

    console.log(`VPC Endpoints created (Secrets Manager, ECR, CloudWatch Logs, S3) in ${hasPrivateSubnets ? 'private' : 'public'} subnets`);
  }
}
