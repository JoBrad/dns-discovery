package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbradley/dns-discovery/internal/discovery"
)

const (
	FormatMarkdown = "markdown"
	FormatJSON     = "json"
	FormatText     = "text"
)

func SaveReportByFormat(baseDir string, res *discovery.DiscoveryResult, format string) (string, error) {
	normalized, err := normalizeFormat(format)
	if err != nil {
		return "", err
	}

	var content string
	var reportName string

	switch normalized {
	case FormatMarkdown:
		content = GenerateMarkdown(res)
		reportName = "report.md"
	case FormatJSON:
		content, err = GenerateJSON(res)
		if err != nil {
			return "", err
		}
		reportName = "report.json"
	case FormatText:
		content = GenerateText(res)
		reportName = "report.txt"
	default:
		return "", fmt.Errorf("unsupported output format %q", format)
	}

	domainDir := filepath.Join(baseDir, res.Domain)
	if err := os.MkdirAll(domainDir, 0755); err != nil {
		return "", err
	}

	reportPath := filepath.Join(domainDir, reportName)
	if err := os.WriteFile(reportPath, []byte(content), 0644); err != nil {
		return "", err
	}

	return reportPath, nil
}

func normalizeFormat(format string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(format))
	if normalized == "" {
		normalized = FormatMarkdown
	}

	switch normalized {
	case FormatMarkdown, FormatJSON, FormatText:
		return normalized, nil
	default:
		return "", fmt.Errorf("invalid output format %q: expected markdown, json, or text", format)
	}
}
