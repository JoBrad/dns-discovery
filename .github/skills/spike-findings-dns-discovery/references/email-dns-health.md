# Email DNS Health

Validate email infrastructure by checking MX, SPF, DMARC, and DKIM records.

## Requirements

- Must check email DNS health: MX records present, SPF (v=spf1), DKIM (_domainkey), DMARC (_dmarc)
- Must produce a health score for email configuration
- Must identify email providers from MX records and DKIM selectors

## How to Build It

**Install:** Use `dnspython` for DNS queries, stdlib `re` for pattern matching.

**MX Records:**
```python
import dns.resolver

def query_mx(domain):
    try:
        answers = dns.resolver.resolve(domain, "MX")
        return sorted([(r.preference, str(r.exchange).rstrip(".")) for r in answers])
    except Exception:
        return []
```

**SPF Detection:**
```python
import re

def check_spf(domain):
    txt_records = query_txt(domain)
    spf_records = [t for t in txt_records if t.startswith("v=spf1")]
    
    if not spf_records:
        return {"present": False, "issues": ["No SPF record found"]}
    
    if len(spf_records) > 1:
        return {"present": True, "issues": [f"Multiple SPF records ({len(spf_records)}) — only one allowed"]}
    
    spf = spf_records[0]
    issues = []
    
    # Check for 'all' mechanism
    if not re.search(r'[~\-\+\?]all\b', spf):
        issues.append("Missing 'all' mechanism")
    
    # Classify policy
    if "~all" in spf:
        qualifier = "softfail (~all) — recommended"
    elif "-all" in spf:
        qualifier = "hardfail (-all) — strictest"
    else:
        qualifier = "unknown"
    
    return {"present": True, "valid": len(issues) == 0, "qualifier": qualifier, "issues": issues}
```

**DMARC Detection:**
```python
def check_dmarc(domain):
    dmarc_domain = f"_dmarc.{domain}"
    txt_records = query_txt(dmarc_domain)
    dmarc_records = [t for t in txt_records if t.startswith("v=DMARC1")]
    
    if not dmarc_records:
        return {"present": False}
    
    dmarc = dmarc_records[0]
    policy_match = re.search(r'\bp=(\w+)', dmarc)
    policy = policy_match.group(1) if policy_match else "none"
    has_rua = "rua=" in dmarc
    has_ruf = "ruf=" in dmarc
    
    return {
        "present": True,
        "policy": policy,  # none, quarantine, reject
        "has_aggregate_reports": has_rua,
        "has_forensic_reports": has_ruf,
    }
```

**DKIM Discovery (Probe-based):**
```python
DKIM_SELECTORS = [
    "google", "selector1", "selector2", "default", "mail",
    "k1", "k2", "s1", "s2", "smtp", "dkim",
    "mandrill", "mailjet", "sendgrid", "sparkpost",
    # ... 27 total
]

def check_dkim(domain):
    found = []
    for selector in DKIM_SELECTORS:
        dkim_domain = f"{selector}._domainkey.{domain}"
        txt_records = query_txt(dkim_domain)
        dkim_records = [t for t in txt_records if "v=DKIM1" in t or "p=" in t]
        if dkim_records:
            record = dkim_records[0]
            key_type_match = re.search(r'k=(\w+)', record)
            key_type = key_type_match.group(1) if key_type_match else "rsa"
            revoked = bool(re.search(r'p=\s*;|p=\s*$', record))
            found.append({"selector": selector, "key_type": key_type, "revoked": revoked})
    
    return {"present": len(found) > 0, "selectors_found": found}
```

## What to Avoid

- Don't assume DKIM is absent if probe finds no selectors — selector names are not standardized; the probe list covers ~90% but not all
- Don't treat multiple SPF records as data, only warn — only one is valid per DNS spec
- Don't forget to join split TXT records (long SPF strings) before regex matching

## Constraints

- DKIM discovery is probe-based, not exhaustive — unknown selectors will be missed
- SPF chains can be arbitrarily deep — follow 1-2 levels for performance, not all chains
- DMARC policy values: `none` (monitoring only), `quarantine` (tag suspicious), `reject` (fail)

## Origin

Spike 004: `email-dns-health`
Source: `.planning/spikes/004-email-dns-health/spike.py`

**Key findings:**
- DKIM probe list covers 27 selectors — covers Microsoft 365, Google, Mandrill, Amazon SES, Zoho, Fastmail
- Selector names reveal provider details (e.g., `mandrill._domainkey` → Mailchimp)
- SPF `include:` chains + DKIM selectors + Spike 001 MX detection paint a complete email provider picture
- Health score: MX ✓/✗, SPF ✓/✗, DMARC ✓/✗, DKIM ✓/✗ (4/4 = all pillars)
