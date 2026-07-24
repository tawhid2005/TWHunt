package sources

import (
	"fmt"
	"twhunt/utils"
)

// Regex-based Sources

type AbuseIPDBSource struct{}
func (s *AbuseIPDBSource) Name() string { return "AbuseIPDB" }
func (s *AbuseIPDBSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://www.abuseipdb.com/whois/%s", domain))
	if err != nil { return nil, err }
	return utils.ExtractAll(string(body), domain), nil
}

type CommonCrawlSource struct{}
func (s *CommonCrawlSource) Name() string { return "CommonCrawl" }
func (s *CommonCrawlSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("http://index.commoncrawl.org/CC-MAIN-2023-50-index?url=*.%s/*&output=json", domain))
	if err != nil { return nil, err }
	return utils.ExtractAll(string(body), domain), nil
}

type NetcraftSource struct{}
func (s *NetcraftSource) Name() string { return "Netcraft" }
func (s *NetcraftSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://searchdns.netcraft.com/?restriction=site+ends+with&host=%s", domain))
	if err != nil { return nil, err }
	return utils.ExtractAll(string(body), domain), nil
}

type RapidDNSSource struct{}
func (s *RapidDNSSource) Name() string { return "RapidDNS" }
func (s *RapidDNSSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://rapiddns.io/s/%s?full=1&down=1", domain))
	if err != nil { return nil, err }
	return utils.ExtractAll(string(body), domain), nil
}

type SiteDossierSource struct{}
func (s *SiteDossierSource) Name() string { return "SiteDossier" }
func (s *SiteDossierSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("http://www.sitedossier.com/parentdomain/%s", domain))
	if err != nil { return nil, err }
	return utils.ExtractAll(string(body), domain), nil
}

type SynapsintSource struct{}
func (s *SynapsintSource) Name() string { return "Synapsint" }
func (s *SynapsintSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://synapsint.com/report.php?domain=%s", domain))
	if err != nil { return nil, err }
	return utils.ExtractAll(string(body), domain), nil
}

type WaybackSource struct{}
func (s *WaybackSource) Name() string { return "Wayback" }
func (s *WaybackSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=*.%s/*&output=json&collapse=urlkey&fl=original", domain))
	if err != nil { return nil, err }
	return utils.ExtractAll(string(body), domain), nil
}

type ShrewdEyeSource struct{}
func (s *ShrewdEyeSource) Name() string { return "ShrewdEye" }
func (s *ShrewdEyeSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://shrewdeye.app/domains/%s.txt", domain))
	if err != nil { return nil, err }
	return utils.ExtractAll(string(body), domain), nil
}

type RiddlerSource struct{}
func (s *RiddlerSource) Name() string { return "Riddler" }
func (s *RiddlerSource) Run(domain string) ([]string, error) {
	body, err := utils.RequestAPI(fmt.Sprintf("https://riddler.io/search?q=pld:%s&view_type=data_table", domain))
	if err != nil { return nil, err }
	return utils.ExtractAll(string(body), domain), nil
}
