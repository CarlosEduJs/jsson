# JSSON v0.0.1 - Initial Release 🚀

First official release of JSSON (JavaScript Simplified Object Notation) - a transpiler that converts simplified syntax to JSON.

## 🎯 Features

- **Simple Syntax**: Write JSON with a cleaner, more readable syntax
- **Include System**: Import and merge other JSSON files with `#include`
- **Comments**: Support for single-line (`//`) and multi-line (`/* */`) comments
- **Type Safety**: Basic type checking during transpilation
- **Cross-Platform**: Binaries available for Windows, Linux, and macOS

## 📦 Installation

Download the appropriate binary for your platform:

- **Windows**: `jsson-v0.0.1-windows-amd64.exe`
- **Linux**: `jsson-v0.0.1-linux-amd64`
- **macOS (Intel)**: `jsson-v0.0.1-darwin-amd64`
- **macOS (Apple Silicon)**: `jsson-v0.0.1-darwin-arm64`

### Linux/macOS
```bash
# Download and make executable
chmod +x jsson-v0.0.1-*
# Move to PATH (optional)
sudo mv jsson-v0.0.1-* /usr/local/bin/jsson
```

### Windows
```powershell
# Rename for easier usage
Rename-Item jsson-v0.0.1-windows-amd64.exe jsson.exe
# Add to PATH or use directly
```

## 🚀 Usage

```bash
jsson -i input.jsson
```

### Options
- `-i <file>`: Input JSSON file (required)
- `-include-merge <mode>`: Include merge strategy: `keep` (default), `overwrite`, or `error`

## 📝 Example

**Input** (`example.jsson`):
```jsson
{
  name: "JSSON",
  version: "0.0.1",
  features: [
    "Simple syntax",
    "Comments support",
    "Include system"
  ]
}
```

**Output**:
```bash
jsson -i example.jsson
```
```json
{
  "name": "JSSON",
  "version": "0.0.1",
  "features": [
    "Simple syntax",
    "Comments support",
    "Include system"
  ]
}
```

## 🎨 VS Code Extension

Install the [JSSON - Simplified JSON Transpiler](https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson) extension for syntax highlighting and language support!

## 📚 Documentation

For more examples and documentation, visit the [GitHub repository](https://github.com/carlosedujs/jsson).

## 🐛 Known Issues

None reported yet! Please report any issues on GitHub.

## 🙏 Acknowledgments

Thanks to everyone who contributed to making JSSON possible!
