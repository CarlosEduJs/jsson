export const employees = [
  {
    name: "Employee 1001",
    email: "emp1001@corp.com",
    department: "Engineering",
    level: "junior",
    salary: 50000,
    id: 1001
  },
  {
    email: "emp1002@corp.com",
    department: "Engineering",
    level: "mid",
    salary: 80000,
    id: 1002,
    name: "Employee 1002"
  },
  {
    level: "senior",
    salary: 120000,
    id: 1003,
    name: "Employee 1003",
    email: "emp1003@corp.com",
    department: "Engineering"
  },
  {
    department: "Sales",
    level: "junior",
    salary: 50000,
    id: 2001,
    name: "Employee 2001",
    email: "emp2001@corp.com"
  },
  {
    department: "Sales",
    level: "mid",
    salary: 80000,
    id: 2002,
    name: "Employee 2002",
    email: "emp2002@corp.com"
  },
  {
    department: "HR",
    level: "mid",
    salary: 80000,
    id: 3001,
    name: "Employee 3001",
    email: "emp3001@corp.com"
  },
  {
    salary: 120000,
    id: 3002,
    name: "Employee 3002",
    email: "emp3002@corp.com",
    department: "HR",
    level: "senior"
  }
] as const;

// Generated types
export type Employees = typeof employees;

