# DNS Discovery Report: github.com

Generated on: Wed, 22 Apr 2026 14:12:03 CDT

## Executive Summary

- **Primary Provider:** AWS Route 53
- **Email Health:** ✅ Healthy (4/4)
- **TLS Health:** 1/1 hosts healthy

## Infrastructure & DNS

### Nameservers (AWS Route 53)
- ns-1707.awsdns-21.co.uk
- dns3.p08.nsone.net
- dns4.p08.nsone.net
- ns-520.awsdns-01.net
- dns2.p08.nsone.net
- ns-421.awsdns-52.com
- dns1.p08.nsone.net
- ns-1283.awsdns-32.org

## Email Security

- **SPF:** `v=spf1 ip4:192.30.252.0/22 include:spf.protection.outlook.com include:_netblocks.google.com include:_netblocks2.google.com include:mail.zendesk.com include:_spf.salesforce.com include:servers.mcsv.net include:mktomail.com include:sendgrid.net ip4:62.253.227.114 ip4:166.78.69.169 ip4:166.78.69.170 ip4:166.78.71.131 ~all` (softfail)
- **DMARC:** `v=DMARC1; p=quarantine; sp=reject; pct=100; rua=mailto:dmarc@github.com; ruf=mailto:dmarc@github.com; fo=1` (p=quarantine)
- **DKIM:** Found selectors: google, selector1, k1, k2, s1, s2

## TLS Health

| Hostname | Reachable | Valid | Expiry | Issuer |
|----------|-----------|-------|--------|--------|
| github.com | ✅ | ✅ | 2026-06-03 | Sectigo Limited |

