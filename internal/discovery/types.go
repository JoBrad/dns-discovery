package discovery

// DNSRecords maps record type name to slice of string representations.
type DNSRecords map[string][]string

// DetectedServices holds services discovered via DNS record pattern matching.
type DetectedServices struct {
	Email                []string
	HostingCDN           []string
	VerificationServices []string
}

// ProviderResult holds NS provider fingerprinting results.
type ProviderResult struct {
	Primary  string
	Counts   map[string]int
	IsSplit  bool
	AllHosts []string
}

// TLSResult holds TLS health check results for a single hostname.
type TLSResult struct {
	Hostname      string
	Reachable     bool
	TLSVersion    string
	CertValid     bool
	CertExpired   bool
	CertExpiry    string // "2026-06-03"
	DaysToExpiry  int
	ExpiryWarning bool   // < 14 days
	Issuer        string
	// ErrorCategory: "EXPIRED", "SELF_SIGNED", "HOSTNAME_MISMATCH", "TIMEOUT", "REFUSED", "DNS_ERROR", ""
	ErrorCategory string
	ErrorDetail   string
}

// SPFResult holds SPF record analysis.
type SPFResult struct {
	Present  bool
	Record   string
	Policy   string // "hardfail", "softfail", "neutral", "pass", "unknown"
	Insecure bool   // +all present
	Issues   []string
}

// DMARCResult holds DMARC record analysis.
type DMARCResult struct {
	Present        bool
	Record         string
	Policy         string // "none", "quarantine", "reject"
	HasAggReports  bool
	HasForensics   bool
}

// DKIMSelector is a discovered DKIM selector.
type DKIMSelector struct {
	Selector string
	KeyType  string
	Revoked  bool
}

// EmailResult holds the full email DNS health evaluation.
type EmailResult struct {
	MXRecords []string // formatted as "priority hostname"
	SPF       SPFResult
	DMARC     DMARCResult
	DKIM      []DKIMSelector
	Score     int    // 0–4
	ScoreText string // "4/4"
}

// DiscoveryResult is the top-level aggregated result for a domain.
type DiscoveryResult struct {
	Domain   string
	DNS      DNSRecords
	Services DetectedServices
	Provider ProviderResult
	TLS      []TLSResult
	Email    EmailResult
}
