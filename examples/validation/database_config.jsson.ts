export const cache = {
  name: "app_cache",
  config: {
    type: "redis",
    port: 6379,
    ssl: false
  },
  host: "redis.internal"
} as const;

export const replication = {
  primary: "app_prod",
  replicas: [
    "replica_1",
    "replica_2"
  ],
  strategy: "round-robin",
  enabled: true
} as const;

export const migrations = {
  enabled: true,
  directory: "./migrations",
  tableName: "schema_migrations"
} as const;

export const databases = [
  {
    environment: "development",
    host: "localhost",
    config: {
      type: "postgres",
      port: 5432,
      ssl: false,
      pool: {
        max: 10,
        idleTimeoutMs: 30000,
        min: 2
      }
    },
    ssl: false,
    name: "app_dev"
  },
  {
    environment: "staging",
    host: "staging-db.internal",
    config: {
      type: "postgres",
      port: 5432,
      ssl: false,
      pool: {
        min: 2,
        max: 10,
        idleTimeoutMs: 30000
      }
    },
    ssl: false,
    name: "app_staging"
  },
  {
    ssl: true,
    name: "app_prod",
    environment: "production",
    host: "db.prod.internal",
    config: {
      ssl: false,
      pool: {
        max: 10,
        idleTimeoutMs: 30000,
        min: 2
      },
      type: "postgres",
      port: 5432
    }
  }
] as const;

// Generated types
export type Cache = typeof cache;
export type Replication = typeof replication;
export type Migrations = typeof migrations;
export type Databases = typeof databases;

