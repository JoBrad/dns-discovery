package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jbradley/dns-discovery/internal/discovery"
	"github.com/jbradley/dns-discovery/internal/report"
)

type OutputFormat string

const (
	OutputMarkdown OutputFormat = "markdown"
	OutputJSON     OutputFormat = "json"
	OutputText     OutputFormat = "text"
)

type RunOptions struct {
	OutputDir   string
	Output      OutputFormat
	Verbose     bool
	LogLocation string
}

func ValidateOutputFormat(value string) (OutputFormat, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch OutputFormat(normalized) {
	case OutputMarkdown, OutputJSON, OutputText:
		return OutputFormat(normalized), nil
	default:
		return "", fmt.Errorf("invalid output format %q: expected one of markdown, json, text", value)
	}
}

type BatchSummary struct {
	Succeeded []string
	Failed    map[string]error
}

var testHookDomainRunner = executeDomain

func RunDomain(domain string, outputDir string) error {
	return testHookDomainRunner(domain, outputDir)
}

func RunBatch(domains []string, outputDir string) BatchSummary {
	return runBatch(domains, outputDir, RunDomain)
}

func (summary BatchSummary) Total() int {
	return len(summary.Succeeded) + len(summary.Failed)
}

func (summary BatchSummary) FailedDomains() []string {
	domains := make([]string, 0, len(summary.Failed))
	for domain := range summary.Failed {
		domains = append(domains, domain)
	}
	sort.Strings(domains)
	return domains
}

func runBatch(domains []string, outputDir string, runner func(string, string) error) BatchSummary {
	summary := BatchSummary{Failed: make(map[string]error)}

	for _, domain := range domains {
		normalized := strings.TrimSpace(strings.ToLower(domain))
		if normalized == "" {
			continue
		}

		if err := runner(normalized, outputDir); err != nil {
			summary.Failed[normalized] = err
			continue
		}

		summary.Succeeded = append(summary.Succeeded, normalized)
	}

	return summary
}

