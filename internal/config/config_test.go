package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

func TestLoadValidJSONDecodesAndNormalizes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	content := `{"output_dir":" reports ","domains":[" github.com ","cloudflare.com"]}`
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
	if len(cfg.Domains) != 2 || cfg.Domains[0] != "github.com" || cfg.Domains[1] != "cloudflare.com" {
		t.Fatalf("unexpected domains: %#v", cfg.Domains)
	}
}

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

func TestResolvePrefersFlagOverConfig(t *testing.T) {
	cfg := Config{OutputDir: "from-config"}

	resolved := Resolve("from-flag", true, cfg)
	if resolved.OutputDir != "from-flag" {
		t.Fatalf("expected flag output dir, got %q", resolved.OutputDir)
	}

	resolved = Resolve(DefaultOutputDir, false, cfg)
	if resolved.OutputDir != "from-config" {
		t.Fatalf("expected config output dir when flag is unchanged, got %q", resolved.OutputDir)
	}
}