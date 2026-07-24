package sources

import (
	"encoding/json"
	"fmt"
	"twhunt/utils"
)

type AlienVaultSource struct{}

func (s *AlienVaultSource) Name() string {
	return "AlienVault"
}

func (s *AlienVaultSource) Run(domain string) ([]string, error) {
	url := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/passive_dns", domain)
	
	body, err := utils.RequestAPI(url)
	if err != nil {
		return nil, err
	}

	var data struct {
		PassiveDns []struct {
			Hostname string `json:"hostname"`
		} `json:"passive_dns"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var subdomains []string
	for _, entry := range data.PassiveDns {
		if entry.Hostname != "" {
			subdomains = append(subdomains, entry.Hostname)
		}
	}

	return subdomains, nil
}
