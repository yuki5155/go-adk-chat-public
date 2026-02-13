import * as cdk from 'aws-cdk-lib';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as logs from 'aws-cdk-lib/aws-logs';
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as route53Targets from 'aws-cdk-lib/aws-route53-targets';
import * as acm from 'aws-cdk-lib/aws-certificatemanager';
import * as path from 'path';
import {
  getEnvironment,
  getCdkDefaultAccount,
  getCdkDefaultRegion,
  extractRootDomain,
  getFrontendDomain
} from 'automation-deploy-template-iac';

// Lambda function configuration
interface LambdaConfig {
  name: string;
  path: string;
  method: string;
  description: string;
  requiresAuth?: boolean;
  streaming?: boolean; // Enable response streaming for SSE
}

(async () => {
  const app = new cdk.App();
  const projectName = app.node.tryGetContext('projectName');
  const environment = getEnvironment(app);
  const rootDomain = app.node.tryGetContext('rootDomain');
  const subdomain = app.node.tryGetContext('subdomain') || 'lambda';  // Default: lambda

  // Validate required parameters
  if (!projectName) {
    console.error('❌ Error: projectName is required');
    console.error('   Use: --context projectName=your-project-name');
    process.exit(1);
  }

  if (!rootDomain) {
    console.error('❌ Error: rootDomain is required (custom domain is mandatory)');
    console.error('   Use: --context rootDomain=yourdomain.com');
    process.exit(1);
  }

  // Construct domain name with custom subdomain (e.g., lambda.dev.yourdomain.com)
  const domain = extractRootDomain(rootDomain);
  const domainName = environment === 'prod'
    ? `${subdomain}.${domain}`
    : `${subdomain}.${environment}.${domain}`;

  const memory = parseInt(app.node.tryGetContext('memory') || '512', 10);
  const timeout = parseInt(app.node.tryGetContext('timeout') || '120', 10); // 120s for chat AI responses
  const stackName = `${projectName}-${environment}-lambda`;

  // Path to Lambda build directory (ZIP files)
  const lambdaBuildPath = path.join(__dirname, '../../backend/build/lambda');

  // Setup environment variables for CORS
  const frontendDomain = rootDomain ? getFrontendDomain(extractRootDomain(rootDomain), environment) : undefined;
  const frontendUrl = frontendDomain ? `https://${frontendDomain}` : 'http://localhost:5173';

  // Cloud environments (dev/staging/prod) use HTTPS, so always use production mode for Secure cookies
  const goEnv = domainName ? 'production' : 'development';

  // Secrets Manager configuration (reuse existing GoogleAuthSecretsStack)
  const secretName = `${projectName}/${environment}/google-auth`;

  // Define all Lambda functions
  const lambdaConfigs: LambdaConfig[] = [
    { name: 'auth-google', path: '/auth/google', method: 'POST', description: 'Google Sign-In' },
    { name: 'auth-refresh', path: '/auth/refresh', method: 'POST', description: 'Token Refresh' },
    { name: 'auth-logout', path: '/auth/logout', method: 'POST', description: 'User Logout' },
    { name: 'get-user', path: '/api/me', method: 'GET', description: 'Get Current User', requiresAuth: true },
    { name: 'admin-dashboard', path: '/admin/dashboard', method: 'GET', description: 'Admin Dashboard', requiresAuth: true },
    { name: 'role-request', path: '/api/role/request', method: 'POST', description: 'Request Role Subscription', requiresAuth: true },
    { name: 'admin-role-requests', path: '/api/admin/role/requests', method: 'GET', description: 'List Pending Role Requests', requiresAuth: true },
    { name: 'admin-role-approve', path: '/api/admin/role/approve', method: 'POST', description: 'Approve Role Request', requiresAuth: true },
    { name: 'admin-role-reject', path: '/api/admin/role/reject', method: 'POST', description: 'Reject Role Request', requiresAuth: true },
    { name: 'admin-users', path: '/api/admin/role/users', method: 'GET', description: 'List Users by Role', requiresAuth: true },
    { name: 'health', path: '/health', method: 'GET', description: 'Health Check' },
    { name: 'hello', path: '/hello', method: 'GET', description: 'Hello Endpoint' },
    // Chat Lambda functions
    { name: 'chat-models', path: '/api/chat/models', method: 'GET', description: 'List AI Models', requiresAuth: true },
    { name: 'chat-threads-create', path: '/api/chat/threads', method: 'POST', description: 'Create Chat Thread', requiresAuth: true },
    { name: 'chat-threads-list', path: '/api/chat/threads', method: 'GET', description: 'List Chat Threads', requiresAuth: true },
    { name: 'chat-threads-get', path: '/api/chat/threads/{id}', method: 'GET', description: 'Get Chat Thread', requiresAuth: true },
    { name: 'chat-threads-delete', path: '/api/chat/threads/{id}', method: 'DELETE', description: 'Delete Chat Thread', requiresAuth: true },
    { name: 'chat-message', path: '/api/chat/threads/{id}/message', method: 'POST', description: 'Send Chat Message', requiresAuth: true },
    { name: 'chat-stream', path: '/api/chat/threads/{id}/stream', method: 'POST', description: 'Stream Chat Message', requiresAuth: true, streaming: true },
  ];

  console.log('=== Lambda Backend Configuration ===');
  console.log(`Deployment Type: ZIP (Go binaries)`);
  console.log(`Build Path: ${lambdaBuildPath}`);
  console.log(`Project Name: ${projectName}`);
  console.log(`Environment: ${environment}`);
  console.log(`Root Domain: ${rootDomain}`);
  console.log(`Subdomain: ${subdomain}`);
  console.log(`API Domain: ${domainName}`);
  console.log(`Frontend URL: ${frontendUrl}`);
  console.log(`GO_ENV: ${goEnv}`);
  console.log(`Secrets Manager: ${secretName}`);
  console.log(`Lambda Memory: ${memory} MB`);
  console.log(`Lambda Timeout: ${timeout} seconds`);
  console.log(`Lambda Functions: ${lambdaConfigs.length}`);

  try {
    const stack = new cdk.Stack(app, stackName, {
      env: {
        account: getCdkDefaultAccount(),
        region: getCdkDefaultRegion()
      },
      tags: {
        Project: projectName,
        Environment: environment,
        StackType: 'lambda',
        CostLevel: 'standard',
        ManagedBy: 'cdk'
      }
    });

    // Import secret from Secrets Manager
    const secret = secretsmanager.Secret.fromSecretNameV2(
      stack,
      'GoogleAuthSecret',
      secretName
    );

    // Lambda execution role with necessary permissions
    const lambdaRole = new iam.Role(stack, 'LambdaExecutionRole', {
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
      managedPolicies: [
        iam.ManagedPolicy.fromAwsManagedPolicyName('service-role/AWSLambdaBasicExecutionRole')
      ]
    });

    // Grant Lambda permission to read secrets
    secret.grantRead(lambdaRole);

    // Grant Lambda permission to access DynamoDB tables
    // Table names follow the pattern: projectName-environment-table-name
    const userRolesTableName = `${projectName}-${environment}-user-roles`;
    const roleRequestsTableName = `${projectName}-${environment}-role-requests`;
    const chatThreadsTableName = `${projectName}-${environment}-chat-threads`;
    const chatSessionsTableName = `${projectName}-${environment}-chat-sessions`;
    const chatEventsTableName = `${projectName}-${environment}-chat-events`;
    const chatMemoriesTableName = `${projectName}-${environment}-chat-memories`;

    // Grant read/write permissions to DynamoDB tables
    lambdaRole.addToPolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      actions: [
        'dynamodb:GetItem',
        'dynamodb:PutItem',
        'dynamodb:UpdateItem',
        'dynamodb:DeleteItem',
        'dynamodb:Query',
        'dynamodb:Scan',
        'dynamodb:BatchGetItem',
        'dynamodb:BatchWriteItem'
      ],
      resources: [
        `arn:aws:dynamodb:${getCdkDefaultRegion()}:${getCdkDefaultAccount()}:table/${userRolesTableName}`,
        `arn:aws:dynamodb:${getCdkDefaultRegion()}:${getCdkDefaultAccount()}:table/${userRolesTableName}/index/*`,
        `arn:aws:dynamodb:${getCdkDefaultRegion()}:${getCdkDefaultAccount()}:table/${roleRequestsTableName}`,
        `arn:aws:dynamodb:${getCdkDefaultRegion()}:${getCdkDefaultAccount()}:table/${roleRequestsTableName}/index/*`,
        `arn:aws:dynamodb:${getCdkDefaultRegion()}:${getCdkDefaultAccount()}:table/${chatThreadsTableName}`,
        `arn:aws:dynamodb:${getCdkDefaultRegion()}:${getCdkDefaultAccount()}:table/${chatThreadsTableName}/index/*`,
        `arn:aws:dynamodb:${getCdkDefaultRegion()}:${getCdkDefaultAccount()}:table/${chatSessionsTableName}`,
        `arn:aws:dynamodb:${getCdkDefaultRegion()}:${getCdkDefaultAccount()}:table/${chatSessionsTableName}/index/*`,
        `arn:aws:dynamodb:${getCdkDefaultRegion()}:${getCdkDefaultAccount()}:table/${chatEventsTableName}`,
        `arn:aws:dynamodb:${getCdkDefaultRegion()}:${getCdkDefaultAccount()}:table/${chatEventsTableName}/index/*`,
        `arn:aws:dynamodb:${getCdkDefaultRegion()}:${getCdkDefaultAccount()}:table/${chatMemoriesTableName}`,
        `arn:aws:dynamodb:${getCdkDefaultRegion()}:${getCdkDefaultAccount()}:table/${chatMemoriesTableName}/index/*`
      ]
    }));

    // Environment variables for all Lambda functions
    const lambdaEnvironment = {
      PORT: '8080',
      GO_ENV: goEnv,
      ALLOWED_ORIGINS: frontendUrl,
      FRONTEND_URL: frontendUrl,
      COOKIE_DOMAIN: `.${domain}`, // e.g., ".mydevportal.com" for cross-subdomain cookies
      GOOGLE_CLIENT_ID: secret.secretValueFromJson('GOOGLE_CLIENT_ID').unsafeUnwrap(),
      GOOGLE_CLIENT_SECRET: secret.secretValueFromJson('GOOGLE_CLIENT_SECRET').unsafeUnwrap(),
      JWT_SECRET: secret.secretValueFromJson('JWT_SECRET').unsafeUnwrap(),
      ROOT_USER_EMAIL: secret.secretValueFromJson('ROOT_USER_EMAIL').unsafeUnwrap(),
      DYNAMODB_USER_ROLES_TABLE: userRolesTableName,
      DYNAMODB_ROLE_REQUESTS_TABLE: roleRequestsTableName,
      // Chat tables
      DYNAMODB_CHAT_THREADS_TABLE: chatThreadsTableName,
      DYNAMODB_CHAT_SESSIONS_TABLE: chatSessionsTableName,
      DYNAMODB_CHAT_EVENTS_TABLE: chatEventsTableName,
      DYNAMODB_CHAT_MEMORIES_TABLE: chatMemoriesTableName,
      // Gemini API
      GOOGLE_AI_API_KEY: secret.secretValueFromJson('GOOGLE_AI_API_KEY').unsafeUnwrap(),
      GEMINI_MODEL: 'gemini-2.0-flash',
    };

    // Create Lambda functions
    const lambdaFunctions = new Map<string, lambda.Function>();

    for (const config of lambdaConfigs) {
      // CloudWatch Logs group
      const logGroup = new logs.LogGroup(stack, `${config.name}LogGroup`, {
        logGroupName: `/aws/lambda/${projectName}-${environment}-lambda-${config.name}`,
        retention: logs.RetentionDays.ONE_MONTH,
        removalPolicy: cdk.RemovalPolicy.DESTROY
      });

      // Path to the ZIP file for this Lambda function
      const functionZipPath = path.join(lambdaBuildPath, `${config.name}.zip`);

      // Create Lambda Function using ZIP file
      const lambdaFunction = new lambda.Function(stack, `${config.name}Function`, {
        functionName: `${projectName}-${environment}-lambda-${config.name}`,
        runtime: lambda.Runtime.PROVIDED_AL2023,  // Custom runtime for Go
        handler: 'bootstrap',  // Go binary name
        code: lambda.Code.fromAsset(functionZipPath),
        memorySize: memory,
        timeout: cdk.Duration.seconds(timeout),
        role: lambdaRole,
        logGroup: logGroup,
        environment: lambdaEnvironment,
        description: `${projectName} ${config.description} (${environment})`,
        architecture: lambda.Architecture.X86_64
      });

      lambdaFunctions.set(config.name, lambdaFunction);
      console.log(`✓ Created Lambda function: ${config.name} from ${functionZipPath}`);
    }

    // REST API Gateway
    const api = new apigateway.RestApi(stack, 'RestApi', {
      restApiName: `${projectName}-${environment}-api`,
      description: `REST API for ${projectName} Lambda backend (${environment})`,
      defaultCorsPreflightOptions: {
        allowOrigins: [frontendUrl],
        allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
        allowHeaders: [
          'Content-Type',
          'Content-Length',
          'Accept-Encoding',
          'X-CSRF-Token',
          'Authorization',
          'accept',
          'origin',
          'Cache-Control',
          'X-Requested-With'
        ],
        allowCredentials: true
      },
      deployOptions: {
        stageName: environment,
        loggingLevel: apigateway.MethodLoggingLevel.INFO,
        dataTraceEnabled: true,
        metricsEnabled: true
      },
      endpointConfiguration: {
        types: [apigateway.EndpointType.REGIONAL]
      }
    });

    // Helper function to get or create API Gateway resource path
    const getOrCreateResource = (api: apigateway.RestApi, path: string): apigateway.IResource => {
      const parts = path.split('/').filter(p => p);
      let resource: apigateway.IResource = api.root;

      for (const part of parts) {
        const existing = resource.getResource(part);
        if (existing) {
          resource = existing;
        } else {
          resource = resource.addResource(part);
        }
      }

      return resource;
    };

    // Wire up each Lambda function to its API Gateway route
    for (const config of lambdaConfigs) {
      const lambdaFunction = lambdaFunctions.get(config.name)!;

      // Configure Lambda integration with optional response streaming
      const lambdaIntegration = new apigateway.LambdaIntegration(lambdaFunction, {
        proxy: true,
        allowTestInvoke: true,
        ...(config.streaming && { responseTransferMode: apigateway.ResponseTransferMode.STREAM }),
      });

      const resource = getOrCreateResource(api, config.path);
      resource.addMethod(config.method, lambdaIntegration);

      console.log(`✓ Mapped ${config.method} ${config.path} → ${config.name}${config.streaming ? ' (streaming)' : ''}`);
    }

    // Custom Domain Configuration (REQUIRED)
    console.log(`Setting up custom domain: ${domainName}`);

    // Lookup existing hosted zone (not import - just lookup)
    const hostedZone = route53.HostedZone.fromLookup(stack, 'HostedZone', {
      domainName: domain
    });

    // Create ACM certificate for custom domain
    const certificate = new acm.Certificate(stack, 'ApiCertificate', {
      domainName: domainName,
      validation: acm.CertificateValidation.fromDns(hostedZone)
    });

    // Create custom domain for API Gateway
    const customDomain = new apigateway.DomainName(stack, 'CustomDomain', {
      domainName: domainName,
      certificate: certificate,
      endpointType: apigateway.EndpointType.REGIONAL,
      securityPolicy: apigateway.SecurityPolicy.TLS_1_2
    });

    // Map custom domain to API Gateway
    new apigateway.BasePathMapping(stack, 'BasePathMapping', {
      domainName: customDomain,
      restApi: api,
      stage: api.deploymentStage
    });

    // Create Route53 A record
    new route53.ARecord(stack, 'ApiAliasRecord', {
      zone: hostedZone,
      recordName: domainName,
      target: route53.RecordTarget.fromAlias(
        new route53Targets.ApiGatewayDomain(customDomain)
      )
    });

    console.log(`✓ Custom domain configured: https://${domainName}`);

    // Stack Outputs
    new cdk.CfnOutput(stack, 'ApiGatewayUrl', {
      value: api.url,
      description: 'API Gateway URL'
    });

    new cdk.CfnOutput(stack, 'ApiGatewayId', {
      value: api.restApiId,
      description: 'API Gateway ID'
    });

    new cdk.CfnOutput(stack, 'LambdaFunctionCount', {
      value: lambdaFunctions.size.toString(),
      description: 'Number of Lambda Functions'
    });

    // Output each Lambda function ARN
    for (const [name, fn] of lambdaFunctions.entries()) {
      new cdk.CfnOutput(stack, `Lambda${name.replace(/-/g, '')}Arn`, {
        value: fn.functionArn,
        description: `${name} Lambda Function ARN`
      });
    }

    // Always output custom domain (it's required)
    new cdk.CfnOutput(stack, 'CustomDomainUrl', {
      value: `https://${domainName}`,
      description: 'Custom Domain URL'
    });

    console.log(`✅ Successfully created ${stackName}`);
    console.log(`✅ Deployed ${lambdaFunctions.size} Lambda functions`);
  } catch (error) {
    console.error('❌ Failed to create Lambda Stack:', (error as Error).message);
    process.exit(1);
  }
})();
