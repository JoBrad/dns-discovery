---
phase: 02
slug: reporting
status: verified
threats_open: 0
asvs_level: 1
created: 2026-04-23
---

# Phase 02 — Security

> Per-phase security contract: threat register, accepted risks, and audit trail.

---

## Trust Boundaries

| Boundary | Description | Data Crossing |
|----------|-------------|---------------|
| DNS responses to markdown report generation | User-controlled DNS record strings are rendered into `report.md` content | Public DNS metadata (record values) |
| Local report output to git workspace | Generated discovery artifacts can be accidentally staged if not ignored | Domain discovery output files under `output/` |

---

## Threat Register

| Threat ID | Category | Component | Disposition | Mitigation | Status |
|-----------|----------|-----------|-------------|------------|--------|
| T-02-01 | Injection | `internal/report/markdown.go` | accept | Markdown output is local artifact consumption only; no HTML rendering or server execution path in phase scope | closed |
| T-02-02 | Information Disclosure | `.gitignore` / `output/` artifact handling | mitigate | `output/` is explicitly gitignored to prevent accidental report data commits | closed |

*Status: open · closed*  
*Disposition: mitigate (implementation required) · accept (documented risk) · transfer (third-party)*

---

## Accepted Risks Log

| Risk ID | Threat Ref | Rationale | Accepted By | Date |
|---------|------------|-----------|-------------|------|
| R-02-01 | T-02-01 | DNS values are rendered to local markdown files only; no executable rendering context in this phase | Project security policy | 2026-04-23 |

*Accepted risks do not resurface in future audit runs.*

---

## Security Audit Trail

| Audit Date | Threats Total | Closed | Open | Run By |
|------------|---------------|--------|------|--------|
| 2026-04-23 | 2 | 2 | 0 | gsd-secure-phase |

---

## Sign-Off

- [x] All threats have a disposition (mitigate / accept / transfer)
- [x] Accepted risks documented in Accepted Risks Log
- [x] `threats_open: 0` confirmed
- [x] `status: verified` set in frontmatter

**Approval:** verified 2026-04-23
