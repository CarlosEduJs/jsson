export const api = {
  timeout: 5000,
  url: "https://api.example.com"
} as const;

export const database = {
  port: 5432,
  host: "localhost"
} as const;

// Generated types
export type Api = typeof api;
export type Database = typeof database;

