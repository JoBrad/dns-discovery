# TLS Health Checks

Connect to a hostname via TLS and verify certificate validity, expiry, and protocol version.

## Requirements

- Must check TLS certificate validity (trusted vs self-signed vs expired)
- Must check TLS certificate expiry with days-until-expiry
- Must check negotiated TLS protocol version (TLSv1.2, TLSv1.3, etc.)
- Must handle non-HTTPS hosts gracefully (timeout, connection refused)

## How to Build It

**Install:** No external dependencies — use stdlib `ssl` + `socket`.

**Pattern:**
```python
import ssl
import socket
from datetime import datetime, timezone

def check_tls(hostname, port=443, timeout=5):
    result = {
        "hostname": hostname,
        "reachable": False,
        "tls_version": None,
        "cert_valid": False,
        "cert_expired": False,
        "cert_expiry": None,
        "days_until_expiry": None,
        "error": None,
    }
    
    ctx = ssl.create_default_context()
    
    try:
        with socket.create_connection((hostname, port), timeout=timeout) as sock:
            with ctx.wrap_socket(sock, server_hostname=hostname) as ssock:
                result["reachable"] = True
                result["tls_version"] = ssock.version()
                
                cert = ssock.getpeercert()
                
                # Parse expiry date
                expiry_str = cert.get("notAfter", "")
                if expiry_str:
                    expiry_dt = datetime.strptime(expiry_str, "%b %d %H:%M:%S %Y %Z")
                    expiry_dt = expiry_dt.replace(tzinfo=timezone.utc)
                    now = datetime.now(timezone.utc)
                    days_left = (expiry_dt - now).days
                    result["cert_expiry"] = expiry_dt.strftime("%Y-%m-%d")
                    result["days_until_expiry"] = days_left
                    result["cert_expired"] = days_left < 0
                
                result["cert_valid"] = not result["cert_expired"]
                
    except ssl.SSLCertVerificationError as e:
        result["reachable"] = True
        result["cert_valid"] = False
        # Map verify_code to specific error
        code = getattr(e, "verify_code", None)
        if code == 10:
            result["error"] = "Certificate EXPIRED"
            result["cert_expired"] = True
        elif code in (18, 19):
            result["error"] = "Self-signed certificate (not trusted)"
        elif code == 62:
            result["error"] = "Hostname mismatch"
        else:
            result["error"] = f"Cert verification failed: {e.reason}"
    
    except socket.timeout:
        result["error"] = f"Timeout after {timeout}s"
    except ConnectionRefusedError:
        result["error"] = "Connection refused (port 443 not open)"
    except socket.gaierror as e:
        result["error"] = f"DNS resolution failed: {e}"
    
    return result
```

## What to Avoid

- Don't ignore OpenSSL error codes — different verify_code values mean different failure modes (expired vs self-signed vs hostname mismatch)
- Don't assume all A records point to HTTPS servers — gracefully handle ConnectionRefused and Timeout
- Don't use a long timeout — 5 seconds is reasonable for health checks

## Constraints

- Port 443 must be open — non-HTTPS services will timeout or refuse
- TLS handshake can reveal the issuer, but SANs and subject CN extraction requires parsing `getpeercert()` carefully
- Certificate chain validation happens inside `ssl.create_default_context()` — can't customize without disabling security

## Origin

Spike 003: `tls-health-check`
Source: `.planning/spikes/003-tls-health-check/spike.py`

**Key findings:**
- Stdlib `ssl` is sufficient — no third-party dependency needed
- OpenSSL verify codes (10, 18/19, 62) map to specific failure modes
- CAA records from DNS spike can be cross-referenced with cert issuer from TLS spike
