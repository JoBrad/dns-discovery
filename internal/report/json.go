package report

import (
	"encoding/json"

	"github.com/jbradley/dns-discovery/internal/discovery"
)

func GenerateJSON(res *discovery.DiscoveryResult) (string, error) {
	payload, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return "", err
	}
	return string(payload) + "\n", nil
}
