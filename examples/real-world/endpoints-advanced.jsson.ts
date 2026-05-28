export const endpoints = [
  {
    auth: "basic",
    rateLimit: 100,
    path: "/api/v1/users",
    methods: [
      "GET",
      "POST",
      "PUT",
      "DELETE"
    ]
  },
  {
    auth: "jwt",
    rateLimit: 1000,
    path: "/api/v2/posts",
    methods: [
      "GET",
      "POST",
      "PUT",
      "DELETE"
    ]
  },
  {
    methods: [
      "GET",
      "POST",
      "PUT",
      "DELETE"
    ],
    auth: "jwt",
    rateLimit: 1000,
    path: "/api/v2/comments"
  },
  {
    auth: "jwt",
    rateLimit: 1000,
    path: "/api/v3/analytics",
    methods: [
      "GET",
      "POST",
      "PUT",
      "DELETE"
    ]
  }
] as const;

// Generated types
export type Endpoints = typeof endpoints;

