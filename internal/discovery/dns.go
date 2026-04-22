package discovery

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"
)

// ErrNXDOMAIN is returned when the queried domain does not exist.
var ErrNXDOMAIN = errors.New("domain does not exist (NXDOMAIN)")

const defaultResolver = "8.8.8.8:53"

var recordTypes = []struct {
	name  string
	qtype uint16
}{
	{"A", dns.TypeA},
	{"AAAA", dns.TypeAAAA},
	{"MX", dns.TypeMX},
	{"NS", dns.TypeNS},
	{"TXT", dns.TypeTXT},
	{"CNAME", dns.TypeCNAME},
	{"SOA", dns.TypeSOA},
	{"CAA", dns.TypeCAA},
	{"SRV", dns.TypeSRV},
}

// Service detection patterns translated from spike 001.
var mxServicePatterns = map[string]string{
	"google.com":                  "Google Workspace",
	"googlemail.com":              "Google Workspace",
	"outlook.com":                 "Microsoft 365",
	"protection.outlook.com":      "Microsoft 365",
	"pphosted.com":                "Proofpoint",
	"mimecast.com":                "Mimecast",
	"mailgun.org":                 "Mailgun",
	"sendgrid.net":                "SendGrid",
	"amazonses.com":               "Amazon SES",
	"zoho.com":                    "Zoho Mail",
	"fastmail.com":                "Fastmail",
	"inbound.cf-emailsecurity.net": "Cloudflare Email Security",
}

var txtServicePatterns = map[string]string{
	"v=spf1":                      "SPF Record",
	"v=DMARC1":                    "DMARC Policy",
	"google-site-verification":    "Google Search Console",
	"MS=ms":                       "Microsoft 365",
	"facebook-domain-verification": "Facebook",
	"docusign":                    "DocuSign",
	"atlassian-domain-verification": "Atlassian",
	"stripe-verification":         "Stripe",
	"adobe-idp-site-verification": "Adobe",
	"apple-domain-verification":   "Apple",
	"ZOOM_verify":                 "Zoom",
}

var cnameServicePatterns = map[string]string{
	"cloudfront.net":  "AWS CloudFront",
	"amazonaws.com":   "AWS",
	"azurewebsites.net": "Azure App Service",
	"fastly.net":      "Fastly CDN",
	"github.io":       "GitHub Pages",
	"netlify.app":     "Netlify",
	"vercel.app":      "Vercel",
	"heroku.com":      "Heroku",
	"pages.dev":       "Cloudflare Pages",
	"workers.dev":     "Cloudflare Workers",
	"shopify.com":     "Shopify",
	"zendesk.com":     "Zendesk",
	"hubspot.com":     "HubSpot",
	"wixdns.net":      "Wix",
}

// queryRaw sends a DNS query over UDP, retrying over TCP if the response is
// truncated. Network/timeout errors return an empty slice with rcode 0.
func queryRaw(domain string, qtype uint16) ([]dns.RR, int, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	m.RecursionDesired = true

	udpClient := &dns.Client{Timeout: 5 * time.Second, Net: "udp"}
	r, _, err := udpClient.Exchange(m, defaultResolver)
	if err != nil {
		return nil, dns.RcodeSuccess, nil
	}

	// Retry over TCP when UDP response was truncated.
	if r.Truncated {
		tcpClient := &dns.Client{Timeout: 5 * time.Second, Net: "tcp"}
		rTCP, _, errTCP := tcpClient.Exchange(m, defaultResolver)
		if errTCP == nil {
			r = rTCP
		}
	}

	return r.Answer, r.Rcode, nil
}

// QueryTXT is exported so the email package can reuse it.
func QueryTXT(domain string) []string {
	answers, rcode, _ := queryRaw(domain, dns.TypeTXT)
	if rcode == dns.RcodeNameError {
		return nil
	}
	var results []string
	for _, ans := range answers {
		if t, ok := ans.(*dns.TXT); ok {
			results = append(results, strings.Join(t.Txt, ""))
		}
	}
	return results
}

// QueryMX is exported so the email package can reuse it.
func QueryMX(domain string) []string {
	answers, rcode, _ := queryRaw(domain, dns.TypeMX)
	if rcode == dns.RcodeNameError {
		return nil
	}
	var results []string
	for _, ans := range answers {
		if mx, ok := ans.(*dns.MX); ok {
			results = append(results, fmt.Sprintf("%d %s", mx.Preference, strings.TrimSuffix(mx.Mx, ".")))
		}
	}
	return results
}

// rrToString converts a DNS resource record to a readable string.
func rrToString(rr dns.RR) string {
	switch v := rr.(type) {
	case *dns.A:
		return v.A.String()
	case *dns.AAAA:
		return v.AAAA.String()
	case *dns.MX:
		return fmt.Sprintf("%d %s", v.Preference, strings.TrimSuffix(v.Mx, "."))
	case *dns.NS:
		return strings.TrimSuffix(v.Ns, ".")
	case *dns.TXT:
		return strings.Join(v.Txt, "")
	case *dns.CNAME:
		return strings.TrimSuffix(v.Target, ".")
	case *dns.SOA:
		return fmt.Sprintf("%s %s serial=%d",
			strings.TrimSuffix(v.Ns, "."),
			strings.TrimSuffix(v.Mbox, "."),
			v.Serial)
	case *dns.CAA:
		return fmt.Sprintf("%d %s \"%s\"", v.Flag, v.Tag, v.Value)
	case *dns.SRV:
		return fmt.Sprintf("%d %d %d %s",
			v.Priority, v.Weight, v.Port,
			strings.TrimSuffix(v.Target, "."))
	default:
		return rr.String()
	}
}

// QueryAllRecords queries all 9 supported DNS record types.
// Returns ErrNXDOMAIN if the domain does not exist.
func QueryAllRecords(domain string) (DNSRecords, error) {
	records := make(DNSRecords)

	for _, rt := range recordTypes {
		answers, rcode, _ := queryRaw(domain, rt.qtype)

		// NXDOMAIN on A query is the reliable signal that the domain itself doesn't exist.
		if rcode == dns.RcodeNameError && rt.qtype == dns.TypeA {
			return nil, ErrNXDOMAIN
		}

		for _, ans := range answers {
			records[rt.name] = append(records[rt.name], rrToString(ans))
		}
	}

	return records, nil
}

// DetectServices scans DNS records for known service patterns.
func DetectServices(records DNSRecords) DetectedServices {
	var svc DetectedServices
	seen := map[string]bool{}

	add := func(category *[]string, name string) {
		if !seen[name] {
			*category = append(*category, name)
			seen[name] = true
		}
	}

	for _, mx := range records["MX"] {
		lower := strings.ToLower(mx)
		for pattern, service := range mxServicePatterns {
			if strings.Contains(lower, pattern) {
				add(&svc.Email, service)
			}
		}
	}

	for _, txt := range records["TXT"] {
		for pattern, service := range txtServicePatterns {
			if strings.Contains(txt, pattern) {
				add(&svc.VerificationServices, service)
			}
		}
	}

	for _, cname := range records["CNAME"] {
		lower := strings.ToLower(cname)
		for pattern, service := range cnameServicePatterns {
			if strings.Contains(lower, pattern) {
				add(&svc.HostingCDN, service)
			}
		}
	}

	return svc
}
