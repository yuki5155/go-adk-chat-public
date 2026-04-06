import * as cdk from 'aws-cdk-lib';
import { BackendStack } from '../lib/backend-stack';
import { CloudformationSdkUtils } from '../lib/cloudformation-utils';
import { RdsRequests, ContainerConfigRequests, CustomSecretsRequests } from '../lib/requests';
import {
  getProjectName,
  getEnvironment,
  getContainerPort,
  getImageTag,
  getCpu,
  getMemory,
  getDesiredCount,
  getCdkDefaultAccount,
  getCdkDefaultRegion,
  STACK_TYPES,
  createStackName,
  createDefaultTags,
  extractRootDomain,
  getBackendDomain,
  getFrontendDomain,
} from '../lib/utils';

(async () => {
    const app = new cdk.App();
    const projectName = app.node.tryGetContext('projectName')
    const environment = getEnvironment(app);
    const rootDomain = app.node.tryGetContext('rootDomain')
    const domainName = rootDomain ? getBackendDomain(extractRootDomain(rootDomain), environment) : undefined;
    const containerPort = getContainerPort(app);
    const imageTag = getImageTag();
    const cpu = getCpu();
    const memory = getMemory();
    const desiredCount = getDesiredCount();
    const stackName = createStackName(projectName, environment, STACK_TYPES.BACKEND);
    const databaseStackName = createStackName(projectName, environment, STACK_TYPES.DATABASE);

    const databaseStack = await CloudformationSdkUtils.create(databaseStackName);
    const isDatabaseStackDeployed = databaseStack.isDeployed;
    const rdsRequests = isDatabaseStackDeployed ? RdsRequests.build(
        databaseStack.getOutputByKey('ClusterEndpoint'),
        databaseStack.getOutputByKey('ClusterPort'),
        databaseStack.getOutputByKey('DatabaseName'),
        databaseStack.getOutputByKey('SecretArn'),
        databaseStack.getOutputByKey('ClusterArn')
    ) : undefined;

    // Setup environment variables for CORS
    const frontendDomain = rootDomain ? getFrontendDomain(extractRootDomain(rootDomain), environment) : undefined;
    const frontendUrl = frontendDomain ? `https://${frontendDomain}` : 'http://localhost:5173';
    
    // Cloud environments (dev/staging/prod) use HTTPS, so always use production mode for Secure cookies
    const goEnv = domainName ? 'production' : 'development';

    // Secrets Manager configuration for Google OAuth and JWT secrets
    // Secret should be created manually in AWS Secrets Manager with the following structure:
    // {
    //   "GOOGLE_CLIENT_ID": "your-client-id.apps.googleusercontent.com",
    //   "GOOGLE_CLIENT_SECRET": "your-client-secret",
    //   "JWT_SECRET": "your-jwt-secret"
    // }
    const secretName = `${projectName}/${environment}/google-auth`;
    
    console.log('=== Backend Environment Configuration ===');
    console.log(`Frontend URL: ${frontendUrl}`);
    console.log(`Backend Domain: ${domainName || 'Not specified'}`);
    console.log(`Environment: ${environment}`);
    console.log(`GO_ENV: ${goEnv} (${domainName ? 'HTTPS/Secure cookies enabled' : 'HTTP/Secure cookies disabled'})`);
    console.log(`Secrets Manager: ${secretName}`);
    
    // Build custom secrets from AWS Secrets Manager
    const customSecretsRequests = CustomSecretsRequests.buildFromName(
      secretName,
      [
        { envVarName: 'GOOGLE_CLIENT_ID', secretKey: 'GOOGLE_CLIENT_ID' },
        { envVarName: 'GOOGLE_CLIENT_SECRET', secretKey: 'GOOGLE_CLIENT_SECRET' },
        { envVarName: 'JWT_SECRET', secretKey: 'JWT_SECRET' },
        { envVarName: 'ROOT_USER_EMAIL', secretKey: 'ROOT_USER_EMAIL' },
        { envVarName: 'GOOGLE_AI_API_KEY', secretKey: 'GOOGLE_AI_API_KEY' },
        { envVarName: 'OPENAI_API_KEY', secretKey: 'OPENAI_API_KEY' },
        { envVarName: 'ANTHROPIC_API_KEY', secretKey: 'ANTHROPIC_API_KEY' },
      ]
    );

    const containerConfigRequests = ContainerConfigRequests.build(
      {
        ALLOWED_ORIGINS: frontendUrl,
        FRONTEND_URL: frontendUrl,
        GO_ENV: goEnv,
        PORT: containerPort.toString()
      },
      customSecretsRequests
    );

    try {
        new BackendStack(app, stackName, {
          projectName,
          environment,
          domainName,
          databaseStackName,
          isDatabaseStackDeployed,
          rdsRequests,
          containerConfigRequests,
          containerPort,
          imageTag,
          cpu,
          memory,
          desiredCount,
          env: {
            account: getCdkDefaultAccount(),
            region: getCdkDefaultRegion()
          },
          tags: createDefaultTags(projectName, environment, STACK_TYPES.BACKEND, 'standard')
        });
    
        console.log(`✅ Successfully created ${stackName}`);
      } catch (error) {
        console.error('❌ Failed to create BackendStack:', (error as Error).message);
        process.exit(1);
      }
    })();