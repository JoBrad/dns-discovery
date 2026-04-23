package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jbradley/dns-discovery/internal/discovery"
)

func TestSaveReportByFormatWritesExpectedExtensions(t *testing.T) {
	baseDir := t.TempDir()
	result := &discovery.DiscoveryResult{Domain: "example.com"}

	tests := []struct {
		name   string
		format string
		suffix string
	}{
		{name: "markdown", format: FormatMarkdown, suffix: ".md"},
		{name: "json", format: FormatJSON, suffix: ".json"},
		{name: "text", format: FormatText, suffix: ".txt"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path, err := SaveReportByFormat(baseDir, result, tc.format)
			if err != nil {
				t.Fatalf("save report: %v", err)
			}
			if !strings.HasSuffix(path, tc.suffix) {
				t.Fatalf("expected suffix %q, got %q", tc.suffix, path)
			}
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("expected report file to exist: %v", err)
			}
			if filepath.Dir(path) != filepath.Join(baseDir, "example.com") {
				t.Fatalf("expected domain directory path, got %q", path)
			}
		})
	}
}

func TestSaveReportByFormatRejectsInvalidFormat(t *testing.T) {
	_, err := SaveReportByFormat(t.TempDir(), &discovery.DiscoveryResult{Domain: "example.com"}, "xml")
	if err == nil {
		t.Fatal("expected invalid format error")
	}
}
