export const numbers = [
  1,
  2,
  3,
  4,
  5
] as const;

export const evens = [
  0,
  2,
  4,
  6,
  8,
  10
] as const;

export const negatives = [
  -5,
  -4,
  -3,
  -2,
  -1,
  0
] as const;

export const countdown = [
  10,
  9,
  8,
  7,
  6,
  5,
  4,
  3,
  2,
  1
] as const;

export const ports = [
  8080,
  8081,
  8082,
  8083,
  8084,
  8085
] as const;

// Generated types
export type Negatives = typeof negatives;
export type Countdown = typeof countdown;
export type Ports = typeof ports;
export type Numbers = typeof numbers;
export type Evens = typeof evens;

