import * as cdk from 'aws-cdk-lib';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as ecr from 'aws-cdk-lib/aws-ecr';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as logs from 'aws-cdk-lib/aws-logs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager';
import * as fs from 'fs';
import * as path from 'path';
import * as yaml from 'js-yaml';
import { Construct } from 'constructs';

interface ToolConfig {
  name: string;
  route: string;
  method: string;
  type: 'zip' | 'container';
  memory?: number;
  timeout?: number;
  s3_access?: boolean;
  environment?: Record<string, string>;
}

interface ToolsManifest {
  tools: ToolConfig[];
}

export interface LambdaToolsStackProps extends cdk.StackProps {
  projectName: string;
  environment: string;
  toolsManifestPath?: string;
}

export class LambdaToolsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: LambdaToolsStackProps) {
    super(scope, id, props);

    const { projectName, environment } = props;

    const manifestPath = props.toolsManifestPath
      ?? path.join(__dirname, '../../lambda-tools/tools.yaml');

    const manifest = yaml.load(fs.readFileSync(manifestPath, 'utf8')) as ToolsManifest;
    const toolConfigs = manifest.tools;

    const needsS3 = toolConfigs.some(t => t.s3_access);

    // ── S3 bucket (created only if at least one tool requires it) ────────────
    let videosBucket: s3.Bucket | undefined;
    if (needsS3) {
      videosBucket = new s3.Bucket(this, 'VideosBucket', {
        bucketName: `${projectName}-${environment}-lambda-tools-videos`,
        removalPolicy: environment === 'prod' ? cdk.RemovalPolicy.RETAIN : cdk.RemovalPolicy.DESTROY,
        autoDeleteObjects: environment !== 'prod',
        lifecycleRules: [{ expiration: cdk.Duration.days(1), prefix: 'combined/' }],
      });
    }

    // ── Shared IAM role ───────────────────────────────────────────────────────
    const lambdaRole = new iam.Role(this, 'ToolsLambdaRole', {
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
      managedPolicies: [
        iam.ManagedPolicy.fromAwsManagedPolicyName('service-role/AWSLambdaBasicExecutionRole'),
      ],
    });

    if (videosBucket) {
      videosBucket.grantReadWrite(lambdaRole);
    }

    // ── API Gateway ───────────────────────────────────────────────────────────
    const api = new apigateway.RestApi(this, 'ToolsApi', {
      restApiName: `${projectName}-${environment}-lambda-tools`,
      description: `Lambda tools API for ${projectName} (${environment})`,
      apiKeySourceType: apigateway.ApiKeySourceType.HEADER,
      deployOptions: {
        stageName: environment,
        loggingLevel: apigateway.MethodLoggingLevel.INFO,
        metricsEnabled: true,
      },
    });

    const apiKey = new apigateway.ApiKey(this, 'ToolsApiKey', {
      apiKeyName: `${projectName}-${environment}-lambda-tools-key`,
      description: `API key for lambda tools (${environment})`,
    });

    const usagePlan = new apigateway.UsagePlan(this, 'ToolsUsagePlan', {
      name: `${projectName}-${environment}-lambda-tools-plan`,
      apiStages: [{ api, stage: api.deploymentStage }],
    });
    usagePlan.addApiKey(apiKey);

    // Store the API key value in Secrets Manager so the Go app can retrieve it
    new secretsmanager.Secret(this, 'ToolsApiKeySecret', {
      secretName: `${projectName}/${environment}/lambda-tools-api-key`,
      description: `API key for lambda tools (${environment})`,
    });

    // ── Create a Lambda + API Gateway route for each tool in the manifest ────
    for (const tool of toolConfigs) {
      const id = tool.name.replace(/-/g, '');
      const toolSourcePath = path.join(__dirname, `../../lambda-tools/${tool.name}`);

      const logGroup = new logs.LogGroup(this, `${id}LogGroup`, {
        logGroupName: `/aws/lambda/${projectName}-${environment}-tools-${tool.name}`,
        retention: logs.RetentionDays.ONE_MONTH,
        removalPolicy: cdk.RemovalPolicy.DESTROY,
      });

      const baseProps = {
        functionName: `${projectName}-${environment}-tools-${tool.name}`,
        memorySize: tool.memory ?? 256,
        timeout: cdk.Duration.seconds(tool.timeout ?? 30),
        role: lambdaRole,
        logGroup,
        description: `${tool.name} tool for ${projectName} (${environment})`,
        environment: {
          ...(tool.s3_access && videosBucket ? { TOOLS_S3_BUCKET: videosBucket.bucketName } : {}),
          ...(tool.environment ?? {}),
        },
      };

      let fn: lambda.IFunction;

      if (tool.type === 'zip') {
        fn = new lambda.Function(this, `${id}Function`, {
          ...baseProps,
          runtime: lambda.Runtime.PYTHON_3_12,
          handler: 'handler.lambda_handler',
          code: lambda.Code.fromAsset(toolSourcePath),
        });
      } else {
        const repoName = `${projectName}-${environment}-tools-${tool.name}`;
        const repo = ecr.Repository.fromRepositoryName(this, `${id}Repo`, repoName);

        fn = new lambda.DockerImageFunction(this, `${id}Function`, {
          ...baseProps,
          code: lambda.DockerImageCode.fromEcr(repo),
        });

        new cdk.CfnOutput(this, `${id}EcrRepo`, {
          value: repo.repositoryUri,
          description: `ECR repository URI for ${tool.name}`,
        });
      }

      const resource = api.root.addResource(tool.route.replace(/^\//, ''));
      resource.addMethod(
        tool.method,
        new apigateway.LambdaIntegration(fn),
        { apiKeyRequired: true },
      );

      new cdk.CfnOutput(this, `${id}Url`, {
        value: `${api.url}${tool.route.replace(/^\//, '')}`,
        description: `${tool.name} tool endpoint`,
      });

      console.log(`✓ Registered tool: ${tool.method} ${tool.route} → ${tool.name} (${tool.type})`);
    }

    // ── Stack outputs ─────────────────────────────────────────────────────────
    new cdk.CfnOutput(this, 'ToolsApiUrl', {
      value: api.url,
      description: 'Lambda tools API Gateway base URL',
    });

    if (videosBucket) {
      new cdk.CfnOutput(this, 'VideosBucketName', {
        value: videosBucket.bucketName,
        description: 'S3 bucket for tool video uploads/outputs',
      });
    }
  }
}
