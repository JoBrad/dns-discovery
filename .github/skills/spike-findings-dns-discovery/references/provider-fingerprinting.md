# Provider Fingerprinting

Identify registrar and nameserver providers by matching NS hostnames against a lookup table.

## Requirements

- Must identify registrar and nameserver providers by friendly name, not just raw hostnames
- Must handle split DNS (multiple NS providers for a single domain)
- Must detect self-hosted DNS (when NS records are subdomains of the queried domain)

## How to Build It

Use a deterministic ordered pattern table and a self-hosted shortcut.

1. Normalize NS hostnames (lowercase, trim trailing dot).
2. Detect self-hosted first: NS host ends with `.<domain>`.
3. Match remaining hosts against ordered provider patterns (~60 entries).
4. Accumulate provider counts and choose primary by max count.
5. Mark split DNS when provider count map has more than one entry.

```go
func identifyNS(nsHost, domain string) string {
    lower := strings.ToLower(strings.TrimSuffix(nsHost, "."))
    if strings.HasSuffix(lower, "."+strings.ToLower(domain)) {
        return "Self-hosted (under " + domain + ")"
    }
    for _, p := range nsPatterns {
        if strings.Contains(lower, p.pattern) {
            return p.provider
        }
    }
    return "Unknown (" + lower + ")"
}

func IdentifyProviders(domain string, nsHosts []string) ProviderResult {
    result := ProviderResult{Counts: make(map[string]int), AllHosts: nsHosts}
    for _, ns := range nsHosts {
        provider := identifyNS(ns, domain)
        result.Counts[provider]++
    }
    for provider, count := range result.Counts {
        if count > result.Counts[result.Primary] || result.Primary == "" {
            result.Primary = provider
        }
    }
    result.IsSplit = len(result.Counts) > 1
    return result
}
```

## What to Avoid

- Do not use exact hostname lookup only; substring patterns are required.
- Do not run generic pattern checks before self-hosted check.
- Do not hide split-provider counts when ties occur.
- Do not map unknown providers to nearest known provider.

## Constraints

- Provider table maintenance is ongoing work as DNS ecosystems change.
- Pattern order matters; keep specific suffixes ahead of generic tokens.
- Unknown results are valid and should remain explicit.

## Origin

Synthesized from spikes: 002

Source files available in:
- `sources/002-ns-registrar-fingerprinting/README.md`
- `sources/002-ns-registrar-fingerprinting/providers.go`
- `sources/002-ns-registrar-fingerprinting/types.go`

**Key findings:**
- Pattern table covers ~60 providers with high accuracy
- Self-hosted detection handles tech giants (Google, Facebook) running their own DNS
- Split DNS detection is clean and useful for understanding redundancy architecture
