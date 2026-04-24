package discovery

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func requireDNSNetwork(t *testing.T) {
	t.Helper()
	_, err := QueryAllRecords("github.com")
	if err != nil {
		t.Skipf("skipping network-dependent test: %v", err)
	}
}

func TestDetectServicesClassifiesPatterns(t *testing.T) {
	records := DNSRecords{
		"MX":    {"10 aspmx.l.google.com"},
		"TXT":   {"google-site-verification=abc123", "v=spf1 include:_spf.google.com ~all"},
		"CNAME": {"abc123.vercel.app"},
	}

	services := DetectServices(records)
	if len(services.Email) == 0 || services.Email[0] != "Google Workspace" {
		t.Fatalf("expected email provider detection, got %#v", services.Email)
	}
	if len(services.HostingCDN) == 0 || services.HostingCDN[0] != "Vercel" {
		t.Fatalf("expected hosting/CDN detection, got %#v", services.HostingCDN)
	}
	if len(services.VerificationServices) == 0 {
		t.Fatalf("expected verification services, got %#v", services.VerificationServices)
	}
}

func TestIdentifyProvidersDetectsSplitAndSelfHosted(t *testing.T) {
	nsHosts := []string{
		"ns-520.awsdns-01.net",
		"dns1.p08.nsone.net",
		"ns1.example.com",
	}

	result := IdentifyProviders("example.com", nsHosts)
	if !result.IsSplit {
		t.Fatal("expected split DNS detection")
	}
	if result.Counts["AWS Route 53"] == 0 {
		t.Fatalf("expected Route 53 count, got %#v", result.Counts)
	}
	if result.Counts["NS1 / IBM NS1 Connect"] == 0 {
		t.Fatalf("expected NS1 count, got %#v", result.Counts)
	}
	if result.Counts["Self-hosted (under example.com)"] == 0 {
		t.Fatalf("expected self-hosted count, got %#v", result.Counts)
	}
}

func TestQueryAllRecordsReturnsDataForKnownDomain(t *testing.T) {
	requireDNSNetwork(t)

	records, err := QueryAllRecords("github.com")
	if err != nil {
		t.Fatalf("query all records: %v", err)
	}
	if len(records["A"]) == 0 {
		t.Fatalf("expected A records for github.com, got %#v", records["A"])
	}
	if len(records["NS"]) == 0 {
		t.Fatalf("expected NS records for github.com, got %#v", records["NS"])
	}
}

func TestQueryAllRecordsReturnsNXDOMAINForInvalidDomain(t *testing.T) {
	requireDNSNetwork(t)

	domain := fmt.Sprintf("phase1-nyquist-%d.invalid", time.Now().UnixNano())
	_, err := QueryAllRecords(domain)
	if err != ErrNXDOMAIN {
		t.Fatalf("expected ErrNXDOMAIN for %s, got %v", domain, err)
	}
}

func TestCheckTLSReturnsMetadataForHealthyHost(t *testing.T) {
	requireDNSNetwork(t)

	result := CheckTLS("github.com")
	if !result.Reachable {
		t.Fatalf("expected github.com reachable TLS endpoint, got %+v", result)
	}
	if result.TLSVersion == "" {
		t.Fatalf("expected TLS version, got %+v", result)
	}
}

func TestCheckTLSClassifiesExpiredEndpoint(t *testing.T) {
	requireDNSNetwork(t)

	result := CheckTLS("expired.badssl.com")
	if result.CertValid {
		t.Fatalf("expected expired.badssl.com to be invalid, got %+v", result)
	}
	allowed := map[string]bool{
		"EXPIRED":      true,
		"CERT_INVALID": true,
		"TLS_ERROR":    true,
	}
	if !allowed[result.ErrorCategory] {
		t.Fatalf("unexpected error category: %q (%+v)", result.ErrorCategory, result)
	}
}

func TestEvaluateEmailHealthReturnsScoreAndPillars(t *testing.T) {
	requireDNSNetwork(t)

	result := EvaluateEmailHealth("github.com")
	if !strings.HasSuffix(result.ScoreText, "/4") {
		t.Fatalf("expected score text to end with /4, got %q", result.ScoreText)
	}
	if result.Score < 0 || result.Score > 4 {
		t.Fatalf("expected score in [0,4], got %d", result.Score)
	}
	if len(result.MXRecords) == 0 {
		t.Fatalf("expected MX records for github.com, got %#v", result.MXRecords)
	}
}
