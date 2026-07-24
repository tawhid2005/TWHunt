package sources

import (
	"encoding/json"
	"fmt"
	"subfinder_clone/utils"
)

type URLScanSource struct{}

func (s *URLScanSource) Name() string {
	return "URLScan"
}

func (s *URLScanSource) Run(domain string) ([]string, error) {
	url := fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s", domain)
	
	body, err := utils.RequestAPI(url)
	if err != nil {
		return nil, err
	}

	var data struct {
		Results []struct {
			Page struct {
				Domain string `json:"domain"`
			} `json:"page"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var subdomains []string
	for _, res := range data.Results {
		if res.Page.Domain != "" {
			subdomains = append(subdomains, res.Page.Domain)
		}
	}

	return subdomains, nil
}
