import { CloudFormationClient, DescribeStacksCommand, Stack } from '@aws-sdk/client-cloudformation';

export class CloudformationSdkUtils {
  public readonly client: CloudFormationClient;
  private stack?: Stack;
  private outputs: Record<string, string>;
  public isDeployed: boolean;

  private constructor(public readonly stackName: string) {
    this.client = new CloudFormationClient({});
    this.isDeployed = false;
    this.outputs = {};
  }

  public static async create(stackName: string): Promise<CloudformationSdkUtils> {
    const instance = new CloudformationSdkUtils(stackName);
    await instance.initialize();
    return instance;
  }

  private async initialize() {
    try {
      const command = new DescribeStacksCommand({ StackName: this.stackName });
      const response = await this.client.send(command);
      this.stack = response.Stacks?.[0];
      if (!this.stack) {
        this.isDeployed = false;
        return;
      }
      const validStatuses = ['CREATE_COMPLETE', 'UPDATE_COMPLETE'];
      this.isDeployed = validStatuses.includes(this.stack?.StackStatus || '');
      this.outputs = this.stack?.Outputs?.reduce((acc, output) => {
        acc[output.OutputKey || ''] = output.OutputValue || '';
        return acc;
      }, {} as Record<string, string>) || {};
    } catch (error) {
      console.log(`Stack '${this.stackName}' not found or error occurred. Proceeding without it.`);
      this.isDeployed = false;
    }
  }

  public getOutputByKey(key: string): string | null {
    if (!this.isDeployed) return null;
    return this.outputs[key] || null;
  }
}
