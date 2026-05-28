export const api = {
  url: "https://api.example.com/v1",
  timeout: 30,
  retries: 3
} as const;

export const final_price = 80 as const;

// Generated types
export type Api = typeof api;
export type Final_price = typeof final_price;

