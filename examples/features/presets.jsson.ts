export const logging = {
  format: "json",
  output: "stdout",
  level: "info"
} as const;

export const api_service = {
  name: "my-api",
  server: {
    maxConnections: 100,
    port: 3000,
    host: "localhost",
    timeout: 30
  },
  logging: {
    level: "debug",
    format: "json",
    output: "stdout"
  }
} as const;

export const dev_server = {
  port: 8080,
  host: "localhost",
  timeout: 30,
  maxConnections: 100
} as const;

export const prod_server = {
  timeout: 60,
  maxConnections: 10000,
  port: 443,
  host: "0.0.0.0"
} as const;

// Generated types
export type Logging = typeof logging;
export type Api_service = typeof api_service;
export type Dev_server = typeof dev_server;
export type Prod_server = typeof prod_server;

