import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { 
  NetworkStack,
  getEnvironment,
  getCostLevel,
  getCdkDefaultAccount,
  getCdkDefaultRegion,
  validateCostLevel,
  STACK_TYPES,
  createStackName,
  createDefaultTags
} from 'automation-deploy-template-iac';

const app = new cdk.App();

// Use helper functions from the published package
const projectName = app.node.tryGetContext('projectName');
const environment = getEnvironment(app);
const costLevel = getCostLevel(app);
const stackName = createStackName(projectName, environment, STACK_TYPES.NETWORK);

// Validation
validateCostLevel(costLevel);

console.log('=== Deploying Network Stack from NPM Package ===');
console.log(`Project Name: ${projectName}`);
console.log(`Environment: ${environment}`);
console.log(`Cost Level: ${costLevel}`);
console.log(`Stack Name: ${stackName}`);

try {
    new NetworkStack(app, stackName, {
        projectName,
        costLevel: costLevel as "minimal" | "standard" | "high-availability",
        environment,
        env: {
            account: getCdkDefaultAccount(),
            region: getCdkDefaultRegion()
        },
        tags: createDefaultTags(projectName, environment, STACK_TYPES.NETWORK, costLevel)
    });
    console.log('Network Stack deployed successfully');
} catch (error) {
    console.error('Error deploying Network Stack:', error);
    process.exit(1);
}