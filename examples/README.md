# JSSON Examples

This directory contains examples demonstrating JSSON features and real-world use cases.

## 📁 Structure

```
examples/
├── basics/          # Beginner-friendly examples (start here!)
├── features/        # Feature demonstrations (maps, templates, presets, etc.)
├── real-world/      # Real-world use cases
├── showcase/        # Version showcases
├── schemas/         # JSON schemas for validation
└── validation/      # Validation examples
```

## 🚀 Getting Started

### Basics (Start Here!)

| File | Description |
|------|-------------|
| [01-hello-world.jsson](basics/01-hello-world.jsson) | Your first JSSON file |
| [02-variables.jsson](basics/02-variables.jsson) | Using variables with `:=` |
| [03-objects.jsson](basics/03-objects.jsson) | Objects and nesting |
| [04-arrays.jsson](basics/04-arrays.jsson) | Arrays in JSSON |
| [05-ranges.jsson](basics/05-ranges.jsson) | Ranges and steps |

### Features

| File | Description |
|------|-------------|
| [map.jsson](features/map.jsson) | Basic map transformations |
| [map-advanced.jsson](features/map-advanced.jsson) | Advanced map usage |
| [template.jsson](features/template.jsson) | Template arrays |
| [presets.jsson](features/presets.jsson) | Reusable presets |
| [ranges.jsson](features/ranges.jsson) | Range expressions |
| [ternary.jsson](features/ternary.jsson) | Conditional expressions |
| [streaming-*.jsson](features/) | Streaming for large datasets |

### Real-World Examples

| File | Description |
|------|-------------|
| [apiconfig.jsson](real-world/apiconfig.jsson) | API configuration |
| [database.jsson](real-world/database.jsson) | Database config |
| [seed.jsson](real-world/seed.jsson) | Database seeding |
| [k8s-deployment.jsson](real-world/k8s-deployment.jsson) | Kubernetes deployment |
| [game_macros_presets.jsson](real-world/game_macros_presets.jsson) | Game configuration |

## 🔧 Running Examples

```bash
# Transpile to JSON
jsson -i examples/basics/01-hello-world.jsson

# Transpile to YAML
jsson -i examples/basics/01-hello-world.jsson -f yaml

# Transpile to TOML
jsson -i examples/basics/01-hello-world.jsson -f toml

# Validate with schema
jsson -i examples/real-world/apiconfig.jsson -schema examples/schemas/api-config.json
```

## 📚 Programmatic Usage

See the [api-usage/](api-usage/) directory for examples of using JSSON as a Go library.
