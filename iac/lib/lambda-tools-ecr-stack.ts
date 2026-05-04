import * as cdk from 'aws-cdk-lib';
import * as ecr from 'aws-cdk-lib/aws-ecr';
import * as fs from 'fs';
import * as path from 'path';
import * as yaml from 'js-yaml';
import { Construct } from 'constructs';

interface ToolConfig {
  name: string;
  type: 'zip' | 'container';
}

interface ToolsManifest {
  tools: ToolConfig[];
}

export interface LambdaToolsECRStackProps extends cdk.StackProps {
  projectName: string;
  environment: string;
  toolsManifestPath?: string;
}

export class LambdaToolsECRStack extends cdk.Stack {
  public readonly repoUris: Record<string, string> = {};

  constructor(scope: Construct, id: string, props: LambdaToolsECRStackProps) {
    super(scope, id, props);

    const { projectName, environment } = props;

    const manifestPath = props.toolsManifestPath
      ?? path.join(__dirname, '../../lambda-tools/tools.yaml');

    const manifest = yaml.load(fs.readFileSync(manifestPath, 'utf8')) as ToolsManifest;

    for (const tool of manifest.tools) {
      if (tool.type !== 'container') continue;

      const id = tool.name.replace(/-/g, '');
      const repoName = `${projectName}-${environment}-tools-${tool.name}`;

      const repo = new ecr.Repository(this, `${id}Repo`, {
        repositoryName: repoName,
        removalPolicy: environment === 'prod' ? cdk.RemovalPolicy.RETAIN : cdk.RemovalPolicy.DESTROY,
        emptyOnDelete: environment !== 'prod',
      });

      this.repoUris[tool.name] = repo.repositoryUri;

      new cdk.CfnOutput(this, `${id}RepoUri`, {
        value: repo.repositoryUri,
        exportName: `${projectName}-${environment}-ecr-${tool.name}`,
        description: `ECR repository URI for ${tool.name}`,
      });
    }
  }
}
