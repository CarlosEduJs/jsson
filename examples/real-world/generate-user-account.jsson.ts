export const users = [
  {
    dept: "Engineering",
    role: "employee",
    id: 101,
    username: "user_101",
    email: "user_101@company.com"
  },
  {
    id: 102,
    username: "user_102",
    email: "user_102@company.com",
    dept: "Engineering",
    role: "employee"
  },
  {
    id: 103,
    username: "user_103",
    email: "user_103@company.com",
    dept: "Engineering",
    role: "employee"
  },
  {
    dept: "Sales",
    role: "employee",
    id: 201,
    username: "user_201",
    email: "user_201@company.com"
  },
  {
    dept: "Sales",
    role: "employee",
    id: 202,
    username: "user_202",
    email: "user_202@company.com"
  }
] as const;

// Generated types
export type Users = typeof users;

