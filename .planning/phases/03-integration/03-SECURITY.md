---
phase: 03
slug: integration
status: verified
threats_open: 0
asvs_level: 1
created: 2026-04-23
---

# Phase 03 — Security

> Per-phase security contract: threat register, accepted risks, and audit trail.

---

## Trust Boundaries

| Boundary | Description | Data Crossing |
|----------|-------------|---------------|
| local config file to runtime settings | Untrusted local JSON config is parsed before execution | `output_dir`, `domains`, and related runtime options |
| CLI/input-file values to batch execution | User-provided domain strings drive DNS/TLS/email network lookups | domain names and execution control flags |
| batch summary to shell exit status | Aggregated per-domain failures influence automation/CI outcomes | success/failure counts and domain-level error text |

---

## Threat Register

| Threat ID | Category | Component | Disposition | Mitigation | Status |
|-----------|----------|-----------|-------------|------------|--------|
| T-03-01 | T | `internal/config/config.go` | mitigate | Strict JSON parsing with unknown-field rejection and normalization (`DisallowUnknownFields`, `TrimSpace`) prevents malformed config tampering from silently altering runtime behavior | closed |
| T-03-02 | D | `cmd/dns-discovery/main.go` | mitigate | Config and input errors fail fast at startup with bounded parsing and explicit error returns before broad network activity | closed |
| T-03-03 | I | CLI error output | accept | Config file path is included in local CLI error messages for debugging; disclosure is low risk in single-user local execution context | closed |
| T-03-04 | T | `internal/app/run.go` | mitigate | Per-domain execution records failures by domain and continues remaining work, preventing one malformed entry from corrupting whole-run state | closed |
| T-03-05 | D | batch execution flow | mitigate | Batch processing executes deterministically one domain at a time without unbounded concurrency or fan-out | closed |
| T-03-06 | R | CLI batch summary and exit handling | mitigate | Deterministic aggregate summaries and non-zero exit on failures provide auditable, automation-safe outcomes | closed |

*Status: open · closed*  
*Disposition: mitigate (implementation required) · accept (documented risk) · transfer (third-party)*

---

## Accepted Risks Log

| Risk ID | Threat Ref | Rationale | Accepted By | Date |
|---------|------------|-----------|-------------|------|
| R-03-01 | T-03-03 | Path disclosure in local CLI error output is acceptable for diagnosability; no multi-tenant or remote exposure in phase scope | Project security policy | 2026-04-23 |

*Accepted risks do not resurface in future audit runs.*

---

## Security Audit Trail

| Audit Date | Threats Total | Closed | Open | Run By |
|------------|---------------|--------|------|--------|
| 2026-04-23 | 6 | 6 | 0 | gsd-secure-phase |

---

## Sign-Off

- [x] All threats have a disposition (mitigate / accept / transfer)
- [x] Accepted risks documented in Accepted Risks Log
- [x] `threats_open: 0` confirmed
- [x] `status: verified` set in frontmatter

**Approval:** verified 2026-04-23
