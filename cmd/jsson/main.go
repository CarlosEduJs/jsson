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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"jsson/internal/transpiler"
	"jsson/internal/validator"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	formatTypeScript = "typescript"
	formatJSON       = "json"
	formatYAML       = "yaml"
)

const (
	Version       = "0.0.6"
	ServerVersion = "0.1.0"
)

func main() {
	// Check for subcommands first
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "serve", "server":
			runServer(os.Args[2:])

			return
		case "help", "-h", "--help":
			printHelp()

			return
		case "version", "-v", "--version":
			fmt.Printf("JSSON v%s\n", Version)

			return
		}
	}

	// Default: run transpiler
	runTranspiler()
}

func printHelp() {
	fmt.Printf(`JSSON v%s - JavaScript Simplified Object Notation

Usage:
  jsson [flags]              Transpile JSSON file
  jsson serve [flags]        Start HTTP server

Transpile Flags:
  -i string              Input JSSON file (required)
  -f string              Output format: json|yaml|toml|typescript (default "json")
  -schema string         Schema file to validate output against (optional)
  -validate-only         Only validate, don't output result
  -stream                Enable streaming mode for large datasets
  -stream-threshold int  Auto-enable streaming threshold (default 10000)
  -include-merge string  Include merge strategy: keep|overwrite|error (default "keep")

Server Flags (jsson serve):
  -port int    Port to listen on (default 8090)
  -cors        Enable CORS for all origins (default true)

Examples:
  jsson -i config.jsson                    # Transpile to JSON
  jsson -i config.jsson -f yaml            # Transpile to YAML  
  jsson -i config.jsson -schema schema.json # Validate output
  jsson serve                              # Start HTTP server
  jsson serve -port 3000                   # Server on port 3000

Documentation: https://docs.jssonlang.tech/
`, Version)
}

// Transpiler Command

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

	// Validate format
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

		for _, msg := range p.Errors() {
			fmt.Println("\t" + msg)
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

	// Schema validation (optional)
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

// Server Command

var (
	serverPort int
	serverCORS bool
)

func runServer(args []string) {
	serverFlags := flag.NewFlagSet("serve", flag.ExitOnError)
	serverFlags.IntVar(&serverPort, "port", 8090, "Port to listen on")
	serverFlags.BoolVar(&serverCORS, "cors", true, "Enable CORS for all origins")
	if err := serverFlags.Parse(args); err != nil {
		log.Fatalf("Failed to parse server flags: %v", err)
	}

	// Routes
	http.HandleFunc("/health", corsMiddleware(healthHandler))
	http.HandleFunc("/version", corsMiddleware(versionHandler))
	http.HandleFunc("/transpile", corsMiddleware(transpileHandler))
	http.HandleFunc("/validate", corsMiddleware(validateHandler))
	http.HandleFunc("/validate-schema", corsMiddleware(validateWithSchemaHandler))

	addr := fmt.Sprintf(":%d", serverPort)

	log.Printf("🚀 JSSON HTTP Server v%s (JSSON v%s)", ServerVersion, Version)
	log.Printf("📡 Listening on http://0.0.0.0%s", addr)
	log.Printf("")
	log.Printf("Endpoints:")
	log.Printf("  POST /transpile        - Transpile JSSON to JSON/YAML/TOML/TypeScript")
	log.Printf("  POST /validate         - Validate JSSON syntax")
	log.Printf("  POST /validate-schema  - Validate transpiled output against schema")
	log.Printf("  GET  /health           - Health check")
	log.Printf("  GET  /version          - Version info")
	log.Printf("")

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

// HTTP Types

type TranspileRequest struct {
	Source          string `json:"source"`
	Format          string `json:"format,omitempty"`
	IncludeMerge    string `json:"include_merge,omitempty"`
	Streaming       bool   `json:"streaming,omitempty"`
	StreamThreshold int64  `json:"stream_threshold,omitempty"`
}

type TranspileResponse struct {
	Success         bool            `json:"success"`
	Output          json.RawMessage `json:"output,omitempty"`
	OutputRaw       string          `json:"output_raw,omitempty"`
	Format          string          `json:"format"`
	Errors          []string        `json:"errors,omitempty"`
	TranspileTimeMs float64         `json:"transpile_time_ms"`
}

type ValidateRequest struct {
	Source string `json:"source"`
}

type ValidateResponse struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors,omitempty"`
}

type ValidateWithSchemaRequest struct {
	Source       string `json:"source"`
	Schema       string `json:"schema"`
	SchemaFormat string `json:"schema_format,omitempty"`
	OutputFormat string `json:"output_format,omitempty"`
}

type ValidateWithSchemaResponse struct {
	Valid           bool              `json:"valid"`
	Errors          []ValidationError `json:"errors,omitempty"`
	Warnings        []ValidationError `json:"warnings,omitempty"`
	TranspiledData  json.RawMessage   `json:"transpiled_data,omitempty"`
	Format          string            `json:"format"`
	SchemaType      string            `json:"schema_type"`
	TranspileTimeMs float64           `json:"transpile_time_ms"`
	ValidateTimeMs  float64           `json:"validate_time_ms"`
}

type ValidationError struct {
	Path       string `json:"path"`
	Message    string `json:"message"`
	SchemaPath string `json:"schema_path,omitempty"`
	Value      string `json:"value,omitempty"`
	Expected   string `json:"expected,omitempty"`
}

type HealthResponse struct {
	Status       string `json:"status"`
	Service      string `json:"service"`
	Version      string `json:"version"`
	JssonVersion string `json:"jsson_version"`
	Timestamp    string `json:"timestamp"`
}

type VersionResponse struct {
	ServerVersion string `json:"server_version"`
	JssonVersion  string `json:"jsson_version"`
	GoVersion     string `json:"go_version"`
}

// HTTP Handlers.
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if serverCORS {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)

			return
		}

		next(w, r)
	}
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	jsonResponse(w, http.StatusOK, HealthResponse{
		Status:       "healthy",
		Service:      "jsson",
		Version:      ServerVersion,
		JssonVersion: Version,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	})
}

