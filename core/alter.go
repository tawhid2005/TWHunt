package core

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

var AlterationWords = []string{
	"api", "dev", "staging", "test", "prod", "beta", "v1", "v2", 
	"admin", "internal", "corp", "mail", "auth", "vpn", "portal",
}

// AlterSubdomains takes existing valid subdomains and applies mutations to find hidden ones
func AlterSubdomains(subdomains []string, domain string, silent bool) []string {
	if !silent {
		fmt.Printf(" %s[*] RUNNING SUBDOMAIN PERMUTATION ENGINE...%s\n", Gold, EndC)
	}

	var results []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 50)

	for _, sub := range subdomains {
		// Only alter the first part of the subdomain (e.g. dev from dev.target.com)
		parts := strings.Split(strings.TrimSuffix(sub, "."+domain), ".")
		if len(parts) == 0 || parts[0] == "" || parts[0] == domain {
			continue
		}
		
		base := parts[0]
		
		for _, word := range AlterationWords {
			mutations := []string{
				fmt.Sprintf("%s-%s.%s", base, word, domain), // dev-api.target.com
				fmt.Sprintf("%s-%s.%s", word, base, domain), // api-dev.target.com
				fmt.Sprintf("%s%s.%s", base, word, domain),  // devapi.target.com
			}

			for _, m := range mutations {
				wg.Add(1)
				sem <- struct{}{}
				go func(target string) {
					defer wg.Done()
					defer func() { <-sem }()
					
					if _, err := net.LookupHost(target); err == nil {
						mu.Lock()
						results = append(results, target)
						mu.Unlock()
					}
				}(m)
			}
		}
	}

	wg.Wait()

	if !silent {
		if len(results) > 0 {
			fmt.Printf(" %s[✓] FOUND %d NEW SUBDOMAINS VIA ALTERATION!%s\n", Mint, len(results), EndC)
		} else {
			fmt.Printf(" %s[-] No new subdomains found via alteration.%s\n", Silver, EndC)
		}
	}
	return results
}
