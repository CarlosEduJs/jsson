package main

import "encoding/json"

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
