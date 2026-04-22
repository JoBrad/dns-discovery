# Project Roadmap

**Project:** DNS Zone Discovery Tool  
**Version:** 1.0  
**Last Updated:** April 22, 2026

## Milestone 1: Core Implementation

### Phase 1: CLI Tool Foundation

**Goal:** Build the core discovery engine with DNS enumeration, provider identification, TLS checks, and email health validation integrated into a working CLI.

**Status:** ✅ Complete

**Plans:** 2 plans

Plan list:
- [x] 01-01-PLAN.md — Build Go CLI foundation, DNS enumeration, and provider fingerprinting
- [x] 01-02-PLAN.md — Add TLS/email health pillars and render executive summary

**Requirement IDs:** CLI-01, ENUM-01, ENUM-02, PROV-01, PROV-02, PROV-03, EMAIL-01, EMAIL-02, EMAIL-03, EMAIL-04, EMAIL-05, TLS-01, TLS-02, TLS-03, TLS-04, OUT-01, OUT-02

**Dependencies:** None — first phase

**Deliverables:**
- Core discovery engine (`internal/discovery/` package)
- CLI interface accepting domain argument
- All health checks integrated
- Output formatting (stdout printing)
- Passes UAT on github.com and cloudflare.com

---

### Phase 2: Reporting & Output

**Goal:** Enrich the Markdown report to match the full stdout output (DNS records table, detected services, split DNS detail, complete email/TLS sections) and ensure generated reports are gitignored.

**Status:** ✅ Complete

**Plans:** 1 plan

Plan list:
- [x] 02-01-PLAN.md — Enrich GenerateMarkdown to cover all stdout sections; add output/ to .gitignore

**Requirement IDs:** RPT-01, RPT-02, RPT-03

**Dependencies:** Phase 1 complete

**Deliverables:**
- Markdown report includes DNS records table (all 9 types)
- Markdown report includes detected services (email, hosting/CDN, verification)
- Markdown report includes split DNS detail
- Markdown report includes MX records with priorities
- Markdown report includes TLS version and days-to-expiry
- `output/` directory is gitignored

---

### Phase 3: Integration & Polish

**Goal:** Add config file support, batch processing, and resilient error handling for multi-domain runs.

**Status:** Built (verify pending)

**Plans:** 2 plans

Plan list:
- [x] 03-01-PLAN.md — Add JSON config loading, defaults, and CLI override precedence
- [x] 03-02-PLAN.md — Add batch domain processing and aggregate error reporting

**Requirement IDs:** CFG-01, CFG-02, BAT-01, BAT-02, ERR-01

**Dependencies:** Phase 2 complete

**Deliverables:**
- Optional config file support via `.dns-discovery.json` or `--config`
- CLI flags override config values predictably
- Batch processing from config or input file
- Per-domain failures do not abort whole batch runs
- Batch runs finish with a clear success/failure summary and non-zero exit on failures

---

## Summary

| Phase | Title | Status | Plans | Deliverable |
|-------|-------|--------|-------|-------------|
| 1 | CLI Tool Foundation | ✅ Complete | 2 | Core discovery engine + CLI |
| 2 | Reporting & Output | ✅ Complete | 1 | Complete Markdown reports |
| 3 | Integration & Polish | Built (verify pending) | 2 | Config, batch mode, error handling |

---

*Initialized from README.md and spike validation results. Phase 1 complete.*
