# Email DNS Health

Validate email infrastructure by checking MX, SPF, DMARC, and DKIM records.

## Requirements

- Must check email DNS health: MX records present, SPF (v=spf1), DKIM (_domainkey), DMARC (_dmarc)
- Must produce a health score for email configuration
- Must identify email providers from MX records and DKIM selectors

## How to Build It

Use existing DNS query helpers plus regex-based policy parsing.

1. Query MX records and format as `priority hostname`.
2. Query TXT at zone apex; detect SPF (`v=spf1`) and evaluate policy/risks.
3. Query TXT at `_dmarc.<domain>` and parse policy/reporting flags.
4. Probe 27 common DKIM selectors at `<selector>._domainkey.<domain>`.
5. Compute 4-pillar score: MX, SPF, DMARC, DKIM.

```go
func EvaluateEmailHealth(domain string) EmailResult {
    result := EmailResult{}
    result.MXRecords = QueryMX(domain)
    result.SPF = checkSPF(domain)
    result.DMARC = checkDMARC(domain)
    result.DKIM = checkDKIM(domain)

    score := 0
    if len(result.MXRecords) > 0 { score++ }
    if result.SPF.Present { score++ }
    if result.DMARC.Present { score++ }
    if len(result.DKIM) > 0 { score++ }

    result.Score = score
    result.ScoreText = string(rune('0'+score)) + "/4"
    return result
}
```

```go
var reSPFAll = regexp.MustCompile(`([~\-\+\?])all\b`)

func checkSPF(domain string) SPFResult {
    txts := QueryTXT(domain)
    // select all v=spf1 records, warn on multiples
    // parse all mechanism and classify hardfail/softfail/neutral/pass
}
```

```go
var reDMARCPol = regexp.MustCompile(`\bp=(\w+)`)

func checkDMARC(domain string) DMARCResult {
    txts := QueryTXT("_dmarc." + domain)
    // find v=DMARC1, parse p= value, check rua=/ruf=
}
```

```go
func checkDKIM(domain string) []DKIMSelector {
    var found []DKIMSelector
    for _, sel := range dkimSelectors { // 27 selectors
        txts := QueryTXT(sel + "._domainkey." + domain)
        // detect v=DKIM1/p=, parse key type, mark revoked on empty p=
    }
    return found
}
```

## What to Avoid

- Do not assume DKIM absence means no DKIM exists; selector probing is partial.
- Do not pass multiple SPF records silently as valid.
- Do not parse SPF/DMARC without regex validation of key fields.
- Do not block full discovery if one email pillar fails.

## Constraints

- DKIM discovery is probe-based and inherently incomplete.
- SPF and TXT records vary widely by provider and can include many includes/verifications.
- DMARC policy semantics remain `none`, `quarantine`, `reject`.

## Origin

Synthesized from spikes: 004

Source files available in:
- `sources/004-email-dns-health/README.md`
- `sources/004-email-dns-health/email.go`
- `sources/004-email-dns-health/types.go`

**Key findings:**
- DKIM probe list covers 27 selectors — covers Microsoft 365, Google, Mandrill, Amazon SES, Zoho, Fastmail
- Selector names reveal provider details (e.g., `mandrill._domainkey` → Mailchimp)
- SPF `include:` chains + DKIM selectors + Spike 001 MX detection paint a complete email provider picture
- Health score: MX ✓/✗, SPF ✓/✗, DMARC ✓/✗, DKIM ✓/✗ (4/4 = all pillars)
