"""
Spike 004: Email DNS Health
Validates MX, SPF, DKIM, and DMARC records for a domain.

Usage: python spike.py <domain>
Example: python spike.py github.com
"""

import dns.resolver
import dns.exception
import sys
import re


# Common DKIM selectors to probe
DKIM_SELECTORS = [
    "google",
    "selector1",  # Microsoft 365
    "selector2",  # Microsoft 365 alternate
    "default",
    "mail",
    "k1",
    "k2",
    "s1",
    "s2",
    "smtp",
    "dkim",
    "email",
    "key1",
    "key2",
    "sig1",
    "sig2",
    "mimecast",
    "mailjet",
    "sendgrid",
    "mandrill",
    "sparkpost",
    "amazonses",
    "protonmail",
    "zoho",
    "fm1",  # Fastmail
    "fm2",
    "fm3",
]


def query_txt(domain: str) -> list[str]:
    """Query TXT records for a domain, returning cleaned strings."""
    try:
        answers = dns.resolver.resolve(domain, "TXT")
        # Join multi-string TXT records (e.g. long SPF split across strings)
        result = []
        for rdata in answers:
            joined = "".join(
                s.decode() if isinstance(s, bytes) else s for s in rdata.strings
            )
            result.append(joined)
        return result
    except (dns.resolver.NXDOMAIN, dns.resolver.NoAnswer, dns.resolver.NoNameservers):
        return []
    except dns.exception.Timeout:
        return []


def query_mx(domain: str) -> list[tuple[int, str]]:
    """Query MX records, returning (priority, hostname) tuples sorted by priority."""
    try:
        answers = dns.resolver.resolve(domain, "MX")
        return sorted([(r.preference, str(r.exchange).rstrip(".")) for r in answers])
    except Exception:
        return []


def check_spf(domain: str) -> dict:
    """Check SPF record exists and is syntactically valid."""
    txt_records = query_txt(domain)
    spf_records = [t for t in txt_records if t.startswith("v=spf1")]

    if not spf_records:
        return {
            "present": False,
            "valid": False,
            "record": None,
            "issues": ["No SPF record found"],
        }

    if len(spf_records) > 1:
        return {
            "present": True,
            "valid": False,
            "record": spf_records[0],
            "issues": [
                f"Multiple SPF records found ({len(spf_records)}) — only one is allowed"
            ],
        }

    spf = spf_records[0]
    issues = []

    # Check it ends with an all mechanism
    if not re.search(r"[~\-\+\?]all\b", spf):
        issues.append("Missing 'all' mechanism at end of SPF record")

    # Warn on softfail vs hardfail
    qualifier = "unknown"
    if "~all" in spf:
        qualifier = "softfail (~all) — recommended"
    elif "-all" in spf:
        qualifier = "hardfail (-all) — strictest"
    elif "+all" in spf:
        issues.append("+all allows anyone to send — very insecure")
        qualifier = "passall (+all) — INSECURE"
    elif "?all" in spf:
        qualifier = "neutral (?all) — no enforcement"

    return {
        "present": True,
        "valid": len(issues) == 0,
        "record": spf,
        "qualifier": qualifier,
        "issues": issues,
    }


def check_dmarc(domain: str) -> dict:
    """Check DMARC record at _dmarc.<domain>."""
    dmarc_domain = f"_dmarc.{domain}"
    txt_records = query_txt(dmarc_domain)
    dmarc_records = [t for t in txt_records if t.startswith("v=DMARC1")]

    if not dmarc_records:
        return {
            "present": False,
            "valid": False,
            "record": None,
            "issues": ["No DMARC record at _dmarc." + domain],
        }

    dmarc = dmarc_records[0]
    issues = []

    # Extract policy
    policy_match = re.search(r"\bp=(\w+)", dmarc)
    policy = policy_match.group(1) if policy_match else "none"

    if policy == "none":
        issues.append("p=none — DMARC is monitoring only, no enforcement")
    elif policy in ("quarantine", "reject"):
        pass  # Good

    # Check reporting addresses
    has_rua = "rua=" in dmarc
    has_ruf = "ruf=" in dmarc
    if not has_rua:
        issues.append("No aggregate report address (rua=) — losing visibility")

    return {
        "present": True,
        "valid": True,
        "record": dmarc,
        "policy": policy,
        "has_aggregate_reports": has_rua,
        "has_forensic_reports": has_ruf,
        "issues": issues,
    }


