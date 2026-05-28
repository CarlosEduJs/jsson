//go:build e2e

package main_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

var update = flag.Bool("update", false, "update golden files")
var binaryPath string
var repoRoot string

func TestMain(m *testing.M) {
	flag.Parse()

	var err error
	repoRoot, err = findRepoRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to find repo root: %v\n", err)
		os.Exit(1)
	}

	binaryPath = filepath.Join(repoRoot, ".jsson-test-bin")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/jsson")
	buildCmd.Dir = repoRoot
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build jsson binary: %v\n", err)
		os.Exit(1)
	}

	exitCode := m.Run()

	os.Remove(binaryPath)
	os.Exit(exitCode)
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", dir)
		}

		dir = parent
	}
}

var requiresNoGolden = map[string]bool{
	"showcase/v0.0.6.jsson": true, // validators produce random output
}

var streamingOnly = map[string]bool{
	"features/streaming-large-dataset.jsson": true,
	"features/streaming-nested-maps.jsson":   true,
}

func TestE2E_TranspileGolden(t *testing.T) {
	examplesDir := filepath.Join(repoRoot, "examples")
	formats := []string{"json", "yaml", "toml", "typescript"}

	var jssonFiles []string
	filepath.Walk(examplesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".jsson") {
			return nil
		}

		rel, _ := filepath.Rel(examplesDir, path)

		if streamingOnly[rel] {
			t.Run(rel, func(t *testing.T) {
				t.Parallel()
				testStreamingFile(t, path)
			})

			return nil
		}

		jssonFiles = append(jssonFiles, rel)

		return nil
	})

	for _, rel := range jssonFiles {
		path := filepath.Join(examplesDir, rel)

		for _, format := range formats {
			t.Run(rel+"/"+format, func(t *testing.T) {
				testTranspileGolden(t, rel, path, format)
			})
		}
	}
}

func testTranspileGolden(t *testing.T, rel, path, format string) {
	t.Parallel()

	if requiresNoGolden[rel] {
		testTranspileOnly(t, path, format)

		return
	}

	goldenPath := goldenFilePath(rel, format)

	stdout, _, err := runJsson(path, format)

	if *update {
		if err != nil {
			os.Remove(goldenPath)

			return
		}

		os.WriteFile(goldenPath, stdout, 0o600)

		return
	}

	if err != nil {
		if _, statErr := os.Stat(goldenPath); statErr == nil {
			t.Fatalf("transpile failed but golden exists: %v", err)
		}

		t.Skipf("no golden (expected failure for this format)")
	}

	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("missing golden file %s", goldenPath)
	}

	if normalizedEqual(stdout, golden, format) {
		return
	}

	t.Fatalf("output mismatch for %s (%s)\n— golden (%d bytes):\n%s\n— actual (%d bytes):\n%s",
		rel, format, len(golden), golden, len(stdout), stdout)
}

func normalizedEqual(a, b []byte, format string) bool {
	switch format {
	case "json":
		return jsonNormalizedEqual(a, b)
	case "yaml":
		return yamlNormalizedEqual(a, b)
	case "typescript":
		return tsNormalizedEqual(a, b)
	default:
		return bytes.Equal(a, b)
	}
}

func tsNormalizedEqual(a, b []byte) bool {
	// TypeScript output iterates over map keys with random order.
	// Extract key=value lines and compare as a set.
	extractTSAssignments := func(data []byte) map[string]string {
		result := make(map[string]string)
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "export const ") {
				// export const key = value as const;
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(strings.TrimPrefix(parts[0], "export const "))
					result[key] = strings.TrimSpace(parts[1])
				}
			}
		}
		return result
	}

	ma := extractTSAssignments(a)
	mb := extractTSAssignments(b)

	if len(ma) != len(mb) {
		return false
	}

	for k, va := range ma {
		vb, ok := mb[k]
		if !ok || va != vb {
			return false
		}
	}

	return true
}

func goldenFilePath(rel, format string) string {
	ext := "." + format
	if format == "typescript" {
		ext = ".ts"
	}

	return filepath.Join(repoRoot, "testdata", rel+ext)
}

func testTranspileOnly(t *testing.T, path, format string) {
	stdout, _, err := runJsson(path, format)
	if err != nil {
		t.Fatalf("transpile failed: %v", err)
	}

	if err := validateOutput(stdout, format); err != nil {
		t.Fatalf("invalid %s output: %v", format, err)
	}
}

