# JSSON

**JavaScript Simplified Object Notation** — A configuration language that transpiles to JSON, YAML, TOML, and TypeScript.

Write config once. Ship it anywhere.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![VS Code Extension](https://img.shields.io/badge/VS%20Code-Extension-blue)](https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson)
[![CI](https://github.com/jsson-lang/jsson/actions/workflows/ci.yml/badge.svg)](https://github.com/jsson-lang/jsson/actions/workflows/ci.yml)

---

## Why JSSON?

JSON is everywhere, but writing it by hand hurts. Missing commas, no comments, no variables, no expressions. YAML is better but has whitespace nightmares and security footguns. TOML is pleasant but limited.

JSSON gives you a modern, human-friendly syntax with **variables, templates, ranges, conditionals, and presets** — and transpiles to the format your tools actually expect.

```javascript
// config.jsson — comments, variables, expressions
base_url := "https://api.example.com"
timeout := 30

server {
  url = base_url + "/v2"
  timeout = timeout
  retries = retries > 5 ? 5 : retries
  active = yes
}
```

```json
// → config.json
{
  "url": "https://api.example.com/v2",
  "timeout": 30,
  "retries": 3,
  "active": true
}
```

---

## Quick Start

### Install

```bash
# Download the latest binary (Linux/macOS)
curl -L https://github.com/jsson-lang/jsson/releases/latest/download/jsson-linux-amd64 -o jsson
chmod +x jsson && sudo mv jsson /usr/local/bin/

# Or build from source
git clone https://github.com/jsson-lang/jsson.git
cd jsson
make build
sudo make install
```

### Your first JSSON file

```javascript
// hello.jsson
message = "Hello, JSSON!"
version = "0.1.0"

app {
  name = "My First App"
  language = "JSSON"
}
```

```bash
jsson -i hello.jsson
```

```json
{
  "message": "Hello, JSSON!",
  "version": "0.1.0",
  "app": {
    "name": "My First App",
    "language": "JSSON"
  }
}
```

---

## Language Tour

### Variables and Expressions

```javascript
base_url := "https://api.example.com"
version := "v1"
timeout := 30
retries := 3

api {
  url = base_url + "/" + version
  timeout = timeout
  retries = retries
}

// Arithmetic
price := 100
discount := 0.2
final_price = price - (price * discount)
```

### Objects and Nesting

```javascript
company {
  name = "TechCorp"

  address {
    street = "123 Main St"
    city = "San Francisco"
    zip = "94105"
  }

  contact {
    email = "hello@techcorp.com"
    phone = "+1 555-0123"
  }
}
```

### Arrays and Ranges

```javascript
// Static arrays
colors = ["red", "green", "blue"]

// Ranges (expanded at transpile time)
ports = 8000..8005
// → [8000, 8001, 8002, 8003, 8004, 8005]

letters = "a".."f"
// → ["a", "b", "c", "d", "e", "f"]
```

### Templates and Maps

Generate repetitive structures with logic:

```javascript
users [
  template { id, role }
  map (u) = {
    id = u.id
    role = u.role
    active = true
  }
  1..100, "admin"
  101..1000, "user"
]
```

→ Produces 1000 user objects with proper IDs and roles.

### Presets

Define reusable configuration blocks:

```javascript
@preset "server-defaults" {
  port = 8080
  host = "localhost"
  timeout = 30
}

// Use with overrides
production = @"server-defaults" {
  port = 443
  host = "0.0.0.0"
  timeout = 60
}
```

### Conditionals

```javascript
env := "production"
debug = env == "development" ? true : false

rate_limit = env == "production" ? 1000 : 100
```

### Validators

Generate realistic data for testing and seeding:

```javascript
@preset "user-gen" {
  id = @uuid
  email = @email
  created_at = @datetime
}

test_user = @"user-gen" {
  email = "admin@test.com"
}
```

Built-in validators: `@uuid`, `@email`, `@datetime`, `@url`, `@ipv4`, `@int`, `@float`, `@string`, and more.

### Includes

Split configs across files:

```javascript
// main.jsson
server {
  port = 3000
}

api = include "api.jsson"
database = include "database.jsson"
```

---

## CLI Reference

```bash
# Transpile to JSON (default)
jsson -i config.jsson

# Pick output format
jsson -i config.jsson -f yaml
jsson -i config.jsson -f toml
jsson -i config.jsson -f typescript

# Read from stdin
echo 'name = "JSSON"' | jsson -i -

# Write to file
jsson -i config.jsson -o output.json

# Minify output
jsson -i config.jsson -m

# Format a JSSON file
jsson fmt config.jsson          # print formatted
jsson fmt -w config.jsson       # write in-place
jsson fmt -check config.jsson   # validate formatting (CI)

# Validate against a schema
jsson -i config.jsson -schema schema.json

# Start HTTP server
jsson serve
jsson serve -port 3000

# Version info
jsson --version
# → JSSON v0.1.0 (commit=abc1234, date=2026-05-28T19:48:26Z)
```

### Output Formats

| Flag | Format | Extensions |
|------|--------|------------|
| `json` (default) | JSON | `.json` |
| `yaml` | YAML | `.yaml`, `.yml` |
| `toml` | TOML | `.toml` |
| `typescript` | TypeScript | `.ts` |

---

## Development

```bash
make build       # Build the binary → bin/jsson
make test        # Run unit tests
make test-e2e    # Run e2e golden file tests
make lint        # go vet + gofmt check
make fmt         # Format all Go source files
make run         # Build and run
make clean       # Remove build artifacts
make install     # Copy to GOPATH/bin
```

The project uses [GoReleaser](https://goreleaser.com) for cross-platform releases. To cut a new release:

```bash
git tag v0.2.0
git push origin v0.2.0
```

The release workflow builds for Linux, macOS, and Windows (amd64 + arm64), creates archives, generates checksums, and publishes a GitHub release.

---

## When to use JSSON

**Good for:**
- Configuration files that need variables and logic
- Generating repetitive config structures (environments, users, deployments)
- Multi-format pipelines (same source → JSON for API + YAML for K8s + TOML for app)
- Test data generation with validators
- Prototyping and scaffolding

**Not ideal for:**
- Machine-to-machine data interchange (stick to JSON)
- Performance-critical parsing (JSSON prioritizes ergonomics over speed)
- Nested data beyond ~5 levels (flatter structures work better)

---

## Documentation

Full documentation at [docs.jssonlang.tech](https://docs.jssonlang.tech)

- [Getting Started](https://docs.jssonlang.tech/docs/core/guides/getting-started)
- [Syntax Reference](https://docs.jssonlang.tech/docs/core/reference/syntax)
- [Presets Guide](https://docs.jssonlang.tech/docs/core/guides/presets)
- [Validators Guide](https://docs.jssonlang.tech/docs/core/guides/validators)
- [CLI Reference](https://docs.jssonlang.tech/docs/cli)
- [HTTP Server API](https://docs.jssonlang.tech/docs/server)
- [VS Code Extension](https://docs.jssonlang.tech/docs/editor)

---

## VS Code Extension

Syntax highlighting, diagnostics, auto-complete, hover previews, and more.

Install from the [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson).

---

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

This project was refactored during a hackathon focused on improving abandonware with AI assistance. The PR tells the full story: typed AST, golden file tests, goreleaser, Makefile, and more — 17 commits, 207 files, ~45k lines changed.

---

## License

MIT © [Carlos Eduardo](https://github.com/carlosedujs)
