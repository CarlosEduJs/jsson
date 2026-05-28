package main

import (
	"context"
	"errors"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"jsson/internal/transpiler"
)

func errsToStrings(errs []error) []string {
	strs := make([]string, len(errs))
	for i, e := range errs {
		strs[i] = e.Error()
	}

	return strs
}

func transpileSource(ctx context.Context, source, format, includeMerge string, streaming bool, streamThreshold int64) (output []byte, errs []string, err error) {
	mergeMode := transpiler.MergeMode(includeMerge)
	if mergeMode == "" {
		mergeMode = transpiler.MergeKeep
	}

	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return nil, errsToStrings(p.Errors()), errors.New("parser errors")
	}

	t := transpiler.New(program, "", mergeMode, "")
	t.SetStreamingMode(streaming, streamThreshold)

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

	if err != nil {
		return output, []string{err.Error()}, err
	}

	return output, errs, err
}

func validateSyntax(source string) (valid bool, errs []string) {
	l := lexer.New(source)
	p := parser.New(l)
	p.ParseProgram()

	errs = errsToStrings(p.Errors())

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