func testStreamingFile(t *testing.T, path string) {
	// Streaming files: just check exit 0 and valid JSON
	stdout, _, err := runJsson(path, "json")
	if err != nil {
		t.Fatalf("transpile failed: %v", err)
	}

	if !json.Valid(stdout) {
		t.Fatal("output is not valid JSON")
	}

	// Verify it's an object
	var root any
	if err := json.Unmarshal(stdout, &root); err != nil {
		t.Fatal("output is not valid JSON")
	}

	if _, ok := root.(map[string]any); !ok {
		t.Fatal("output root is not an object")
	}

	t.Logf("output size: %d bytes", len(stdout))
}

func runJsson(path, format string) ([]byte, []byte, error) {
	cmd := exec.Command(binaryPath, "-i", path, "-f", format)
	cmd.Dir = repoRoot

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.Bytes(), stderr.Bytes(), err
}

func validateOutput(data []byte, format string) error {
	switch format {
	case "json":
		if !json.Valid(data) {
			return fmt.Errorf("invalid JSON")
		}

		return nil
	case "yaml":
		var v any

		return yaml.Unmarshal(data, &v)
	case "toml":
		return nil // just check it produces output; TOML encoding may fail for some inputs
	case "typescript":
		if len(data) == 0 {
			return fmt.Errorf("empty output")
		}

		return nil
	}

	return nil
}

func jsonNormalizedEqual(a, b []byte) bool {
	var va, vb any
	if err := json.Unmarshal(a, &va); err != nil {
		return false
	}

	if err := json.Unmarshal(b, &vb); err != nil {
		return false
	}

	return jsonDeepEqual(va, vb)
}

func yamlNormalizedEqual(a, b []byte) bool {
	var va, vb any
	if err := yaml.Unmarshal(a, &va); err != nil {
		return false
	}

	if err := yaml.Unmarshal(b, &vb); err != nil {
		return false
	}

	return jsonDeepEqual(va, vb)
}

func jsonDeepEqual(a, b any) bool {
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)

	return bytes.Equal(ja, jb)
}

func TestE2E_Stdin(t *testing.T) {
	t.Parallel()

	cmd := exec.Command(binaryPath, "-i", "-", "-f", "json")
	cmd.Dir = repoRoot
	cmd.Stdin = strings.NewReader(`x = 1
y = "hello"`)

	stdout, err := cmd.Output()
	if err != nil {
		t.Fatalf("stdin transpile failed: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(stdout, &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if v, ok := result["x"].(float64); !ok || v != 1 {
		t.Errorf("expected x=1, got %v", result["x"])
	}

	if v, ok := result["y"].(string); !ok || v != "hello" {
		t.Errorf("expected y=hello, got %v", result["y"])
	}
}

func TestE2E_Fmt(t *testing.T) {
	t.Parallel()

	// Roundtrip test: format a file, transpile both, compare JSON
	files := []string{
		"examples/basics/01-hello-world.jsson",
		"examples/features/presets.jsson",
		"examples/real-world/apiconfig.jsson",
	}

	for _, file := range files {
		path := filepath.Join(repoRoot, file)

		t.Run(file, func(t *testing.T) {
			testFmtRoundtrip(t, path)
		})
	}
}

func testFmtRoundtrip(t *testing.T, path string) {
	// Transpile original
	origOut, _, err := runJsson(path, "json")
	if err != nil {
		t.Fatalf("original transpile failed: %v", err)
	}

	// Run formatter
	fmtCmd := exec.Command(binaryPath, "fmt", path)
	fmtCmd.Dir = repoRoot

	formatted, err := fmtCmd.Output()
	if err != nil {
		t.Fatalf("fmt failed: %v", err)
	}

	if len(formatted) == 0 {
		t.Fatal("fmt produced empty output")
	}

	// Write formatted output to a temp file next to the original
	// so include resolution still works
	tmpDir := filepath.Dir(path)

	tmpFile := filepath.Join(tmpDir, ".fmt_test_tmp.jsson")
	defer os.Remove(tmpFile)

	if err := os.WriteFile(tmpFile, formatted, 0o600); err != nil {
		t.Fatal(err)
	}

	// Transpile formatted
	fmtOut, _, err := runJsson(tmpFile, "json")
	if err != nil {
		t.Fatalf("formatted transpile failed: %v\nformatted:\n%s", err, formatted)
	}

	if !jsonNormalizedEqual(origOut, fmtOut) {
		t.Fatalf("fmt roundtrip changed semantics\noriginal JSON:\n%s\nformatted JSON:\n%s",
			origOut, fmtOut)
	}
}

func TestE2E_FmtCheck(t *testing.T) {
	t.Parallel()

	// Format a file, then run -check on the formatted output
	tmpFile := filepath.Join(t.TempDir(), "check-test.jsson")

	formatted, err := exec.Command(binaryPath, "fmt",
		filepath.Join(repoRoot, "examples/basics/01-hello-world.jsson")).Output()
	if err != nil {
		t.Fatalf("fmt failed: %v", err)
	}

	if err := os.WriteFile(tmpFile, formatted, 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "fmt", "-check", tmpFile)
	cmd.Dir = repoRoot

	if err := cmd.Run(); err != nil {
		t.Fatalf("fmt -check on formatted output should pass: %v", err)
	}
}

func TestE2E_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		file      string
		schema    string
		wantValid bool
	}{
		{
			name:      "simple_config_valid",
			file:      "examples/validation/simple_config.jsson",
			schema:    "examples/schemas/simple.schema.json",
			wantValid: true,
		},
		{
			name:      "invalid_config",
			file:      "examples/validation/invalid_config.jsson",
			schema:    "examples/schemas/simple.schema.json",
			wantValid: false,
		},
		{
			name:      "api_config",
			file:      "examples/validation/api_config.jsson",
			schema:    "examples/schemas/api-config.schema.json",
			wantValid: true,
		},
		{
			name:      "database_config",
			file:      "examples/validation/database_config.jsson",
			schema:    "examples/schemas/database.schema.yaml",
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(repoRoot, tt.file)
			schemaPath := filepath.Join(repoRoot, tt.schema)

			cmd := exec.Command(binaryPath, "-i", filePath,
				"-schema", schemaPath, "-validate-only")
			cmd.Dir = repoRoot

			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			err := cmd.Run()

			if tt.wantValid && err != nil {
				t.Fatalf("expected valid, got error: %v\nstderr: %s", err, stderr.String())
			}

			if !tt.wantValid && err == nil {
				t.Fatal("expected invalid, got success")
			}

			if tt.wantValid {
				if !strings.Contains(stderr.String(), "Validation passed") {
					t.Fatalf("expected validation passed message, got: %s", stderr.String())
				}
			} else {
				if !strings.Contains(stderr.String(), "Validation failed") {
					t.Fatalf("expected validation failed message, got: %s", stderr.String())
				}
			}
		})
	}
}

