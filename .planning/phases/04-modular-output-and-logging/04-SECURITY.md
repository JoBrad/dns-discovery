---
phase: 04
slug: modular-output-and-logging
status: verified
threats_open: 0
asvs_level: 1
created: 2026-04-23
---

# Phase 04 — Security

> Per-phase security contract: threat register, accepted risks, and audit trail.

---

## Trust Boundaries

| Boundary | Description | Data Crossing |
|----------|-------------|---------------|
| discovery result → renderer | Untrusted/variable scan data transformed into output artifacts | DNS record values, TLS metadata, provider names |
| renderer output → filesystem | Generated content written to user-controlled output locations | Rendered report files (markdown/json/text) |
| domain-level errors → stderr/log sink | Failure information crosses user-visible and persistent log channels | Error messages with domain context |

---

## Threat Register

| Threat ID | Category | Component | Disposition | Mitigation | Status |
|-----------|----------|-----------|-------------|------------|--------|
| T-04-04 | Tampering | internal/report/output.go | mitigate | `ValidateOutputFormat` and `SaveReportByFormat` use explicit `switch` statements; default case returns a structured error — no implicit fallback to unknown flavors | closed |
| T-04-05 | Repudiation | internal/app/run.go | mitigate | `RunDiscovery` appends to `summary.Succeeded` and `summary.Failed` in input-slice order, producing deterministic per-domain audit records | closed |
| T-04-06 | Information Disclosure | stderr/log error surfacing | mitigate | Errors wrapped with `fmt.Errorf` include domain and operation context (e.g. `"✗ %s: %v"`, `"discovery failed for %s"`); raw stack traces not emitted | closed |

*Status: open · closed*
*Disposition: mitigate (implementation required) · accept (documented risk) · transfer (third-party)*

---

## Accepted Risks Log

No accepted risks.

---

## Security Audit Trail

| Audit Date | Threats Total | Closed | Open | Run By |
|------------|---------------|--------|------|--------|
| 2026-04-23 | 3 | 3 | 0 | gsd-secure-phase |

---

## Sign-Off

- [x] All threats have a disposition (mitigate / accept / transfer)
- [x] Accepted risks documented in Accepted Risks Log
- [x] `threats_open: 0` confirmed
- [x] `status: verified` set in frontmatter

**Approval:** verified 2026-04-23
