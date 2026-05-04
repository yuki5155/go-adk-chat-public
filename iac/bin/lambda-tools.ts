import * as cdk from 'aws-cdk-lib';
import { LambdaToolsStack } from '../lib/lambda-tools-stack';
import { getEnvironment, getCdkDefaultAccount, getCdkDefaultRegion } from '../lib/utils';

(async () => {
  const app = new cdk.App();
  const projectName = app.node.tryGetContext('projectName');
  const environment = getEnvironment(app);

  if (!projectName) {
    console.error('❌ Error: projectName is required');
    console.error('   Use: --context projectName=your-project-name');
    process.exit(1);
  }

  const stackName = `${projectName}-${environment}-lambda-tools`;

  new LambdaToolsStack(app, stackName, {
    projectName,
    environment,
    env: {
      account: getCdkDefaultAccount(),
      region: getCdkDefaultRegion(),
    },
    tags: {
      Project: projectName,
      Environment: environment,
      StackType: 'lambda-tools',
      ManagedBy: 'cdk',
    },
  });

  console.log(`✅ LambdaToolsStack: ${stackName}`);
})();
