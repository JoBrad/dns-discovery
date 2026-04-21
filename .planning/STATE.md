# Project State

**Last Updated:** April 21, 2026  
**Current Status:** Initialized — Ready for Phase 1 Planning

## Completed Work

### Spike Validation (✓ All Complete)

- ✓ Spike 001: DNS Record Enumeration — VALIDATED
- ✓ Spike 002: NS Registrar Fingerprinting — VALIDATED
- ✓ Spike 003: TLS Health Check — VALIDATED
- ✓ Spike 004: Email DNS Health — VALIDATED

**Spike Findings Packaged:** `.github/skills/spike-findings-dns-discovery/SKILL.md`

Findings include:
- 4 implementation reference guides (references/)
- 4 source spike files (sources/)
- Stack decisions (dnspython, stdlib ssl/re, Python 3.14.2)
- Tested domains (github.com, cloudflare.com, google.com, badssl.com)

## Current Decisions

### Stack (Locked)

- **Language:** Python 3.14.2
- **CLI Framework:** `click` (pending decision in Phase 1 planning)
- **DNS:** `dnspython` v2.6+ (VALIDATED ✓ in all 4 spikes)
- **TLS:** stdlib `ssl` (VALIDATED ✓ in spike 003)
- **Pattern Matching:** stdlib `re` (VALIDATED ✓ in spikes 001, 004)
- **Output Format:** Markdown (per README.md)

### Architecture Pattern

Four-pillar discovery architecture (from spikes):
1. **DNS Enumeration** — Query all 9 record types
2. **Provider Identification** — Map NS to friendly names
3. **TLS Health** — Check certs and protocol versions
4. **Email DNS Health** — MX/SPF/DMARC/DKIM validation

Each pillar is independently testable; orchestrate in `axeman/core.py`.

## Pending Decisions

- **CLI Framework:** `click` vs custom argparse vs other (to be decided in Phase 1 planning)
- **Output Location:** `.output/` directory name confirmed, but structure TBD in Phase 2
- **Batch Processing:** Deferred to Phase 3

## Project Todos

- [ ] Phase 1: Implement core discovery engine
- [ ] Phase 1: Integrate all 4 pillars into CLI
- [ ] Phase 1: Verify on github.com and cloudflare.com (UAT)
- [ ] Phase 2: Markdown report generation
- [ ] Phase 2: Output directory management
- [ ] Phase 3: Config file support
- [ ] Phase 3: Batch mode

## Known Constraints

- DKIM discovery is probe-based (27 common selectors) — some non-standard selectors may be missed
- Split DNS detection depends on NS record count breakdown
- Non-HTTPS A records timeout gracefully but provide limited info
- TLS checks timeout after 5 seconds (configurable)

## Blockers

None currently.

## Context Files

- `.planning/PROJECT.md` — Project vision and structure
- `.planning/REQUIREMENTS.md` — Phase 1 requirements
- `.planning/ROADMAP.md` — Full roadmap with phases
- `.github/skills/spike-findings-dns-discovery/` — Packaged spike learnings

---

*Project is ready for Phase 1 planning. Next step: `/gsd-plan-phase 1`*
