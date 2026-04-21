# Project: DNS Zone Discovery Tool

**Status:** Initialized  
**Created:** April 21, 2026  
**Version:** 1.0

## Vision

Build a DNS zone discovery tool that, given a domain name, produces an executive summary of the zone's configuration. Output includes registrar and nameserver identification (with friendly provider names), a list of configured services/hosts/redirects, and health checks for each service — validating email DNS (MX/SPF/DKIM/DMARC), TLS certificate validity, and minimum TLS 1.2+ for web-facing targets.

## Core Deliverable

A CLI-first Python tool (`dns-discovery <domain>`) that outputs a Markdown report to the `output/` directory with:

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

- **Language**: Python 3.14.2
- **CLI Framework**: `click` (for CLI interaction)
- **DNS**: `dnspython` v2.6+ (proven across all 4 spikes)
- **TLS**: stdlib `ssl` (no external dependency)
- **Pattern Matching**: stdlib `re`
- **Output**: Markdown to `output/` directory, eventual HTML via mk-docs

## Project Structure

```
.
├── main.py                    # Entry point (CLI setup)
├── axeman/
│   ├── __init__.py
│   ├── core.py               # Main discovery orchestration
│   ├── certlib.py            # TLS/certificate logic
│   └── [other modules]       # DNS, provider, email modules
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

*Initialized via project setup. Next: Phase 1 planning.*
