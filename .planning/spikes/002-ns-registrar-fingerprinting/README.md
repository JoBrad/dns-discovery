---
spike: "002"
name: ns-registrar-fingerprinting
type: standard
validates: "Given a set of NS records, when matched against a provider lookup table, then a friendly provider name is returned with high accuracy for common providers"
verdict: VALIDATED
related: ["001"]
tags: [dns, ns, registrar, fingerprinting, provider-detection]
---

# Spike 002: NS/Registrar Fingerprinting

## What This Validates
Given NS hostnames for a domain, when pattern-matched against a lookup table of known provider substrings, then a friendly provider name is returned. Special case: when all NS records are subdomains of the queried domain, it's identified as "Self-hosted."

## Research
No library needed — pure string pattern matching against a curated lookup table. The table covers ~60 providers including cloud DNS (Route 53, Cloudflare, Azure, Google Cloud DNS), registrar-bundled DNS (GoDaddy, Namecheap, Name.com, IONOS, etc.), and specialist DNS providers (NS1, UltraDNS, Dyn, easyDNS).

## How to Run
```
uv run python spike.py <domain>
# Examples:
uv run python spike.py github.com       # split DNS
uv run python spike.py cloudflare.com   # self-brand
uv run python spike.py google.com       # self-hosted
```

## What to Expect
- Each NS hostname mapped to a provider name
- Primary provider identified (most NS records)
- Split DNS flagged with provider breakdown

## Investigation Trail

**Iteration 1 — Basic table matching:** Pattern substring match works for all tested providers. github.com correctly identified as NS1+Route53 split DNS.

**Edge case discovered — Self-hosted DNS:** `google.com` uses `ns[1-4].google.com`. These don't match any provider pattern since Google Corp uses their own infrastructure. Added self-hosted detection: if all NS records are subdomains of the queried domain, report "Self-hosted (under domain)".

**Split DNS behaviour:** github.com has 4 NS1 + 4 Route53 nameservers — a real active/active split for redundancy. The primary provider defaults to NS1 (alphabetically first via `most_common(1)` tie-break). In the real build, when counts are equal, surfacing all providers is better than picking one.

**Coverage gap:** "Unknown" will appear for truly bespoke infrastructure (private nameservers not under the queried domain). This is acceptable — the real build should surface the raw NS hostname when the provider is Unknown.

## Results
**Verdict: VALIDATED ✓**

- Pattern table covers the vast majority of real-world providers tested
- Self-hosted detection handles large tech companies running own DNS
- Split DNS detection is clean and useful
- Graceful Unknown fallback works correctly
- Table completeness is a maintenance concern, not a blocking one
