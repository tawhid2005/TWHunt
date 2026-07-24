package core

import (
	"fmt"
	"net/http"
	"sync"
)

var VulnPaths = []string{
	"/.env",
	"/.git/config",
	"/phpinfo.php",
	"/.DS_Store",
	"/server-status",
	"/wp-config.php.bak",
}

// ProbeVulns checks common sensitive endpoints on live subdomains
func ProbeVulns(subdomains []string, silent bool) map[string][]string {
	if !silent {
		fmt.Printf(" %s[*] PROBING FOR LOW-HANGING VULNERABILITIES...%s\n", Gold, EndC)
	}

	results := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20) // lower concurrency to avoid immediate bans

	for _, sub := range subdomains {
		for _, path := range VulnPaths {
			wg.Add(1)
			sem <- struct{}{}
			go func(s, p string) {
				defer wg.Done()
				defer func() { <-sem }()

				url := "http://" + s + p
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("User-Agent", "Mozilla/5.0")
				
				resp, err := jsClient.Do(req)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				// If it returns 200 and it's not a generic 404 disguised as 200
				if resp.StatusCode == 200 && resp.ContentLength > 0 && resp.ContentLength < 100000 {
					mu.Lock()
					results[s] = append(results[s], url)
					mu.Unlock()
				}
			}(sub, path)
		}
	}

	wg.Wait()

	if !silent {
		if len(results) > 0 {
			fmt.Printf(" %s[!] FOUND POTENTIAL EXPOSED FILES ON %d HOSTS!%s\n", Coral, len(results), EndC)
		} else {
			fmt.Printf(" %s[✓] NO EXPOSED FILES FOUND.%s\n", Mint, EndC)
		}
	}
	return results
}