func versionHandler(w http.ResponseWriter, _ *http.Request) {
	jsonResponse(w, http.StatusOK, VersionResponse{
		ServerVersion: ServerVersion,
		JssonVersion:  Version,
		GoVersion:     "1.21+",
	})
}

func transpileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed. Use POST.",
		})

		return
	}

	start := time.Now()

	var req TranspileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		elapsed := float64(time.Since(start).Microseconds()) / 1000
		jsonResponse(w, http.StatusBadRequest, TranspileResponse{
			Success:         false,
			Errors:          []string{"Invalid JSON: " + err.Error()},
			Format:          "json",
			TranspileTimeMs: elapsed,
		})

		return
	}

	if req.Source == "" {
		elapsed := float64(time.Since(start).Microseconds()) / 1000
		jsonResponse(w, http.StatusBadRequest, TranspileResponse{
			Success:         false,
			Errors:          []string{"Source is required"},
			Format:          req.Format,
			TranspileTimeMs: elapsed,
		})

		return
	}

	output, errs, err := transpileSource(req.Source, req.Format, req.IncludeMerge, req.Streaming, req.StreamThreshold)

	elapsed := float64(time.Since(start).Microseconds()) / 1000

	format := req.Format
	if format == "" {
		format = formatJSON
	}

	if err != nil {
		jsonResponse(w, http.StatusOK, TranspileResponse{
			Success:         false,
			Errors:          errs,
			Format:          format,
			TranspileTimeMs: elapsed,
		})

		return
	}

	response := TranspileResponse{
		Success:         true,
		Format:          format,
		TranspileTimeMs: elapsed,
	}

	if format == formatJSON {
		if json.Valid(output) {
			response.Output = json.RawMessage(output)
		} else {
			response.OutputRaw = string(output)
		}
	} else {
		response.OutputRaw = string(output)
	}

	jsonResponse(w, http.StatusOK, response)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed. Use POST.",
		})

		return
	}

	var req ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, ValidateResponse{
			Valid:  false,
			Errors: []string{"Invalid JSON: " + err.Error()},
		})

		return
	}

	valid, errs := validateSyntax(req.Source)

	jsonResponse(w, http.StatusOK, ValidateResponse{
		Valid:  valid,
		Errors: errs,
	})
}

func validateWithSchemaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed. Use POST.",
		})

		return
	}

	var req ValidateWithSchemaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, ValidateWithSchemaResponse{
			Valid:  false,
			Errors: []ValidationError{{Path: "$", Message: "Invalid JSON request: " + err.Error()}},
		})

		return
	}

	if req.Source == "" {
		jsonResponse(w, http.StatusBadRequest, ValidateWithSchemaResponse{
			Valid:  false,
			Errors: []ValidationError{{Path: "$", Message: "Source is required"}},
		})

		return
	}

	if req.Schema == "" {
		jsonResponse(w, http.StatusBadRequest, ValidateWithSchemaResponse{
			Valid:  false,
			Errors: []ValidationError{{Path: "$", Message: "Schema is required"}},
		})

		return
	}

	outputFormat := req.OutputFormat
	if outputFormat == "" {
		outputFormat = formatJSON
	}

	transpileStart := time.Now()
	output, transpileErrors, err := transpileSource(req.Source, outputFormat, "keep", false, 10000)
	transpileTime := float64(time.Since(transpileStart).Microseconds()) / 1000

	if err != nil {
		jsonResponse(w, http.StatusOK, ValidateWithSchemaResponse{
			Valid:           false,
			Errors:          convertToValidationErrors(transpileErrors),
			Format:          outputFormat,
			TranspileTimeMs: transpileTime,
		})

		return
	}

	validateStart := time.Now()
	v := validator.New()

	var (
		schema     *validator.Schema
		schemaType string
	)

	schemaFormat := req.SchemaFormat
	switch schemaFormat {
	case "":
		var detectedFormat string

		schema, detectedFormat, err = v.LoadSchemaAuto(req.Schema)
		schemaType = detectedFormat + "-schema"
	case "yaml":
		schema, err = v.LoadSchemaFromYAML(req.Schema)
		schemaType = "yaml-schema"
	default:
		schema, err = v.LoadSchemaFromJSON(req.Schema)
		schemaType = "json-schema"
	}

	if err != nil {
		jsonResponse(w, http.StatusOK, ValidateWithSchemaResponse{
			Valid:           false,
			Errors:          []ValidationError{{Path: "$schema", Message: "Invalid schema: " + err.Error()}},
			Format:          outputFormat,
			TranspileTimeMs: transpileTime,
		})

		return
	}

	result := v.Validate(output, schema, outputFormat)
	validateTime := float64(time.Since(validateStart).Microseconds()) / 1000

	var validationErrors []ValidationError

	for _, e := range result.Errors {
		var valueStr string
		if e.Value != nil {
			valueStr = fmt.Sprintf("%v", e.Value)
		}

		validationErrors = append(validationErrors, ValidationError{
			Path:       e.Path,
			Message:    e.Message,
			SchemaPath: e.SchemaPath,
			Value:      valueStr,
			Expected:   e.Expected,
		})
	}

	response := ValidateWithSchemaResponse{
		Valid:           result.Valid,
		Errors:          validationErrors,
		Format:          outputFormat,
		SchemaType:      schemaType,
		TranspileTimeMs: transpileTime,
		ValidateTimeMs:  validateTime,
	}

	if result.Valid && outputFormat == "json" {
		if json.Valid(output) {
			response.TranspiledData = json.RawMessage(output)
		}
	}

	jsonResponse(w, http.StatusOK, response)
}

// Core Functions

func transpileSource(source, format, includeMerge string, streaming bool, streamThreshold int64) (output []byte, errs []string, err error) {
	if format == "" {
		format = "json"
	}

	if includeMerge == "" {
		includeMerge = "keep"
	}

	if streamThreshold == 0 {
		streamThreshold = 10000
	}

	format = strings.ToLower(format)
	if format == "ts" {
		format = "typescript"
	}

	validFormats := map[string]bool{
		"json": true, "yaml": true, "toml": true, "typescript": true,
	}
	if !validFormats[format] {
		return nil, nil, fmt.Errorf("invalid format: %s", format)
	}

	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return nil, p.Errors(), errors.New("parser errors")
	}

	t := transpiler.New(program, "", includeMerge, "")
	t.SetStreamingMode(streaming, streamThreshold)

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

	if err != nil {
		return output, []string{err.Error()}, err
	}

	return output, errs, err
}

func validateSyntax(source string) (valid bool, errs []string) {
	l := lexer.New(source)
	p := parser.New(l)
	p.ParseProgram()

	errs = p.Errors()

	return len(errs) == 0, errs
}

func convertToValidationErrors(errs []string) []ValidationError {
	result := make([]ValidationError, 0, len(errs))
	for _, e := range errs {
		result = append(result, ValidationError{
			Path:    "$",
			Message: e,
		})
	}

	return result
}
