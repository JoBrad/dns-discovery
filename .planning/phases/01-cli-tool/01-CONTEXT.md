# Phase 1: CLI Tool Foundation - Context

**Gathered:** April 21, 2026  
**Status:** Ready for planning  
**Source:** README.md + Spike findings from spike-findings-dns-discovery skill

<domain>
## Phase Boundary

Build a working CLI tool that accepts a domain name and produces a discovery report showing:
1. DNS zone overview (all 9 record types)
2. Registrar and nameserver identification with friendly provider names
3. Configured services detected from DNS records
4. Health checks for each service category (email DNS health: MX/SPF/DMARC/DKIM; TLS: cert validity, expiry, version)
5. Executive summary with key findings

Deliverable: Executable `dns-discovery <domain>` that works on real domains (github.com, cloudflare.com tested).

## Architecture Pattern (From Spikes)

Four independent discovery pillars, each proven in spike work:

1. **DNS Enumeration** (Spike 001)
   - Query all 9 record types (A, AAAA, MX, NS, TXT, CNAME, SOA, CAA, SRV)
   - Detect services from patterns (MX → email, TXT → SaaS/SPF, CNAME → hosting)
   - Handle NXDOMAIN vs NoAnswer vs Timeout cleanly

2. **Provider Fingerprinting** (Spike 002)
   - Map NS hostnames to friendly provider names (~60 providers in table)
   - Detect self-hosted NS (when NS is subdomain of queried domain)
   - Identify split DNS (multiple providers managing one domain's NS)

3. **TLS Health Checks** (Spike 003)
   - Connect on port 443, verify certificate
   - Distinguish failures: expired (code 10), self-signed (18/19), hostname-mismatch (62)
   - Extract TLS version, cert issuer, days-until-expiry
   - Handle non-HTTPS A records gracefully (timeout/refused)

4. **Email DNS Health** (Spike 004)
   - Query MX, SPF (TXT v=spf1), DMARC (_dmarc.<domain>), DKIM (probe 27 selectors)
   - Classify SPF policy: softfail (~all) vs hardfail (-all)
   - Extract DMARC policy: none/quarantine/reject with reporting status
   - Produce 4-pillar health score (MX/SPF/DMARC/DKIM)

</domain>

<decisions>
## Locked Implementation Decisions

These came from spike validation and are NON-NEGOTIABLE:

### Stack Choices
- **DNS library:** `dnspython` v2.6+ (proven in all 4 spikes, clean exception handling, no external resolvers needed)
- **TLS:** stdlib `ssl` only (no third-party library, full chain validation, OpenSSL error codes work)
- **Pattern matching:** stdlib `re` (sufficient for SPF/DMARC/DKIM syntax)
- **Python version:** 3.14.2 (asyncio.Runner support, modern syntax)
- **Output format:** Markdown to stdout and `output/` directory (per README.md)

### Provider Table
- **Exact provider list:** The ~60-entry provider table from Spike 002 is the canonical source
- **Order matters:** Most-specific patterns first (e.g., "azure-dns.com" before "dns.com")
- **Self-hosted detection logic:** If `ns_lower.endswith("." + domain.lower())` → "Self-hosted (under domain)"
- **Unknown fallback:** Surface raw NS hostname when no pattern matches

### DKIM Probe List
- **Exact selectors:** The 27-selector list from Spike 004 (google, selector1/2, default, mail, k1/k2, s1/s2, dkim, etc.)
- **Coverage:** ~90% of real deployments; some non-standard selectors may be missed (acceptable)

### Service Detection Patterns
- **MX patterns:** 12 common email providers (Google Workspace, Microsoft 365, SendGrid, Mailgun, Zoho, Fastmail, etc.)
- **TXT patterns:** 11 SaaS verification services (Google Search Console, Microsoft 365, Facebook, Stripe, Zoom, etc.)
- **CNAME patterns:** 14 hosting/CDN providers (CloudFront, GitHub Pages, Netlify, Vercel, Heroku, etc.)

### Email Health Scoring
- **4-pillar model:** MX present, SPF present, DMARC present, DKIM selectors found
- **Score format:** "{score}/4" (e.g., github.com → "4/4")
- **SPF policy:** Warn if softfail (~all) vs hardfail (-all); catch insecure +all
- **DMARC policy:** Extract p=none|quarantine|reject; note missing reporting (rua)

### TLS Validation
- **OpenSSL verify codes:**
  - Code 10 → "Certificate EXPIRED"
  - Code 18/19 → "Self-signed certificate (not trusted)"
  - Code 62 → "Hostname mismatch"
- **Expiry warning threshold:** <14 days
- **Timeout:** 5 seconds (reasonable for health checks, not interactive)
- **Days-until-expiry:** Calculate from notAfter field in UTC

### Error Handling
- **NXDOMAIN:** Terminal error (domain doesn't exist) → exit with message
- **NoAnswer/NoNameservers:** Graceful → record type not configured, continue
- **Timeout:** Graceful → record type not queried, continue
- **TLS failures:** Categorized (not fatal) → specific error message, continue to next check
- **Non-HTTPS ports:** Timeout/refused → no TLS info available, continue

### CLI Interface
- **Argument:** Domain name (required)
- **Output:** Printed to stdout in formatted sections
- **Exit code:** 0 if domain exists and checks complete, 1 if domain doesn't exist

## the agent's Discretion

Decisions to be made during implementation:

### CLI Framework
- **Options:** `click` vs custom `argparse` vs other
- **Recommendation from README:** "click library to help make the CLI usage easy"
- **Decision:** To be finalized during Phase 1 planning based on simplicity/features needed

### Module Organization
- **Current structure:** `axeman/` package with `__init__.py`, `core.py`, `certlib.py`
- **Decision:** Organize remaining modules (dns.py, providers.py, email.py, output.py) during planning

### Output Structure (Phase 1)
- **Scope:** Print to stdout with readable formatting (sections, emoji, summaries)
- **Deferred:** File output to `output/` directory deferred to Phase 2

### Progress Display
- **Optional:** Spinner, progress bars, or simple print-as-you-go? (Agent's choice during implementation)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Spike Findings (MANDATORY)
- `.github/skills/spike-findings-dns-discovery/SKILL.md` — Master requirements, stack choices, features overview
  - `references/dns-enumeration.md` — How to query all record types and detect services
  - `references/provider-fingerprinting.md` — NS pattern matching, provider table, self-hosted detection
  - `references/tls-health.md` — Certificate validation, OpenSSL error codes, expiry calculation
  - `references/email-dns-health.md` — MX/SPF/DMARC/DKIM queries and validation rules

### Project Context
- `.planning/PROJECT.md` — Project vision and deliverables
- `.planning/REQUIREMENTS.md` — Phase 1 requirements and UAT criteria
- `.planning/ROADMAP.md` — Full phase breakdown
- `README.md` — User-facing vision (CLI-first, Markdown output, eventual HTML via mk-docs)

</canonical_refs>

<specifics>
## Specific Requirements from README

1. **CLI-first interaction:** User runs `dns-discovery example.com` and gets results
2. **Executive summary:** Provide quick overview at the top
3. **Zone info:** List zone name, registrar, nameservers with friendly provider names
4. **Configured services:** Hosts, redirects, SaaS tools detected from records
5. **Health checks per service:** 
   - Email: MX, SPF, DMARC, DKIM validation
   - Web/CDN: TLS cert valid, not expired, TLS 1.2+
6. **Readable output:** Not raw record dumps (use spike reference patterns)
7. **Library option:** Can be loaded as library (not just CLI) — design for modularity

## Test Domains (From Spike Validation)

**github.com** (complex case):
- DNS: Full zone with CAA records
- NS: Split between NS1 (4 records) and Route53 (4 records)
- Services: Email, Google Workspace, Microsoft 365, Zendesk, Salesforce, Mailchimp in SPF chains
- Email Health: 4/4 score, 6 DKIM selectors (google, selector1, k1, k2, s1, s2)
- TLS: ✓ TLSv1.3, Sectigo Limited, valid 43 days

**cloudflare.com** (complex case):
- DNS: Full zone with CAA records
- NS: Self-branded Cloudflare nameservers (self-hosted)
- Services: Cloudflare Email Security, Mailchimp via Mandrill
- Email Health: 4/4 score, 3 DKIM selectors (k1, s1, mandrill)
- TLS: ✓ TLSv1.3, Google, valid 50 days

Expected output on these domains should match spike validation results exactly.

</specifics>

<deferred>
## Deferred to Future Phases

- **File output:** Writing to `output/` directory (Phase 2)
- **HTML generation:** mk-docs configuration and HTML rendering (Phase 2)
- **Config files:** Domain lists, filter options, custom selectors (Phase 3)
- **Batch processing:** Multiple domains in one run (Phase 3)
- **Interactive mode:** Menu-driven discovery (Phase 3+)
- **Distribution:** PyPI package, Docker image (Phase 4+)
- **Async/concurrent checks:** Parallel health checks (Phase 3, performance)

</deferred>

---

*Phase 1 Context ready for planning. Next: Create detailed PLAN.md with task breakdown, dependencies, and verification criteria.*
