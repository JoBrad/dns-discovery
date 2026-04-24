# Project: DNS Zone Discovery Tool

**Status:** Initialized  
**Created:** April 21, 2026  
**Version:** 1.0

## Vision

Build a DNS zone discovery tool that, given a domain name, produces an executive summary of the zone's configuration. Output includes registrar and nameserver identification (with friendly provider names), a list of configured services/hosts/redirects, and health checks for each service — validating email DNS (MX/SPF/DKIM/DMARC), TLS certificate validity, and minimum TLS 1.2+ for web-facing targets.

## What This Is

This project is a CLI-first Go application that performs DNS discovery and produces actionable reports per domain. It is designed for practical operations use: fast domain checks, repeatable outputs, and clear health signals for DNS, email, and TLS posture.

## Core Value

- Turn raw DNS data into operator-friendly summaries
- Detect provider posture (including split DNS) quickly
- Surface email and TLS risks in one execution flow
- Generate saved artifacts for auditability and follow-up

## Requirements

- Accept a domain argument and support batch execution from config/input files
- Discover key DNS record types and infer associated services/providers
- Evaluate email DNS health (MX, SPF, DMARC, DKIM) and TLS certificate health
- Produce report output in supported formats under the configured output directory
- Provide reliable error handling so one domain failure does not abort whole batch runs

## Core Deliverable

A CLI-first Go tool (`dns-discovery <domain>`) that outputs a Markdown report to the `output/` directory with:

1. **Zone Overview** — Quick executive summary
2. **Registrar & Nameservers** — Friendly provider names + split DNS detection
3. **Configured Services** — Hosts, redirects, SaaS tools detected from records
4. **Health Checks** — Per-service validation (email DNS health, TLS cert validity, TLS version)
5. **Executive Summary** — Key findings and alerts

## Foundation

### Spike Validation (Complete)

All core features validated through 4 spikes:
- **Spike 001**: DNS enumeration — All 9 record types, service detection
- **Spike 002**: NS fingerprinting — Provider detection, split DNS
- **Spike 003**: TLS health checks — Cert validation, version, expiry
- **Spike 004**: Email DNS health — MX/SPF/DKIM/DMARC scoring

Findings packaged in `.github/skills/spike-findings-dns-discovery/SKILL.md`

### Stack (Validated)

- **Language**: Go 1.23+
- **CLI Framework**: `cobra` (modern CLI framework for Go)
- **DNS**: `miekg/dns` Go library (equivalent to dnspython in functionality)
- **TLS**: `crypto/tls` stdlib (no external dependency)
- **Pattern Matching**: `regexp` stdlib package
- **Output**: Markdown to `output/` directory, eventual HTML via mk-docs

## Project Structure

```
.
├── main.go                    # Entry point (CLI setup)
├── cmd/
│   └── dns-discovery/         # CLI command package
│       └── main.go            # Entry point
├── internal/
│   ├── discovery/             # Core discovery logic
│   │   ├── dns.go             # DNS enumeration
│   │   ├── providers.go       # Provider fingerprinting
│   │   ├── tls.go             # TLS health checks
│   │   └── email.go           # Email DNS health validation
│   └── output/                # Output formatting
├── output/                    # Generated reports
├── .planning/
│   ├── PROJECT.md (this file)
│   ├── ROADMAP.md
│   ├── REQUIREMENTS.md
│   ├── STATE.md
│   ├── config.json
│   └── phases/
│       ├── 01-cli-tool/
│       ├── 02-reporting/
│       ├── 03-integration/
│       └── ...
├── go.mod                     # Go module definition
├── go.sum                     # Go dependency checksums
└── README.md
```

## Roadmap

Phase breakdown (detailed in ROADMAP.md):

1. **Phase 1: CLI Tool Foundation** — Core discovery engine (DNS, provider fingerprinting, TLS checks, email health)
2. **Phase 2: Reporting & Output** — Markdown formatting, output directory, HTML generation
3. **Phase 3: Integration & Polish** — Config file support, batch mode, error handling
4. **Phase 4+** — Advanced features, performance optimization, distribution

## Success Criteria

When Phase 1 is complete:
- `dns-discovery example.com` produces a structured discovery output in memory
- All 4 spike validations (DNS, providers, TLS, email) fully integrated
- Health checks working for real domains (github.com, cloudflare.com tested)
- CLI accepts domain argument, outputs results
- Code follows project conventions from spike work

## Key Dates

- **Spike Validation**: Completed April 21, 2026
- **Project Initialization**: April 21, 2026
- **Phase 1 Planning**: April 21, 2026
- **Phase 1 Implementation**: In progress

---

## Current State

**Shipped:** v1.0 — 2026-04-23  
**Status:** Milestone complete. 34/34 requirements satisfied.  
**Codebase:** 3,342 lines Go across `cmd/`, `internal/discovery`, `internal/app`, `internal/report`, `internal/config`  
**Stack:** Go 1.23+, cobra, miekg/dns, crypto/tls stdlib  

### Validated Requirements (v1.0)

- ✓ CLI accepts domain argument and batch sources — v1.0
- ✓ DNS enumeration (9 record types), provider fingerprinting (~60), split DNS detection — v1.0
- ✓ Email DNS health scoring (MX/SPF/DMARC/DKIM 4-pillar) with 27-selector DKIM probing — v1.0
- ✓ TLS cert validation, expiry warnings, TLS version detection — v1.0
- ✓ Markdown/JSON/text output flavors with configurable output dir and log location — v1.0
- ✓ Unified RunDiscovery orchestration with deterministic batch summaries — v1.0

## Next Milestone Goals

To be defined via `/gsd-new-milestone`.

---

*Last updated: 2026-04-23 after v1.0 milestone*
