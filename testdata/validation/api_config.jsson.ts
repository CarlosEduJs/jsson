export const auth = {
  type: "jwt",
  expiresIn: "7d",
  issuer: "api.example.com"
} as const;

export const rateLimit = {
  windowMs: 60000,
  maxRequests: 100,
  message: "Too many requests, please try again later"
} as const;

export const version = "1.0.0" as const;

export const servers = [
  {
    host: "development.example.com",
    port: 3000,
    protocol: "http",
    name: "dev-server",
    environment: "development"
  },
  {
    protocol: "http",
    name: "stage-server",
    environment: "staging",
    host: "staging.example.com",
    port: 4000
  },
  {
    port: 443,
    protocol: "https",
    name: "prod-server",
    environment: "production",
    host: "api.example.com"
  }
] as const;

export const endpoints = [
  {
    path: "/api/users",
    method: "GET",
    requiresAuth: true,
    rateLimit: 1000
  },
  {
    path: "/api/users",
    method: "POST",
    requiresAuth: true,
    rateLimit: 100
  },
  {
    path: "/api/posts",
    method: "GET",
    requiresAuth: false,
    rateLimit: 1000
  },
  {
    path: "/api/posts",
    method: "POST",
    requiresAuth: true,
    rateLimit: 100
  }
] as const;

// Generated types
export type Auth = typeof auth;
export type RateLimit = typeof rateLimit;
export type Version = typeof version;
export type Servers = typeof servers;
export type Endpoints = typeof endpoints;

