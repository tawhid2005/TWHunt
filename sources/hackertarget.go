package sources

import (
	"fmt"
	"twhunt/utils"
)

type HackerTargetSource struct{}

func (s *HackerTargetSource) Name() string {
	return "HackerTarget"
}

func (s *HackerTargetSource) Run(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain)
	
	body, err := utils.RequestAPI(url)
	if err != nil {
		return nil, err
	}

	return utils.ExtractAll(string(body), domain), nil
}
