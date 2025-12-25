# JSSON v0.0.6 Release Notes - Presets, Validators, HTTP Server & LSP

## New Features

### Presets System
Define reusable configuration blocks with `@preset` and apply them with `@use`:

```javascript
@preset "api" {
  timeout = 30
  retries = 3
}

service = @use "api" {
  endpoint = "/v1"
}
```

### Validators
Auto-generate realistic data with built-in validators:

```javascript
user {
  id = @uuid
  email = @email
  createdAt = @datetime
  website = @url
  ip = @ipv4
}
```

**Available validators:** `@uuid`, `@email`, `@datetime`, `@date`, `@url`, `@ipv4`, `@ipv6`, `@filepath`, `@int(min, max)`, `@float(min, max)`, `@bool`

### Built-in HTTP Server
Serve your configurations directly via HTTP with zero configuration:

```bash
# Start the server
jsson --serve -i config.jsson

# Custom port
jsson --serve --port 3000 -i config.jsson

# Fetch from any service
curl http://localhost:8080/config
```

**Features:**
- Production-ready lightweight server
- Hot-reload on file changes
- CORS support for web clients
- Multiple output formats (JSON, YAML, TOML)
- Perfect for distributing configs across infrastructure

### Boolean Literals
Use `yes/no` and `on/off` as more readable alternatives to `true/false`:

```javascript
config {
  enabled = yes
  debug = on
  production = true
}
```

### Output Formatting
Control JSON output format:

```bash
# Minified (no whitespace)
jsson -i config.jsson --minify

# Custom indentation
jsson -i config.jsson --indent 4
```

### Language Server Protocol (LSP)
Full LSP support for VS Code extension:
- Real-time diagnostics
- Auto-complete for keywords and presets
- Go-to-definition for presets
- Hover information
- Syntax highlighting

---

## Bug Fixes

### Int/Float Comparison
Fixed comparison operators to correctly handle mixed int/float comparisons:

```javascript
// Now works correctly
value = 10
result = value > 9.5  // true
```

### String Keys with Special Characters
Fixed parser to accept string keys starting with `/` or other special characters:

```javascript
// Now works
api {
  "/users" = "GET"
  "/auth/login" = "POST"
}
```

### Bare Identifier Handling
Improved error messages for undefined variables. Variables now have priority over bare identifiers:

```javascript
// Correct
name = "Alice"

// Error: variable 'Alice' not found
// Suggestion: Did you mean to use a string? Try: name = "Alice"
name = Alice
```

---

## Documentation

- Complete documentation site at [docs.jssonlang.tech](https://docs.jssonlang.tech)
- New sections: Patterns & Anti-Patterns, Editor/LSP, HTTP Server
- Real-world examples updated with v0.0.6 features
- Contributing guide (CONTRIBUTING.md)

---

## Internal Changes

- Refactored transpiler into modular components (`eval_*.go`, `format_*.go`)
- Added 1300+ comprehensive tests (edge cases, presets, validators)
- Public API in `pkg/` directory
- Build scripts for Linux (`build.sh`) and Windows (`build.ps1`)
- CI/CD with GitHub Actions

---

## Installation

### CLI

```bash
# Linux AMD64
curl -L https://github.com/jsson-lang/jsson/releases/download/v0.0.6/jsson-v0.0.6-linux-amd64 -o jsson
chmod +x jsson && sudo mv jsson /usr/local/bin/

# macOS AMD64 (Intel)
curl -L https://github.com/jsson-lang/jsson/releases/download/v0.0.6/jsson-v0.0.6-darwin-amd64 -o jsson
chmod +x jsson && sudo mv jsson /usr/local/bin/

# macOS ARM64 (Apple Silicon)
curl -L https://github.com/jsson-lang/jsson/releases/download/v0.0.6/jsson-v0.0.6-darwin-arm64 -o jsson
chmod +x jsson && sudo mv jsson /usr/local/bin/

# Windows
# Download jsson-v0.0.6-windows-amd64.exe from releases
```

### VS Code Extension

Install from the [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson).

---

## Full Changelog

**Full Changelog**: https://github.com/jsson-lang/jsson/compare/v0.0.5...v0.0.6

---

## What's Next?

JSSON v0.0.6 is the last pre-alpha release. The next version will be **v1.0.0** with:
- Frozen specification (no breaking changes)
- Production-ready stability
- Official announcement

Stay tuned! :)
