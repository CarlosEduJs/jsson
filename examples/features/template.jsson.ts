export const users = [
  {
    name: "João",
    age: 19,
    job: "Student",
    height: 1.75
  },
  {
    name: "Maria",
    age: 25,
    job: "Teacher",
    height: 1.65
  },
  {
    job: "Doctor",
    height: 1.8,
    name: "Pedro",
    age: 30
  },
  {
    name: "Ana",
    age: 22,
    job: "Nurse",
    height: 1.68
  }
] as const;

// Generated types
export type Users = typeof users;

