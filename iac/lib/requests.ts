// ============================================================================
// RdsRequests
// ============================================================================

export class RdsRequests {
  constructor(
    public readonly clusterEndpoint: string | null,
    public readonly clusterPort: string | null,
    public readonly databaseName: string | null,
    public readonly secretArn: string | null,
    public readonly clusterArn: string | null,
  ) {}

  static build(
    clusterEndpoint: string | null,
    clusterPort: string | null,
    databaseName: string | null,
    secretArn: string | null,
    clusterArn: string | null,
  ): RdsRequests {
    return new RdsRequests(clusterEndpoint, clusterPort, databaseName, secretArn, clusterArn);
  }
}

// ============================================================================
// CustomSecretsRequests
// ============================================================================

export interface SecretKeyMapping {
  envVarName: string;
  secretKey: string;
}

export interface SecretConfiguration {
  secretName?: string;
  secretArn?: string;
  keyMappings: SecretKeyMapping[];
}

export class CustomSecretsRequests {
  constructor(public readonly secretConfigurations: SecretConfiguration[]) {
    this.validate();
  }

  private validate() {
    for (const config of this.secretConfigurations) {
      if (!config.secretName && !config.secretArn) {
        throw new Error('Each secret configuration requires either secretName or secretArn');
      }
      if (!config.keyMappings || config.keyMappings.length === 0) {
        throw new Error('Each secret configuration must have at least one key mapping');
      }
      for (const mapping of config.keyMappings) {
        if (!mapping.envVarName?.trim()) throw new Error('Secret key mapping must have a valid envVarName');
        if (!mapping.secretKey?.trim()) throw new Error('Secret key mapping must have a valid secretKey');
      }
    }
  }

  hasSecrets(): boolean {
    return this.secretConfigurations.length > 0 &&
      this.secretConfigurations.some(c => c.keyMappings.length > 0);
  }

  static buildFromName(secretName: string, keyMappings: SecretKeyMapping[]): CustomSecretsRequests {
    return new CustomSecretsRequests([{ secretName, keyMappings }]);
  }

  static buildFromArn(secretArn: string, keyMappings: SecretKeyMapping[]): CustomSecretsRequests {
    return new CustomSecretsRequests([{ secretArn, keyMappings }]);
  }

  static buildFromMultiple(secretConfigurations: SecretConfiguration[]): CustomSecretsRequests {
    return new CustomSecretsRequests(secretConfigurations);
  }

  static buildEmpty(): CustomSecretsRequests {
    return new CustomSecretsRequests([]);
  }
}

// ============================================================================
// ContainerConfigRequests
// ============================================================================

export class ContainerConfigRequests {
  constructor(
    public readonly customEnvironmentVariables?: Record<string, string>,
    public readonly customSecretsRequests?: CustomSecretsRequests,
    public readonly overrideDefaults: boolean = false,
  ) {}

  hasCustomEnvironmentVariables(): boolean {
    return !!(this.customEnvironmentVariables && Object.keys(this.customEnvironmentVariables).length > 0);
  }

  hasCustomSecrets(): boolean {
    return !!(this.customSecretsRequests && this.customSecretsRequests.hasSecrets());
  }

  hasCustomConfiguration(): boolean {
    return this.hasCustomEnvironmentVariables() || this.hasCustomSecrets();
  }

  static build(
    customEnvironmentVariables?: Record<string, string>,
    customSecretsRequests?: CustomSecretsRequests,
    overrideDefaults?: boolean,
  ): ContainerConfigRequests {
    return new ContainerConfigRequests(customEnvironmentVariables, customSecretsRequests, overrideDefaults);
  }

  static buildEmpty(): ContainerConfigRequests {
    return new ContainerConfigRequests();
  }
}
