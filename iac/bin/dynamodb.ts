#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib';
import { DynamoDBStack } from '../lib/dynamodb-stack';
import {
  getEnvironment,
  getCdkDefaultAccount,
  getCdkDefaultRegion,
} from '../lib/utils';

// ============================================================================
// Main
// ============================================================================
const app = new cdk.App();
const projectName = app.node.tryGetContext('projectName');
const environment = getEnvironment(app);

if (!projectName) {
  console.error('❌ Error: projectName is required');
  console.error('   Use: --context projectName=your-project-name');
  process.exit(1);
}

const stackName = `${projectName}-${environment}-dynamodb`;

console.log('=========================');
console.log('DynamoDB Stack Deployment');
console.log('=========================');
console.log(`Project Name: ${projectName}`);
console.log(`Environment: ${environment}`);
console.log(`Stack Name: ${stackName}`);
console.log(`Region: ${getCdkDefaultRegion()}`);
console.log(`Account: ${getCdkDefaultAccount()}`);
console.log('=========================');

try {
  new DynamoDBStack(app, stackName, {
    projectName,
    environment,
    env: {
      account: getCdkDefaultAccount(),
      region: getCdkDefaultRegion()
    },
    tags: {
      Project: projectName,
      Environment: environment,
      StackType: 'dynamodb',
      CostLevel: 'standard',
      ManagedBy: 'cdk'
    }
  });

  console.log(`✅ Successfully created ${stackName}`);
  console.log('');
  console.log('To deploy this stack:');
  console.log(`  npx cdk deploy ${stackName} \\`);
  console.log(`    --context projectName=${projectName} \\`);
  console.log(`    --context environment=${environment}`);
  console.log('');
  console.log('Or use the Makefile:');
  console.log(`  make dynamodb-deploy ENV=${environment}`);
} catch (error) {
  console.error('❌ Failed to create DynamoDBStack:', error instanceof Error ? error.message : String(error));
  process.exit(1);
}
