#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { FrontendStack, CertificateStack } from '../lib/frontend-stack';
import {
  getEnvironment,
  getCostLevel,
  getCdkDefaultAccount,
  getCdkDefaultRegion,
  validateCostLevel,
  STACK_TYPES,
  createStackName,
  extractRootDomain,
} from '../lib/utils';

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

if (domainName) {
  console.log('Auto-detection enabled for certificate and hosted zone');
} else {
  console.log('No domain specified, using CloudFront default domain');
}

try {
  // CloudFront requires ACM certificate in us-east-1.
  // Create a dedicated CertificateStack in us-east-1 when a custom domain is provided,
  // then pass the certificate reference to the main FrontendStack via crossRegionReferences.
  let certificate;
  if (domainName) {
    const certStackName = `${stackName}-cert-us-east-1`;
    const certStack = new CertificateStack(app, certStackName, {
      domainName,
      env: {
        account: getCdkDefaultAccount(),
        region: 'us-east-1',
      },
      crossRegionReferences: true,
    });
    certificate = certStack.certificate;
    console.log(`Certificate stack created: ${certStackName} (us-east-1)`);
  }

  new FrontendStack(app, stackName, {
    projectName,
    environment,
    costLevel: costLevel as 'minimal' | 'standard' | 'high-availability',
    domainName,
    certificate,
    env: {
      account: getCdkDefaultAccount(),
      region: getCdkDefaultRegion()
    },
    crossRegionReferences: true // allows referencing the us-east-1 certificate
  });

  console.log(`✅ Successfully created ${stackName} with cost level: ${costLevel}`);
} catch (error) {
  console.error('❌ Failed to create FrontendStack:', (error as Error).message);
  process.exit(1);
}
