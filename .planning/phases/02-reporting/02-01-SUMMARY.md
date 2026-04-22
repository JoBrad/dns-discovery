# Phase 02-reporting / Plan 01 — SUMMARY

**Status:** ✅ Complete  
**Commit:** f471b7f  
**Date:** April 22, 2026

## What Was Built

### Task 1: Enrich GenerateMarkdown
Rewrote `internal/report/markdown.go` `GenerateMarkdown` function to add all sections
present in stdout but missing from the file report:

- **Executive Summary:** Added Split DNS line when `Provider.IsSplit` is true
- **Infrastructure & DNS:** Added Split DNS Providers table (sorted, provider → NS count)
- **DNS Records section (new):** All 9 record types in canonical order (A/AAAA/MX/NS/TXT/CNAME/SOA/CAA/SRV), `—` for absent types, `<br>` join for multi-value types
- **Detected Services section (new):** Email Providers, Hosting & CDN, Verification & SaaS subsections (each rendered only when non-empty)
- **Email Security — MX Records subsection (new):** Priority/Host table; "No MX records configured" when empty
- **TLS Health table:** Added TLS version and Days columns; ⚠️ prefix when `ExpiryWarning`; ErrorCategory shown for unreachable hosts
- Added `sortedKeys` helper for deterministic map iteration

### Task 2: Gitignore output/
- Added `output/` to `.gitignore`
- Ran `git rm --cached output/` to untrack 3 previously committed report files

## Verification

- `go build ./...` — ✅ clean
- `github.com` report: 6 `##` headings, Split DNS table present, TLS shows TLSv1.3/42d
- `cloudflare.com` report: 6 `##` headings, no Split DNS section (single provider)
- `git status` — output/ not tracked

## Requirements Closed

- RPT-01 ✅ — DNS Records table in report.md
- RPT-02 ✅ — Detected Services section in report.md  
- RPT-03 ✅ — Split DNS, MX priorities, TLS version/days all present

## Files Modified

- `internal/report/markdown.go` — enriched GenerateMarkdown, added sortedKeys, added sort import
- `.gitignore` — added output/
- `output/{cloudflare,github,google}.com/report.md` — deleted from git tracking (now gitignored)
