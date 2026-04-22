package discovery

import (
	"regexp"
	"strings"
)

// dkimSelectors is the canonical 27-selector probe list from Spike 004.
var dkimSelectors = []string{
	"google",
	"selector1",  // Microsoft 365
	"selector2",  // Microsoft 365 alternate
	"default",
	"mail",
	"k1",
	"k2",
	"s1",
	"s2",
	"smtp",
	"dkim",
	"email",
	"key1",
	"key2",
	"sig1",
	"sig2",
	"mimecast",
	"mailjet",
	"sendgrid",
	"mandrill",
	"sparkpost",
	"amazonses",
	"protonmail",
	"zoho",
	"fm1", // Fastmail
	"fm2",
	"fm3",
}

var (
	reSPFAll    = regexp.MustCompile(`([~\-\+\?])all\b`)
	reDMARCPol  = regexp.MustCompile(`\bp=(\w+)`)
	reDKIMKey   = regexp.MustCompile(`k=(\w+)`)
	reDKIMEmpty = regexp.MustCompile(`p=\s*[;$]|p=\s*$`)
)

// checkSPF extracts and evaluates the v=spf1 TXT record.
func checkSPF(domain string) SPFResult {
	txts := QueryTXT(domain)
	var spfRecords []string
	for _, t := range txts {
		if strings.HasPrefix(t, "v=spf1") {
			spfRecords = append(spfRecords, t)
		}
	}

	if len(spfRecords) == 0 {
		return SPFResult{Present: false}
	}

	result := SPFResult{Present: true, Record: spfRecords[0]}

	if len(spfRecords) > 1 {
		result.Issues = append(result.Issues,
			"Multiple SPF records found — only one is valid per RFC 7208")
	}

	spf := spfRecords[0]
	m := reSPFAll.FindStringSubmatch(spf)
	if m == nil {
		result.Policy = "unknown"
		result.Issues = append(result.Issues, "Missing 'all' mechanism")
	} else {
		switch m[1] {
		case "-":
			result.Policy = "hardfail"
		case "~":
			result.Policy = "softfail"
		case "?":
			result.Policy = "neutral"
		case "+":
			result.Policy = "pass"
			result.Insecure = true
			result.Issues = append(result.Issues, "+all is insecure — permits any sender")
		}
	}

	return result
}

// checkDMARC queries the _dmarc subdomain and evaluates the policy.
func checkDMARC(domain string) DMARCResult {
	txts := QueryTXT("_dmarc." + domain)
	var dmarcRecord string
	for _, t := range txts {
		if strings.HasPrefix(t, "v=DMARC1") {
			dmarcRecord = t
			break
		}
	}

	if dmarcRecord == "" {
		return DMARCResult{Present: false}
	}

	result := DMARCResult{Present: true, Record: dmarcRecord}

	if m := reDMARCPol.FindStringSubmatch(dmarcRecord); m != nil {
		result.Policy = m[1]
	} else {
		result.Policy = "none"
	}

	result.HasAggReports = strings.Contains(dmarcRecord, "rua=")
	result.HasForensics = strings.Contains(dmarcRecord, "ruf=")

	return result
}

// checkDKIM probes the canonical 27 selectors and returns found keys.
func checkDKIM(domain string) []DKIMSelector {
	var found []DKIMSelector
	for _, sel := range dkimSelectors {
		dkimDomain := sel + "._domainkey." + domain
		txts := QueryTXT(dkimDomain)
		for _, t := range txts {
			if strings.Contains(t, "v=DKIM1") || strings.Contains(t, "p=") {
				ds := DKIMSelector{Selector: sel, KeyType: "rsa"}
				if m := reDKIMKey.FindStringSubmatch(t); m != nil {
					ds.KeyType = m[1]
				}
				ds.Revoked = reDKIMEmpty.MatchString(t)
				found = append(found, ds)
				break
			}
		}
	}
	return found
}

// EvaluateEmailHealth runs all four email DNS health checks and scores them.
func EvaluateEmailHealth(domain string) EmailResult {
	result := EmailResult{}

	result.MXRecords = QueryMX(domain)
	result.SPF = checkSPF(domain)
	result.DMARC = checkDMARC(domain)
	result.DKIM = checkDKIM(domain)

	// 4-pillar scoring: MX / SPF / DMARC / DKIM
	score := 0
	if len(result.MXRecords) > 0 {
		score++
	}
	if result.SPF.Present {
		score++
	}
	if result.DMARC.Present {
		score++
	}
	if len(result.DKIM) > 0 {
		score++
	}

	result.Score = score
	result.ScoreText = string(rune('0'+score)) + "/4"

	return result
}
