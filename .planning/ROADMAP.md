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

**Status:** In Planning

**Plans:** 1 plan

Plan list:
- [ ] 02-01-PLAN.md — Enrich GenerateMarkdown to cover all stdout sections; add output/ to .gitignore

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

### Phase 3: Integration & Polish (Planned)

**Goal:** Add config file support, batch processing, advanced error handling.

---

## Summary

| Phase | Title | Status | Plans | Deliverable |
|-------|-------|--------|-------|-------------|
| 1 | CLI Tool Foundation | ✅ Complete | 2 | Core discovery engine + CLI |
| 2 | Reporting & Output | In Planning | 1 | Complete Markdown reports |
| 3 | Integration & Polish | Planned | — | Config, batch mode, error handling |

---

*Initialized from README.md and spike validation results. Phase 1 complete.*
