package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"jsson/internal/validator"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	log.Printf("🛑 Press Ctrl+C to stop the server")

	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}

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

	output, errs, err := transpileSource(r.Context(), req.Source, req.Format, req.IncludeMerge, req.Streaming, req.StreamThreshold)

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
	output, transpileErrors, err := transpileSource(r.Context(), req.Source, outputFormat, "keep", false, 10000)
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
