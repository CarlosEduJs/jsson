## What's New in {{ .Tag }}

### New Features

<!-- List new features here -->

### Bug Fixes

<!-- List bug fixes here -->

### Documentation

<!-- List documentation improvements here -->

### ⚡ Performance

<!-- List performance improvements here -->

### Internal Changes

<!-- List refactoring, tooling, etc. here -->

---

## Installation

### CLI

```bash
# Linux AMD64
curl -L https://github.com/jsson-lang/jsson/releases/download/{{ .Tag }}/jsson-{{ .Tag }}-linux-amd64 -o jsson
chmod +x jsson && sudo mv jsson /usr/local/bin/

# macOS AMD64 (Intel)
curl -L https://github.com/jsson-lang/jsson/releases/download/{{ .Tag }}/jsson-{{ .Tag }}-darwin-amd64 -o jsson
chmod +x jsson && sudo mv jsson /usr/local/bin/

# macOS ARM64 (Apple Silicon)
curl -L https://github.com/jsson-lang/jsson/releases/download/{{ .Tag }}/jsson-{{ .Tag }}-darwin-arm64 -o jsson
chmod +x jsson && sudo mv jsson /usr/local/bin/
```

### VS Code Extension

Install from the [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson).

---

## Full Changelog

**Full Changelog**: https://github.com/jsson-lang/jsson/compare/{{ .PreviousTag }}...{{ .Tag }}
