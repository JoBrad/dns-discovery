# Project Roadmap

**Project:** DNS Zone Discovery Tool  
**Version:** 1.0  
**Last Updated:** April 21, 2026

## Milestone 1: Core Implementation

### Phase 1: CLI Tool Foundation

**Goal:** Build the core discovery engine with DNS enumeration, provider identification, TLS checks, and email health validation integrated into a working CLI.

**Status:** Ready for Planning

**Plans:** 1 plan planned

Plan list:
- [ ] 01-01-PLAN.md — Implement core discovery logic

**Requirement IDs:** CLI-01, ENUM-01, ENUM-02, PROV-01, PROV-02, PROV-03, EMAIL-01, EMAIL-02, EMAIL-03, EMAIL-04, EMAIL-05, TLS-01, TLS-02, TLS-03, TLS-04, OUT-01, OUT-02

**Dependencies:** None — first phase

**Deliverables:**
- Core discovery engine (`internal/discovery/` package and supporting modules)
- CLI interface accepting domain argument
- All health checks integrated
- Output formatting (stdout printing)
- Passes UAT on github.com and cloudflare.com

---

### Phase 2: Reporting & Output (Planned)

**Goal:** Convert discovery results into formatted Markdown reports, write to `output/` directory, and set up HTML generation via mk-docs.

**Status:** Planned

**Deliverables:**
- Markdown report generation
- `output/` directory management
- mk-docs configuration for HTML (future)

---

### Phase 3: Integration & Polish (Planned)

**Goal:** Add config file support, batch processing, advanced error handling.

---

## Summary

| Phase | Title | Status | Plans | Deliverable |
|-------|-------|--------|-------|-------------|
| 1 | CLI Tool Foundation | Ready for Planning | 1 | Core discovery engine + CLI |
| 2 | Reporting & Output | Planned | — | Markdown reports + output |
| 3 | Integration & Polish | Planned | — | Config, batch mode, error handling |

---

*Initialized from README.md and spike validation results. Phase 1 ready for detailed planning.*
