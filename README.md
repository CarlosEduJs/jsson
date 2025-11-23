# JSSON VS Code Extension

Syntax highlighting for JSSON (JSON Simplified Object Notation) files.

## Features

✨ **Syntax Highlighting** for all JSSON features:
- Keywords: `include`, `template`, `map`, `step`
- Comments (`//`)
- Strings (quoted and bare identifiers)
- Numbers (integers and floats)
- Booleans (`true`, `false`)
- Operators (`=`, `+`, `-`, `*`, `/`, `..`, `.`)
- Auto-closing brackets and braces

## What is JSSON?

JSSON is a human-friendly syntax that transpiles to JSON. It removes the pain of writing JSON manually by:

✅ Eliminating quotes for keys  
✅ Removing trailing commas  
✅ Adding templates for arrays  
✅ Supporting ranges (`1..10`)  
✅ Including other files  
✅ Reducing repetition  

## Example

**JSSON:**
```jsson
// Simple and clean!
users [
  template { name, age }
  
  João, 19
  Maria, 25
  Pedro, 30
]

ports = 8080..8085

include "database.jsson"
```

**Transpiles to JSON:**
```json
{
  "users": [
    { "name": "João", "age": 19 },
    { "name": "Maria", "age": 25 },
    { "name": "Pedro", "age": 30 }
  ],
  "ports": [8080, 8081, 8082, 8083, 8084, 8085]
}
```

## Installation Local

Install from the VS Code Marketplace:

1. Open VS Code
2. Press `Ctrl+P` (or `Cmd+P` on Mac)
3. Type: `ext install carlosedujs.jsson`
4. Press Enter

Or search for "JSSON" in the Extensions view (`Ctrl+Shift+X`).

## Usage

1. Create a file with `.jsson` extension
2. Start writing JSSON syntax
3. Enjoy automatic syntax highlighting!

## Learn More

- [JSSON GitHub Repository](https://github.com/carlosedujs/jsson)
- [JSSON Documentation](https://github.com/carlosedujs/jsson#readme)

## Release Notes Extension JSSON

### 0.0.1

Initial release of JSSON syntax highlighting:
- Basic syntax highlighting for all JSSON keywords
- Support for comments, strings, numbers, and operators
- Auto-closing pairs for brackets and braces
- Language configuration for better editing experience

## Contributing

Found a bug or want to contribute? Visit the [GitHub repository](https://github.com/carlosedujs/jsson).

## License

MIT

---

**Enjoy coding with JSSON!** 🚀
