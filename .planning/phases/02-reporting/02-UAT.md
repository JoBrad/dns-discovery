---
status: complete
phase: 02-reporting
source: [".planning/phases/02-reporting/02-01-SUMMARY.md"]
started: 2026-04-22T19:40:16Z
updated: 2026-04-22T19:45:21Z
---

## Current Test

[testing complete]

## Tests

### 1. Markdown report includes DNS Records table with all 9 record types
expected: report.md contains a DNS Records table with rows for A/AAAA/MX/NS/TXT/CNAME/SOA/CAA/SRV.
result: pass

### 2. Markdown report includes Detected Services sections when services are found
expected: report.md shows Detected Services with relevant subsections (Email Providers, Hosting & CDN, Verification & SaaS) when data exists.
result: pass

### 3. Markdown report includes split DNS detail for split-provider domains
expected: For github.com, report.md includes Split DNS summary and Split DNS Providers table with provider counts.
result: pass

### 4. Email Security section includes MX priority/host table and health context
expected: report.md includes MX Records priority/host table plus SPF/DMARC/DKIM and email health score.
result: pass

### 5. TLS Health table includes TLS version and days-to-expiry fields
expected: report.md TLS table includes TLS and Days columns, with readable values for reachable valid hosts.
result: pass

### 6. Generated output files are not tracked by git
expected: output/ artifacts are ignored in git status after runs.
result: pass

## Summary

total: 6
passed: 6
issues: 0
pending: 0
skipped: 0
blocked: 0

## Gaps

