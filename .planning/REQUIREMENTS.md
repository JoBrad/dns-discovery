# Project Requirements

**Last Updated:** April 22, 2026

## Phase Requirements

### Phase 1: CLI Tool Foundation

**Phase Goal:** Build the core discovery engine with DNS enumeration, provider identification, and health checks integrated into a working CLI.

**Requirement IDs for Phase 1:**

| ID | Category | Requirement | UAT |
|---|---|---|---|
| CLI-01 | Interaction | CLI accepts domain argument (`dns-discovery example.com`) | Command runs without error, accepts 1 argument |
| ENUM-01 | DNS Enumeration | Query all 9 DNS record types (A, AAAA, MX, NS, TXT, CNAME, SOA, CAA, SRV) | All types queried, None missing |
| ENUM-02 | Service Detection | Detect services from MX/TXT/CNAME records (email providers, SaaS tools, hosting) | Service patterns recognized in output |
| PROV-01 | Nameserver ID | Map NS records to friendly provider names (~60 providers supported) | github.com → NS1 + Route53, cloudflare.com → Cloudflare |
| PROV-02 | Self-Hosted Detection | Detect self-hosted NS (when NS is subdomain of queried domain) | google.com → "Self-hosted (under google.com)" |
| PROV-03 | Split DNS Detection | Identify when multiple providers manage NS for one domain | github.com → "Split: NS1 (4), Route53 (4)" |
| EMAIL-01 | MX Records | Query and list MX records with priority | github.com shows 5 MX records |
| EMAIL-02 | SPF Validation | Query v=spf1 TXT record, validate syntax, classify policy (softfail/hardfail) | SPF record found, policy identified |
| EMAIL-03 | DMARC Validation | Query _dmarc.<domain>, extract policy (none/quarantine/reject) | DMARC present, policy extracted |
| EMAIL-04 | DKIM Discovery | Probe 27 common DKIM selectors, identify found keys | At least 1 DKIM selector found on test domains |
| EMAIL-05 | Email Health Score | Produce 4-pillar health score (MX/SPF/DMARC/DKIM) | github.com scores 4/4, cloudflare.com scores 4/4 |
| TLS-01 | TLS Check | Connect to A/CNAME targets on port 443, verify cert | github.com → TLSv1.3, valid |
| TLS-02 | Cert Validity | Distinguish expired, self-signed, hostname mismatch | badssl.com endpoints caught all 3 modes |
| TLS-03 | Expiry Warning | Calculate days until expiry, warn if <14 days | github.com shows "valid until 2026-06-03 (43d)" |
| TLS-04 | Graceful Non-HTTPS | Handle A records pointing to non-HTTPS services | Timeout/refused handled without crash |
| OUT-01 | Output Structure | Organize discovery results in memory (not yet written to disk) | All checks complete, results available |
| OUT-02 | Readable Summary | Present zone overview in human-readable format | Print to stdout with sections and formatting |

---

### Phase 2: Reporting & Output

**Phase Goal:** Enrich the Markdown report file to match the full information presented on stdout, and gitignore generated output.

**Requirement IDs for Phase 2:**

| ID | Category | Requirement | UAT |
|---|---|---|---|
| RPT-01 | Report Completeness | Markdown report includes DNS records table (all 9 types with values) | `report.md` has DNS Records section with A/MX/NS/TXT etc. |
| RPT-02 | Report Completeness | Markdown report includes detected services (email providers, hosting/CDN, verification SaaS) | `report.md` has Detected Services section matching stdout |
| RPT-03 | Report Completeness | Markdown report includes split DNS detail, MX records with priorities, TLS version and days-to-expiry | All three fields present in generated report for github.com |

---

### Phase 3: Integration & Polish

**Phase Goal:** Add file-based configuration, batch execution across multiple domains, and clearer error handling for single and multi-domain runs.

**Requirement IDs for Phase 3:**

| ID | Category | Requirement | UAT |
|---|---|---|---|
| CFG-01 | Configuration | Tool loads settings from `.dns-discovery.json` or an explicit `--config` path | Running with config file applies configured output dir and domain list |
| CFG-02 | Configuration | CLI flags override config values and defaults remain stable when config is absent | `--output-dir` beats config `output_dir`; command still works without config |
| BAT-01 | Batch Mode | Tool processes multiple domains from config `domains:` list or `--input-file` | One invocation generates reports for all listed domains |
| BAT-02 | Batch Resilience | Failure on one domain does not abort the whole batch | Mixed valid/invalid domain file still processes valid domains and reports failures |
| ERR-01 | Error Handling | Invalid config/input errors are specific and actionable; batch runs end with aggregate success/failure summary | Bad config shows file/field context; batch run prints counts and exits non-zero if any failures |

---

## Validation Approach

**End-of-Phase Verification:** Run tool on github.com and cloudflare.com, verify all 4 categories (DNS, Providers, TLS, Email) produce correct output matching spike validation results.

**Test Domains:**
- `github.com` — Complex: split DNS, multiple email tools, 6 DKIM selectors
- `cloudflare.com` — Complex: self-branded NS, multiple SaaS services, Mandrill/Mailchimp

**Key Assertions (from spike work):**
- github.com: NS1 + Route53 split, email score 4/4, TLS ✓
- cloudflare.com: Cloudflare self-branded, email score 4/4, TLS ✓

---

*Derived from spike validation results and README vision. See spike findings in `.github/skills/spike-findings-dns-discovery/` for implementation patterns.*
