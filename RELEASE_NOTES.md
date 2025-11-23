# 🚀 **JSSON v0.1.0 – The Human JSON Release**

JSSON (**JavaScript Simplified Object Notation**) is a modern, human-friendly syntax that transpiles directly to JSON — with templates, ranges, includes, maps, and more.

This is the **first full-featured release** of the language.

---

# ✨ **🎯 Key Features**

### 🔶 **1. Template Arrays (⚡ Killer Feature)**

Generate structured JSON objects using simple row-based syntax:

```jsson
users [
  template { name, age, job }
  João, 19, Student
  Maria, 25, Teacher
]
```

Output:

```json
{
  "users": [
    { "name": "João", "age": 19, "job": "Student" },
    { "name": "Maria", "age": 25, "job": "Teacher" }
  ]
}
```

---

### 🔶 **2. Map Transformer**

```jsson
routes [
  template { path, method }

  map (item) = {
    path = "/api/" + item.path
    method = item.method
  }

  users, GET
  posts, POST
]
```

---

### 🔶 **3. Includes (modularization)**

Supports relative file includes:

```jsson
include "./config/database.jsson"
```

Circular include detection ✔
Include cache ✔

---

### 🔶 **4. Ranges and Step Support**

```jsson
ports = [ 8080..8085 ]
even = [ 0..10 step 2 ]
```

---

### 🔶 **5. Literal Types**

* Strings
* Integers
* **Floats (NEW!)**
* Booleans
* Objects
* Arrays
* Identifiers
* Member access (`obj.key`)
* String concatenation (`"a" + b`)

---

### 🔶 **6. Clean Syntax**

No braces required for arrays/objects inside templates.
Readable. Minimal. Fast to write.

---

### 🔶 **7. Fully JSON-Accurate Output**

The transpiler guarantees 100% valid JSON output.

---

### 🔶 **8. Wizard/Goblin Error Messages™**

Fun, descriptive error reporting:

```
Syntax wizard: line 3 col 12 — expected '}' — wizard can't find the closing brace
```

---

# 📦 **Installation**

Download the binary for your OS:

* **Windows:** `jsson-v0.1.0-windows-amd64.exe`
* **Linux:** `jsson-v0.1.0-linux-amd64`
* **macOS (Intel):** `jsson-v0.1.0-darwin-amd64`
* **macOS (Apple Silicon):** `jsson-v0.1.0-darwin-arm64`

### Linux/macOS

```bash
chmod +x jsson-v0.1.0-*
sudo mv jsson-v0.1.0-* /usr/local/bin/jsson
```

### Windows

```powershell
Rename-Item jsson-v0.1.0-windows-amd64.exe jsson.exe
```

---

# 🚀 **Usage**

```bash
jsson -i input.jsson
```

---

# 📘 **Example**

**example.jsson**

```jsson
app {
  name = "JSSON"
  version = "0.1.0"

  ports = [ 8080..8083 ]

  authors [
    template { name, role }
    Carlos, Creator
    João, Contributor
  ]
}
```

Output is valid JSON.

---

# 🎨 **VS Code Extension**

Official syntax highlighting + language support:

👉 [https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson](https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson)

---

# 📚 Documentation

Docs & playground:
👉 [https://github.com/carlosedujs/jsson](https://github.com/carlosedujs/jsson)

---

# 🐛 Known Issues

* No known issues after float fix
* Please report anything unexpected in GitHub issues

---

# 🙏 Acknowledgments

Thanks to everyone helping shape this language.
Special thanks to the wizards, goblins and gremlins of the parser.