def check_dkim(domain: str, selectors: list[str] = DKIM_SELECTORS) -> dict:
    """
    Probe common DKIM selectors at <selector>._domainkey.<domain>.
    Returns all found selectors with their key info.
    """
    found = []
    probed = 0

    for selector in selectors:
        dkim_domain = f"{selector}._domainkey.{domain}"
        probed += 1
        txt_records = query_txt(dkim_domain)
        dkim_records = [
            t for t in txt_records if "v=DKIM1" in t or "k=rsa" in t or "p=" in t
        ]
        if dkim_records:
            record = dkim_records[0]
            # Extract key type
            key_type_match = re.search(r"k=(\w+)", record)
            key_type = key_type_match.group(1) if key_type_match else "rsa"
            # Check if key is revoked (empty p=)
            revoked = bool(re.search(r"p=\s*;|p=\s*$", record))
            found.append(
                {
                    "selector": selector,
                    "dkim_domain": dkim_domain,
                    "key_type": key_type,
                    "revoked": revoked,
                }
            )

    return {
        "present": len(found) > 0,
        "selectors_found": found,
        "selectors_probed": probed,
    }


def check_email_health(domain: str) -> dict:
    print(f"\n{'=' * 60}")
    print(f"  Email DNS Health: {domain}")
    print(f"{'=' * 60}\n")

    # 1. MX Records
    mx_records = query_mx(domain)
    print(f"  MX Records:")
    if mx_records:
        for priority, host in mx_records:
            print(f"    [{priority:3}] {host}")
        print(f"    → {len(mx_records)} mail server(s) configured")
    else:
        print(f"    ✗ No MX records — domain cannot receive email")

    # 2. SPF
    print(f"\n  SPF Record:")
    spf = check_spf(domain)
    if spf["present"]:
        short = spf["record"][:80] + "..." if len(spf["record"]) > 80 else spf["record"]
        print(f"    ✓ Present: {short}")
        print(f"    Policy: {spf.get('qualifier', 'unknown')}")
        for issue in spf["issues"]:
            print(f"    ⚠ {issue}")
    else:
        for issue in spf["issues"]:
            print(f"    ✗ {issue}")

    # 3. DMARC
    print(f"\n  DMARC Record:")
    dmarc = check_dmarc(domain)
    if dmarc["present"]:
        print(f"    ✓ Present: p={dmarc.get('policy', '?')}")
        print(
            f"    Aggregate reports: {'Yes' if dmarc.get('has_aggregate_reports') else 'No'}"
        )
        for issue in dmarc["issues"]:
            print(f"    ⚠ {issue}")
    else:
        for issue in dmarc["issues"]:
            print(f"    ✗ {issue}")

    # 4. DKIM
    print(f"\n  DKIM (probing {len(DKIM_SELECTORS)} selectors):")
    dkim = check_dkim(domain)
    if dkim["present"]:
        for sel in dkim["selectors_found"]:
            status = "✗ REVOKED" if sel["revoked"] else "✓"
            print(f"    {status} selector={sel['selector']}  ({sel['key_type']})")
            print(f"       {sel['dkim_domain']}")
    else:
        print(f"    ✗ No DKIM selectors found (probed {dkim['selectors_probed']})")
        print(f"    ℹ DKIM may use a non-standard selector not in the probe list")

    # Summary
    has_mx = bool(mx_records)
    has_spf = spf["present"]
    has_dmarc = dmarc["present"]
    has_dkim = dkim["present"]

    score = sum([has_mx, has_spf, has_dmarc, has_dkim])
    print(f"\n  {'=' * 40}")
    print(f"  Email Health Score: {score}/4")
    print(f"    MX:    {'✓' if has_mx else '✗'}")
    print(f"    SPF:   {'✓' if has_spf else '✗'}")
    print(f"    DMARC: {'✓' if has_dmarc else '✗'}")
    print(f"    DKIM:  {'✓' if has_dkim else '✗ (not found in probe list)'}")
    print()

    return {"mx": mx_records, "spf": spf, "dmarc": dmarc, "dkim": dkim}


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python spike.py <domain>")
        print("Example: python spike.py github.com")
        sys.exit(1)
    check_email_health(sys.argv[1])
