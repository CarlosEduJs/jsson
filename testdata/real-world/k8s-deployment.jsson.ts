export const configMaps = [
  {
    cacheTTL: "3600",
    featureFlags: "stable",
    apiVersion: "v1",
    kind: "ConfigMap",
    name: "app-config-prod",
    namespace: "prod",
    databaseHost: "postgres-prod.internal",
    redisHost: "redis-prod.internal"
  },
  {
    redisHost: "redis-staging.internal",
    cacheTTL: "60",
    featureFlags: "experimental",
    apiVersion: "v1",
    kind: "ConfigMap",
    name: "app-config-staging",
    namespace: "staging",
    databaseHost: "postgres-staging.internal"
  },
  {
    namespace: "dev",
    databaseHost: "postgres-dev.internal",
    redisHost: "redis-dev.internal",
    cacheTTL: "60",
    featureFlags: "experimental",
    apiVersion: "v1",
    kind: "ConfigMap",
    name: "app-config-dev"
  }
] as const;

export const deployments = [
  {
    kind: "Deployment",
    memoryRequest: "512Mi",
    logLevel: "error",
    name: "api-prod",
    image: "myregistry/api:stable",
    memoryLimit: "1Gi",
    environment: "prod",
    apiVersion: "apps/v1",
    namespace: "prod",
    replicas: 5,
    cpuLimit: "1000m",
    cpuRequest: "500m"
  },
  {
    name: "web-prod",
    namespace: "prod",
    image: "myregistry/web:stable",
    memoryRequest: "512Mi",
    cpuRequest: "500m",
    apiVersion: "apps/v1",
    memoryLimit: "1Gi",
    logLevel: "error",
    replicas: 3,
    cpuLimit: "1000m",
    environment: "prod",
    kind: "Deployment"
  },
  {
    environment: "prod",
    apiVersion: "apps/v1",
    cpuRequest: "500m",
    cpuLimit: "1000m",
    kind: "Deployment",
    name: "worker-prod",
    namespace: "prod",
    replicas: 4,
    image: "myregistry/worker:stable",
    memoryLimit: "1Gi",
    logLevel: "error",
    memoryRequest: "512Mi"
  },
  {
    namespace: "staging",
    cpuLimit: "500m",
    kind: "Deployment",
    apiVersion: "apps/v1",
    memoryRequest: "256Mi",
    cpuRequest: "250m",
    memoryLimit: "512Mi",
    environment: "staging",
    replicas: 2,
    image: "myregistry/api:latest",
    logLevel: "debug",
    name: "api-staging"
  },
  {
    name: "web-staging",
    namespace: "staging",
    environment: "staging",
    logLevel: "debug",
    image: "myregistry/web:latest",
    cpuRequest: "250m",
    memoryLimit: "512Mi",
    replicas: 2,
    cpuLimit: "500m",
    apiVersion: "apps/v1",
    kind: "Deployment",
    memoryRequest: "256Mi"
  },
  {
    namespace: "staging",
    memoryLimit: "512Mi",
    environment: "staging",
    image: "myregistry/worker:latest",
    apiVersion: "apps/v1",
    replicas: 1,
    memoryRequest: "256Mi",
    kind: "Deployment",
    name: "worker-staging",
    cpuRequest: "250m",
    cpuLimit: "500m",
    logLevel: "debug"
  },
  {
    memoryLimit: "512Mi",
    namespace: "dev",
    memoryRequest: "256Mi",
    apiVersion: "apps/v1",
    replicas: 1,
    image: "myregistry/api:latest",
    cpuRequest: "250m",
    cpuLimit: "500m",
    environment: "dev",
    kind: "Deployment",
    logLevel: "debug",
    name: "api-dev"
  },
  {
    memoryRequest: "256Mi",
    kind: "Deployment",
    name: "web-dev",
    namespace: "dev",
    cpuLimit: "500m",
    logLevel: "debug",
    image: "myregistry/web:latest",
    cpuRequest: "250m",
    memoryLimit: "512Mi",
    environment: "dev",
    apiVersion: "apps/v1",
    replicas: 1
  },
  {
    name: "worker-dev",
    replicas: 1,
    image: "myregistry/worker:latest",
    memoryRequest: "256Mi",
    cpuLimit: "500m",
    logLevel: "debug",
    apiVersion: "apps/v1",
    kind: "Deployment",
    namespace: "dev",
    cpuRequest: "250m",
    memoryLimit: "512Mi",
    environment: "dev"
  }
] as const;

export const services = [
  {
    kind: "Service",
    targetPort: 8080,
    protocol: "TCP",
    port: 8080,
    environment: "prod",
    app: "api",
    apiVersion: "v1",
    name: "api-prod",
    namespace: "prod",
    type: "LoadBalancer"
  },
  {
    app: "web",
    kind: "Service",
    protocol: "TCP",
    targetPort: 3000,
    environment: "prod",
    name: "web-prod",
    apiVersion: "v1",
    namespace: "prod",
    type: "LoadBalancer",
    port: 3000
  },
  {
    type: "ClusterIP",
    name: "api-staging",
    namespace: "staging",
    app: "api",
    apiVersion: "v1",
    kind: "Service",
    port: 8080,
    targetPort: 8080,
    protocol: "TCP",
    environment: "staging"
  },
  {
    targetPort: 3000,
    apiVersion: "v1",
    kind: "Service",
    name: "web-staging",
    namespace: "staging",
    port: 3000,
    protocol: "TCP",
    type: "ClusterIP",
    app: "web",
    environment: "staging"
  },
  {
    environment: "dev",
    apiVersion: "v1",
    kind: "Service",
    namespace: "dev",
    type: "ClusterIP",
    targetPort: 8080,
    port: 8080,
    protocol: "TCP",
    name: "api-dev",
    app: "api"
  },
  {
    port: 3000,
    protocol: "TCP",
    environment: "dev",
    kind: "Service",
    namespace: "dev",
    type: "ClusterIP",
    app: "web",
    apiVersion: "v1",
    name: "web-dev",
    targetPort: 3000
  }
] as const;

// Generated types
export type ConfigMaps = typeof configMaps;
export type Deployments = typeof deployments;
export type Services = typeof services;

