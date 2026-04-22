package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jbradley/dns-discovery/internal/discovery"
)

// GenerateMarkdown generates a Markdown report from discovery results.
func GenerateMarkdown(res *discovery.DiscoveryResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# DNS Discovery Report: %s\n\n", res.Domain))
	sb.WriteString(fmt.Sprintf("Generated on: %s\n\n", time.Now().Format(time.RFC1123)))

	// Executive Summary
	sb.WriteString("## Executive Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Primary Provider:** %s\n", res.Provider.Primary))
	
	emailStatus := "✅ Healthy"
	if res.Email.Score < 3 {
		emailStatus = "⚠️ Issues Found"
	}
	sb.WriteString(fmt.Sprintf("- **Email Health:** %s (%s)\n", emailStatus, res.Email.ScoreText))
	
	tlsTotal := len(res.TLS)
	tlsHealthy := 0
	for _, t := range res.TLS {
		if t.Reachable && t.CertValid && !t.CertExpired {
			tlsHealthy++
		}
	}
	sb.WriteString(fmt.Sprintf("- **TLS Health:** %d/%d hosts healthy\n\n", tlsHealthy, tlsTotal))

	// Infrastructure & DNS
	sb.WriteString("## Infrastructure & DNS\n\n")
	sb.WriteString(fmt.Sprintf("### Nameservers (%s)\n", res.Provider.Primary))
	for _, ns := range res.Provider.AllHosts {
		sb.WriteString(fmt.Sprintf("- %s\n", ns))
	}
	sb.WriteString("\n")

	if len(res.Services.HostingCDN) > 0 {
		sb.WriteString("### Hosting & CDN\n")
		for _, s := range res.Services.HostingCDN {
			sb.WriteString(fmt.Sprintf("- %s\n", s))
		}
		sb.WriteString("\n")
	}

	// Email Security
	sb.WriteString("## Email Security\n\n")
	if res.Email.SPF.Present {
		sb.WriteString(fmt.Sprintf("- **SPF:** `%s` (%s)\n", res.Email.SPF.Record, res.Email.SPF.Policy))
		for _, issue := range res.Email.SPF.Issues {
			sb.WriteString(fmt.Sprintf("  - ⚠️ %s\n", issue))
		}
	} else {
		sb.WriteString("- **SPF:** ❌ Not found\n")
	}

	if res.Email.DMARC.Present {
		sb.WriteString(fmt.Sprintf("- **DMARC:** `%s` (p=%s)\n", res.Email.DMARC.Record, res.Email.DMARC.Policy))
	} else {
		sb.WriteString("- **DMARC:** ❌ Not found\n")
	}

	if len(res.Email.DKIM) > 0 {
		selectors := []string{}
		for _, s := range res.Email.DKIM {
			selectors = append(selectors, s.Selector)
		}
		sb.WriteString(fmt.Sprintf("- **DKIM:** Found selectors: %s\n", strings.Join(selectors, ", ")))
	} else {
		sb.WriteString("- **DKIM:** ❌ No selectors found\n")
	}
	sb.WriteString("\n")

	// TLS Health
	sb.WriteString("## TLS Health\n\n")
	sb.WriteString("| Hostname | Reachable | Valid | Expiry | Issuer |\n")
	sb.WriteString("|----------|-----------|-------|--------|--------|\n")
	for _, t := range res.TLS {
		reachable := "✅"
		if !t.Reachable {
			reachable = "❌"
		}
		valid := "✅"
		if !t.CertValid {
			valid = "❌"
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n", 
			t.Hostname, reachable, valid, t.CertExpiry, t.Issuer))
	}
	sb.WriteString("\n")

	return sb.String()
}

// SaveReport writes the Markdown report to the output directory.
func SaveReport(baseDir string, res *discovery.DiscoveryResult) (string, error) {
	domainDir := filepath.Join(baseDir, res.Domain)
	if err := os.MkdirAll(domainDir, 0755); err != nil {
		return "", err
	}

	content := GenerateMarkdown(res)
	reportPath := filepath.Join(domainDir, "report.md")
	if err := os.WriteFile(reportPath, []byte(content), 0644); err != nil {
		return "", err
	}

	return reportPath, nil
}
