package sources

import (
	"encoding/json"
	"fmt"
	"strings"
	"subfinder_clone/utils"
)

type CertSpotterSource struct{}

func (s *CertSpotterSource) Name() string {
	return "CertSpotter"
}

func (s *CertSpotterSource) Run(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&include_subdomains=true&expand=dns_names", domain)
	
	body, err := utils.RequestAPI(url)
	if err != nil {
		return nil, err
	}

	var data []struct {
		DnsNames []string `json:"dns_names"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var subdomains []string
	for _, entry := range data {
		for _, name := range entry.DnsNames {
			name = strings.TrimPrefix(name, "*.")
			if name != "" {
				subdomains = append(subdomains, name)
			}
		}
	}

	return subdomains, nil
}
