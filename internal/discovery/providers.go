package discovery

import (
	"strings"
)

// nsPatterns is the canonical ~60-entry provider table from Spike 002.
// Order matters: most-specific patterns first.
var nsPatterns = []struct {
	pattern  string
	provider string
}{
	// Cloud DNS
	{"awsdns", "AWS Route 53"},
	{"cloudflare.com", "Cloudflare"},
	{"nsone.net", "NS1 / IBM NS1 Connect"},
	{"googledomains.com", "Google Domains"},
	{"dns.google", "Google Cloud DNS"},
	{"azure-dns.com", "Azure DNS"},
	{"azure-dns.net", "Azure DNS"},
	{"azure-dns.org", "Azure DNS"},
	{"azure-dns.info", "Azure DNS"},
	{"digitaloceandns.com", "DigitalOcean DNS"},
	{"dnsimple.com", "DNSimple"},
	{"dnsimple.net", "DNSimple"},
	{"dynect.net", "Dyn / Oracle Cloud DNS"},
	{"ultradns.net", "UltraDNS"},
	{"ultradns.com", "UltraDNS"},
	{"ultradns.org", "UltraDNS"},
	{"ultradns.biz", "UltraDNS"},
	{"edns.biz", "Akamai Edge DNS"},
	{"akamai.net", "Akamai"},
	{"akam.net", "Akamai"},
	// Registrar-bundled DNS
	{"domaincontrol.com", "GoDaddy"},
	{"secureserver.net", "GoDaddy"},
	{"registrar-servers.com", "Namecheap"},
	{"namecheaphosting.com", "Namecheap"},
	{"name.com", "Name.com"},
	{"hover.com", "Hover"},
	{"namesilo.com", "NameSilo"},
	{"enom.com", "eNom"},
	{"networksolutions.com", "Network Solutions"},
	{"name-services.com", "Network Solutions"},
	{"web.com", "Web.com"},
	{"dotster.com", "Dotster"},
	{"register.com", "Register.com"},
	{"gkg.net", "GKG / TuCows"},
	{"tucows.com", "TuCows"},
	{"1and1.com", "IONOS / 1&1"},
	{"ionos.com", "IONOS"},
	{"ui-dns.com", "IONOS"},
	{"ui-dns.de", "IONOS"},
	{"ui-dns.biz", "IONOS"},
	{"ui-dns.org", "IONOS"},
	{"hichina.com", "Alibaba Cloud (Aliyun)"},
	{"alibabadns.com", "Alibaba Cloud (Aliyun)"},
	{"bluehost.com", "Bluehost"},
	{"hostgator.com", "HostGator"},
	{"siteground.net", "SiteGround"},
	{"inmotionhosting.com", "InMotion Hosting"},
	{"dreamhost.com", "DreamHost"},
	{"linode.com", "Linode / Akamai"},
	{"he.net", "Hurricane Electric"},
	{"afraid.org", "FreeDNS (afraid.org)"},
	{"cloudns.net", "ClouDNS"},
	{"pointhq.com", "PointHQ"},
	{"buddyns.com", "BuddyNS"},
	{"constellix.com", "Constellix"},
	{"easydns.com", "easyDNS"},
	{"easydns.net", "easyDNS"},
	{"rage4.com", "Rage4"},
	{"zendesk.com", "Zendesk"},
	{"shopify.com", "Shopify"},
	{"squarespace.com", "Squarespace"},
	{"parkingcrew.net", "ParkingCrew (parked domain)"},
	{"sedoparking.com", "Sedo Parking (parked domain)"},
	{"above.com", "Above.com (parked domain)"},
	{"huaweicloud.com", "Huawei Cloud DNS"},
}

// identifyNS returns the friendly provider name for a single NS hostname.
func identifyNS(nsHost, domain string) string {
	lower := strings.ToLower(strings.TrimSuffix(nsHost, "."))

	// Self-hosted: NS is a subdomain of the queried domain.
	if strings.HasSuffix(lower, "."+strings.ToLower(domain)) {
		return "Self-hosted (under " + domain + ")"
	}

	for _, p := range nsPatterns {
		if strings.Contains(lower, p.pattern) {
			return p.provider
		}
	}

	return "Unknown (" + lower + ")"
}

// IdentifyProviders maps NS hostnames to friendly provider names.
// It detects self-hosted DNS, split DNS, and unknown providers.
func IdentifyProviders(domain string, nsHosts []string) ProviderResult {
	result := ProviderResult{
		Counts:   make(map[string]int),
		AllHosts: nsHosts,
	}

	for _, ns := range nsHosts {
		provider := identifyNS(ns, domain)
		result.Counts[provider]++
	}

	// Primary = provider with the most NS records.
	maxCount := 0
	for provider, count := range result.Counts {
		if count > maxCount || (count == maxCount && result.Primary == "") {
			maxCount = count
			result.Primary = provider
		}
	}

	result.IsSplit = len(result.Counts) > 1

	return result
}
