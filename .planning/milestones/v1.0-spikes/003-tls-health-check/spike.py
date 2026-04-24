"""
Spike 003: TLS Health Check
Checks TLS certificate validity, expiry, and minimum protocol version for a hostname.

Usage: python spike.py <hostname> [<hostname2> ...]
Example: python spike.py github.com cloudflare.com expired.badssl.com
"""

import ssl
import socket
import sys
from datetime import datetime, timezone
from typing import Optional


def check_tls(hostname: str, port: int = 443, timeout: int = 5) -> dict:
    """
    Connect to hostname:port via TLS and inspect the certificate.
    Returns a dict with health check results.
    """
    result = {
        "hostname": hostname,
        "port": port,
        "reachable": False,
        "tls_version": None,
        "cert_valid": False,
        "cert_expired": False,
        "cert_expiry": None,
        "days_until_expiry": None,
        "cert_issuer": None,
        "cert_subject": None,
        "error": None,
    }

    ctx = ssl.create_default_context()

    try:
        with socket.create_connection((hostname, port), timeout=timeout) as sock:
            with ctx.wrap_socket(sock, server_hostname=hostname) as ssock:
                result["reachable"] = True
                result["tls_version"] = ssock.version()

                cert = ssock.getpeercert()

                # Parse expiry
                expiry_str = cert.get("notAfter", "")
                if expiry_str:
                    expiry_dt = datetime.strptime(expiry_str, "%b %d %H:%M:%S %Y %Z")
                    expiry_dt = expiry_dt.replace(tzinfo=timezone.utc)
                    now = datetime.now(timezone.utc)
                    days_left = (expiry_dt - now).days
                    result["cert_expiry"] = expiry_dt.strftime("%Y-%m-%d")
                    result["days_until_expiry"] = days_left
                    result["cert_expired"] = days_left < 0

                # If we got here, SSL validation passed (cert is trusted + not expired by ssl module)
                result["cert_valid"] = not result["cert_expired"]

                # Extract issuer
                issuer_dict = dict(x[0] for x in cert.get("issuer", []))
                result["cert_issuer"] = issuer_dict.get("organizationName", "Unknown")

                # Extract subject CN
                subject_dict = dict(x[0] for x in cert.get("subject", []))
                result["cert_subject"] = subject_dict.get("commonName", "Unknown")

    except ssl.SSLCertVerificationError as e:
        result["reachable"] = True  # We connected, but cert is invalid
        result["cert_valid"] = False
        # Distinguish common failure modes via verify_code
        code = getattr(e, "verify_code", None)
        if code == 10:
            result["error"] = "Certificate EXPIRED"
            result["cert_expired"] = True
        elif code in (18, 19):
            result["error"] = "Self-signed certificate (not trusted)"
        elif code == 62:
            result["error"] = "Hostname mismatch (cert issued for different domain)"
        else:
            result["error"] = (
                f"Cert verification failed: {getattr(e, 'verify_message', e.reason)}"
            )
    except ssl.SSLError as e:
        result["reachable"] = True
        result["cert_valid"] = False
        result["error"] = f"SSL error: {e}"
    except ConnectionRefusedError:
        result["error"] = "Connection refused (port 443 not open)"
    except socket.timeout:
        result["error"] = f"Timeout after {timeout}s"
    except socket.gaierror as e:
        result["error"] = f"DNS resolution failed: {e}"
    except OSError as e:
        result["error"] = f"Network error: {e}"

    return result


def tls_verdict(result: dict) -> tuple[str, str]:
    """Returns (status_emoji, summary) for a result."""
    if result["error"] and not result["reachable"]:
        return "✗", result["error"]
    if not result["cert_valid"]:
        return "✗", result.get("error") or "Certificate invalid"
    if result["cert_expired"]:
        return "✗", "Certificate EXPIRED"
    if result["days_until_expiry"] is not None and result["days_until_expiry"] < 14:
        return "⚠", f"Certificate expiring in {result['days_until_expiry']} days"

    tls = result.get("tls_version", "")
    if tls in ("TLSv1", "TLSv1.1") or tls is None:
        return "⚠", f"Weak TLS: {tls}"

    return (
        "✓",
        f"TLS {tls}, cert valid until {result['cert_expiry']} ({result['days_until_expiry']}d)",
    )


def print_result(result: dict):
    emoji, summary = tls_verdict(result)
    print(f"\n  {emoji} {result['hostname']}")
    print(f"     {summary}")
    if result["cert_issuer"]:
        print(f"     Issuer:  {result['cert_issuer']}")
    if result["cert_subject"]:
        print(f"     Subject: {result['cert_subject']}")
    if result["tls_version"]:
        print(f"     TLS:     {result['tls_version']}")


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python spike.py <hostname> [<hostname2> ...]")
        print("Example: python spike.py github.com cloudflare.com")
        sys.exit(1)

    hostnames = sys.argv[1:]

    print(f"\n{'=' * 60}")
    print(
        f"  TLS Health Check ({len(hostnames)} host{'s' if len(hostnames) > 1 else ''})"
    )
    print(f"{'=' * 60}")

    results = [check_tls(h) for h in hostnames]
    for r in results:
        print_result(r)

    ok = sum(1 for r in results if tls_verdict(r)[0] == "✓")
    warn = sum(1 for r in results if tls_verdict(r)[0] == "⚠")
    fail = sum(1 for r in results if tls_verdict(r)[0] == "✗")
    print(f"\n  {'=' * 40}")
    print(f"  Results: {ok} OK  {warn} WARN  {fail} FAIL")
    print()
