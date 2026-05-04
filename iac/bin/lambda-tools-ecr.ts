import * as cdk from 'aws-cdk-lib';
import { LambdaToolsECRStack } from '../lib/lambda-tools-ecr-stack';
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

  const stackName = `${projectName}-${environment}-lambda-tools-ecr`;

  new LambdaToolsECRStack(app, stackName, {
    projectName,
    environment,
    env: {
      account: getCdkDefaultAccount(),
      region: getCdkDefaultRegion(),
    },
    tags: {
      Project: projectName,
      Environment: environment,
      StackType: 'lambda-tools-ecr',
      ManagedBy: 'cdk',
    },
  });

  console.log(`✅ LambdaToolsECRStack: ${stackName}`);
})();
