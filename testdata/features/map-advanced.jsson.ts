export const products = [
  {
    name: "Laptop Pro 15",
    price: 1299.99,
    category: "electronics",
    currency: "USD",
    inStock: true,
    url: "/products/1",
    id: 1
  },
  {
    name: "Wireless Mouse",
    price: 29.99,
    category: "accessories",
    currency: "USD",
    inStock: true,
    url: "/products/2",
    id: 2
  },
  {
    id: 3,
    name: "USB-C Cable",
    price: 12.99,
    category: "accessories",
    currency: "USD",
    inStock: true,
    url: "/products/3"
  },
  {
    price: 399.99,
    category: "electronics",
    currency: "USD",
    inStock: true,
    url: "/products/4",
    id: 4,
    name: "Monitor 27"
  },
  {
    price: 89.99,
    category: "accessories",
    currency: "USD",
    inStock: true,
    url: "/products/5",
    id: 5,
    name: "Keyboard Mechanical"
  }
] as const;

// Generated types
export type Products = typeof products;

