# dns-discovery

CLI tool for DNS zone discovery and lightweight security posture checks.

Given one or more domains, the tool:

- Queries core DNS records (`A`, `AAAA`, `MX`, `NS`, `TXT`, `CNAME`, `SOA`, `CAA`, `SRV`)
- Fingerprints nameserver providers with friendly names (for example, GoDaddy, Cloudflare, Route 53)
- Detects common services from DNS patterns (email, hosting/CDN, verification/SaaS)
- Evaluates email DNS health (MX, SPF, DMARC, DKIM) with a 4-point score
- Performs an HTTPS/TLS certificate check on the domain
- Writes a Markdown report per domain under an output directory

## Current Implementation Status

This repository is currently implemented in Go and uses Cobra for the CLI.

- Entry point: `cmd/dns-discovery/main.go`
- Domain discovery pipeline: `internal/app/run.go`
- DNS and service detection: `internal/discovery/`
- Markdown report generation: `internal/report/markdown.go`

## Requirements

- Go 1.26+

## Build

```bash
go build ./cmd/dns-discovery
```

This produces a `dns-discovery` binary in the current directory.

## CLI Usage

```text
Discover DNS configuration, provider fingerprint, TLS health,
and email DNS health for any domain.

Usage:
  dns-discovery [domain] [flags]

Flags:
      --config string       Path to JSON config file
  -h, --help                help for dns-discovery
      --input-file string   Path to newline-delimited domains file
  -o, --output-dir string   Directory to save reports (default "output")
```

## Example Usage

### 1) Single domain

```bash
go run ./cmd/dns-discovery github.com
```

Writes report to:

```text
output/github.com/report.md
```

### 2) Single domain with custom output directory

```bash
go run ./cmd/dns-discovery github.com --output-dir reports
```

Writes report to:

```text
reports/github.com/report.md
```

### 3) Batch mode from input file

Create a newline-delimited file:

```text
github.com
cloudflare.com
not-a-real-domain.invalid
```

Run:

```bash
go run ./cmd/dns-discovery --input-file domain_list.txt
```

Behavior:

- Continues processing all domains
- Prints a batch summary with succeeded and failed domains
- Returns a non-zero exit code if any domain fails

### 4) Batch mode from config

If no domain argument and no `--input-file` are provided, the tool reads `domains` from config.

Create `.dns-discovery.json` in the working directory:

```json
{
  "output_dir": "output",
  "domains": [
    "github.com",
    "cloudflare.com"
  ]
}
```

Run:

```bash
go run ./cmd/dns-discovery
```

### 5) Explicit config file path

```bash
go run ./cmd/dns-discovery --config test_config.json
```

## Config File

JSON schema currently supported:

```json
{
  "output_dir": "output",
  "domains": ["example.com", "example.org"]
}
```

Notes:

- Unknown JSON fields are rejected
- Empty/whitespace values are normalized or rejected as appropriate
- CLI `--output-dir` overrides `output_dir` from config

## What the Report Contains

Each domain report is Markdown and currently includes:

- Executive summary
- Nameserver provider analysis (including split-DNS detection)
- DNS records table
- Detected services
- Email security section (MX/SPF/DMARC/DKIM and score)
- TLS health table

Example generated report:

- `output/github.com/report.md`

## Implementation Notes and Limitations

- DNS queries currently use resolver `8.8.8.8:53`
- TLS check is performed once per domain hostname (not per every A/AAAA target)
- Provider and service detection are pattern-based fingerprints

## Run Tests

```bash
go test ./...
```

Contributor guidance (including test comment conventions) is documented in `CONTRIBUTING.md`.

See: [CONTRIBUTING.md](CONTRIBUTING.md)

