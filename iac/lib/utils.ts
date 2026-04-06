import * as cdk from 'aws-cdk-lib';

export const STACK_TYPES = {
  NETWORK: 'Network',
  DATABASE: 'Database',
  BACKEND: 'Backend',
  FRONTEND: 'Frontend',
  SECRETS: 'Secrets',
} as const;

export type StackType = typeof STACK_TYPES[keyof typeof STACK_TYPES];

export function getEnvironment(app?: cdk.App, defaultEnv: string = 'dev'): string {
  if (app) {
    return app.node.tryGetContext('environment') || process.env.ENVIRONMENT || defaultEnv;
  }
  return process.env.ENVIRONMENT || defaultEnv;
}

export function getCostLevel(app?: cdk.App): string {
  if (app) {
    return app.node.tryGetContext('costLevel') || process.env.COST_LEVEL || 'standard';
  }
  return process.env.COST_LEVEL || 'standard';
}

export function getCdkDefaultAccount(): string | undefined {
  return process.env.CDK_DEFAULT_ACCOUNT;
}

export function getCdkDefaultRegion(defaultRegion: string = 'ap-northeast-1'): string {
  return process.env.CDK_DEFAULT_REGION || defaultRegion;
}

export const VALID_COST_LEVELS = ['minimal', 'standard', 'high-availability'] as const;

export function validateCostLevel(costLevel: string): void {
  if (!VALID_COST_LEVELS.includes(costLevel as any)) {
    throw new Error(`Invalid costLevel: ${costLevel}. Valid options are: ${VALID_COST_LEVELS.join(', ')}`);
  }
}

export const createStackName = (projectName: string, environment: string, stackType: StackType): string => {
  return `${projectName}-${environment}-${stackType}Stack`;
};

export const createDefaultTags = (
  projectName: string,
  environment: string,
  stackType: StackType,
  costLevel: string,
  additionalTags?: Record<string, string>
) => {
  const baseTags = {
    Project: projectName,
    Environment: environment,
    StackType: stackType,
    CostLevel: costLevel,
  };
  return additionalTags ? { ...baseTags, ...additionalTags } : baseTags;
};

export function extractRootDomain(domain: string): string {
  const parts = domain.split('.');
  return parts.slice(-2).join('.');
}

export function getBackendDomain(rootDomain: string, environment: string): string {
  switch (environment) {
    case 'dev': return `devapi.${rootDomain}`;
    case 'stg': return `stgapi.${rootDomain}`;
    case 'prod': return `api.${rootDomain}`;
    default: return `api.${rootDomain}`;
  }
}

export function getFrontendDomain(rootDomain: string, environment: string): string {
  switch (environment) {
    case 'dev': return `dev.${rootDomain}`;
    case 'stg': return `stg.${rootDomain}`;
    case 'prod': return rootDomain;
    default: return rootDomain;
  }
}

export function getProjectName(): string {
  const name = process.env.PROJECT_NAME;
  if (!name) throw new Error('PROJECT_NAME environment variable is required');
  return name;
}

export function getContainerPort(app?: cdk.App, defaultPort: number = 8000): number {
  if (app) {
    const ctx = app.node.tryGetContext('containerPort');
    if (ctx !== undefined) return typeof ctx === 'number' ? ctx : parseInt(ctx.toString());
  }
  return process.env.CONTAINER_PORT ? parseInt(process.env.CONTAINER_PORT) : defaultPort;
}

export function getImageTag(): string {
  return process.env.IMAGE_TAG || 'latest';
}

export function getCpu(): number {
  return process.env.CPU ? parseInt(process.env.CPU) : 256;
}

export function getMemory(): number {
  return process.env.MEMORY ? parseInt(process.env.MEMORY) : 512;
}

export function getDesiredCount(): number {
  return process.env.DESIRED_COUNT ? parseInt(process.env.DESIRED_COUNT) : 1;
}
