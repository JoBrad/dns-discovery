package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jbradley/dns-discovery/internal/discovery"
)

// dnsRecordOrder defines the canonical display order for DNS record types.
var dnsRecordOrder = []string{"A", "AAAA", "MX", "NS", "TXT", "CNAME", "SOA", "CAA", "SRV"}

// GenerateMarkdown generates a Markdown report from discovery results.
func GenerateMarkdown(res *discovery.DiscoveryResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# DNS Discovery Report: %s\n\n", res.Domain))
	sb.WriteString(fmt.Sprintf("Generated on: %s\n\n", time.Now().Format(time.RFC1123)))

	// --- Executive Summary ---
	sb.WriteString("## Executive Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Primary Provider:** %s\n", res.Provider.Primary))
	if res.Provider.IsSplit {
		providerNames := sortedKeys(res.Provider.Counts)
		parts := make([]string, 0, len(providerNames))
		for _, p := range providerNames {
			parts = append(parts, fmt.Sprintf("%s (%d)", p, res.Provider.Counts[p]))
		}
		sb.WriteString(fmt.Sprintf("- **Split DNS:** %s\n", strings.Join(parts, ", ")))
	}
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

	// --- Infrastructure & DNS ---
	sb.WriteString("## Infrastructure & DNS\n\n")
	sb.WriteString(fmt.Sprintf("### Nameservers (%s)\n", res.Provider.Primary))
	for _, ns := range res.Provider.AllHosts {
		sb.WriteString(fmt.Sprintf("- %s\n", ns))
	}
	sb.WriteString("\n")
	if res.Provider.IsSplit {
		sb.WriteString("### Split DNS Providers\n\n")
		sb.WriteString("| Provider | NS Records |\n")
		sb.WriteString("|----------|------------|\n")
		for _, p := range sortedKeys(res.Provider.Counts) {
			sb.WriteString(fmt.Sprintf("| %s | %d |\n", p, res.Provider.Counts[p]))
		}
		sb.WriteString("\n")
	}

	// --- DNS Records ---
	sb.WriteString("## DNS Records\n\n")
	sb.WriteString("| Type | Records |\n")
	sb.WriteString("|------|---------|\n")
	for _, rtype := range dnsRecordOrder {
		vals := res.DNS[rtype]
		if len(vals) == 0 {
			sb.WriteString(fmt.Sprintf("| %s | — |\n", rtype))
		} else {
			sb.WriteString(fmt.Sprintf("| %s | %s |\n", rtype, strings.Join(vals, "<br>")))
		}
	}
	sb.WriteString("\n")

	// --- Detected Services ---
	hasServices := len(res.Services.Email) > 0 || len(res.Services.HostingCDN) > 0 || len(res.Services.VerificationServices) > 0
	if hasServices {
		sb.WriteString("## Detected Services\n\n")
		if len(res.Services.Email) > 0 {
			sb.WriteString("### Email Providers\n")
			for _, s := range res.Services.Email {
				sb.WriteString(fmt.Sprintf("- %s\n", s))
			}
			sb.WriteString("\n")
		}
		if len(res.Services.HostingCDN) > 0 {
			sb.WriteString("### Hosting & CDN\n")
			for _, s := range res.Services.HostingCDN {
				sb.WriteString(fmt.Sprintf("- %s\n", s))
			}
			sb.WriteString("\n")
		}
		if len(res.Services.VerificationServices) > 0 {
			sb.WriteString("### Verification & SaaS\n")
			for _, s := range res.Services.VerificationServices {
				sb.WriteString(fmt.Sprintf("- %s\n", s))
			}
			sb.WriteString("\n")
		}
	}

	// --- Email Security ---
	sb.WriteString("## Email Security\n\n")

	// MX Records subsection
	if len(res.Email.MXRecords) > 0 {
		sb.WriteString("### MX Records\n\n")
		sb.WriteString("| Priority | Host |\n")
		sb.WriteString("|----------|------|\n")
		for _, mx := range res.Email.MXRecords {
			parts := strings.SplitN(mx, " ", 2)
			if len(parts) == 2 {
				sb.WriteString(fmt.Sprintf("| %s | %s |\n", parts[0], parts[1]))
			} else {
				sb.WriteString(fmt.Sprintf("| — | %s |\n", mx))
			}
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString("**MX Records:** No MX records configured\n\n")
	}

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
		selectors := make([]string, 0, len(res.Email.DKIM))
		for _, s := range res.Email.DKIM {
			selectors = append(selectors, s.Selector)
		}
		sb.WriteString(fmt.Sprintf("- **DKIM:** Found selectors: %s\n", strings.Join(selectors, ", ")))
	} else {
		sb.WriteString("- **DKIM:** ❌ No selectors found\n")
	}
	sb.WriteString(fmt.Sprintf("\n**Email Health Score:** %s\n\n", res.Email.ScoreText))

	// --- TLS Health ---
	sb.WriteString("## TLS Health\n\n")
	sb.WriteString("| Hostname | Reachable | Valid | TLS | Expiry | Days | Issuer |\n")
	sb.WriteString("|----------|-----------|-------|-----|--------|------|--------|\n")
	for _, t := range res.TLS {
		reachable := "✅"
		if !t.Reachable {
			reachable = "❌"
		}
		valid := "✅"
		if !t.CertValid {
			valid = "❌"
		}
		tlsVer := t.TLSVersion
		if tlsVer == "" {
			tlsVer = "—"
		}
		expiry := t.CertExpiry
		if expiry == "" {
			expiry = "—"
		}
		days := "—"
		if t.Reachable && t.CertValid {
			if t.ExpiryWarning {
				days = fmt.Sprintf("⚠️ %dd", t.DaysToExpiry)
			} else {
				days = fmt.Sprintf("%dd", t.DaysToExpiry)
			}
		} else if !t.Reachable && t.ErrorCategory != "" {
			days = t.ErrorCategory
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |\n",
			t.Hostname, reachable, valid, tlsVer, expiry, days, t.Issuer))
	}
	sb.WriteString("\n")

	return sb.String()
}

// sortedKeys returns a sorted slice of keys from a map[string]int.
func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SaveReport writes the Markdown report to the output directory.
func SaveReport(baseDir string, res *discovery.DiscoveryResult) (string, error) {
	return SaveReportByFormat(baseDir, res, FormatMarkdown)
}
