export const version = "0.0.6" as const;

export const app = {
  name: "My First App",
  language: "JSSON"
} as const;

export const message = "Hello, JSSON!" as const;

// Generated types
export type Message = typeof message;
export type Version = typeof version;
export type App = typeof app;

