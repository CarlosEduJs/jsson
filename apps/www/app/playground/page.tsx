"use client";

import { Metadata } from "next";
import { useTranspiler } from "@/hooks/use-transpiler";
import { JSSONEditor } from "@/components/playground/editor";
import { OutputViewer } from "@/components/playground/output-viewer";
import { Button } from "@/components/ui/button";
import Logo from "@/components/shared/logo";
import Link from "next/link";

export const metadata: Metadata = {
  title: "Playground - JSSON",
  description: "Playground for JSSON - JavaScript Simplified Object Notation",
};

const DEFAULT_CODE = `// Welcome to the JSSON Playground!
// Try editing this code to see the magic happen.

server {
  port = 8080
  host = "localhost"
  debug = true
  
  // Database configuration
  database {
    type = "postgres"
    url = "postgres://user:pass@localhost:5432/db"
  }
}

// Generate some users
users [
  template { id, name, role }
  
  map (u) = {
    id = u.id
    name = u.name
    role = u.role
    active = true
  }
  
  1, "Alice", "admin"
  2, "Bob", "user"
  3, "Charlie", "user"
]`;

export default function PlaygroundPage() {
  const { code, setCode, output, error, isWasmLoaded } =
    useTranspiler(DEFAULT_CODE);

  return (
    <div className="flex flex-col h-screen bg-background overflow-hidden">
      <header className="flex items-center justify-between px-6 py-3 border-b border-border bg-background/80 backdrop-blur-md z-10">
        <div className="flex items-center gap-4">
          <Logo size="md" />
          <div className="h-6 w-px bg-border" />
          <span className="text-sm font-medium text-muted-foreground">
            Playground
          </span>
          {!isWasmLoaded && (
            <span className="text-xs text-muted-foreground animate-pulse">
              Loading compiler...
            </span>
          )}
        </div>
        <Link href={"https://jsson-docs.vercel.app/"}>
          <Button size="sm">Go to docs</Button>
        </Link>
      </header>

      <main className="flex-1 flex gap-4 p-4 min-h-0">
        <div className="flex-1 min-w-0">
          <JSSONEditor value={code} onChange={(val) => setCode(val || "")} />
        </div>
        <div className="flex-1 min-w-0">
          <OutputViewer jssonCode={code} error={error} />
        </div>
      </main>
    </div>
  );
}
