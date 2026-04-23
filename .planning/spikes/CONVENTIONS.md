# Spike Conventions

Patterns and stack choices established across spike sessions. New spikes follow these unless the question requires otherwise.

## Stack

Language: Go 1.26+
CLI framework: `spf13/cobra`
DNS library: `miekg/dns`
TLS/SSL: stdlib `crypto/tls`, `crypto/x509`, `net`
Pattern matching: stdlib `regexp` and `strings`

## Structure

Production code layout used for spike-derived features:
- `cmd/dns-discovery/main.go` for CLI argument/config orchestration
- `internal/discovery/` for DNS, provider, TLS, and email checks
- `internal/app/run.go` for pipeline execution and console output
- `internal/report/markdown.go` for report generation

## Patterns

DNS queries:
- Query fixed record set: A, AAAA, MX, NS, TXT, CNAME, SOA, CAA, SRV
- Retry truncated UDP responses over TCP
- Treat NXDOMAIN on A query as terminal domain-not-found
- Keep missing record types as empty/not-configured, not fatal

Service and provider detection:
- Pattern-match MX/TXT/CNAME values for service detection
- Keep provider pattern tables ordered from most-specific to least-specific
- Keep explicit Unknown fallbacks for unmapped providers
- Surface split-DNS when multiple provider fingerprints are present

TLS checks:
- Use `tls.DialWithDialer` with `ServerName`
- Extract cert expiry/issuer from peer cert
- Classify failures into stable categories: EXPIRED, SELF_SIGNED, HOSTNAME_MISMATCH, TIMEOUT, REFUSED, DNS_ERROR, TLS_ERROR
- Keep timeout short (5s)

Email checks:
- Score 4 pillars: MX, SPF, DMARC, DKIM
- Probe 27 common DKIM selectors (no global enumeration exists)
- Parse SPF/DMARC with regex and emit human-readable policy findings

## Tools and Libraries

Preferred:
- `github.com/miekg/dns` for DNS transport and record decoding
- `github.com/spf13/cobra` for CLI flag/args handling
- Go stdlib crypto/net/time/regexp for security and parsing

Avoid:
- Hard-coding provider assumptions outside lookup tables
- Failing the whole run on single-record or single-domain errors in batch mode
- Introducing extra dependencies when stdlib already covers TLS/parsing needs

## Testing Domains

Use these for live validation:
- `github.com` for split NS and rich email/service footprint
- `cloudflare.com` for self-branded NS and strict email policies
- `google.com` for self-hosted DNS patterns
- badssl endpoints for TLS negatives:
  - `expired.badssl.com`
  - `self-signed.badssl.com`
  - `wrong.host.badssl.com`
