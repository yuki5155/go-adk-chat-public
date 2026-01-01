#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib';
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager';
import { Construct } from 'constructs';
import {
  getEnvironment,
  getCdkDefaultAccount,
  getCdkDefaultRegion,
  STACK_TYPES,
  createStackName,
  createDefaultTags
} from 'automation-deploy-template-iac';

/**
 * Google Auth Secrets Stack Props
 */
interface GoogleAuthSecretsStackProps extends cdk.StackProps {
  projectName: string;
  environment: string;
}

/**
 * Google Auth Secrets Stack
 * Creates secrets in AWS Secrets Manager for Google OAuth and JWT configuration
 */
class GoogleAuthSecretsStack extends cdk.Stack {
  public readonly googleAuthSecret: secretsmanager.Secret;
  public readonly secretName: string;
  public readonly secretArn: string;

  constructor(scope: Construct, id: string, props: GoogleAuthSecretsStackProps) {
    super(scope, id, props);

    const { projectName, environment } = props;

    // Create secret name following project naming convention
    // Format: {projectName}/{environment}/google-auth
    this.secretName = `${projectName}/${environment}/google-auth`;

    // ============================================================================
    // Google Auth Secret
    // Contains Google OAuth credentials and JWT secret
    // ============================================================================
    this.googleAuthSecret = new secretsmanager.Secret(this, 'GoogleAuthSecret', {
      secretName: this.secretName,
      description: `Google OAuth and JWT secrets for ${projectName} ${environment} environment`,
      secretStringValue: cdk.SecretValue.unsafePlainText(JSON.stringify({
        // Google OAuth credentials - REPLACE THESE VALUES after deployment
        GOOGLE_CLIENT_ID: 'your-client-id.apps.googleusercontent.com',
        GOOGLE_CLIENT_SECRET: 'your-client-secret',
        
        // JWT Secret for token signing - REPLACE THIS VALUE after deployment
        JWT_SECRET: 'replace-with-secure-random-secret',
      })),
      removalPolicy: cdk.RemovalPolicy.RETAIN, // Keep secrets on stack deletion for safety
    });

    this.secretArn = this.googleAuthSecret.secretArn;

    // Console output
    console.log('=========================');
    console.log('✅ Google Auth Secrets Stack');
    console.log('=========================');
    console.log(`Project Name: ${projectName}`);
    console.log(`Environment: ${environment}`);
    console.log(`Secret Name: ${this.secretName}`);
    console.log('');
    console.log('📝 Secret Keys:');
    console.log('   - GOOGLE_CLIENT_ID');
    console.log('   - GOOGLE_CLIENT_SECRET');
    console.log('   - JWT_SECRET');
    console.log('');
    console.log('⚠️  IMPORTANT: After deployment, update the secret values:');
    console.log(`   aws secretsmanager put-secret-value \\`);
    console.log(`     --secret-id "${this.secretName}" \\`);
    console.log(`     --secret-string '{"GOOGLE_CLIENT_ID":"...","GOOGLE_CLIENT_SECRET":"...","JWT_SECRET":"..."}'`);
    console.log('=========================');
  }

  /**
   * Get the stack name for secrets stack
   */
  static getStackName(projectName: string, environment: string): string {
    return createStackName(projectName, environment, STACK_TYPES.SECRETS);
  }
}

// ============================================================================
// Main
// ============================================================================
const app = new cdk.App();
const projectName = app.node.tryGetContext('projectName');
const environment = getEnvironment(app);

if (!projectName) {
  console.error('❌ Error: projectName is required. Use --context projectName=your-project-name');
  process.exit(1);
}

const stackName = GoogleAuthSecretsStack.getStackName(projectName, environment);

console.log('=========================');
console.log('Secrets Stack Deployment');
console.log('=========================');
console.log(`Project Name: ${projectName}`);
console.log(`Environment: ${environment}`);
console.log(`Stack Name: ${stackName}`);
console.log('=========================');

try {
  new GoogleAuthSecretsStack(app, stackName, {
    projectName,
    environment,
    env: {
      account: getCdkDefaultAccount(),
      region: getCdkDefaultRegion()
    },
    tags: createDefaultTags(projectName, environment, STACK_TYPES.SECRETS, 'standard')
  });

  console.log(`✅ Successfully created ${stackName}`);
  console.log('');
  console.log('To deploy this stack:');
  console.log(`  npx cdk deploy ${stackName} --context projectName=${projectName} --context environment=${environment}`);
} catch (error) {
  console.error('❌ Failed to create SecretsStack:', error instanceof Error ? error.message : String(error));
  process.exit(1);
}
