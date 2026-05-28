export const name = "my-service" as const;

export const port = 8080 as const;

export const enabled = true as const;

export const tags = [
  "api",
  "production",
  "v1"
] as const;

// Generated types
export type Name = typeof name;
export type Port = typeof port;
export type Enabled = typeof enabled;
export type Tags = typeof tags;

