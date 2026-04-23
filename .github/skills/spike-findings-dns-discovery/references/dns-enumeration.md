# DNS Enumeration

Query all DNS record types for a domain and detect configured services from record patterns.

## Requirements

- Must work from a single domain input
- Must enumerate all standard record types (A, AAAA, MX, NS, TXT, CNAME, SOA, CAA, SRV)
- Must detect services from record patterns (email provider from MX, hosting from CNAME, SaaS tools from TXT)
- Must produce a readable output grouped by record type

## How to Build It

Use `github.com/miekg/dns` and keep one canonical record list and resolver path.

1. Define record types once and iterate in fixed order.
2. Query each type with recursive DNS (`RecursionDesired=true`) against the configured resolver.
3. Retry truncated UDP responses over TCP.
4. Convert resource records to clean, human-readable strings by record type.
5. Treat NXDOMAIN on `A` query as terminal domain-not-found; treat missing answers as not configured.
6. Run service detection after enumeration using MX/TXT/CNAME lookup maps.

```go
var recordTypes = []struct {
    name  string
    qtype uint16
}{
    {"A", dns.TypeA}, {"AAAA", dns.TypeAAAA}, {"MX", dns.TypeMX},
    {"NS", dns.TypeNS}, {"TXT", dns.TypeTXT}, {"CNAME", dns.TypeCNAME},
    {"SOA", dns.TypeSOA}, {"CAA", dns.TypeCAA}, {"SRV", dns.TypeSRV},
}

func QueryAllRecords(domain string) (DNSRecords, error) {
    records := make(DNSRecords)

    for _, rt := range recordTypes {
        answers, rcode, _ := queryRaw(domain, rt.qtype)
        if rcode == dns.RcodeNameError && rt.qtype == dns.TypeA {
            return nil, ErrNXDOMAIN
        }
        for _, ans := range answers {
            records[rt.name] = append(records[rt.name], rrToString(ans))
        }
    }
    return records, nil
}
```

```go
func DetectServices(records DNSRecords) DetectedServices {
    var svc DetectedServices
    seen := map[string]bool{}

    add := func(category *[]string, name string) {
        if !seen[name] {
            *category = append(*category, name)
            seen[name] = true
        }
    }

    for _, mx := range records["MX"] {
        lower := strings.ToLower(mx)
        for pattern, service := range mxServicePatterns {
            if strings.Contains(lower, pattern) {
                add(&svc.Email, service)
            }
        }
    }
    return svc
}
```

## What to Avoid

- Do not collapse NXDOMAIN and empty answer into one outcome.
- Do not skip TCP fallback for truncated DNS responses.
- Do not parse record strings generically when typed resource records are available.
- Do not run service matching without case normalization.

## Constraints

- Resolver is currently fixed to `8.8.8.8:53` in implementation.
- Service detection is pattern-based and only as complete as lookup maps.
- TXT records can hold multiple provider/verifier signals and can be long.

## Origin

Synthesized from spikes: 001

Source files available in:
- `sources/001-dns-record-enumeration/README.md`
- `sources/001-dns-record-enumeration/dns.go`
- `sources/001-dns-record-enumeration/types.go`

**Key findings:**
- All 9 record types query cleanly with `miekg/dns`
- SPF `include:` chains reveal SaaS tool usage
- CAA records provide cert issuer policy (cross-reference with TLS spike)
