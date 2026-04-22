package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jbradley/dns-discovery/internal/discovery"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dns-discovery <domain>",
	Short: "DNS zone discovery tool",
	Long: `Discover DNS configuration, provider fingerprint, TLS health,
and email DNS health for any domain.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(args[0])
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(domain string) error {
	domain = strings.TrimSpace(strings.ToLower(domain))

	sep := strings.Repeat("=", 60)

	fmt.Printf("\n%s\n", sep)
	fmt.Printf("  DNS Zone Discovery: %s\n", domain)
	fmt.Printf("%s\n\n", sep)

	// ── DNS enumeration ──────────────────────────────────────────
	fmt.Println("Querying DNS records...")
	records, err := discovery.QueryAllRecords(domain)
	if err != nil {
		return fmt.Errorf("  ERROR: %w", err)
	}

	services := discovery.DetectServices(records)

	// ── Provider fingerprinting ───────────────────────────────────
	nsHosts := records["NS"]
	providerResult := discovery.IdentifyProviders(domain, nsHosts)

	// ── TLS check (A record targets) ─────────────────────────────
	var tlsResults []discovery.TLSResult
	for _, aRecord := range records["A"] {
		tlsResults = append(tlsResults, discovery.CheckTLS(domain))
		_ = aRecord
		break // One TLS check per domain (use the domain name, not raw IP)
	}

	// ── Email health ──────────────────────────────────────────────
	emailResult := discovery.EvaluateEmailHealth(domain)

	// ═══════════════════════════════════════════════════════════════
	// REPORT
	// ═══════════════════════════════════════════════════════════════

	// Nameserver Providers
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

	// DNS Records
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

	// Detected Services
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

	// TLS Health
	fmt.Printf("\n%s\n  TLS Health\n%s\n", sep, sep)
	if len(tlsResults) == 0 {
		fmt.Println("  No A records found — TLS check skipped")
	}
	for _, t := range tlsResults {
		fmt.Printf("  %s:\n", t.Hostname)
		if t.ErrorCategory != "" {
			fmt.Printf("     ✗  %s — %s\n", t.ErrorCategory, t.ErrorDetail)
		} else if t.CertValid {
			warning := ""
			if t.ExpiryWarning {
				warning = fmt.Sprintf("  ⚠  Expires soon!")
			}
			fmt.Printf("     ✓  Valid | %s | Expires: %s (%dd) | Issuer: %s%s\n",
				t.TLSVersion, t.CertExpiry, t.DaysToExpiry, t.Issuer, warning)
		}
	}

	// Email DNS Health
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

	// Executive Summary
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
		fmt.Printf(" — missing:")
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
		fmt.Printf(" %s", strings.Join(missing, ", "))
	}
	fmt.Println()

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