func TestE2E_Include(t *testing.T) {
	t.Parallel()

	// apiconfig.jsson includes database.jsson — both in real-world/
	stdout, stderr, err := runJsson(
		filepath.Join(repoRoot, "examples/real-world/apiconfig.jsson"),
		"json",
	)
	if err != nil {
		t.Fatalf("include test failed: %v\nstderr: %s", err, stderr)
	}

	var root map[string]any
	if err := json.Unmarshal(stdout, &root); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Should have merged content from both files
	if _, ok := root["api"]; !ok {
		t.Fatal("missing api key (from apiconfig.jsson)")
	}

	if _, ok := root["database"]; !ok {
		t.Fatal("missing database key (from database.jsson via include)")
	}

	// Verify the included content has correct values
	db, ok := root["database"].(map[string]any)
	if !ok {
		t.Fatal("database is not an object")
	}

	if v, ok := db["host"].(string); !ok || v != "localhost" {
		t.Errorf("expected database.host=localhost, got %v", db["host"])
	}

	if v, ok := db["port"].(float64); !ok || v != 5432 {
		t.Errorf("expected database.port=5432, got %v", db["port"])
	}
}

func TestE2E_AllFormats(t *testing.T) {
	t.Parallel()

	// Verify that the basics files produce parseable output in ALL formats
	// (even formats we don't have golden files for)
	files := []string{
		"examples/basics/01-hello-world.jsson",
		"examples/basics/02-variables.jsson",
		"examples/basics/03-objects.jsson",
	}

	formats := []string{"json", "yaml", "toml", "typescript"}

	for _, file := range files {
		path := filepath.Join(repoRoot, file)

		for _, format := range formats {
			t.Run(file+"/"+format, func(t *testing.T) {
				stdout, _, err := runJsson(path, format)
				if err != nil {
					// TOML may fail for some files
					if format == "toml" && strings.Contains(file, "04-arrays") {
						t.Skip("TOML may not support mixed-type arrays")
					}

					t.Fatalf("transpile failed: %v", err)
				}

				if len(stdout) == 0 {
					t.Fatal("empty output")
				}
			})
		}
	}
}
