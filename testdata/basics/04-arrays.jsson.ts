export const colors = [
  "red",
  "green",
  "blue"
] as const;

export const numbers = [
  1,
  2,
  3,
  4,
  5
] as const;

export const users = [
  {
    name: "Alice",
    age: 25
  },
  {
    age: 30,
    name: "Bob"
  },
  {
    name: "Charlie",
    age: 35
  }
] as const;

export const mixed = [
  1,
  "two",
  true,
  {
    nested: "object"
  }
] as const;

// Generated types
export type Users = typeof users;
export type Mixed = typeof mixed;
export type Colors = typeof colors;
export type Numbers = typeof numbers;

