# Provider Fingerprinting

Identify registrar and nameserver providers by matching NS hostnames against a lookup table.

## Requirements

- Must identify registrar and nameserver providers by friendly name, not just raw hostnames
- Must handle split DNS (multiple NS providers for a single domain)
- Must detect self-hosted DNS (when NS records are subdomains of the queried domain)

## How to Build It

**Install:** No external dependencies — pure string pattern matching.

**Pattern:**
```python
import dns.resolver

# Ordered from most-specific to least-specific
NS_PATTERNS = [
    ("awsdns", "AWS Route 53"),
    ("cloudflare.com", "Cloudflare"),
    ("nsone.net", "NS1 / IBM NS1 Connect"),
    ("googledomains.com", "Google Domains"),
    ("azure-dns.com", "Azure DNS"),
    ("digitaloceandns.com", "DigitalOcean DNS"),
    # ... ~60 patterns total
]

def identify_provider(ns_hostname, query_domain):
    ns_lower = ns_hostname.rstrip(".").lower()
    # Check self-hosted first
    if ns_lower.endswith("." + query_domain.lower()):
        return f"Self-hosted (under {query_domain})"
    # Check pattern table
    for pattern, provider in NS_PATTERNS:
        if pattern in ns_lower:
            return provider
    return "Unknown"

# Query and fingerprint
answers = dns.resolver.resolve(domain, "NS")
ns_hosts = [str(r).rstrip(".") for r in answers]
providers = Counter(identify_provider(ns, domain) for ns in ns_hosts)
primary = providers.most_common(1)[0][0]
is_split = len(providers) > 1
```

## What to Avoid

- Don't rely on exact NS hostname matching — pattern matching is more robust
- Don't assume the provider with most NS records is the only one — split DNS is real (e.g., github.com: NS1 + Route53)
- Don't guess for "Unknown" providers — surface the raw hostname when unknown

## Constraints

- The lookup table is not exhaustive — new DNS providers require manual entry
- "Unknown" will appear for truly bespoke infrastructure or small/new providers
- Order matters: most-specific patterns must come before less-specific ones (e.g., check "azure-dns.com" before "dns.com")

## Origin

Spike 002: `ns-registrar-fingerprinting`
Source: `.planning/spikes/002-ns-registrar-fingerprinting/spike.py`

**Key findings:**
- Pattern table covers ~60 providers with high accuracy
- Self-hosted detection handles tech giants (Google, Facebook) running their own DNS
- Split DNS detection is clean and useful for understanding redundancy architecture
