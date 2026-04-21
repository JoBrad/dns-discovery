---
spike: "001"
name: dns-record-enumeration
type: standard
validates: "Given a domain, when queried for all common record types, then a complete zone picture is returned with service detection from record patterns"
verdict: VALIDATED
related: ["002", "003", "004"]
tags: [dns, enumeration, python, dnspython, services]
---

# Spike 001: DNS Record Enumeration

## What This Validates
Given a domain, when we query A, AAAA, MX, NS, TXT, CNAME, SOA, CAA, SRV using `dnspython`, then we receive a structured zone picture and can detect configured services (email provider, hosting, SaaS tools) from record patterns.

## Research
`dnspython` (`dns.resolver`) is the de-facto Python DNS library. It supports all standard record types, raises clean typed exceptions (`NXDOMAIN`, `NoAnswer`, `NoNameservers`, `Timeout`), and requires no external resolvers. No competing approaches considered — there is no credible Python alternative for DNS resolution.

## How to Run
```
uv run python spike.py <domain>
# Example: uv run python spike.py github.com
```

## What to Expect
- All record types queried; unsupported types show "not configured"
- Records printed grouped by type with counts
- Detected services listed by category (email, hosting, verification services)
- Summary with totals and boolean flags for IPv6/MX/SPF

## Investigation Trail

**Iteration 1 — Basic enumeration:** Confirmed `dnspython` resolves all 9 record types cleanly. NXDOMAIN and NoAnswer exceptions are distinct and handleable.

**Iteration 2 — Service detection:** Added pattern-matching against MX values (email providers), TXT record prefixes (SPF, DMARC, SaaS verifications), and CNAME suffixes (CDN/hosting). Tested on `github.com` and `cloudflare.com`.

**Surprising finding — SPF as service signal:** github.com's SPF `include:` directives reveal Zendesk, Salesforce, Mailchimp, SendGrid in use — all services configured to send email on their behalf. This is richer signal than MX alone for understanding a domain's SaaS footprint.

**Surprising finding — CAA records:** github.com has 7 CAA records listing permitted certificate authorities (DigiCert, Let's Encrypt, Sectigo). This is valuable for the TLS health check spike — we can cross-reference the cert issuer against the CAA policy.

**Edge case tested:** SPF records are sometimes split across multiple TXT strings by DNS (the `"v=spf1 ..." "27.114 ..."` split). `dnspython` returns them as separate strings in the same record; the real build will need to join them.

## Results
**Verdict: VALIDATED ✓**

- `dnspython` reliably queries all 9 record types
- Service detection from MX + TXT patterns is high-signal and low-false-positive
- CNAME-based hosting detection requires subdomains (not apex), which is expected
- SPF `include:` chains are a rich secondary signal for SaaS tool detection
- CAA records provide cert issuer policy — useful cross-reference for TLS spike
- Split TXT records (long SPF) need joining in the real build