func executeDomain(domain string, outputDir string) error {
	sep := strings.Repeat("=", 60)

	fmt.Printf("\n%s\n", sep)
	fmt.Printf("  DNS Zone Discovery: %s\n", domain)
	fmt.Printf("%s\n\n", sep)

	// DNS enumeration
	fmt.Println("Querying DNS records...")
	records, err := discovery.QueryAllRecords(domain)
	if err != nil {
		return fmt.Errorf("  ERROR: %w", err)
	}

	services := discovery.DetectServices(records)

	// Provider fingerprinting
	nsHosts := records["NS"]
	providerResult := discovery.IdentifyProviders(domain, nsHosts)

	// TLS check (A record targets)
	var tlsResults []discovery.TLSResult
	for _, aRecord := range records["A"] {
		tlsResults = append(tlsResults, discovery.CheckTLS(domain))
		_ = aRecord
		break // One TLS check per domain (use the domain name, not raw IP)
	}

	// Email health
	emailResult := discovery.EvaluateEmailHealth(domain)

	finalResult := &discovery.DiscoveryResult{
		Domain:   domain,
		DNS:      records,
		Services: services,
		Provider: providerResult,
		TLS:      tlsResults,
		Email:    emailResult,
	}

	// Report
	fmt.Printf("\n%s\n  Nameserver Providers\n%s\n", sep, sep)
	fmt.Printf("  Primary provider: %s\n", providerResult.Primary)
	if providerResult.IsSplit {
		fmt.Println("  ⚠  Split DNS detected:")
		for provider, count := range providerResult.Counts {
			fmt.Printf("     • %s (%d NS records)\n", provider, count)
		}
	}
	if len(nsHosts) > 0 {
		fmt.Println("\n  Nameservers:")
		for _, ns := range nsHosts {
			fmt.Printf("     %s\n", ns)
		}
	}

	fmt.Printf("\n%s\n  DNS Records\n%s\n", sep, sep)
	recordOrder := []string{"A", "AAAA", "MX", "NS", "TXT", "CNAME", "SOA", "CAA", "SRV"}
	for _, rtype := range recordOrder {
		recs := records[rtype]
		if len(recs) == 0 {
			fmt.Printf("  %-6s  — not configured\n", rtype)
			continue
		}
		fmt.Printf("  %-6s  (%d record(s)):\n", rtype, len(recs))
		for _, r := range recs {
			fmt.Printf("           %s\n", r)
		}
	}

	fmt.Printf("\n%s\n  Detected Services\n%s\n", sep, sep)
	anyServices := false
	if len(services.Email) > 0 {
		fmt.Println("  Email Providers:")
		for _, s := range services.Email {
			fmt.Printf("     • %s\n", s)
		}
		anyServices = true
	}
	if len(services.HostingCDN) > 0 {
		fmt.Println("  Hosting / CDN:")
		for _, s := range services.HostingCDN {
			fmt.Printf("     • %s\n", s)
		}
		anyServices = true
	}
	if len(services.VerificationServices) > 0 {
		fmt.Println("  Verification / SaaS:")
		for _, s := range services.VerificationServices {
			fmt.Printf("     • %s\n", s)
		}
		anyServices = true
	}
	if !anyServices {
		fmt.Println("  None detected")
	}

	fmt.Printf("\n%s\n  Email DNS Health\n%s\n", sep, sep)
	if len(emailResult.MXRecords) > 0 {
		fmt.Printf("  MX Records (%d):\n", len(emailResult.MXRecords))
		for _, mx := range emailResult.MXRecords {
			fmt.Printf("     %s\n", mx)
		}
	} else {
		fmt.Println("  MX: ✗ No MX records found")
	}
	printSPF(emailResult.SPF)
	printDMARC(emailResult.DMARC)
	printDKIM(emailResult.DKIM)

	scoreIcon := "✗"
	if emailResult.Score == 4 {
		scoreIcon = "✓"
	}
	fmt.Printf("\n  Email Health Score: %s %s\n", emailResult.ScoreText, scoreIcon)

	fmt.Printf("\n%s\n  TLS Health\n%s\n", sep, sep)
	if len(tlsResults) == 0 {
		fmt.Println("  ✗ No A records found to check TLS")
	} else {
		for _, tr := range tlsResults {
			status := "✓"
			if !tr.CertValid || tr.CertExpired || tr.ExpiryWarning {
				status = "⚠"
			}
			if !tr.Reachable {
				status = "✗"
			}
			fmt.Printf("  %s %-20s (expires: %s, issuer: %s)\n", status, tr.Hostname, tr.CertExpiry, tr.Issuer)
			if tr.ErrorDetail != "" {
				fmt.Printf("      Error: %s (%s)\n", tr.ErrorCategory, tr.ErrorDetail)
			}
		}
	}

	fmt.Printf("\n%s\n  Executive Summary\n%s\n", sep, sep)
	fmt.Printf("  Domain:    %s\n", domain)
	fmt.Printf("  DNS:       %d record type(s) configured\n", len(records))
	fmt.Printf("  Provider:  %s", providerResult.Primary)
	if providerResult.IsSplit {
		fmt.Printf(" (split DNS — %d providers)", len(providerResult.Counts))
	}
	fmt.Println()
	if len(tlsResults) > 0 {
		t := tlsResults[0]
		if t.CertValid {
			fmt.Printf("  TLS:       ✓ Valid (%s, %dd until expiry)\n", t.TLSVersion, t.DaysToExpiry)
		} else {
			fmt.Printf("  TLS:       ✗ %s\n", t.ErrorCategory)
		}
	}
	fmt.Printf("  Email:     %s", emailResult.ScoreText)
	if emailResult.Score < 4 {
		missing := []string{}
		if len(emailResult.MXRecords) == 0 {
			missing = append(missing, "MX")
		}
		if !emailResult.SPF.Present {
			missing = append(missing, "SPF")
		}
		if !emailResult.DMARC.Present {
			missing = append(missing, "DMARC")
		}
		if len(emailResult.DKIM) == 0 {
			missing = append(missing, "DKIM")
		}
		fmt.Printf(" — missing: %s", strings.Join(missing, ", "))
	}
	fmt.Println()

	reportPath, err := report.SaveReport(outputDir, finalResult)
	if err != nil {
		fmt.Printf("\n  ⚠  Failed to save report: %v\n", err)
	} else {
		fmt.Printf("\n  ✓ Report saved to: %s\n", reportPath)
	}

	fmt.Printf("\n%s\n\n", sep)
	return nil
}

func printSPF(spf discovery.SPFResult) {
	if !spf.Present {
		fmt.Println("  SPF:  ✗ No SPF record found")
		return
	}
	icon := "✓"
	if spf.Insecure || len(spf.Issues) > 0 {
		icon = "⚠"
	}
	fmt.Printf("  SPF:  %s %s (%s)\n", icon, spf.Record, spf.Policy)
	for _, issue := range spf.Issues {
		fmt.Printf("          ⚠  %s\n", issue)
	}
}

func printDMARC(dmarc discovery.DMARCResult) {
	if !dmarc.Present {
		fmt.Println("  DMARC: ✗ No DMARC record found")
		return
	}
	icon := "✓"
	if dmarc.Policy == "none" {
		icon = "⚠"
	}
	rua := ""
	if !dmarc.HasAggReports {
		rua = " (no rua= reporting)"
	}
	fmt.Printf("  DMARC: %s p=%s%s\n", icon, dmarc.Policy, rua)
}

func printDKIM(selectors []discovery.DKIMSelector) {
	if len(selectors) == 0 {
		fmt.Println("  DKIM:  ✗ No DKIM selectors found (probed 27 common selectors)")
		return
	}
	names := make([]string, 0, len(selectors))
	for _, s := range selectors {
		name := s.Selector
		if s.Revoked {
			name += "(revoked)"
		}
		names = append(names, name)
	}
	fmt.Printf("  DKIM:  ✓ %d selector(s) found: %s\n", len(selectors), strings.Join(names, ", "))
}
