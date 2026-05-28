export const routes = [
  {
    upstreamPort: 8080,
    allowMethods: [
      "GET",
      "POST",
      "PUT",
      "DELETE",
      "PATCH"
    ],
    authType: "jwt",
    authRequired: true,
    strategy: "sliding-window",
    service: "auth",
    version: 2,
    healthCheck: "/health",
    upstreamHost: "auth-service.internal",
    protocol: "http",
    requestsPerMinute: 100,
    cacheEnabled: false,
    tokenExpiry: 7200,
    allowHeaders: [
      "Content-Type",
      "Authorization"
    ],
    cacheTTL: 300,
    maxAge: 3600,
    burstSize: 100,
    path: "/api/v2/auth",
    id: "auth-v2"
  },
  {
    burstSize: 10,
    tokenExpiry: 7200,
    cacheEnabled: false,
    cacheTTL: 300,
    service: "payment",
    healthCheck: "/health",
    path: "/api/v2/payments",
    protocol: "http",
    id: "payment-v2",
    version: 2,
    strategy: "sliding-window",
    allowHeaders: [
      "Content-Type",
      "Authorization"
    ],
    allowMethods: [
      "GET",
      "POST",
      "PUT",
      "DELETE",
      "PATCH"
    ],
    maxAge: 3600,
    authType: "jwt",
    requestsPerMinute: 50,
    upstreamHost: "payment-service.internal",
    upstreamPort: 8080,
    authRequired: true
  },
  {
    strategy: "sliding-window",
    upstreamHost: "catalog-service.internal",
    authRequired: true,
    protocol: "http",
    allowHeaders: [
      "Content-Type",
      "Authorization"
    ],
    version: 3,
    tokenExpiry: 3600,
    authType: "jwt",
    path: "/api/v3/products",
    requestsPerMinute: 1000,
    burstSize: 100,
    maxAge: 3600,
    id: "catalog-v3",
    cacheEnabled: true,
    healthCheck: "/health",
    allowMethods: [
      "GET",
      "POST",
      "PUT",
      "DELETE",
      "PATCH"
    ],
    cacheTTL: 300,
    service: "catalog",
    upstreamPort: 8080
  },
  {
    version: 3,
    maxAge: 3600,
    cacheEnabled: false,
    allowMethods: [
      "GET",
      "POST",
      "PUT",
      "DELETE",
      "PATCH"
    ],
    allowHeaders: [
      "Content-Type",
      "Authorization"
    ],
    requestsPerMinute: 1000,
    path: "/api/v3/search",
    authRequired: true,
    cacheTTL: 300,
    id: "search-v3",
    upstreamPort: 8080,
    protocol: "http",
    healthCheck: "/health",
    tokenExpiry: 3600,
    upstreamHost: "search-service.internal",
    authType: "jwt",
    strategy: "sliding-window",
    burstSize: 100,
    service: "search"
  },
  {
    authRequired: true,
    burstSize: 100,
    strategy: "sliding-window",
    cacheEnabled: false,
    service: "users",
    upstreamHost: "users-service.internal",
    version: 2,
    cacheTTL: 300,
    allowMethods: [
      "GET",
      "POST",
      "PUT",
      "DELETE",
      "PATCH"
    ],
    healthCheck: "/health",
    id: "users-v2",
    upstreamPort: 8080,
    requestsPerMinute: 1000,
    allowHeaders: [
      "Content-Type",
      "Authorization"
    ],
    path: "/api/v2/users",
    protocol: "http",
    maxAge: 3600,
    authType: "jwt",
    tokenExpiry: 7200
  },
  {
    cacheTTL: 300,
    cacheEnabled: false,
    requestsPerMinute: 1000,
    burstSize: 100,
    upstreamHost: "orders-service.internal",
    id: "orders-v2",
    service: "orders",
    strategy: "sliding-window",
    path: "/api/v2/orders",
    authRequired: true,
    version: 2,
    allowMethods: [
      "GET",
      "POST",
      "PUT",
      "DELETE",
      "PATCH"
    ],
    healthCheck: "/health",
    authType: "jwt",
    allowHeaders: [
      "Content-Type",
      "Authorization"
    ],
    protocol: "http",
    upstreamPort: 8080,
    maxAge: 3600,
    tokenExpiry: 7200
  },
  {
    version: 1,
    upstreamPort: 8080,
    tokenExpiry: 7200,
    maxAge: 3600,
    healthCheck: "/health",
    path: "/api/v1/notify",
    allowHeaders: [
      "Content-Type",
      "Authorization"
    ],
    burstSize: 100,
    protocol: "http",
    authRequired: true,
    cacheEnabled: false,
    cacheTTL: 300,
    upstreamHost: "notifications-service.internal",
    requestsPerMinute: 1000,
    service: "notifications",
    authType: "basic",
    strategy: "sliding-window",
    allowMethods: [
      "GET",
      "POST",
      "PUT",
      "DELETE",
      "PATCH"
    ],
    id: "notifications-v1"
  },
  {
    requestsPerMinute: 1000,
    path: "/api/v3/analytics",
    cacheTTL: 300,
    upstreamHost: "analytics-service.internal",
    maxAge: 3600,
    protocol: "http",
    cacheEnabled: false,
    upstreamPort: 8080,
    healthCheck: "/health",
    authType: "jwt",
    strategy: "sliding-window",
    allowMethods: [
      "GET",
      "POST",
      "PUT",
      "DELETE",
      "PATCH"
    ],
    burstSize: 100,
    authRequired: true,
    version: 3,
    tokenExpiry: 3600,
    allowHeaders: [
      "Content-Type",
      "Authorization"
    ],
    id: "analytics-v3",
    service: "analytics"
  }
] as const;

export const loadBalancers = [
  {
    stickySession: true,
    errorThreshold: 50,
    volumeThreshold: 10,
    service: "auth",
    circuitTimeout: 30,
    healthCheckPath: "/health",
    circuitBreakerEnabled: true,
    strategy: "least-connections",
    healthCheckTimeout: 5,
    healthCheckInterval: 30,
    unhealthyThreshold: 3,
    healthyThreshold: 2
  },
  {
    healthCheckPath: "/health",
    healthCheckTimeout: 5,
    errorThreshold: 50,
    volumeThreshold: 10,
    strategy: "round-robin",
    circuitTimeout: 30,
    healthCheckInterval: 30,
    unhealthyThreshold: 3,
    healthyThreshold: 2,
    circuitBreakerEnabled: true,
    service: "payment",
    stickySession: false
  },
  {
    service: "catalog",
    healthyThreshold: 2,
    volumeThreshold: 10,
    healthCheckPath: "/health",
    healthCheckInterval: 30,
    healthCheckTimeout: 5,
    unhealthyThreshold: 3,
    stickySession: false,
    circuitBreakerEnabled: true,
    errorThreshold: 50,
    circuitTimeout: 30,
    strategy: "ip-hash"
  },
  {
    volumeThreshold: 10,
    strategy: "weighted-round-robin",
    healthCheckInterval: 30,
    errorThreshold: 50,
    circuitTimeout: 30,
    service: "search",
    stickySession: false,
    healthCheckPath: "/health",
    unhealthyThreshold: 3,
    healthyThreshold: 2,
    healthCheckTimeout: 5,
    circuitBreakerEnabled: true
  }
] as const;

// Generated types
export type Routes = typeof routes;
export type LoadBalancers = typeof loadBalancers;

