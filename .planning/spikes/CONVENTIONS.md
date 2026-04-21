# Spike Conventions

Patterns and stack choices established across spike sessions. New spikes follow these unless the question requires otherwise.

## Stack

**Language:** Python 3.14+
**DNS library:** `dnspython` — comprehensive, actively maintained, handles all standard record types
**TLS/SSL:** stdlib `ssl` + stdlib `socket` — no third-party needed
**Pattern matching:** stdlib `re` for SPF/DMARC parsing

## Structure

Each spike lives in `.planning/spikes/NNN-descriptive-name/`:
- `spike.py` — executable spike code
- `README.md` — frontmatter + investigation trail + results

## Patterns

**DNS queries:** Use `dns.resolver.resolve(domain, record_type)` directly. Handle exceptions:
- `dns.resolver.NXDOMAIN` — domain doesn't exist
- `dns.resolver.NoAnswer` — record type not configured
- `dns.resolver.NoNameservers` — resolution failed
- Return empty list on error (except NXDOMAIN, which is terminal)

**Service detection:** Pattern-match against known provider substrings. Keep lookup tables as ordered lists (most-specific first). Gracefully fall back to "Unknown" for unmapped providers.

**TLS connections:** Use `ssl.create_default_context()` for full chain validation. Extract certs with `ssock.getpeercert()`. Map `ssl.SSLCertVerificationError.verify_code` to human-readable failure reasons:
- Code 10 → expired
- Code 18/19 → self-signed/untrusted
- Code 62 → hostname mismatch

**TXT record parsing:** TXT records can be split across multiple DNS strings (e.g. long SPF). Join them before regex matching.

**DKIM discovery:** Probe 20-30 common selectors (google, selector1, selector2, k1, k2, s1, s2, default, mail, mandrill, etc.). No way to enumerate, so probe is inherently incomplete but covers 90%+ of real deployments.

## Tools & Libraries

**dnspython** — v2.6+
- Reliable, no external dependencies beyond network
- Clean exception types
- Supports all standard + obscure record types (CAA, SRV, etc.)

**Avoid:**
- External resolver services (too slow, rate-limited)
- Sync/async complexity (all queries are fast enough for sync)
- Regex complexity beyond `r'pattern'` (keep it readable)

## Testing Domains

**Use these for live validation:**
- `github.com` — split NS (NS1 + Route53), Microsoft 365 email, 6 DKIM selectors
- `cloudflare.com` — self-branded NS, Cloudflare Email Security, strict DMARC p=reject, Mandrill selectors
- `google.com` — self-hosted NS (ns1-4.google.com), Google Workspace, Google Cloud DNS
- `badssl.com` endpoints:
  - `expired.badssl.com` — expired cert
  - `self-signed.badssl.com` — untrusted cert
  - `wrong.host.badssl.com` — hostname mismatch

These cover the full range of real-world scenarios.
