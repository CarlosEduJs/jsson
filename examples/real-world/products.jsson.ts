export const products = [
  {
    price: 999.99,
    hasDiscount: true,
    finalPrice: 899.991,
    inStock: true,
    id: 1,
    name: "Laptop"
  },
  {
    price: 29.99,
    hasDiscount: false,
    finalPrice: 29.99,
    inStock: true,
    id: 2,
    name: "Mouse"
  },
  {
    finalPrice: 79.99,
    inStock: true,
    id: 3,
    name: "Keyboard",
    price: 79.99,
    hasDiscount: false
  },
  {
    hasDiscount: true,
    finalPrice: 269.99100000000004,
    inStock: true,
    id: 4,
    name: "Monitor",
    price: 299.99
  },
  {
    hasDiscount: true,
    finalPrice: 134.991,
    inStock: true,
    id: 5,
    name: "Headphones",
    price: 149.99
  }
] as const;

// Generated types
export type Products = typeof products;

