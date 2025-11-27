"use client";

import { motion } from "motion/react";
import { Check } from "lucide-react";
import { useState } from "react";
import { cn } from "@/lib/utils";
import { CodeBlock } from "../shared/code-block";

const codeExamples = [
  {
    title: "Basic Configuration",
    jsson: `// Simple configuration
api {
  url = "https://api.example.com"
  timeout = 5000
  retries = 3
  
  headers {
    content_type = "application/json"
    auth_token = "bearer-token-123"
  }
}`,
    json: `{
  "api": {
    "url": "https://api.example.com",
    "timeout": 5000,
    "retries": 3,
    "headers": {
      "content_type": "application/json",
      "auth_token": "bearer-token-123"
    }
  }
}`,
    yaml: `api:
    headers:
        auth_token: bearer-token-123
        content_type: application/json
    retries: 3
    timeout: 5000
    url: https://api.example.com`,
    toml: `[api]
  retries = 3
  timeout = 5000
  url = "https://api.example.com"
  [api.headers]
    auth_token = "bearer-token-123"
    content_type = "application/json"`,
    typescript: `export const api = {
  url: "https://api.example.com",
  timeout: 5000,
  retries: 3,
  headers: {
    content_type: "application/json",
    auth_token: "bearer-token-123"
  }
} as const;

// Generated types
export type Api = typeof api;`,
  },
  {
    title: "Multi-Region Deployment",
    jsson: `servers [
  template { id, region, tier }
  
  map (s) = {
    id = "srv-" + s.id
    region = s.region
    tier = s.tier
    ip = "10." + (s.id / 100) + ".0." + (s.id % 100)
    replicas = s.tier == "prod" ? 5 : 2
  }
  
  // US East - Production
  100..109, "us-east-1", prod
  
  // US West - Staging
  200..204, "us-west-2", staging
  
  // EU - Production
  300..314, "eu-central-1", prod
  
  // APAC - Dev
  400..402, "ap-south-1", dev
]`,
    json: `
{
  "servers": [
    {
      "id": "srv-100",
      "region": "us-east-1",
      "tier": "prod",
      "ip": "10.1.0.0",
      "replicas": 5
    },
    {
      "id": "srv-101",
      "region": "us-east-1",
      "tier": "prod",
      "ip": "10.1.0.1",
      "replicas": 5
    },
    {
      "id": "srv-102",
      "region": "us-east-1",
      "tier": "prod",
      "ip": "10.1.0.2",
      "replicas": 5
    },
    //... more servers
    {
      "id": "srv-401",
      "ip": "10.4.01.0.1",
      "region": "ap-south-1",
      "replicas": 2,
      "tier": "dev"
    },
    {
      "id": "srv-402",
      "ip": "10.4.02.0.2",
      "region": "ap-south-1",
      "replicas": 2,
      "tier": "dev"
    }
  ]
}`,
    yaml: `
servers:
    - id: srv-100
      ip: 10.1.0.0
      region: us-east-1
      replicas: 5
      tier: prod
    - id: srv-101
      ip: 10.1.01.0.1
      region: us-east-1
      replicas: 5
      tier: prod
    - id: srv-102
      ip: 10.1.02.0.2
      region: us-east-1
      replicas: 5
      tier: prod
	//... more servers
    - id: srv-401
      ip: 10.4.01.0.1
      region: ap-south-1
      replicas: 2
      tier: dev
    - id: srv-402
      ip: 10.4.02.0.2
      region: ap-south-1
      replicas: 2
      tier: dev
    `,
    toml: `
[[servers]]
  id = "srv-100"
  ip = "10.1.0.0"
  region = "us-east-1"
  replicas = 5
  tier = "prod"

[[servers]]
  id = "srv-101"
  ip = "10.1.01.0.1"
  region = "us-east-1"
  replicas = 5
  tier = "prod"

[[servers]]
  id = "srv-102"
  ip = "10.1.02.0.2"
  region = "us-east-1"
  replicas = 5
  tier = "prod"

//... more servers

[[servers]]
  id = "srv-401"
  ip = "10.4.01.0.1"
  region = "ap-south-1"
  replicas = 2
  tier = "dev"

[[servers]]
  id = "srv-402"
  ip = "10.4.02.0.2"
  region = "ap-south-1"
  replicas = 2
  tier = "dev"
`,
    typescript: `
export const servers = [
  {
    tier: "prod",
    ip: "10.1.0.0",
    replicas: 5,
    id: "srv-100",
    region: "us-east-1"
  },
  {
    region: "us-east-1",
    tier: "prod",
    ip: "10.1.01.0.1",
    replicas: 5,
    id: "srv-101"
  },
  {
    id: "srv-102",
    region: "us-east-1",
    tier: "prod",
    ip: "10.1.02.0.2",
    replicas: 5
  },
  //... more servers
  {
    id: "srv-400",
    region: "ap-south-1",
    tier: "dev",
    ip: "10.4.0.0",
    replicas: 2
  },
  {
    id: "srv-401",
    region: "ap-south-1",
    tier: "dev",
    ip: "10.4.01.0.1",
    replicas: 2
  },
  {
    id: "srv-402",
    region: "ap-south-1",
    tier: "dev",
    ip: "10.4.02.0.2",
    replicas: 2
  }
] as const;

// Generated types
export type Servers = typeof servers;

    `,
  },
  {
    title: "Dynamic Resources",
    jsson: `
resources [
template { name, type, region }

  map (res) = {
    // Auto-generate standardized ID
    id = res.type + "-" + res.name + "-" + res.region
    
    name = res.name
    type = res.type
    region = res.region
    
    tags {
      managed_by = "jsson"
      env = "production"
    }
  }

  "web-server", ec2, "us-east-1" 
  "db-primary", rds, "us-west-2"
]`,
    json: `{
"resources": [
    {
      "id": "ec2-web-server-us-east-1",
      "name": "web-server",
      "type": "ec2",
      "region": "us-east-1",
      "tags": {
        "managed_by": "jsson",
        "env": "production"
      }
    },
    {
      "id": "rds-db-primary-us-west-2",
      "name": "db-primary",
      "type": "rds",
      "region": "us-west-2",
      "tags": {
        "managed_by": "jsson",
        "env": "production"
      }
    }
  ]
}`,
    yaml: `
resources:
    - id: ec2-web-server-us-east-1
      name: web-server
      region: us-east-1
      tags:
        env: production
        managed_by: jsson
      type: ec2
    - id: rds-db-primary-us-west-2
      name: db-primary
      region: us-west-2
      tags:
        env: production
        managed_by: jsson
      type: rds    
    `,
    toml: `
[[resources]]
    id = "ec2-web-server-us-east-1"
    name = "web-server"
    region = "us-east-1"
    tags = { env = "production", managed_by = "jsson" }
    type = "ec2"

[[resources]]
    id = "rds-db-primary-us-west-2"
    name = "db-primary"
    region = "us-west-2"
    tags = { env = "production", managed_by = "jsson" }
    type = "rds"
    `,
    typescript: `
export const resources = [
  {
    id: "ec2-web-server-us-east-1",
    name: "web-server",
    region: "us-east-1",
    tags: {
      env: "production",
      managed_by: "jsson"
    },
    type: "ec2"
  },
  {
    id: "rds-db-primary-us-west-2",
    name: "db-primary",
    region: "us-west-2",
    tags: {
      env: "production",
      managed_by: "jsson"
    },
    type: "rds"
  }
] as const;

// Generated types
export type Resources = typeof resources;
    `,
  },
];

