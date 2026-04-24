# Spike Wrap-Up Summary

**Date:** April 23, 2026
**Spikes processed:** 4
**Feature areas:** 4
**Skill output:** `./.github/skills/spike-findings-dns-discovery/`

## Processed Spikes

| # | Name | Type | Verdict | Feature Area | Key Finding |
|---|------|------|---------|--------------|-------------|
| 001 | dns-record-enumeration | standard | VALIDATED | DNS Enumeration | All 9 record types query cleanly; SPF `include:` chains reveal SaaS tools |
| 002 | ns-registrar-fingerprinting | standard | VALIDATED ✓ | Provider Fingerprinting | Pattern table covers ~60 providers; self-hosted & split DNS detection |
| 003 | tls-health-check | standard | VALIDATED | TLS Health Checks | Go stdlib TLS/X509 checks distinguish cert failure and connectivity modes |
| 004 | email-dns-health | standard | VALIDATED ✓ | Email DNS Health | 27-selector DKIM probe, selector names reveal providers, 4-pillar health score |

## Key Findings

### Stack Choices (Locked)
- DNS: `miekg/dns`
- TLS: Go stdlib `crypto/tls`, `crypto/x509`, `net`
- Pattern matching: Go stdlib `regexp` and `strings`
- CLI and execution: Go + Cobra

### Design Decisions Validated

1. **Service Detection** — Pattern matching on MX/TXT/CNAME records is high-signal
   - MX: Email provider detection
   - TXT: SPF/DMARC/SaaS verification services
   - CNAME: CDN/hosting detection
   - SPF `include:` chains reveal deeper SaaS footprint beyond MX

2. **Provider Fingerprinting** — Substring matching against ordered pattern table
   - Most-specific patterns first (e.g., "azure-dns.com" before "dns.com")
   - ~60 providers covers vast majority of real-world deployments
   - Self-hosted detection: if NS is subdomain of queried domain, flag as "Self-hosted (under domain)"
   - Split DNS: when multiple providers manage NS records for redundancy

3. **TLS Health** — Stdlib `ssl` sufficient for all scenarios
   - `tls.DialWithDialer` + `ServerName` performs verification during handshake
   - Error mapping: EXPIRED, SELF_SIGNED, HOSTNAME_MISMATCH, TIMEOUT, REFUSED, DNS_ERROR, TLS_ERROR
   - Days-until-expiry calculation for proactive warnings (<14d threshold)
   - Graceful handling of non-HTTPS hosts (timeout, connection refused)

4. **Email DNS Health** — Four-pillar validation
   - **MX**: Must be present for email delivery
   - **SPF**: Must have `v=spf1` record with `all` mechanism
   - **DMARC**: At `_dmarc.<domain>` with policy (none/quarantine/reject) and reporting
   - **DKIM**: Probe-based discovery of ~27 common selectors (covers Microsoft 365, Google, Mandrill, Amazon SES, Zoho, Fastmail, etc.)
   - Selector names reveal email provider (e.g., `mandrill._domainkey` → Mailchimp)

### Constraints & Gotchas

- **TXT record joining**: Long SPF records are split across multiple TXT strings in DNS response; must join before pattern matching
- **DKIM probe-based**: No way to enumerate all DKIM selectors; probe covers ~90%, some may be missed
- **Split DNS is real**: github.com has NS1 + Route53 nameservers for redundancy; tool must detect and report both
- **Non-HTTP hosts**: A records may point to mail servers, game servers, etc.; TLS check times out gracefully without blocking
- **CAA records cross-reference**: CAA policy (permitted issuers) can be compared against actual cert issuer for validation

### Refresh Notes

- Existing spike findings were re-packaged against current production code paths:
   - `internal/discovery/dns.go`
   - `internal/discovery/providers.go`
   - `internal/discovery/tls.go`
   - `internal/discovery/email.go`
- Skill `sources/` now stores per-spike folders containing the spike README plus matching Go source snapshots.

## Validated Test Domains

| Domain | DNS | NS Provider | Email | TLS |
|--------|-----|-------------|-------|-----|
| github.com | Full zone | NS1 + Route53 (split) | 4/4 (6 DKIM) | ✓ TLSv1.3, Sectigo |
| cloudflare.com | Full zone | Cloudflare (self-branded) | 4/4 (3 DKIM incl. Mandrill) | ✓ TLSv1.3, Google |
| google.com | Full zone | Self-hosted | — | ✓ TLSv1.3, Google |
| badssl.com endpoints | — | — | — | ✗ (expired, self-signed, hostname mismatch) |

## Next Steps

The spike findings skill is ready for implementation planning. The 4 reference files in `references/` provide step-by-step implementation blueprints with:
- Requirements (non-negotiable design decisions)
- How to build it (code patterns and recipes)
- What to avoid (gotchas and dead ends)
- Constraints (library limitations, version requirements)
- Origin (spike source files available in `sources/`)

All foundation work complete. Ready for:
- `/gsd-plan-phase 1` — Begin full implementation of the CLI tool
- `/gsd-spike` (frontier mode) — Identify additional spikes if needed
- Direct implementation if ready to proceed
