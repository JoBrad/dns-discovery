---
phase: 04
slug: modular-output-and-logging
status: complete
nyquist_compliant: true
wave_0_complete: true
created: 2026-04-23
---

# Phase 04 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none |
| **Quick run command** | `go test ./internal/config ./internal/app ./internal/report ./cmd/dns-discovery` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~15 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/config ./internal/app ./internal/report ./cmd/dns-discovery`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 60 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 04-01-01 | 01 | 1 | MOD-01 | — | Output enum validated; unsupported values rejected fast before scan execution. | unit | `go test ./internal/config ./cmd/dns-discovery -run 'TestLoad\|TestResolve'` | ✅ | ✅ green |
| 04-01-02 | 01 | 1 | LOG-01, LOG-02, LOG-03 | — | Log location resolves with CLI > config > default precedence; verbose flag exposed at app boundary. | unit | `go test ./internal/app -run 'TestValidateOutput'` && `go build ./...` | ✅ | ✅ green |
| 04-02-01 | 02 | 2 | MOD-02, MOD-03 | T-04-05 | Deterministic ordered success/failure from unified RunDiscovery; per-domain failures do not abort batch. | unit | `go test ./internal/app ./cmd/dns-discovery -run 'TestRunDiscovery'` | ✅ | ✅ green |
| 04-02-02 | 02 | 2 | MOD-04, ERR-02 | T-04-04 / T-04-06 | Explicit output enum dispatch; render/write errors surfaced with domain context to stderr and log sink. | unit | `go test ./internal/report ./internal/app ./cmd/dns-discovery && go build ./...` | ✅ | ✅ green |

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
