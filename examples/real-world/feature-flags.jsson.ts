export const serviceConfig = [
  {
    metricsEnabled: true,
    logFormat: "json",
    environment: "prod",
    databaseHost: "db-prod.internal",
    profilingEnabled: false,
    metricsInterval: 60,
    cacheTTL: 3600,
    serviceName: "api",
    cacheHost: "redis-prod.internal",
    logLevel: "error",
    logDestination: "cloudwatch",
    databasePort: 5432,
    maxConnections: 100,
    sslEnabled: true,
    maxMemory: "2gb",
    connectionTimeout: 30000,
    cachePort: 6379
  },
  {
    logDestination: "stdout",
    environment: "staging",
    cachePort: 6379,
    metricsEnabled: true,
    maxConnections: 20,
    serviceName: "api",
    metricsInterval: 300,
    databasePort: 5432,
    connectionTimeout: 10000,
    sslEnabled: false,
    logLevel: "debug",
    maxMemory: "512mb",
    cacheHost: "redis-staging.internal",
    databaseHost: "db-staging.internal",
    cacheTTL: 300,
    logFormat: "json",
    profilingEnabled: false
  },
  {
    cacheTTL: 300,
    databaseHost: "db-dev.internal",
    databasePort: 5432,
    logFormat: "json",
    serviceName: "api",
    cacheHost: "redis-dev.internal",
    metricsEnabled: true,
    logLevel: "debug",
    maxConnections: 20,
    profilingEnabled: true,
    maxMemory: "512mb",
    metricsInterval: 300,
    environment: "dev",
    cachePort: 6379,
    sslEnabled: false,
    connectionTimeout: 10000,
    logDestination: "stdout"
  },
  {
    serviceName: "worker",
    environment: "prod",
    cachePort: 6379,
    maxMemory: "2gb",
    cacheTTL: 3600,
    profilingEnabled: false,
    databaseHost: "db-prod.internal",
    logDestination: "cloudwatch",
    metricsInterval: 60,
    cacheHost: "redis-prod.internal",
    connectionTimeout: 30000,
    logFormat: "json",
    maxConnections: 100,
    sslEnabled: true,
    logLevel: "error",
    databasePort: 5432,
    metricsEnabled: true
  },
  {
    databasePort: 5432,
    metricsInterval: 300,
    maxMemory: "512mb",
    cachePort: 6379,
    databaseHost: "db-staging.internal",
    cacheHost: "redis-staging.internal",
    logFormat: "json",
    environment: "staging",
    sslEnabled: false,
    logDestination: "stdout",
    maxConnections: 20,
    cacheTTL: 300,
    logLevel: "debug",
    serviceName: "worker",
    profilingEnabled: false,
    connectionTimeout: 10000,
    metricsEnabled: true
  },
  {
    cacheTTL: 300,
    connectionTimeout: 10000,
    logDestination: "stdout",
    profilingEnabled: true,
    metricsInterval: 300,
    logFormat: "json",
    maxConnections: 20,
    logLevel: "debug",
    cachePort: 6379,
    sslEnabled: false,
    metricsEnabled: true,
    databaseHost: "db-dev.internal",
    databasePort: 5432,
    maxMemory: "512mb",
    environment: "dev",
    serviceName: "worker",
    cacheHost: "redis-dev.internal"
  }
] as const;

export const featureFlags = [
  {
    rolloutPercentage: 0,
    rolloutStrategy: "gradual",
    name: "new-dashboard",
    environment: "prod",
    enabled: false,
    key: "new-dashboard_prod",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team"
  },
  {
    rolloutPercentage: 100,
    rolloutStrategy: "gradual",
    name: "new-dashboard",
    environment: "staging",
    enabled: true,
    key: "new-dashboard_staging",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team"
  },
  {
    modifiedBy: "devops-team",
    rolloutPercentage: 100,
    rolloutStrategy: "gradual",
    name: "new-dashboard",
    environment: "dev",
    enabled: true,
    key: "new-dashboard_dev",
    lastModified: "2025-11-25T09:00:00Z"
  },
  {
    environment: "prod",
    enabled: true,
    key: "dark-mode_prod",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team",
    rolloutPercentage: 100,
    rolloutStrategy: "gradual",
    name: "dark-mode"
  },
  {
    enabled: true,
    key: "dark-mode_staging",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team",
    rolloutPercentage: 100,
    rolloutStrategy: "gradual",
    name: "dark-mode",
    environment: "staging"
  },
  {
    name: "dark-mode",
    environment: "dev",
    enabled: true,
    key: "dark-mode_dev",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team",
    rolloutPercentage: 100,
    rolloutStrategy: "gradual"
  },
  {
    rolloutStrategy: "gradual",
    name: "crypto-payments",
    environment: "prod",
    enabled: false,
    key: "crypto-payments_prod",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team",
    rolloutPercentage: 0
  },
  {
    rolloutPercentage: 100,
    rolloutStrategy: "gradual",
    name: "crypto-payments",
    environment: "staging",
    enabled: true,
    key: "crypto-payments_staging",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team"
  },
  {
    environment: "dev",
    enabled: true,
    key: "crypto-payments_dev",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team",
    rolloutPercentage: 100,
    rolloutStrategy: "gradual",
    name: "crypto-payments"
  },
  {
    key: "installments_prod",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team",
    rolloutPercentage: 100,
    rolloutStrategy: "gradual",
    name: "installments",
    environment: "prod",
    enabled: true
  },
  {
    enabled: true,
    key: "installments_staging",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team",
    rolloutPercentage: 100,
    rolloutStrategy: "gradual",
    name: "installments",
    environment: "staging"
  },
  {
    rolloutStrategy: "gradual",
    name: "installments",
    environment: "dev",
    enabled: true,
    key: "installments_dev",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team",
    rolloutPercentage: 100
  },
  {
    key: "ai-recommendations_prod",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team",
    rolloutPercentage: 0,
    rolloutStrategy: "gradual",
    name: "ai-recommendations",
    environment: "prod",
    enabled: false
  },
  {
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team",
    rolloutPercentage: 100,
    rolloutStrategy: "gradual",
    name: "ai-recommendations",
    environment: "staging",
    enabled: true,
    key: "ai-recommendations_staging"
  },
  {
    rolloutStrategy: "gradual",
    name: "ai-recommendations",
    environment: "dev",
    enabled: true,
    key: "ai-recommendations_dev",
    lastModified: "2025-11-25T09:00:00Z",
    modifiedBy: "devops-team",
    rolloutPercentage: 100
  }
] as const;

// Generated types
export type ServiceConfig = typeof serviceConfig;
export type FeatureFlags = typeof featureFlags;

