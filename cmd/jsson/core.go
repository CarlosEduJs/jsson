package main

import (
	"errors"
	"fmt"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"jsson/internal/transpiler"
	"strings"
)

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
