# TLS Health Checks

Connect to a hostname via TLS and verify certificate validity, expiry, and protocol version.

## Requirements

- Must check TLS certificate validity (trusted vs self-signed vs expired)
- Must check TLS certificate expiry with days-until-expiry
- Must check negotiated TLS protocol version (TLSv1.2, TLSv1.3, etc.)
- Must handle non-HTTPS hosts gracefully (timeout, connection refused)

## How to Build It

Use Go stdlib only (`crypto/tls`, `crypto/x509`, `net`, `time`).

1. Dial `hostname:443` with a short timeout (5s).
2. Use `tls.DialWithDialer` with `ServerName` set for hostname verification.
3. On success, extract TLS version and peer certificate metadata.
4. Compute days-to-expiry and warning window (<14 days).
5. On failure, classify errors into stable categories for reporting.

```go
func CheckTLS(hostname string) TLSResult {
    result := TLSResult{Hostname: hostname}
    dialer := &net.Dialer{Timeout: 5 * time.Second}

    conn, err := tls.DialWithDialer(dialer, "tcp", hostname+":443", &tls.Config{ServerName: hostname})
    if err == nil {
        defer conn.Close()
        result.Reachable = true
        result.CertValid = true

        state := conn.ConnectionState()
        result.TLSVersion = tlsVersionName(state.Version)
        if len(state.PeerCertificates) > 0 {
            cert := state.PeerCertificates[0]
            result.CertExpiry = cert.NotAfter.UTC().Format("2006-01-02")
            result.DaysToExpiry = int(time.Until(cert.NotAfter).Hours() / 24)
            result.CertExpired = result.DaysToExpiry < 0
            result.ExpiryWarning = result.DaysToExpiry >= 0 && result.DaysToExpiry < 14
        }
        return result
    }

    // classify x509 and network failures
    // EXPIRED, SELF_SIGNED, HOSTNAME_MISMATCH, TIMEOUT, REFUSED, DNS_ERROR, TLS_ERROR
    return classifyTLSError(result, err)
}
```

## What to Avoid

- Do not skip `ServerName` in TLS config or hostname checks break.
- Do not treat all TLS failures as generic handshake errors.
- Do not hard-fail discovery on non-HTTPS endpoints.
- Do not use long dial timeouts in batch scans.

## Constraints

- Port 443 must be reachable for successful health checks.
- Certificate validity is evaluated by system trust roots.
- TLS version names should be normalized for reporting consistency.

## Origin

Synthesized from spikes: 003

Source files available in:
- `sources/003-tls-health-check/README.md`
- `sources/003-tls-health-check/tls.go`
- `sources/003-tls-health-check/types.go`

**Key findings:**
- Go stdlib TLS/X509 stack is sufficient — no third-party dependency needed
- Error categories can be mapped to actionable outcomes
- CAA records from DNS spike can be cross-referenced with cert issuer from TLS spike
