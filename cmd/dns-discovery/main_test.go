package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Verifies plain text domain lists are accepted and normalized to lowercase.
// This confirms file-based batch input behavior for valid newline-delimited
// domains.
func TestLoadDomainsFromFileAcceptsPlainTextDomains(t *testing.T) {
	path := filepath.Join(t.TempDir(), "domains.txt")
	content := "github.com\nCloudFlare.com\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write domain file: %v", err)
	}

	domains, err := loadDomainsFromFile(path)
	if err != nil {
		t.Fatalf("load domains: %v", err)
	}

	if len(domains) != 2 {
		t.Fatalf("expected 2 domains, got %d", len(domains))
	}
	if domains[0] != "github.com" || domains[1] != "cloudflare.com" {
		t.Fatalf("unexpected normalized domains: %#v", domains)
	}
}

// Verifies the input domain file rejects entries containing invalid characters.
// This enforces plain-text domain formatting and prevents malformed tokens from
// being accepted into batch execution.
func TestLoadDomainsFromFileRejectsInvalidCharacters(t *testing.T) {
	path := filepath.Join(t.TempDir(), "domains.txt")
	content := "github.com\nbad_domain!.com\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write domain file: %v", err)
	}

	_, err := loadDomainsFromFile(path)
	if err == nil {
		t.Fatal("expected invalid domain character error")
	}
	if !strings.Contains(err.Error(), "invalid domain") {
		t.Fatalf("expected invalid domain error, got %v", err)
	}
}

// Verifies the input domain file rejects entries containing invalid characters.
// This enforces plain-text domain formatting and prevents malformed tokens from
// being accepted into batch execution.
func TestLoadDomainsFromInvalidFileFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "domains.json")
	content := "[\n  \"github.com\",\n  \"bad_domain!.com\"\n]"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write domain file: %v", err)
	}

	_, err := loadDomainsFromFile(path)
	if err == nil {
		t.Fatal("expected invalid file type")
	}
	if !strings.Contains(err.Error(), "must have .txt extension") {
		t.Fatalf("expected invalid file type error, got %v", err)
	}
}

// Verifies non-text bytes are rejected instead of being treated as valid domains.
// This keeps --input-file constrained to plain text, newline-delimited domain
// entries.
func TestLoadDomainsFromFileRejectsNonTextContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "domains.txt")
	if err := os.WriteFile(path, []byte{0xff, 0xfe, 0xfd}, 0644); err != nil {
		t.Fatalf("write domain file: %v", err)
	}

	_, err := loadDomainsFromFile(path)
	if err == nil {
		t.Fatal("expected non-text input file error")
	}
	if !strings.Contains(err.Error(), "invalid domain") {
		t.Fatalf("expected invalid domain error, got %v", err)
	}
}
