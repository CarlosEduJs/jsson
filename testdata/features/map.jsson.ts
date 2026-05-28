export const routes = [
  {
    method: "GET",
    path: "/api/users"
  },
  {
    path: "/api/posts",
    method: "POST"
  },
  {
    method: "DELETE",
    path: "/api/comments"
  }
] as const;

// Generated types
export type Routes = typeof routes;

