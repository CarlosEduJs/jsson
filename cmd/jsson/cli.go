package main

import (
	"flag"
	"fmt"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"jsson/internal/transpiler"
	"jsson/internal/validator"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func runTranspiler() {
	inputPtr := flag.String("i", "", "Input JSSON file")
	formatPtr := flag.String("f", "json", "Output format: json|yaml|toml|typescript")
	mergeMode := flag.String("include-merge", "keep", "Include merge strategy: keep|overwrite|error")
	streamingPtr := flag.Bool("stream", false, "Enable streaming mode for large datasets")
	streamThreshold := flag.Int64("stream-threshold", 10000, "Auto-enable streaming for ranges larger than N items")
	schemaPtr := flag.String("schema", "", "Schema file (JSON/YAML) to validate output against")
	validateOnly := flag.Bool("validate-only", false, "Only validate, don't output transpiled result")
	minifyPtr := flag.Bool("m", false, "Minify output (no whitespace)")
	minifyLong := flag.Bool("minify", false, "Minify output (no whitespace)")
	indentPtr := flag.Int("indent", 2, "Number of spaces for indentation (default 2)")

	flag.Parse()

	if *inputPtr == "" {
		fmt.Println("Please provide an input file with -i")
		fmt.Println("Use 'jsson help' for usage information")
		os.Exit(1)
	}

	format := strings.ToLower(*formatPtr)
	validFormats := map[string]bool{
		"json": true, "yaml": true, "toml": true,
		"typescript": true, "ts": true,
	}

	if !validFormats[format] {
		fmt.Printf("Invalid format: %s. Must be json, yaml, toml or typescript\n", *formatPtr)
		os.Exit(1)
	}

	if format == "ts" {
		format = formatTypeScript
	}

	data, err := os.ReadFile(*inputPtr)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	absInput, err := filepath.Abs(*inputPtr)
	if err != nil {
		fmt.Printf("Error resolving input path: %v\n", err)
		os.Exit(1)
	}

	baseDir := filepath.Dir(absInput)

	l := lexer.New(string(data))
	l.SetSourceFile(absInput)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Parser errors:")

		for _, err := range p.Errors() {
			fmt.Println("\t" + err.Error())
		}

		os.Exit(1)
	}

	t := transpiler.New(program, baseDir, *mergeMode, absInput)
	t.SetStreamingMode(*streamingPtr, *streamThreshold)

	minify := *minifyPtr || *minifyLong
	t.SetOutputFormat(minify, *indentPtr)

	startTime := time.Now()

	var output []byte

	switch format {
	case formatJSON:
		output, err = t.Transpile()
	case formatYAML:
		output, err = t.TranspileToYAML()
	case "toml":
		output, err = t.TranspileToTOML()
	case formatTypeScript:
		output, err = t.TranspileToTypeScript()
	}

	elapsed := time.Since(startTime)

	if err != nil {
		fmt.Printf("Transpilation error: %v\n", err)
		os.Exit(1)
	}

	if *schemaPtr != "" {
		schemaData, err := os.ReadFile(*schemaPtr)
		if err != nil {
			fmt.Printf("Error reading schema file: %v\n", err)
			os.Exit(1)
		}

		v := validator.New()

		schema, schemaFormat, err := v.LoadSchemaAuto(string(schemaData))
		if err != nil {
			fmt.Printf("Error parsing schema: %v\n", err)
			os.Exit(1)
		}

		result := v.Validate(output, schema, format)

		if !result.Valid {
			fmt.Fprintf(os.Stderr, "\n❌ Validation failed against schema (%s format):\n", schemaFormat)

			for _, verr := range result.Errors {
				fmt.Fprintf(os.Stderr, "  • %s: %s\n", verr.Path, verr.Message)

				if verr.Value != nil {
					fmt.Fprintf(os.Stderr, "    Got: %v\n", verr.Value)
				}

				if verr.Expected != "" {
					fmt.Fprintf(os.Stderr, "    Expected: %s\n", verr.Expected)
				}
			}

			os.Exit(2)
		}

		fmt.Fprintf(os.Stderr, "✓ Validation passed against schema\n")

		if *validateOnly {
			fmt.Fprintf(os.Stderr, "✓ Compiled and validated in %v\n", elapsed)
			os.Exit(0)
		}
	}

	if !*validateOnly {
		fmt.Println(string(output))
	}

	fmt.Fprintf(os.Stderr, "✓ Compiled in %v\n", elapsed)
}
