# DNS Enumeration

Query all DNS record types for a domain and detect configured services from record patterns.

## Requirements

- Must work from a single domain input
- Must enumerate all standard record types (A, AAAA, MX, NS, TXT, CNAME, SOA, CAA, SRV)
- Must detect services from record patterns (email provider from MX, hosting from CNAME, SaaS tools from TXT)
- Must produce a readable output grouped by record type

## How to Build It

**Install:** Add `dnspython` to dependencies (already in pyproject.toml for this project).

**Pattern:**
```python
import dns.resolver
import dns.exception

RECORD_TYPES = ["A", "AAAA", "MX", "NS", "TXT", "CNAME", "SOA", "CAA", "SRV"]

def query_records(domain, record_type):
    try:
        answers = dns.resolver.resolve(domain, record_type)
        return [str(r) for r in answers]
    except dns.resolver.NXDOMAIN:
        return None  # domain doesn't exist — terminal
    except (dns.resolver.NoAnswer, dns.resolver.NoNameservers):
        return []  # record type not configured
    except dns.exception.Timeout:
        return []  # timeout

# Query all types
for rtype in RECORD_TYPES:
    result = query_records(domain, rtype)
    if result is None:
        print(f"ERROR: Domain does not exist")
        sys.exit(1)
    elif result:
        print(f"{rtype}: {len(result)} records")
        for r in result:
            print(f"  {r}")
```

**Service Detection Pattern:**
```python
SERVICE_PATTERNS = {
    "mx": {
        "google.com": "Google Workspace",
        "protection.outlook.com": "Microsoft 365",
        # ... more patterns
    },
    "txt": {
        "v=spf1": "SPF Record",
        "google-site-verification": "Google Search Console",
        # ... more patterns
    },
    "cname": {
        "cloudflare.com": "Cloudflare",
        "github.io": "GitHub Pages",
        # ... more patterns
    },
}

def detect_services(records):
    detected = {}
    for mx in records.get("MX", []):
        for pattern, service in SERVICE_PATTERNS["mx"].items():
            if pattern in mx.lower():
                detected.setdefault("email", set()).add(service)
    # ... similar for TXT and CNAME
    return detected
```

## What to Avoid

- Don't query external DNS services — use `dns.resolver.resolve()` directly (uses system resolvers)
- Don't treat `NoAnswer` and `NXDOMAIN` as the same — they mean different things
- Don't concatenate TXT records naively — long SPF records are split across multiple TXT strings; join them before regex matching

## Constraints

- `dnspython` exception types are typed, which is good (clean error handling)
- CAA and SRV records are less commonly configured but still worth querying
- Service detection is inherently limited to known patterns — new services require table maintenance

## Origin

Spike 001: `dns-record-enumeration`
Source: `.planning/spikes/001-dns-record-enumeration/spike.py`

**Key findings:**
- All 9 record types query cleanly with `dnspython`
- SPF `include:` chains reveal SaaS tool usage
- CAA records provide cert issuer policy (cross-reference with TLS spike)
