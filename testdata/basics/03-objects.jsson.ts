export const user = {
  active: true,
  name: "João",
  age: 25
} as const;

export const company = {
  name: "TechCorp",
  address: {
    city: "São Paulo",
    country: "Brazil",
    street: "Rua Principal, 123"
  },
  contact: {
    email: "contact@techcorp.com",
    phone: "+55 11 1234-5678"
  }
} as const;

// Generated types
export type User = typeof user;
export type Company = typeof company;

