export const odds = [
  1,
  3,
  5,
  7,
  9
] as const;

export const doubled = [
  2,
  4,
  6,
  8,
  10
] as const;

export const squares = [
  1,
  4,
  9,
  16,
  25
] as const;

export const ports = [
  8080,
  8081,
  8082,
  8083,
  8084,
  8085
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

export const evens = [
  0,
  2,
  4,
  6,
  8,
  10
] as const;

// Generated types
export type Ports = typeof ports;
export type Countdown = typeof countdown;
export type Evens = typeof evens;
export type Odds = typeof odds;
export type Doubled = typeof doubled;
export type Squares = typeof squares;

