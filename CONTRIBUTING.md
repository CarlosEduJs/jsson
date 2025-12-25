# Contributing to JSSON

Thank you for your interest in contributing to JSSON! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)

---

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help maintain a welcoming environment

---

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Basic understanding of compilers/transpilers

### Fork and Clone

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/jsson.git
cd jsson

# Add upstream remote
git remote add upstream https://github.com/jsson-lang/jsson.git
```

---

## Development Setup

### Build from Source

```bash
# Build CLI
go build -o jsson ./cmd/jsson

# Build LSP Server
go build -o jsson-lsp ./cmd/lsp

# Build WASM (for playground)
GOOS=js GOARCH=wasm go build -o jsson.wasm ./cmd/wasm
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/parser
go test ./internal/transpiler
```

---

## Project Structure

```
jsson/
├── cmd/
│   ├── jsson/       # CLI entry point
│   ├── lsp/         # LSP server
│   └── wasm/        # WASM build
├── internal/
│   ├── lexer/       # Tokenization
│   ├── parser/      # AST generation
│   ├── transpiler/  # Code generation
│   ├── lsp/         # Language server
│   └── validator/   # Schema validation
├── pkg/             # Public API
├── examples/        # Example .jsson files
└── scripts/         # Build scripts
```

---

## How to Contribute

### Reporting Bugs

**Before submitting:**
- Check existing issues
- Verify it's reproducible in the latest version

**Include:**
- JSSON version (`jsson version`)
- Operating system
- Minimal reproducible example
- Expected vs actual behavior

**Example:**
```markdown
**JSSON Version:** v0.0.6
**OS:** Ubuntu 22.04

**Input (.jsson):**
```jsson
config { name = test }
```

**Expected:** `{"config": {"name": "test"}}`
**Actual:** Error: ...
```

### Suggesting Features

- Open an issue with `[Feature Request]` prefix
- Describe the use case
- Provide examples of how it would work
- Explain why it's valuable

---

## Coding Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Run `gofmt` before committing
- Use meaningful variable names
- Add comments for exported functions

**Example:**
```go
// TranspileToJSON converts a JSSON AST to JSON bytes.
// Returns an error if the AST contains invalid nodes.
func TranspileToJSON(program *ast.Program) ([]byte, error) {
    // Implementation
}
```

### File Naming

- `snake_case` for test files: `presets_test.go`
- `camelCase` for regular files: `transpiler.go`
- Group related functionality: `eval_*.go`, `format_*.go`

### Error Handling

- Return errors, don't panic (except for truly unrecoverable situations)
- Wrap errors with context: `fmt.Errorf("failed to parse: %w", err)`
- Use custom error types when appropriate

---

## Testing

### Test Coverage

- **Minimum 70% coverage** for new code
- **100% coverage** for critical paths (parser, transpiler core)

### Test Structure

```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "basic case",
            input:    `name = "test"`,
            expected: `{"name":"test"}`,
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### What to Test

- ✅ Happy path
- ✅ Edge cases
- ✅ Error conditions
- ✅ Boundary values
- ✅ Integration between components

---

## Commit Guidelines

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(parser): add support for @preset directive

Implements preset definitions and @use syntax for reusable
configuration blocks.

Closes #123
```

```
fix(transpiler): correct int/float comparison in maps

Previously, comparing integers with floats would fail.
Now both types are properly coerced.

Fixes #456
```

### Commit Best Practices

- One logical change per commit
- Write clear, descriptive messages
- Reference issues when applicable
- Keep commits atomic and reversible

---

## Pull Request Process

### Before Submitting

1. **Sync with upstream:**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run tests:**
   ```bash
   go test ./...
   ```

3. **Format code:**
   ```bash
   gofmt -w .
   ```

4. **Update documentation** if needed

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Added tests for new functionality
- [ ] All tests pass
- [ ] Manually tested

## Checklist
- [ ] Code follows project style
- [ ] Self-reviewed code
- [ ] Commented complex logic
- [ ] Updated documentation
```

### Review Process

- Maintainers will review within 3-5 days
- Address feedback promptly
- Be open to suggestions
- Squash commits if requested

---

## Release Process

### Versioning

JSSON follows [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

### Release Checklist

1. Update version in `cmd/jsson/main.go`
2. Update `CHANGELOG.md`
3. Run full test suite
4. Build binaries for all platforms
5. Create GitHub release with binaries
6. Update documentation site

---

## Areas Needing Help

### High Priority

- [ ] Performance optimization for large files
- [ ] More real-world examples
- [ ] Integration tests
- [ ] Fuzzing tests

### Medium Priority

- [ ] Additional validators (`@phone`, `@country`, etc.)
- [ ] Better error messages with suggestions
- [ ] VS Code extension improvements

### Documentation

- [ ] Video tutorials
- [ ] Migration guides from other formats
- [ ] Best practices guide

---

## Questions?

- Open a [Discussion](https://github.com/jsson-lang/jsson/discussions)
- Join our community (Discord/Slack - coming soon)
- Email: [maintainer email]

---

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to JSSON! :)
