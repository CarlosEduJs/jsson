export const templates = {
  empty: [
  ],
  single: [
    {
      name: "Alice",
      id: 1
    }
  ],
  withMap: [
    {
      id: 10,
      value: "test!",
      extra: true
    },
    {
      id: 20,
      value: "demo!",
      extra: true
    }
  ]
} as const;

// Generated types
export type Templates = typeof templates;

