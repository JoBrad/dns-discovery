package report

import (
	"strings"
	"testing"

	"github.com/jbradley/dns-discovery/internal/discovery"
)

func testDiscoveryResult() *discovery.DiscoveryResult {
	return &discovery.DiscoveryResult{
		Domain: "example.com",
		DNS: discovery.DNSRecords{
			"A":     {"203.0.113.10"},
			"MX":    {"10 mail.example.com", "20 backup.example.com"},
			"NS":    {"ns1.example.net", "ns2.example.net"},
			"TXT":   {"v=spf1 include:_spf.example.com ~all"},
			"SOA":   {"ns1.example.net hostmaster.example.com 1 7200 3600 1209600 3600"},
			"CNAME": {},
		},
		Services: discovery.DetectedServices{
			Email:                []string{"Google Workspace"},
			HostingCDN:           []string{"Cloudflare"},
			VerificationServices: []string{"Stripe"},
		},
		Provider: discovery.ProviderResult{
			Primary: "Route53",
			Counts:  map[string]int{"Route53": 2, "NS1": 2},
			IsSplit: true,
			AllHosts: []string{
				"ns1.example.net",
				"ns2.example.net",
			},
		},
		TLS: []discovery.TLSResult{
			{
				Hostname:      "example.com",
				Reachable:     true,
				TLSVersion:    "TLSv1.3",
				CertValid:     true,
				CertExpiry:    "2026-12-31",
				DaysToExpiry:  30,
				ExpiryWarning: false,
				Issuer:        "Example CA",
			},
		},
		Email: discovery.EmailResult{
			MXRecords: []string{"10 mail.example.com", "20 backup.example.com"},
			SPF: discovery.SPFResult{
				Present: true,
				Record:  "v=spf1 include:_spf.example.com ~all",
				Policy:  "softfail",
			},
			DMARC: discovery.DMARCResult{
				Present: true,
				Record:  "v=DMARC1; p=quarantine",
				Policy:  "quarantine",
			},
			DKIM:      []discovery.DKIMSelector{{Selector: "google"}},
			Score:     4,
			ScoreText: "4/4",
		},
	}
}

func TestGenerateMarkdownIncludesDNSRecordsTable(t *testing.T) {
	content := GenerateMarkdown(testDiscoveryResult())

	checks := []string{
		"## DNS Records",
		"| Type | Records |",
		"| A | 203.0.113.10 |",
		"| MX | 10 mail.example.com<br>20 backup.example.com |",
		"| AAAA | — |",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Fatalf("expected markdown to contain %q", check)
		}
	}
}

func TestGenerateMarkdownIncludesDetectedServices(t *testing.T) {
	content := GenerateMarkdown(testDiscoveryResult())

	checks := []string{
		"## Detected Services",
		"### Email Providers",
		"- Google Workspace",
		"### Hosting & CDN",
		"- Cloudflare",
		"### Verification & SaaS",
		"- Stripe",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Fatalf("expected markdown to contain %q", check)
		}
	}
}

func TestGenerateMarkdownIncludesSplitDNSMXAndTLSFields(t *testing.T) {
	content := GenerateMarkdown(testDiscoveryResult())

	checks := []string{
		"- **Split DNS:**",
		"### Split DNS Providers",
		"| Provider | NS Records |",
		"## Email Security",
		"### MX Records",
		"| Priority | Host |",
		"| 10 | mail.example.com |",
		"## TLS Health",
		"| Hostname | Reachable | Valid | TLS | Expiry | Days | Issuer |",
		"| example.com | ✅ | ✅ | TLSv1.3 | 2026-12-31 | 30d | Example CA |",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Fatalf("expected markdown to contain %q", check)
		}
	}
}
