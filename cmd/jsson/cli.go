package main

import (
	"context"
	"flag"
	"fmt"
	"io"
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
	inputPtr := flag.String("i", "", "Input JSSON file (use - for stdin)")
	outputPtr := flag.String("o", "", "Output file (default stdout)")
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
		fmt.Fprintln(os.Stderr, "Please provide an input file with -i")
		fmt.Fprintln(os.Stderr, "Use 'jsson help' for usage information")

		os.Exit(1)
	}

	format := strings.ToLower(*formatPtr)
	validFormats := map[string]bool{
		"json": true, "yaml": true, "toml": true,
		"typescript": true, "ts": true,
	}

	if !validFormats[format] {
		fmt.Fprintf(os.Stderr, "Invalid format: %s. Must be json, yaml, toml or typescript\n", *formatPtr)

		os.Exit(1)
	}

	if format == "ts" {
		format = formatTypeScript
	}

	var (
		data     []byte
		absInput string
		baseDir  string
	)

	if *inputPtr == "-" {
		var err error

		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)

			os.Exit(1)
		}

		absInput = "stdin"

		baseDir, err = os.Getwd()
		if err != nil {
			baseDir = "."
		}
	} else {
		var err error

		data, err = os.ReadFile(*inputPtr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)

			os.Exit(1)
		}

		absInput, err = filepath.Abs(*inputPtr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving input path: %v\n", err)

			os.Exit(1)
		}

		baseDir = filepath.Dir(absInput)
	}

	l := lexer.New(string(data))

	if absInput != "stdin" {
		l.SetSourceFile(absInput)
	}

	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Fprintln(os.Stderr, "Parser errors:")

		for _, e := range p.Errors() {
			fmt.Fprintln(os.Stderr, "\t"+e.Error())
		}

		os.Exit(1)
	}

	t := transpiler.New(program, baseDir, transpiler.MergeMode(*mergeMode), absInput)
	t.SetStreamingMode(*streamingPtr, *streamThreshold)

	minify := *minifyPtr || *minifyLong
	t.SetOutputFormat(minify, *indentPtr)

	startTime := time.Now()

	var (
		output []byte
		err    error
	)

	ctx := context.Background()

	switch format {
	case formatJSON:
		output, err = t.Transpile(ctx)
	case formatYAML:
		output, err = t.TranspileToYAML(ctx)
	case "toml":
		output, err = t.TranspileToTOML(ctx)
	case formatTypeScript:
		output, err = t.TranspileToTypeScript(ctx)
	}

	elapsed := time.Since(startTime)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Transpilation error: %v\n", err)

		os.Exit(1)
	}

	if *schemaPtr != "" {
		schemaData, err := os.ReadFile(*schemaPtr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading schema file: %v\n", err)

			os.Exit(1)
		}

		v := validator.New()

		schema, schemaFormat, err := v.LoadSchemaAuto(string(schemaData))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing schema: %v\n", err)

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
		if *outputPtr != "" {
			if err := os.WriteFile(*outputPtr, output, 0o600); err != nil { //nolint:gosec // user-specified path
				fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)

				os.Exit(1)
			}
		} else {
			fmt.Println(string(output))
		}
	}

	fmt.Fprintf(os.Stderr, "✓ Compiled in %v\n", elapsed)
}
