package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jbradley/dns-discovery/internal/discovery"
)

func GenerateText(res *discovery.DiscoveryResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("DNS Discovery Report: %s\n", res.Domain))
	sb.WriteString(strings.Repeat("=", 60) + "\n")
	sb.WriteString(fmt.Sprintf("Primary Provider: %s\n", res.Provider.Primary))
	sb.WriteString(fmt.Sprintf("Email Health: %s\n", res.Email.ScoreText))
	sb.WriteString(fmt.Sprintf("TLS Hosts: %d\n\n", len(res.TLS)))

	sb.WriteString("Nameservers:\n")
	for _, ns := range res.Provider.AllHosts {
		sb.WriteString(fmt.Sprintf("- %s\n", ns))
	}
	sb.WriteString("\n")

	sb.WriteString("DNS Records:\n")
	for _, rtype := range dnsRecordOrder {
		records := res.DNS[rtype]
		if len(records) == 0 {
			sb.WriteString(fmt.Sprintf("%s: (none)\n", rtype))
			continue
		}
		sb.WriteString(fmt.Sprintf("%s:\n", rtype))
		for _, record := range records {
			sb.WriteString(fmt.Sprintf("  - %s\n", record))
		}
	}
	sb.WriteString("\n")

	hasServices := len(res.Services.Email) > 0 || len(res.Services.HostingCDN) > 0 || len(res.Services.VerificationServices) > 0
	if hasServices {
		sb.WriteString("Detected Services:\n")
		writeServiceSection(&sb, "Email", res.Services.Email)
		writeServiceSection(&sb, "Hosting/CDN", res.Services.HostingCDN)
		writeServiceSection(&sb, "Verification", res.Services.VerificationServices)
		sb.WriteString("\n")
	}

	sb.WriteString("TLS Summary:\n")
	for _, t := range res.TLS {
		status := "ok"
		if !t.Reachable {
			status = t.ErrorCategory
		} else if !t.CertValid || t.CertExpired {
			status = "invalid"
		} else if t.ExpiryWarning {
			status = "expiring_soon"
		}
		sb.WriteString(fmt.Sprintf("- %s: %s (issuer=%s, expiry=%s)\n", t.Hostname, status, t.Issuer, t.CertExpiry))
	}

	if res.Provider.IsSplit {
		sb.WriteString("\nSplit DNS Providers:\n")
		providerNames := sortedProviderNames(res.Provider.Counts)
		for _, name := range providerNames {
			sb.WriteString(fmt.Sprintf("- %s (%d)\n", name, res.Provider.Counts[name]))
		}
	}

	return sb.String()
}

func writeServiceSection(sb *strings.Builder, heading string, values []string) {
	if len(values) == 0 {
		return
	}
	sb.WriteString(fmt.Sprintf("%s:\n", heading))
	for _, value := range values {
		sb.WriteString(fmt.Sprintf("  - %s\n", value))
	}
}

func sortedProviderNames(counts map[string]int) []string {
	names := make([]string, 0, len(counts))
	for name := range counts {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
