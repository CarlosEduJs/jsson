export const servers = [
  {
    tier: "prod",
    replicas: 5,
    id: "srv-101",
    region: "us-east-1"
  },
  {
    id: "srv-102",
    region: "us-east-1",
    tier: "prod",
    replicas: 5
  },
  {
    tier: "staging",
    replicas: 2,
    id: "srv-201",
    region: "us-west-2"
  },
  {
    tier: "staging",
    replicas: 2,
    id: "srv-202",
    region: "us-west-2"
  },
  {
    replicas: 5,
    id: "srv-301",
    region: "eu-central-1",
    tier: "prod"
  },
  {
    region: "ap-south-1",
    tier: "dev",
    replicas: 2,
    id: "srv-401"
  }
] as const;

// Generated types
export type Servers = typeof servers;

