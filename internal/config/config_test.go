package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Verifies that loading a missing config path returns an error that includes
// the original file path. This ensures callers get actionable context for
// troubleshooting bad config locations.
func TestLoadMissingPathIncludesFileName(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected missing file error")
	}
	if !strings.Contains(err.Error(), path) {
		t.Fatalf("expected error to include path %q, got %v", path, err)
	}
}

// Verifies that valid JSON config is decoded and normalized by trimming
// whitespace from output_dir and domain entries. This confirms load-time
// normalization behavior expected by runtime config resolution.
func TestLoadValidJSONDecodesAndNormalizes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	content := `{"output_dir":" reports ","output":" JSON ","log_location":" logs/custom.log ","domains":[" github.com ","cloudflare.com"]}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.OutputDir != "reports" {
		t.Fatalf("expected trimmed output dir, got %q", cfg.OutputDir)
	}
	if cfg.Output != "json" {
		t.Fatalf("expected normalized output value, got %q", cfg.Output)
	}
	if cfg.LogLocation != "logs/custom.log" {
		t.Fatalf("expected normalized log location, got %q", cfg.LogLocation)
	}
	if len(cfg.Domains) != 2 || cfg.Domains[0] != "github.com" || cfg.Domains[1] != "cloudflare.com" {
		t.Fatalf("unexpected domains: %#v", cfg.Domains)
	}
}

// Verifies that invalid JSON field types are reported as parse errors and that
// the error includes both the config path and failing field. This keeps parse
// failures explicit and easy to diagnose.
func TestLoadInvalidJSONReturnsParseError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	content := `{"output_dir":true,"domains":["github.com"]}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected invalid JSON type error")
	}
	if !strings.Contains(err.Error(), "parse config") {
		t.Fatalf("expected parse config error, got %v", err)
	}
	if !strings.Contains(err.Error(), path) {
		t.Fatalf("expected error to include path %q, got %v", path, err)
	}
	if !strings.Contains(err.Error(), "output_dir") {
		t.Fatalf("expected error to mention output_dir, got %v", err)
	}
}

// Verifies malformed JSON syntax is rejected and surfaced as a parse error.
// This ensures config files must be structurally valid JSON, not just
// semantically correct field values.
func TestLoadMalformedJSONReturnsParseError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	content := `{"output_dir":"output","domains":["github.com"]`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected malformed JSON parse error")
	}
	if !strings.Contains(err.Error(), "parse config") {
		t.Fatalf("expected parse config error, got %v", err)
	}
}

// Verifies plain text or other non-JSON content is rejected by config loading.
// This keeps the config format strict and prevents accidental acceptance of
// arbitrary text files.
func TestLoadNonJSONContentReturnsParseError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	content := `this is not json`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected non-JSON parse error")
	}
	if !strings.Contains(err.Error(), "parse config") {
		t.Fatalf("expected parse config error, got %v", err)
	}
}

// Verifies a config.txt file with non-JSON content is rejected as invalid
// config format. This protects against accidentally treating plain text
// files as valid runtime configuration.
func TestLoadTxtConfigReturnsParseError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.txt")
	content := `output_dir=output`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected invalid file type error for config.txt")
	}
	if !strings.Contains(err.Error(), "must have .json extension") {
		t.Fatalf("expected invalid file type error, got %v", err)
	}
	if !strings.Contains(err.Error(), path) {
		t.Fatalf("expected error to include path %q, got %v", path, err)
	}
}

// Verifies Resolve precedence rules for output_dir selection. A user-provided
// flag value must override config, while an unchanged flag should preserve the
// config value.
func TestResolvePrefersFlagOverConfig(t *testing.T) {
	cfg := Config{OutputDir: "from-config", Output: "json", LogLocation: "logs/from-config.log"}

	resolved := Resolve("from-flag", true, "text", true, "logs/from-flag.log", true, cfg)
	if resolved.OutputDir != "from-flag" {
		t.Fatalf("expected flag output dir, got %q", resolved.OutputDir)
	}
	if resolved.Output != "text" {
		t.Fatalf("expected flag output format, got %q", resolved.Output)
	}
	if resolved.LogLocation != "logs/from-flag.log" {
		t.Fatalf("expected flag log location, got %q", resolved.LogLocation)
	}

	resolved = Resolve(DefaultOutputDir, false, DefaultOutput, false, DefaultLogLocation, false, cfg)
	if resolved.OutputDir != "from-config" {
		t.Fatalf("expected config output dir when flag is unchanged, got %q", resolved.OutputDir)
	}
	if resolved.Output != "json" {
		t.Fatalf("expected config output when output flag is unchanged, got %q", resolved.Output)
	}
	if resolved.LogLocation != "logs/from-config.log" {
		t.Fatalf("expected config log location when log flag is unchanged, got %q", resolved.LogLocation)
	}
}
