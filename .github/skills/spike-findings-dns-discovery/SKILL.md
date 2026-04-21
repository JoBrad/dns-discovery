---
name: spike-findings-dns-discovery
description: Implementation blueprint from spike experiments. Requirements, proven patterns, and verified knowledge for building the DNS zone discovery tool. Auto-loaded during implementation work.
---

<context>
## Project: dns-discovery

Build a DNS zone discovery tool that, given a domain name, produces an executive summary of the zone's configuration. Output includes registrar and nameserver identification (with friendly provider names), a list of configured services/hosts/redirects, and health checks for each service — validating email DNS (MX/SPF/DKIM/DMARC), TLS certificate validity, and minimum TLS 1.2+ for web-facing targets.

Spike sessions wrapped: April 21, 2026
</context>

<requirements>
## Non-Negotiable Design Decisions

These requirements emerged from spike validation and **must** be honored in the real implementation:

- Must work from a single domain input (e.g., `example.com`)
- Must identify registrar and nameserver providers by friendly name, not just raw hostnames
- Must check email DNS health: MX records present, SPF (v=spf1), DKIM (_domainkey), DMARC (_dmarc)
- Must check TLS health for A/CNAME targets: valid cert, not expired, TLS 1.2+
- Must produce a readable output (not just raw record dumps)
- Use `dnspython` for all DNS queries (proven in 4 spikes)
- Use stdlib `ssl` for TLS validation (no third-party TLS library)
- Use stdlib `re` for pattern matching (SPF, DMARC, DKIM)
- Gracefully handle non-HTTPS hosts (timeout, connection refused)
- Support split DNS detection (multiple providers for one domain)
- DKIM discovery via probe of ~27 common selectors (no enumeration API exists)
- Service detection from record patterns (MX, TXT, CNAME) with known provider lookup tables
</requirements>

<findings_index>
## Feature Areas

| Area | Reference | Key Spike Finding |
|------|-----------|-------------------|
| DNS Enumeration | references/dns-enumeration.md | All 9 record types (A, AAAA, MX, NS, TXT, CNAME, SOA, CAA, SRV) query cleanly with `dnspython`; SPF `include:` chains reveal SaaS tools |
| Provider Fingerprinting | references/provider-fingerprinting.md | Pattern table covers ~60 DNS providers; self-hosted detection handles tech giants; split DNS detection identifies redundancy architecture |
| TLS Health Checks | references/tls-health.md | Stdlib `ssl` sufficient; OpenSSL verify codes distinguish expired vs self-signed vs hostname mismatch; days-until-expiry enables proactive warnings |
| Email DNS Health | references/email-dns-health.md | Probe-based DKIM discovery (27 selectors) covers 90%+ of real deployments; selector names reveal provider (mandrill→Mailchimp); 4-pillar score (MX/SPF/DMARC/DKIM) |

## Source Files

Original spike source files are preserved in `sources/` for complete reference and testing.

### Validated Domains

All spikes tested on:
- **github.com** — Complex: split DNS (NS1 + Route53), multiple email providers, 6 DKIM selectors
- **cloudflare.com** — Complex: self-branded NS records, multiple SaaS services, Mandrill/Mailchimp detected
- **google.com** — Self-hosted DNS infrastructure
- **badssl.com endpoints** — TLS failure modes (expired, self-signed, hostname mismatch)
</findings_index>

<metadata>
## Processed Spikes

- 001-dns-record-enumeration (VALIDATED)
- 002-ns-registrar-fingerprinting (VALIDATED)
- 003-tls-health-check (VALIDATED)
- 004-email-dns-health (VALIDATED)

**All 4 spikes VALIDATED ✓** — No blocking issues, all proof-of-concepts working with real domains.
</metadata>
