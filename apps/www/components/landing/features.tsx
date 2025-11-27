"use client";

import { motion } from "motion/react";
import {
  Zap,
  Copy,
  Repeat,
  Calculator,
  FileJson,
  ArrowRightLeft,
} from "lucide-react";

const features = [
  {
    icon: ArrowRightLeft,
    title: "Multi-Format Output",
    description:
      "Write once, transpile to JSON, YAML, TOML, or TypeScript. One source of truth for all your config formats.",
  },
  {
    icon: Copy,
    title: "Templates",
    description:
      "Define reusable object blueprints once. Eliminate repetitive JSON by instantiating templates with different inputs.",
  },
  {
    icon: ArrowRightLeft,
    title: "Maps",
    description:
      "Transform values declaratively. Use map() to generate dynamic objects, computed fields, and derived structures.",
  },
  {
    icon: Repeat,
    title: "Smart Ranges",
    description:
      "Generate large datasets instantly with ranges like 1..5000 or stepped sequences like 0..10 step 2.",
  },
  {
    icon: Calculator,
    title: "Expressions & Logic",
    description:
      "Supports math operations, string concatenation, ternary conditions, and computed properties directly in your config.",
  },
  {
    icon: FileJson,
    title: "Includes & Composition",
    description:
      "Split big configs into multiple files. Import and merge them with include, supporting keep, overwrite, and error modes.",
  },
  {
    icon: Zap,
    title: "Native Types",
    description:
      "First-class support for strings, ints, floats, booleans, arrays, and objects. No more JSON quoting hell.",
  },
  {
    icon: Calculator,
    title: "Inline & Minifiable",
    description:
      "JSSON doesn't break when minified. Whitespace-independent and fully inline-safe for prompts and compact configs.",
  },
  {
    icon: Zap,
    title: "High-Scale Generation",
    description:
      "Generate thousands or millions of structured records with a few lines. Perfect for seeds, mocks and IA datasets.",
  },
  {
    icon: Zap,
    title: "LLM-Friendly by Nature",
    description:
      "Designed without intending to optimize tokens, yet performs exceptionally well when used in AI prompts.",
  },
];

export function Features() {
  return (
    <section className="py-24 bg-muted/30">
      <div className="container mx-auto px-4 md:px-6">
        <div className="text-center mb-16">
          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
            className="text-3xl font-bold tracking-tight md:text-4xl lg:text-5xl"
          >
            Why JSSON?
          </motion.h2>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5, delay: 0.1 }}
            className="mt-4 text-lg text-muted-foreground max-w-2xl mx-auto"
          >
            Stop repeating yourself in JSON. Upgrade to a language designed for
            modern application configuration.
          </motion.p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((feature, index) => (
            <motion.div
              key={index}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5, delay: index * 0.1 }}
              className="group relative overflow-hidden rounded-2xl border border-border bg-card p-8 hover:shadow-lg transition-all duration-300 hover:-translate-y-1"
            >
              <div className="absolute inset-0 bg-linear-to-br from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
              <div className="relative z-10">
                <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors duration-300">
                  <feature.icon className="h-6 w-6" />
                </div>
                <h3 className="text-xl font-semibold mb-2">{feature.title}</h3>
                <p className="text-muted-foreground">{feature.description}</p>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
