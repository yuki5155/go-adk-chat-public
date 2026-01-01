#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import {
  FrontendStack,
  getEnvironment,
  getCostLevel,
  getCdkDefaultAccount,
  getCdkDefaultRegion,
  validateCostLevel,
  STACK_TYPES,
  createStackName,
  extractRootDomain
} from 'automation-deploy-template-iac';

const app = new cdk.App();

// Get project name from CDK context (passed via --context projectName=...)
const projectName = app.node.tryGetContext('projectName');
const environment = getEnvironment(app);
const costLevel = getCostLevel(app);
const rootDomain = app.node.tryGetContext('rootDomain') || process.env.ROOT_DOMAIN;
const domainName = rootDomain ? `${environment}.${extractRootDomain(rootDomain)}` : undefined;

const stackName = createStackName(projectName, environment, STACK_TYPES.FRONTEND);

// Validation
validateCostLevel(costLevel);

console.log('=== Deploying Frontend Stack ===');
console.log(`Project Name: ${projectName}`);
console.log(`Environment: ${environment}`);
console.log(`Cost Level: ${costLevel}`);
console.log(`Stack Name: ${stackName}`);
console.log(`Root Domain: ${rootDomain || 'Not specified'}`);
console.log(`Frontend Domain: ${domainName || 'Not specified (CloudFront default domain will be used)'}`);

// Auto-detection method
let autoDetectionMethod: string;
if (domainName) {
  autoDetectionMethod = 'auto-detect';
  console.log('🔍 Auto-detection enabled for certificate and hosted zone');
  console.log('⚠️ Certificate ARN and Hosted Zone ID will be auto-detected at deployment time');
  console.log('⚠️ If auto-detection fails, the deployment will use CloudFront default domain');
} else {
  autoDetectionMethod = 'default-domain';
  console.log('ℹ️ No domain specified, using CloudFront default domain');
}

try {
  // Create Frontend Stack using the published npm package
  new FrontendStack(app, stackName, {
    projectName,
    environment,
    costLevel: costLevel as 'minimal' | 'standard' | 'high-availability',
    domainName,
    autoDetectResources: !!domainName,
    env: {
      account: getCdkDefaultAccount(),
      region: getCdkDefaultRegion()
    },
    crossRegionReferences: true // us-east-1 certificate reference
  });

  console.log(`✅ Successfully created ${stackName} with cost level: ${costLevel}`);
} catch (error) {
  console.error('❌ Failed to create FrontendStack:', (error as Error).message);
  process.exit(1);
}
