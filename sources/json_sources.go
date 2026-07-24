package sources

import (
	"encoding/json"
	"fmt"
	"strings"
	"subfinder_clone/utils"
)

// JSON-based Sources

type AnubisSource struct{}
func (s *AnubisSource) Name() string { return "Anubis" }
func (s *AnubisSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://jldc.me/anubis/subdomains/%s", domain))
	if err != nil { return nil, err }
	var data []string
	if err := json.Unmarshal(body, &data); err != nil { return nil, err }
	var subs []string
	for _, sub := range data {
		if strings.HasSuffix(sub, domain) { subs = append(subs, sub) }
	}
	return subs, nil
}

type BeVigilSource struct{}
func (s *BeVigilSource) Name() string { return "BeVigil" }
func (s *BeVigilSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://osint.bevigil.com/api/%s/subdomains/", domain))
	if err != nil { return nil, err }
	var data struct { Subdomains []string `json:"subdomains"` }
	if err := json.Unmarshal(body, &data); err != nil { return nil, err }
	return data.Subdomains, nil
}

type BufferOverSource struct{}
func (s *BufferOverSource) Name() string { return "BufferOver" }
func (s *BufferOverSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://dns.bufferover.run/dns?q=.%s", domain))
	if err != nil { return nil, err }
	var data struct { FDNS_A []string `json:"FDNS_A"` }
	if err := json.Unmarshal(body, &data); err != nil { return nil, err }
	var subs []string
	for _, entry := range data.FDNS_A {
		parts := strings.Split(entry, ",")
		if len(parts) > 1 { subs = append(subs, parts[1]) }
	}
	return subs, nil
}

type FullHuntSource struct{}
func (s *FullHuntSource) Name() string { return "FullHunt" }
func (s *FullHuntSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://fullhunt.io/api/v1/domain/%s/subdomains", domain))
	if err != nil { return nil, err }
	var data struct { Hosts []string `json:"hosts"` }
	if err := json.Unmarshal(body, &data); err != nil { return nil, err }
	return data.Hosts, nil
}

type OmnisintSource struct{}
func (s *OmnisintSource) Name() string { return "Omnisint" }
func (s *OmnisintSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://sonar.omnisint.io/all/%s", domain))
	if err != nil { return nil, err }
	var data []string
	if err := json.Unmarshal(body, &data); err != nil { return nil, err }
	return data, nil
}

type SubdomainCenterSource struct{}
func (s *SubdomainCenterSource) Name() string { return "SubCenter" }
func (s *SubdomainCenterSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://api.subdomain.center/?domain=%s", domain))
	if err != nil { return nil, err }
	var data []string
	if err := json.Unmarshal(body, &data); err != nil { return nil, err }
	return data, nil
}

type VirusTotalSource struct{}
func (s *VirusTotalSource) Name() string { return "VirusTotal" }
func (s *VirusTotalSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://www.virustotal.com/ui/domains/%s/subdomains?limit=40", domain))
	if err != nil { return nil, err }
	var data struct {
		Data []struct { ID string `json:"id"` } `json:"data"`
	}
	if err := json.Unmarshal(body, &data); err != nil { return nil, err }
	var subs []string
	for _, entry := range data.Data {
		subs = append(subs, entry.ID)
	}
	return subs, nil
}

type ThreatCrowdSource struct{}
func (s *ThreatCrowdSource) Name() string { return "ThreatCrowd" }
func (s *ThreatCrowdSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://ci-www.threatcrowd.org/searchApi/v2/domain/report/?domain=%s", domain))
	if err != nil { return nil, err }
	var data struct { Subdomains []string `json:"subdomains"` }
	if err := json.Unmarshal(body, &data); err != nil { return nil, err }
	var subs []string
	for _, sub := range data.Subdomains {
		if strings.HasSuffix(sub, domain) { subs = append(subs, sub) }
	}
	return subs, nil
}

type ThreatMinerSource struct{}
func (s *ThreatMinerSource) Name() string { return "ThreatMiner" }
func (s *ThreatMinerSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://api.threatminer.org/v2/domain.php?q=%s&rt=5", domain))
	if err != nil { return nil, err }
	var data struct { Results []string `json:"results"` }
	if err := json.Unmarshal(body, &data); err != nil { return nil, err }
	var subs []string
	for _, sub := range data.Results {
		if strings.HasSuffix(sub, domain) { subs = append(subs, sub) }
	}
	return subs, nil
}

type ShodanCTSource struct{}
func (s *ShodanCTSource) Name() string { return "ShodanCT" }
func (s *ShodanCTSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://ctl.shodan.io/api/v1/domain/%s/hostnames", domain))
	if err != nil { return nil, err }
	var data []string
	if err := json.Unmarshal(body, &data); err != nil { return nil, err }
	return data, nil
}
