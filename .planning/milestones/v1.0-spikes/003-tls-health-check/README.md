---
spike: "003"
name: tls-health-check
type: standard
validates: "Given an A/CNAME target hostname, when connected via TLS, then cert validity, expiry, TLS version, and failure reason are determined using stdlib ssl only"
verdict: VALIDATED
related: ["001", "004"]
tags: [tls, ssl, health-check, certificates, python, stdlib]
---

# Spike 003: TLS Health Check

## What This Validates
Given a hostname, when we open a TLS connection on port 443, then we can determine: is the cert trusted? Is it expired? What TLS version? And for failures — is it expired, self-signed, or a hostname mismatch?

## Research
Python stdlib `ssl` module is sufficient — no third-party library needed. `ssl.create_default_context()` performs full chain validation. `ssock.version()` returns the negotiated TLS version string ("TLSv1.2", "TLSv1.3"). `ssl.SSLCertVerificationError.verify_code` maps to OpenSSL error codes that distinguish failure modes.

Verified test endpoints from badssl.com:
- `expired.badssl.com` — intentionally expired cert
- `self-signed.badssl.com` — self-signed cert
- `wrong.host.badssl.com` — cert issued for different hostname

## How to Run
```
uv run python spike.py <hostname> [<hostname2> ...]
# Examples:
uv run python spike.py github.com cloudflare.com
uv run python spike.py expired.badssl.com self-signed.badssl.com wrong.host.badssl.com
```

## What to Expect
- ✓ = TLS version, cert valid, expiry date, issuer
- ⚠ = Expiring soon (<14d) or weak TLS (1.0/1.1)
- ✗ = Specific failure reason (expired / self-signed / hostname mismatch / unreachable)

## Investigation Trail

**Iteration 1 — Basic connection:** `ssl.create_default_context().wrap_socket()` works cleanly. `ssock.version()` returns "TLSv1.3" for modern sites.

**Edge case — cert parsing:** `getpeercert()` returns expiry as `"Jun  3 00:00:00 2026 GMT"` format, parsed with `strptime`. Days-until-expiry calculation works.

**Key finding — error discrimination:** All three badssl.com failure modes initially showed "CERTIFICATE_VERIFY_FAILED". Added OpenSSL verify code lookup:
- Code 10 = expired
- Code 18/19 = self-signed / untrusted root
- Code 62 = hostname mismatch

This produces specific, actionable error messages.

**Finding — CAA cross-reference:** Spike 001 found github.com has CAA records allowing Sectigo. TLS spike confirms cert is issued by Sectigo Limited. The real build can cross-reference these.

**Non-HTTP hosts:** A records pointing at non-web servers (e.g., mail servers, game servers) will hit ConnectionRefused or Timeout. Both are handled gracefully with specific error messages.

## Results
**Verdict: VALIDATED ✓**

- Stdlib `ssl` handles all tested scenarios — no third-party dependency needed
- All three cert failure modes produce distinct human-readable messages
- TLSv1.2 vs TLSv1.3 detection works
- Days-until-expiry enables proactive expiry warnings (<14d threshold)
- Non-HTTPS hosts handled gracefully (timeout/refused)
- github.com: TLSv1.3, issued by Sectigo (consistent with CAA policy from Spike 001)
