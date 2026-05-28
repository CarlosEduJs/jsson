/*
JSSON - JavaScript Simplified Object Notation
==============================================

A human-friendly syntax that transpiles to JSON, YAML, TOML, and TypeScript.

Usage:

	jsson [command] [flags]

Commands:

	(default)    Transpile JSSON file to output format
	serve        Start HTTP server for API access

Transpile Flags:

	-i string    Input JSSON file (required)
	-f string    Output format: json|yaml|toml|typescript (default "json")
	-schema      Schema file to validate output against (optional)
	-validate-only  Only validate, don't output result
	-stream      Enable streaming mode for large datasets
	-stream-threshold  Auto-enable streaming threshold (default 10000)

Server Flags:

	-port int    Port to listen on (default 8090)
	-cors        Enable CORS (default true)

Examples:

	# Transpile to JSON
	jsson -i config.jsson

	# Transpile to YAML
	jsson -i config.jsson -f yaml

	# Validate against schema
	jsson -i config.jsson -schema schema.json

	# Start HTTP server
	jsson serve

	# Start server on custom port
	jsson serve -port 3000
*/
package main

import (
	"fmt"
	"os"

	"jsson/internal/version"
)

const (
	formatTypeScript = "typescript"
	formatJSON       = "json"
	formatYAML       = "yaml"
)

const (
	ServerVersion = "0.1.0"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "serve", "server":
			runServer(os.Args[2:])

			return
		case "fmt", "format":
			os.Args = append([]string{os.Args[0] + " fmt"}, os.Args[2:]...)
			runFormatter()

			return
		case "help", "-h", "--help":
			printHelp()

			return
		case "version", "-v", "--version":
			fmt.Printf("JSSON v%s (commit=%s, date=%s)\n", version.Version, version.Commit, version.Date)

			return
		}
	}

	runTranspiler()
}

func printHelp() {
	fmt.Printf(`JSSON v%s - JavaScript Simplified Object Notation

Usage:
  jsson [flags]              Transpile JSSON file
  jsson serve [flags]        Start HTTP server
  jsson fmt [flags] <file>   Format JSSON file

Commands:
  serve, server       Start HTTP server
  fmt, format         Format JSSON files (use -w to write in-place)
  help                Show this help
  version             Show version

Transpile Flags:
  -i string              Input JSSON file (required, use - for stdin)
  -o string              Output file (default: stdout)
  -f string              Output format: json|yaml|toml|typescript (default "json")
  -m                     Minify output (no whitespace)
  -minify                Minify output (no whitespace)
  -indent int            Number of spaces for indentation (default 2)
  -schema string         Schema file to validate output against
  -validate-only         Only validate, don't output result
  -stream                Enable streaming mode for large datasets
  -stream-threshold int  Auto-enable streaming threshold (default 10000)
  -include-merge string  Include merge strategy: keep|overwrite|error (default "keep")

Server Flags (jsson serve):
  -port int    Port to listen on (default 8090)
  -cors        Enable CORS for all origins (default true)

Examples:
  jsson -i config.jsson                     # Transpile to JSON
  jsson -i config.jsson -f yaml             # Transpile to YAML
  jsson -i config.jsson -o output.json      # Write to file
  cat config.jsson | jsson -i -             # Read from stdin
  jsson -i config.jsson -schema schema.json # Validate output
  jsson serve                               # Start HTTP server
  jsson serve -port 3000                    # Server on port 3000

Documentation: https://docs.jssonlang.tech/
`, version.Version)
}
