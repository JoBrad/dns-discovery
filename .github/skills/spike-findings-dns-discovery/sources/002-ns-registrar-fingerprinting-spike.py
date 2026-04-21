"""
Spike 002: NS/Registrar Fingerprinting
Maps nameserver hostnames to friendly provider names using pattern matching.

Usage: python spike.py <domain>
Example: python spike.py github.com
"""

import dns.resolver
import dns.exception
import sys
from collections import Counter

# Ordered from most-specific to least-specific
# Each entry: (substring_pattern, provider_name)
NS_PATTERNS = [
    # Cloud DNS
    ("awsdns", "AWS Route 53"),
    ("cloudflare.com", "Cloudflare"),
    ("nsone.net", "NS1 / IBM NS1 Connect"),
    ("googledomains.com", "Google Domains"),
    ("dns.google", "Google Cloud DNS"),
    ("azure-dns.com", "Azure DNS"),
    ("azure-dns.net", "Azure DNS"),
    ("azure-dns.org", "Azure DNS"),
    ("azure-dns.info", "Azure DNS"),
    ("digitaloceandns.com", "DigitalOcean DNS"),
    ("dnsimple.com", "DNSimple"),
    ("dnsimple.net", "DNSimple"),
    ("dynect.net", "Dyn / Oracle Cloud DNS"),
    ("ultradns.net", "UltraDNS"),
    ("ultradns.com", "UltraDNS"),
    ("ultradns.org", "UltraDNS"),
    ("ultradns.biz", "UltraDNS"),
    ("edns.biz", "Akamai Edge DNS"),
    ("akamai.net", "Akamai"),
    ("akam.net", "Akamai"),
    # Registrar-bundled DNS
    ("domaincontrol.com", "GoDaddy"),
    ("secureserver.net", "GoDaddy"),
    ("registrar-servers.com", "Namecheap"),
    ("namecheaphosting.com", "Namecheap"),
    ("name.com", "Name.com"),
    ("hover.com", "Hover"),
    ("namesilo.com", "NameSilo"),
    ("enom.com", "eNom"),
    ("networksolutions.com", "Network Solutions"),
    ("name-services.com", "Network Solutions"),
    ("web.com", "Web.com"),
    ("dotster.com", "Dotster"),
    ("register.com", "Register.com"),
    ("gkg.net", "GKG / TuCows"),
    ("tucows.com", "TuCows"),
    ("1and1.com", "IONOS / 1&1"),
    ("ionos.com", "IONOS"),
    ("ui-dns.com", "IONOS"),
    ("ui-dns.de", "IONOS"),
    ("ui-dns.biz", "IONOS"),
    ("ui-dns.org", "IONOS"),
    ("hichina.com", "Alibaba Cloud (Aliyun)"),
    ("alibabadns.com", "Alibaba Cloud (Aliyun)"),
    ("bluehost.com", "Bluehost"),
    ("hostgator.com", "HostGator"),
    ("siteground.net", "SiteGround"),
    ("inmotionhosting.com", "InMotion Hosting"),
    ("dreamhost.com", "DreamHost"),
    ("linode.com", "Linode / Akamai"),
    ("he.net", "Hurricane Electric"),
    ("afraid.org", "FreeDNS (afraid.org)"),
    ("cloudns.net", "ClouDNS"),
    ("pointhq.com", "PointHQ"),
    ("buddyns.com", "BuddyNS"),
    ("constellix.com", "Constellix"),
    ("easydns.com", "easyDNS"),
    ("easydns.net", "easyDNS"),
    ("rage4.com", "Rage4"),
    ("zendesk.com", "Zendesk"),
    ("shopify.com", "Shopify"),
    ("squarespace.com", "Squarespace"),
    ("parkingcrew.net", "ParkingCrew (parked domain)"),
    ("sedoparking.com", "Sedo Parking (parked domain)"),
    ("above.com", "Above.com (parked domain)"),
    ("huaweicloud.com", "Huawei Cloud DNS"),
]


def identify_provider(ns_hostname: str, query_domain: str = "") -> str:
    """Return a friendly provider name for a nameserver hostname."""
    ns_lower = ns_hostname.rstrip(".").lower()
    # Check if NS is a subdomain of the queried domain (self-hosted)
    if query_domain and ns_lower.endswith("." + query_domain.lower()):
        return f"Self-hosted (under {query_domain})"
    for pattern, provider in NS_PATTERNS:
        if pattern in ns_lower:
            return provider
    return "Unknown"


def fingerprint_nameservers(domain: str) -> dict:
    """Query NS records and identify the DNS provider."""
    try:
        answers = dns.resolver.resolve(domain, "NS")
        ns_hosts = [str(r).rstrip(".") for r in answers]
    except dns.resolver.NXDOMAIN:
        print(f"ERROR: '{domain}' does not exist")
        sys.exit(1)
    except Exception as e:
        print(f"ERROR: Could not resolve NS records: {e}")
        sys.exit(1)

    provider_counts = Counter(identify_provider(ns, domain) for ns in ns_hosts)
    # The dominant provider is the one with the most nameservers
    primary_provider = provider_counts.most_common(1)[0][0]
    is_split = len(provider_counts) > 1

    return {
        "domain": domain,
        "nameservers": ns_hosts,
        "providers": dict(provider_counts),
        "primary_provider": primary_provider,
        "is_split_dns": is_split,
    }


def print_result(result: dict):
    print(f"\n{'=' * 60}")
    print(f"  NS Fingerprinting: {result['domain']}")
    print(f"{'=' * 60}\n")

    print("  Nameservers:")
    for ns in sorted(result["nameservers"]):
        provider = identify_provider(ns, result["domain"])
        print(f"    {ns}")
        print(f"      → {provider}")

    print(f"\n  Primary DNS Provider: {result['primary_provider']}")

    if result["is_split_dns"]:
        print(f"\n  ⚠ Split DNS detected:")
        for provider, count in result["providers"].items():
            print(f"    • {provider} ({count} NS records)")

    print()


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python spike.py <domain>")
        print("Examples: python spike.py github.com")
        print("          python spike.py cloudflare.com")
        sys.exit(1)

    domain = sys.argv[1].strip().lower()
    result = fingerprint_nameservers(domain)
    print_result(result)
