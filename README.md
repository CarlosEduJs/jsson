# JSSON

[![JSSON Banner](https://jssonlang.tech/og-image.png)](https://jssonlang.tech)

**JavaScript Simplified Object Notation** — A universal configuration meta-format.

Write once, transpile to **JSON**, **YAML**, **TOML**, or **TypeScript**.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![VS Code Extension](https://img.shields.io/badge/VS%20Code-Extension-blue)](https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson)

---

## Quick Example

**JSSON:**
```javascript
@preset "api" {
  timeout = 30
  retries = 3
}

users [
  template { id, role }
  map (u) = @use "api" {
    id = @uuid
    email = @email
    role = u.role
    active = yes
  }
  1..5, "admin"
  6..100, "user"
]
```

**Output (100 users):**
```json
{
  "users": [
    {
      "timeout": 30,
      "retries": 3,
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "email": "user_kx7m@example.com",
      "role": "admin",
      "active": true
    }
    // ... 99 more
  ]
}
```

---

## Features

- **Logic-First**: Variables, ranges, maps, conditionals, arithmetic
- **Presets**: Reusable configuration templates with `@preset` and `@use`
- **Validators**: Auto-generate UUIDs, emails, dates with `@uuid`, `@email`, `@datetime`
- **Multi-Format**: Transpile to JSON, YAML, TOML, TypeScript
- **VS Code Extension**: Full LSP support (syntax highlighting, diagnostics, auto-complete)
- **HTTP Server**: Built-in REST API for integration
- **Schema Validation**: Validate output against JSON Schema
- **Streaming**: Handle millions of records efficiently

---

## Installation

### CLI

```bash
# Download from releases
curl -L https://github.com/carlosedujs/jsson/releases/latest/download/jsson-linux-amd64 -o jsson
chmod +x jsson && sudo mv jsson /usr/local/bin/

# Or build from source
go build -o jsson ./cmd/jsson
```

### VS Code Extension

Install from the [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson).

---

## Usage

```bash
# Transpile to JSON
jsson -i config.jsson > config.json

# Transpile to YAML
jsson -i config.jsson -f yaml > config.yaml

# Minified output
jsson -i config.jsson -m > config.min.json

# Start HTTP server
jsson serve
```

---

## Documentation

📚 **Full documentation:** [docs.jssonlang.tech](https://docs.jssonlang.tech)

- [Getting Started](https://docs.jssonlang.tech/docs/core/guides/getting-started)
- [Syntax Reference](https://docs.jssonlang.tech/docs/core/reference/syntax)
- [Presets Guide](https://docs.jssonlang.tech/docs/core/guides/presets)
- [Validators Guide](https://docs.jssonlang.tech/docs/core/guides/validators)
- [Patterns & Anti-Patterns](https://docs.jssonlang.tech/docs/core/patterns)
- [CLI Reference](https://docs.jssonlang.tech/docs/cli)
- [HTTP Server API](https://docs.jssonlang.tech/docs/server)
- [VS Code Extension](https://docs.jssonlang.tech/docs/editor)

---

## What's New in v0.0.6

- **Presets**: Reusable configuration blocks with `@preset` and `@use`
- **Validators**: Auto-generate data with `@uuid`, `@email`, `@datetime`, `@url`, `@ipv4`, etc.
- **Boolean Literals**: Use `yes`/`no` and `on`/`off` as alternatives to `true`/`false`
- **Minify Flag**: `--minify` for compact JSON output
- **Bug Fixes**: Int/float comparison, bare identifier handling

See the [Changelog](https://docs.jssonlang.tech/docs/changelog) for full details.

---

## Contributing

Contributions are welcome! Please read the [Contributing Guide](CONTRIBUTING.md).

---

## License

MIT © [Carlos Eduardo](https://github.com/carlosedujs)
