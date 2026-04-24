---
phase: 01
slug: cli-tool
status: verified
nyquist_compliant: true
wave_0_complete: true
created: 2026-04-23
---

# Phase 01 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | go.mod |
| **Quick run command** | `go test ./cmd/dns-discovery ./internal/discovery` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~20 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./cmd/dns-discovery ./internal/discovery`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 30 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 01-01-01 | 01 | 1 | CLI-01 | T-01-01 | Domain input is normalized and invalid/empty inputs are rejected before discovery calls | unit | `go test ./cmd/dns-discovery -run 'TestResolveDomainsPrefersPositionalArgument|TestResolveDomainsUsesConfigDomainsWhenNoArgsOrFile'` | ✅ | ✅ green |
| 01-01-02 | 01 | 1 | ENUM-01 | T-01-02 | DNS enumeration queries supported record sets and yields structured outputs | integration | `go test ./internal/discovery -run TestQueryAllRecordsReturnsDataForKnownDomain` | ✅ | ✅ green |
| 01-01-03 | 01 | 1 | ENUM-02 | T-01-03 | Service detection maps MX/TXT/CNAME patterns into service categories | unit | `go test ./internal/discovery -run TestDetectServicesClassifiesPatterns` | ✅ | ✅ green |
| 01-01-04 | 01 | 1 | PROV-01 | T-01-05 | NS hosts map to known friendly providers | unit | `go test ./internal/discovery -run TestIdentifyProvidersDetectsSplitAndSelfHosted` | ✅ | ✅ green |
| 01-01-05 | 01 | 1 | PROV-02 | T-01-05 | NS under queried domain is labeled self-hosted | unit | `go test ./internal/discovery -run TestIdentifyProvidersDetectsSplitAndSelfHosted` | ✅ | ✅ green |
| 01-01-06 | 01 | 1 | PROV-03 | T-01-05 | Multi-provider NS sets are reported as split DNS | unit | `go test ./internal/discovery -run TestIdentifyProvidersDetectsSplitAndSelfHosted` | ✅ | ✅ green |
| 01-02-01 | 02 | 2 | TLS-01 | T-01-07 | Healthy TLS endpoint returns protocol/certificate metadata | integration | `go test ./internal/discovery -run TestCheckTLSReturnsMetadataForHealthyHost` | ✅ | ✅ green |
| 01-02-02 | 02 | 2 | TLS-02 | T-01-07 | Expired/invalid TLS endpoints are categorized without process crash | integration | `go test ./internal/discovery -run TestCheckTLSClassifiesExpiredEndpoint` | ✅ | ✅ green |
| 01-02-03 | 02 | 2 | TLS-03 | T-01-09 | TLS metadata includes expiry-day signals for certificate lifecycle monitoring | integration | `go test ./internal/discovery -run TestCheckTLSReturnsMetadataForHealthyHost` | ✅ | ✅ green |
| 01-02-04 | 02 | 2 | TLS-04 | T-01-09 | Non-ideal TLS conditions are surfaced non-fatally and execution continues | integration | `go test ./internal/discovery -run TestCheckTLSClassifiesExpiredEndpoint` | ✅ | ✅ green |
| 01-02-05 | 02 | 2 | EMAIL-01 | T-01-08 | Email evaluation includes MX record discovery | integration | `go test ./internal/discovery -run TestEvaluateEmailHealthReturnsScoreAndPillars` | ✅ | ✅ green |
| 01-02-06 | 02 | 2 | EMAIL-02 | T-01-08 | SPF policy parsing contributes to email posture | integration | `go test ./internal/discovery -run TestEvaluateEmailHealthReturnsScoreAndPillars` | ✅ | ✅ green |
| 01-02-07 | 02 | 2 | EMAIL-03 | T-01-08 | DMARC presence/policy contributes to email posture | integration | `go test ./internal/discovery -run TestEvaluateEmailHealthReturnsScoreAndPillars` | ✅ | ✅ green |
| 01-02-08 | 02 | 2 | EMAIL-04 | T-01-08 | DKIM selector probing contributes to email posture | integration | `go test ./internal/discovery -run TestEvaluateEmailHealthReturnsScoreAndPillars` | ✅ | ✅ green |
| 01-02-09 | 02 | 2 | EMAIL-05 | T-01-08 | Email score reports normalized {0..4}/4 pillar health | integration | `go test ./internal/discovery -run TestEvaluateEmailHealthReturnsScoreAndPillars` | ✅ | ✅ green |
| 01-02-10 | 02 | 2 | OUT-01 | T-01-10 | Discovery result contracts remain machine-consumable across pillars | integration | `go test ./internal/discovery -run 'TestQueryAllRecordsReturnsDataForKnownDomain|TestEvaluateEmailHealthReturnsScoreAndPillars'` | ✅ | ✅ green |
| 01-02-11 | 02 | 2 | OUT-02 | T-01-10 | End-user receives readable, sectioned output without fatal interruption | integration | `go test ./cmd/dns-discovery -run 'TestResolveDomainsPrefersPositionalArgument|TestResolveDomainsUsesConfigDomainsWhenNoArgsOrFile'` | ✅ | ✅ green |

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
- [x] Feedback latency < 30s
- [x] `nyquist_compliant: true` set in frontmatter

**Approval:** approved 2026-04-23

## Validation Audit 2026-04-23

| Metric | Count |
|--------|-------|
| Gaps found | 6 |
| Resolved | 6 |
| Escalated | 0 |
