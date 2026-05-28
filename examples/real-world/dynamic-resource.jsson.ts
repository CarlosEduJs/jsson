export const resources = [
  {
    region: "us-east-1",
    tags: {
      managed_by: "jsson",
      env: "production"
    },
    id: "ec2-web-server-us-east-1",
    name: "web-server",
    type: "ec2"
  },
  {
    tags: {
      managed_by: "jsson",
      env: "production"
    },
    id: "rds-db-primary-us-west-2",
    name: "db-primary",
    type: "rds",
    region: "us-west-2"
  },
  {
    name: "cache-cluster",
    type: "redis",
    region: "eu-central-1",
    tags: {
      managed_by: "jsson",
      env: "production"
    },
    id: "redis-cache-cluster-eu-central-1"
  }
] as const;

// Generated types
export type Resources = typeof resources;

