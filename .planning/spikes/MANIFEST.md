# Spike Manifest

## Idea
Build a DNS zone discovery tool that, given a domain name, produces an executive summary of the zone's configuration. Output includes registrar and nameserver identification (with friendly provider names), a list of configured services/hosts/redirects, and health checks for each service — validating email DNS (MX/SPF/DKIM/DMARC), TLS certificate validity, and minimum TLS 1.2+ for web-facing targets.

## Requirements
- Must work from a single domain input (e.g., `example.com`)
- Must identify registrar and nameserver providers by friendly name, not just raw hostnames
- Must check email DNS health: MX records present, SPF (v=spf1), DKIM (_domainkey), DMARC (_dmarc)
- Must check TLS health for A/CNAME targets: valid cert, not expired, TLS 1.2+
- Must produce a readable output (not just raw record dumps)

## Spikes

| # | Name | Type | Validates | Verdict | Tags |
|---|------|------|-----------|---------|------|
| 001 | dns-record-enumeration | standard | Given a domain, when queried for all record types, then a complete zone picture is returned | VALIDATED ✓ | dns, enumeration, python |
| 002 | ns-registrar-fingerprinting | standard | Given NS records, when matched against a lookup table, then a friendly provider name is returned | VALIDATED ✓ | dns, ns, registrar, fingerprinting |
| 003 | tls-health-check | standard | Given an A/CNAME target, when connected via TLS, then cert validity and TLS version are determined | VALIDATED ✓ | tls, ssl, health-check, certificates |
| 004 | email-dns-health | standard | Given a domain with MX records, when TXT and _domainkey subdomains are inspected, then SPF/DKIM/DMARC health is identified | VALIDATED ✓ | dns, email, spf, dkim, dmarc |
