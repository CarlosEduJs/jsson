export const test2 = "hello" as const;

export const test3 = [
  1,
  2
] as const;

export const test1 = 1 as const;

// Generated types
export type Test2 = typeof test2;
export type Test3 = typeof test3;
export type Test1 = typeof test1;

