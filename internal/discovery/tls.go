package discovery

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	tlsPort    = 443
	tlsTimeout = 5 * time.Second
	expiryWarnDays = 14
)

// CheckTLS connects to hostname:443 and evaluates certificate health.
// Non-HTTPS hosts are handled gracefully — errors are categorised, not fatal.
func CheckTLS(hostname string) TLSResult {
	result := TLSResult{Hostname: hostname}

	dialer := &net.Dialer{Timeout: tlsTimeout}
	conn, err := tls.DialWithDialer(dialer, "tcp",
		fmt.Sprintf("%s:%d", hostname, tlsPort),
		&tls.Config{ServerName: hostname},
	)

	if err == nil {
		defer conn.Close()
		result.Reachable = true
		result.CertValid = true

		state := conn.ConnectionState()
		result.TLSVersion = tlsVersionName(state.Version)

		if len(state.PeerCertificates) > 0 {
			cert := state.PeerCertificates[0]
			result.CertExpiry = cert.NotAfter.UTC().Format("2006-01-02")
			days := int(time.Until(cert.NotAfter).Hours() / 24)
			result.DaysToExpiry = days
			result.CertExpired = days < 0
			result.ExpiryWarning = days >= 0 && days < expiryWarnDays
			result.Issuer = certIssuer(cert)
		}
		return result
	}

	// Classify the connection/verification error.
	result.Reachable = true // We got a response (even if a TLS error)

	var certInvalid x509.CertificateInvalidError
	var unknownAuth x509.UnknownAuthorityError
	var hostnameErr x509.HostnameError
	var netErr net.Error

	switch {
	case errors.As(err, &certInvalid):
		if certInvalid.Reason == x509.Expired {
			result.ErrorCategory = "EXPIRED"
			result.ErrorDetail = "Certificate has expired"
			result.CertExpired = true
			// Still try to extract expiry from the bad cert
			if certInvalid.Cert != nil {
				result.CertExpiry = certInvalid.Cert.NotAfter.UTC().Format("2006-01-02")
				result.DaysToExpiry = int(time.Until(certInvalid.Cert.NotAfter).Hours() / 24)
				result.Issuer = certIssuer(certInvalid.Cert)
			}
		} else {
			result.ErrorCategory = "CERT_INVALID"
			result.ErrorDetail = fmt.Sprintf("Certificate invalid: %v", certInvalid.Reason)
		}

	case errors.As(err, &unknownAuth):
		result.ErrorCategory = "SELF_SIGNED"
		result.ErrorDetail = "Self-signed certificate (not trusted by system)"

	case errors.As(err, &hostnameErr):
		result.ErrorCategory = "HOSTNAME_MISMATCH"
		result.ErrorDetail = fmt.Sprintf("Certificate hostname mismatch (cert is for %s)", hostnameErr.Certificate.Subject.CommonName)
		result.Issuer = certIssuer(hostnameErr.Certificate)

	case errors.As(err, &netErr) && netErr.Timeout():
		result.Reachable = false
		result.ErrorCategory = "TIMEOUT"
		result.ErrorDetail = fmt.Sprintf("Connection timed out after %.0fs", tlsTimeout.Seconds())

	default:
		errStr := err.Error()
		if strings.Contains(errStr, "connection refused") {
			result.Reachable = false
			result.ErrorCategory = "REFUSED"
			result.ErrorDetail = "Port 443 not open (not an HTTPS server)"
		} else if strings.Contains(errStr, "no such host") || strings.Contains(errStr, "lookup") {
			result.Reachable = false
			result.ErrorCategory = "DNS_ERROR"
			result.ErrorDetail = "DNS resolution failed"
		} else {
			result.ErrorCategory = "TLS_ERROR"
			result.ErrorDetail = err.Error()
		}
	}

	return result
}

func tlsVersionName(v uint16) string {
	switch v {
	case tls.VersionTLS10:
		return "TLSv1.0"
	case tls.VersionTLS11:
		return "TLSv1.1"
	case tls.VersionTLS12:
		return "TLSv1.2"
	case tls.VersionTLS13:
		return "TLSv1.3"
	default:
		return fmt.Sprintf("TLS 0x%04x", v)
	}
}

func certIssuer(cert *x509.Certificate) string {
	if len(cert.Issuer.Organization) > 0 {
		return cert.Issuer.Organization[0]
	}
	if cert.Issuer.CommonName != "" {
		return cert.Issuer.CommonName
	}
	return "Unknown"
}
