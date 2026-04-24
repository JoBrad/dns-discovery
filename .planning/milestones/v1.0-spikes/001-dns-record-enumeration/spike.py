"""
Spike 001: DNS Record Enumeration
Queries all common DNS record types for a domain and detects configured services.

Usage: python spike.py <domain>
Example: python spike.py github.com
"""

import dns.resolver
import dns.exception
import sys

RECORD_TYPES = ["A", "AAAA", "MX", "NS", "TXT", "CNAME", "SOA", "CAA", "SRV"]

# Service detection patterns
SERVICE_PATTERNS = {
    "mx": {
        "google.com": "Google Workspace",
        "googlemail.com": "Google Workspace",
        "outlook.com": "Microsoft 365",
        "protection.outlook.com": "Microsoft 365",
        "pphosted.com": "Proofpoint",
        "mimecast.com": "Mimecast",
        "mailgun.org": "Mailgun",
        "sendgrid.net": "SendGrid",
        "amazonses.com": "Amazon SES",
        "zoho.com": "Zoho Mail",
        "fastmail.com": "Fastmail",
        "inbound.cf-emailsecurity.net": "Cloudflare Email Security",
    },
    "txt": {
        "v=spf1": "SPF Record",
        "v=DMARC1": "DMARC Policy",
        "google-site-verification": "Google Search Console",
        "MS=ms": "Microsoft 365",
        "facebook-domain-verification": "Facebook",
        "docusign": "DocuSign",
        "atlassian-domain-verification": "Atlassian",
        "stripe-verification": "Stripe",
        "adobe-idp-site-verification": "Adobe",
        "apple-domain-verification": "Apple",
        "ZOOM_verify": "Zoom",
    },
    "cname": {
        "cloudfront.net": "AWS CloudFront",
        "amazonaws.com": "AWS",
        "azurewebsites.net": "Azure App Service",
        "fastly.net": "Fastly CDN",
        "github.io": "GitHub Pages",
        "netlify.app": "Netlify",
        "vercel.app": "Vercel",
        "heroku.com": "Heroku",
        "pages.dev": "Cloudflare Pages",
        "workers.dev": "Cloudflare Workers",
        "shopify.com": "Shopify",
        "zendesk.com": "Zendesk",
        "hubspot.com": "HubSpot",
        "wixdns.net": "Wix",
    },
}


def query_records(domain, record_type):
    try:
        answers = dns.resolver.resolve(domain, record_type)
        return [str(r) for r in answers]
    except dns.resolver.NXDOMAIN:
        return None
    except (dns.resolver.NoAnswer, dns.resolver.NoNameservers):
        return []
    except dns.exception.Timeout:
        return []
    except Exception:
        return []


def detect_services(records):
    detected = {}
    for mx in records.get("MX", []):
        for pattern, service in SERVICE_PATTERNS["mx"].items():
            if pattern in mx.lower():
                detected.setdefault("email", set()).add(service)
    for txt in records.get("TXT", []):
        txt_clean = txt.strip('"')
        for pattern, service in SERVICE_PATTERNS["txt"].items():
            if pattern in txt_clean:
                detected.setdefault("verification_services", set()).add(service)
    for cname in records.get("CNAME", []):
        for pattern, service in SERVICE_PATTERNS["cname"].items():
            if pattern in cname.lower():
                detected.setdefault("hosting", set()).add(service)
    return {k: sorted(v) for k, v in detected.items()}


def enumerate_dns(domain):
    print(f"\n{'=' * 60}")
    print(f"  DNS Zone Enumeration: {domain}")
    print(f"{'=' * 60}\n")

    records = {}
    for record_type in RECORD_TYPES:
        result = query_records(domain, record_type)
        if result is None:
            if record_type == "A":
                print(f"ERROR: '{domain}' does not exist (NXDOMAIN)")
                sys.exit(1)
        elif result:
            records[record_type] = result
            print(f"{record_type:8} ({len(result)} records):")
            for r in result:
                print(f"         {r}")
        else:
            print(f"{record_type:8} — not configured")

    services = detect_services(records)
    print(f"\n{'=' * 60}")
    print("  Detected Services")
    print(f"{'=' * 60}")
    if services:
        for category, service_list in services.items():
            print(f"\n  {category.replace('_', ' ').title()}:")
            for svc in service_list:
                print(f"    • {svc}")
    else:
        print("  None detected")

    print(f"\n{'=' * 60}")
    print("  Summary")
    print(f"{'=' * 60}")
    print(f"  Record types found: {len(records)}/{len(RECORD_TYPES)}")
    print(f"  Total records:      {sum(len(v) for v in records.values())}")
    print(f"  Has IPv6 (AAAA):    {'Yes' if 'AAAA' in records else 'No'}")
    print(f"  Has email (MX):     {'Yes' if 'MX' in records else 'No'}")
    has_spf = any("v=spf1" in t for t in records.get("TXT", []))
    print(f"  Has SPF:            {'Yes' if has_spf else 'No'}")
    print()

    return records


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python spike.py <domain>")
        sys.exit(1)
    enumerate_dns(sys.argv[1])