export function CodeDemo() {
  const [activeTab, setActiveTab] = useState(0);
  const [outputFormat, setOutputFormat] = useState<
    "json" | "yaml" | "toml" | "typescript"
  >("json");

  const getOutput = (example: (typeof codeExamples)[0]) => {
    switch (outputFormat) {
      case "yaml":
        return example.yaml || example.json;
      case "toml":
        return example.toml || example.json;
      case "typescript":
        return example.typescript || example.json;
      default:
        return example.json;
    }
  };

  return (
    <section className="py-24 overflow-hidden">
      <div className="container mx-auto px-4 md:px-6">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          <div>
            <motion.h2
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              className="text-3xl font-bold tracking-tight md:text-4xl lg:text-5xl mb-6"
            >
              Write Logic, <br />
              <span className="text-primary">
                Get JSON, YAML, TOML, TypeScript.
              </span>
            </motion.h2>
            <motion.p
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.1 }}
              className="text-lg text-muted-foreground mb-8"
            >
              JSSON brings the power of a real programming language to your
              configuration files. No more copy-pasting or manual error-prone
              editing.
            </motion.p>

            <div className="space-y-4">
              {[
                "Native variables & constants",
                "Arithmetic & Conditional Logic",
                "Templates & Maps for complex arrays",
                "Smart Ranges with steps",
                "Modular configuration (include)",
              ].map((item, i) => (
                <motion.div
                  key={i}
                  initial={{ opacity: 0, x: -20 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: 0.2 + i * 0.1 }}
                  className="flex items-center gap-3"
                >
                  <div className="flex h-6 w-6 items-center justify-center rounded-full bg-primary/10 text-primary">
                    <Check className="h-3.5 w-3.5" />
                  </div>
                  <span className="text-muted-foreground">{item}</span>
                </motion.div>
              ))}
            </div>
          </div>

          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            whileInView={{ opacity: 1, scale: 1 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
            className="relative"
          >
            <div className="relative rounded-xl border border-border bg-card overflow-hidden">
              <div className="flex border-b border-border bg-muted/30">
                {codeExamples.map((example, i) => (
                  <button
                    key={i}
                    onClick={() => setActiveTab(i)}
                    className={cn(
                      "px-6 py-3 text-sm font-medium transition-colors border-r border-border",
                      activeTab === i
                        ? "bg-card text-foreground"
                        : "bg-transparent text-muted-foreground hover:text-foreground"
                    )}
                  >
                    {example.title}
                  </button>
                ))}
              </div>

              <div className="grid grid-cols-2 divide-x divide-border">
                <div className="p-0">
                  <div className="px-4 py-2 text-xs font-mono text-muted-foreground border-b border-border bg-muted/10 flex justify-between">
                    <span>input.jsson</span>
                    <span className="text-primary">JSSON</span>
                  </div>
                  <div className="h-[400px] overflow-hidden bg-card">
                    <CodeBlock
                      code={codeExamples[activeTab].jsson}
                      language="jsson"
                      className="p-4 h-full overflow-auto"
                    />
                  </div>
                </div>
                <div className="p-0 bg-muted/5">
                  <div className="px-4 py-2 text-xs font-mono text-muted-foreground border-b border-border bg-muted/10 flex justify-between items-center">
                    <span>
                      output.
                      {outputFormat === "typescript" ? "ts" : outputFormat}
                    </span>
                    <div className="flex gap-1">
                      {(["json", "yaml", "toml", "typescript"] as const).map(
                        (format) => (
                          <button
                            key={format}
                            onClick={() => setOutputFormat(format)}
                            className={cn(
                              "px-2 py-0.5 text-[10px] rounded transition-colors uppercase font-semibold",
                              outputFormat === format
                                ? "bg-primary text-primary-foreground"
                                : "bg-muted/50 text-muted-foreground hover:bg-muted hover:text-foreground"
                            )}
                          >
                            {format === "typescript" ? "TS" : format}
                          </button>
                        )
                      )}
                    </div>
                  </div>
                  <div className="h-[400px] overflow-hidden">
                    <CodeBlock
                      code={getOutput(codeExamples[activeTab])}
                      language="json"
                      className="p-4 h-full overflow-auto"
                    />
                  </div>
                </div>
              </div>
            </div>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
