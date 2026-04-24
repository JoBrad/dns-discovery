---
phase: 03
slug: integration
status: complete
nyquist_compliant: true
wave_0_complete: true
created: 2026-04-23
---

# Phase 03 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none |
| **Quick run command** | `go test ./internal/config ./internal/app ./cmd/dns-discovery` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~15 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/config ./internal/app ./cmd/dns-discovery`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 60 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 03-01-01 | 01 | 1 | CFG-01 | T-03-01 / T-03-02 | Reject malformed config and fail fast with file-aware parse errors. | unit | `go test ./internal/config -run 'TestLoad'` | ✅ | ✅ green |
| 03-01-02 | 01 | 1 | CFG-02 | T-03-01 / — | Enforce precedence `flags > config > defaults` without unsafe side effects. | unit | `go test ./internal/config -run 'TestResolvePrefersFlagOverConfig' && go build ./...` | ✅ | ✅ green |
| 03-02-01 | 02 | 2 | BAT-02 | T-03-04 / T-03-05 | Continue batch execution after per-domain failures and isolate errors by domain. | unit | `go test ./internal/app -run 'TestRunBatch'` | ✅ | ✅ green |
| 03-02-02 | 02 | 2 | BAT-01, ERR-01 | T-03-04 / T-03-06 | Accept valid `.txt` domain input, reject invalid entries, and preserve actionable failures. | unit | `go test ./cmd/dns-discovery -run 'TestLoadDomainsFrom' && go test ./internal/config` | ✅ | ✅ green |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

Existing infrastructure covers all phase requirements.

---

## Manual-Only Verifications

All phase behaviors have automated verification.

---

## Validation Sign-Off

- [x] All tasks have `<automated>` verify or Wave 0 dependencies
- [x] Sampling continuity: no 3 consecutive tasks without automated verify
- [x] Wave 0 covers all MISSING references
- [x] No watch-mode flags
- [x] Feedback latency < 60s
- [x] `nyquist_compliant: true` set in frontmatter

**Approval:** approved 2026-04-23

## Validation Audit 2026-04-23 (Re-audit)

| Metric | Count |
|--------|-------|
| Gaps found | 0 |
| Resolved | 0 |
| Escalated | 0 |
