# DNS Discovery Report: cloudflare.com

Generated on: Wed, 22 Apr 2026 14:12:05 CDT

## Executive Summary

- **Primary Provider:** Self-hosted (under cloudflare.com)
- **Email Health:** ✅ Healthy (4/4)
- **TLS Health:** 1/1 hosts healthy

## Infrastructure & DNS

### Nameservers (Self-hosted (under cloudflare.com))
- ns4.cloudflare.com
- ns5.cloudflare.com
- ns6.cloudflare.com
- ns3.cloudflare.com
- ns7.cloudflare.com

## Email Security

- **SPF:** `v=spf1 ip4:199.15.212.0/22 ip4:173.245.48.0/20 include:_spf.google.com include:spf1.mcsv.net include:spf.mandrillapp.com include:mail.zendesk.com include:stspg-customer.com include:_spf.salesforce.com -all` (hardfail)
- **DMARC:** `v=DMARC1; p=reject; pct=100; rua=mailto:rua@cloudflare.com,mailto:cloudflare@dmarc.area1reports.com; ruf=mailto:cloudflare@dmarc.area1reports.com` (p=reject)
- **DKIM:** Found selectors: k1, s1, mandrill

## TLS Health

| Hostname | Reachable | Valid | Expiry | Issuer |
|----------|-----------|-------|--------|--------|
| cloudflare.com | ✅ | ✅ | 2026-06-10 | Google Trust Services |

