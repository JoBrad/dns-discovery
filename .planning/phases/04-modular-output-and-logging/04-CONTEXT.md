# Phase 4: Modular output and logging - Context

**Gathered:** 2026-04-23
**Status:** Ready for planning

<domain>
## Phase Boundary

Modularize output generation and logging behavior without changing DNS discovery scope. Phase 4 delivers:
- Output selection via `output` with `markdown` (default), `json`, or `text`
- Unified run orchestration for single and batch flows in `internal/app/run.go`
- Logging controls (default error-level logging, verbose option, configurable log location)
- Report rendering in `internal/report` keyed by selected output flavor

Out of scope: adding new DNS checks, changing discovery algorithms, or changing provider/email/TLS logic.

</domain>

<decisions>
## Implementation Decisions

### Output Contract
- **D-01:** Use one term: `output`. Do not introduce a separate top-level `format` concept.
- **D-02:** `output` values are `markdown`, `json`, `text`; default is `markdown`.
- **D-03:** `output` defines the format of generated output file(s), including batch runs.

### Logging Policy
- **D-04:** Logging level and logging destination are separate concerns.
- **D-05:** Default log level is `error`.
- **D-06:** Add a `verbose` flag that enables verbose logging.
- **D-07:** Default log file location is `logs/` under repo root (same directory level as `go.mod`).
- **D-08:** `logLocation` can be set via config file or CLI argument with precedence: CLI > config > default.
- **D-09:** Error output is written to both log output and stderr.

### Runtime Orchestration Shape
- **D-10:** Consolidate execution into a single app entrypoint (`RunDiscovery`) for single and batch behavior.
- **D-11:** Replace public split between `RunDomain` and `RunBatch`; `RunDiscovery` accepts output selection and orchestrates scan + reporting.
- **D-12:** `BatchSummary` must support deterministic ordering for both succeeded and failed results (ordered lists for both keys).
- **D-13:** Batch output must preserve per-domain failure visibility.

### Package Boundaries
- **D-14:** Keep discovery layer in `internal/discovery` (no package relocation).
- **D-15:** Put output flavor/render functions in `internal/report`.
- **D-16:** `internal/app/run.go` should orchestrate obtaining scan results and passing them to report functions, not embed full report presentation logic.

### the agent's Discretion
- Exact Go type names and field shape for `DomainSummary` and failure entries as long as they satisfy D-12/D-13.
- Exact logger implementation (`slog` vs wrapper) as long as it satisfies D-04 through D-09.
- Exact file naming conventions for non-markdown outputs if deterministic and documented.

</decisions>

<specifics>
## Specific Ideas

- `internal/app/run.go` is the primary refactor target for this phase.
- `executeDomain` should produce summary data useful for output/report generation rather than printing directly.
- Existing markdown output behavior remains the baseline and must continue to work when `output=markdown`.

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Current runtime and CLI flow
- `cmd/dns-discovery/main.go` — current CLI flags, config loading, batch orchestration path
- `internal/app/run.go` — current execution/reporting/logging coupling to refactor

### Existing output implementation
- `internal/report/markdown.go` — current markdown renderer and save path behavior

### Requirements and roadmap context
- `.planning/ROADMAP.md` — Phase 4 section and sequencing
- `.planning/REQUIREMENTS.md` — existing requirement style and IDs for prior phases
- `.planning/phases/01-cli-tool/01-CONTEXT.md` — prior locked context conventions and canonical refs style

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/report/markdown.go`: complete markdown renderer and file writer already in place.
- `internal/config/config.go`: precedence resolution utility can be extended for log location/output config.

### Established Patterns
- CLI precedence model already used for output directory and config values.
- Batch execution currently aggregates successes/failures and returns summary metadata from app layer.

### Integration Points
- `cmd/dns-discovery/main.go` is the control point for new CLI flags/args (`output`, `verbose`, log location).
- `internal/app/run.go` is where execution orchestration can shift from print-driven to report/logger-driven flow.
- `internal/report` is the extension point for additional output flavors.

</code_context>

<deferred>
## Deferred Ideas

- Advanced logging backends/rotation policy (beyond local file + stdout behavior).
- Output transport beyond local files (S3, HTTP streaming) unless required by this phase plan.

</deferred>

---

*Phase: 04-modular-output-and-logging*
*Context gathered: 2026-04-23*
