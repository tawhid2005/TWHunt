package sources

import (
	"encoding/json"
	"fmt"
	"strings"
	"subfinder_clone/utils"
)

type CrtshSource struct{}

func (s *CrtshSource) Name() string {
	return "crt.sh"
}

func (s *CrtshSource) Run(domain string) ([]string, error) {
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	
	body, err := utils.RequestAPI(url)
	if err != nil {
		return nil, err
	}

	var results []struct {
		NameValue string `json:"name_value"`
	}

	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}

	var subdomains []string
	for _, res := range results {
		names := strings.Split(res.NameValue, "\n")
		for _, name := range names {
			name = strings.TrimSpace(name)
			name = strings.TrimPrefix(name, "*.")
			if name != "" {
				subdomains = append(subdomains, name)
			}
		}
	}

	return subdomains, nil
}
