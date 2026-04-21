---
spike: "004"
name: email-dns-health
type: standard
validates: "Given a domain with MX records, when TXT and _domainkey subdomains are inspected, then SPF/DKIM/DMARC health is correctly identified with specific policy details"
verdict: VALIDATED
related: ["001"]
tags: [dns, email, spf, dkim, dmarc, dns-txt]
---

# Spike 004: Email DNS Health

## What This Validates
Given a domain, when we query MX records, SPF (TXT `v=spf1`), DMARC (TXT at `_dmarc.<domain>`), and DKIM (probing common selectors at `<selector>._domainkey.<domain>`), then we can produce a health score with specific findings about email infrastructure.

## Research
No library needed — pure DNS queries with regex pattern matching. SPF/DMARC are straightforward TXT queries. DKIM discovery requires probing 20-30 common selectors since there's no way to enumerate all selectors via DNS. Tested selectors cover Microsoft 365 (`selector1/selector2`), Google (`google`), Mandrill (`mandrill`), Zoho, Fastmail, Amazon SES, etc.

## How to Run
```
uv run python spike.py <domain>
# Examples:
uv run python spike.py github.com       # email health score 4/4
uv run python spike.py cloudflare.com   # email health score 4/4
```

## What to Expect
- MX record count and server names
- SPF policy: softfail (~all) vs hardfail (-all) with warnings
- DMARC policy: p=none (monitoring), p=quarantine, p=reject with report address status
- DKIM selectors found via probe with key type (rsa/ed25519) and revocation status
- Summary health score: MX ✓/✗, SPF ✓/✗, DMARC ✓/✗, DKIM ✓/✗

## Investigation Trail

**Iteration 1 — Basic queries:** MX and SPF records query cleanly. DMARC at `_dmarc.<domain>` works for both github.com and cloudflare.com.

**Edge case — DKIM selector discovery:** Unlike SPF/DMARC, there's no single "where is DKIM" query. Must probe common selector names. github.com has 6 selectors (google, selector1, k1, k2, s1, s2). cloudflare.com has 3 (k1, s1, mandrill).

**Key finding — Selector patterns reveal provider:** The `mandrill._domainkey.cloudflare.com` selector indicates Mandrill/Mailchimp is used for transactional email. Spike 001 also detected Mandrill's IP ranges in SPF. This cross-reference validates tool chain.

**SPF policy analysis:** github.com uses softfail (~all) — recommended for gradual enforcement. cloudflare.com uses hardfail (-all) — strictest, rejects unauthenticated mail. Both have extensive SPF `include:` chains revealing third-party email services.

**DMARC implementation:** Both domains use aggregate reporting (rua=) for feedback. github.com: p=quarantine. cloudflare.com: p=reject — very strict. Neither has forensic reporting (ruf=) enabled.

**Health scoring:** Both domains score 4/4 because they implement all four pillars. In practice, many domains skip DMARC (p=none) or have incomplete DKIM.

## Results
**Verdict: VALIDATED ✓**

- MX/SPF/DMARC queries via dnspython work perfectly
- DKIM probe list of 27 selectors covers real-world deployments
- Health scoring is intuitive (4/4 = all pillars implemented)
- Selector discovery reveals provider details (e.g. "mandrill" → Mailchimp)
- SPF policy parsing distinguishes softfail vs hardfail
- No false positives on selector revocation detection (empty `p=`)
- Cross-reference opportunity: SPF `include:` values + DKIM selectors + Spike 001 service detection confirm consistent email provider picture
