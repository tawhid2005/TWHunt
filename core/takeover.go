package core

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

var TakeoverSignatures = map[string]string{
	"github.io":              "GitHub Pages",
	"herokuapp.com":          "Heroku",
	"s3.amazonaws.com":       "AWS S3",
	"s3-website":             "AWS S3",
	"azurewebsites.net":      "Azure",
	"cloudapp.net":           "Azure",
	"trafficmanager.net":     "Azure TrafficManager",
	"ghost.io":               "Ghost",
	"bitbucket.io":           "Bitbucket",
	"pantheonsite.io":        "Pantheon",
	"shopify.com":            "Shopify",
	"myshopify.com":          "Shopify",
	"zendesk.com":            "Zendesk",
	"readme.io":              "Readme.io",
	"wordpress.com":          "WordPress",
	"surge.sh":               "Surge.sh",
	"strikingly.com":         "Strikingly",
	"strikinglydns.com":      "Strikingly",
	"fastly.net":             "Fastly",
	"helpscoutdocs.com":      "Help Scout",
	"statuspage.io":          "StatusPage",
	"bounceme.net":           "No-IP",
	"cargocollective.com":    "Cargo Collective",
	"cloudfront.net":         "CloudFront",
	"netlify.com":            "Netlify",
	"pages.dev":              "Cloudflare Pages",
}

// CheckTakeovers scans a list of subdomains for potential CNAME takeovers
func CheckTakeovers(subdomains []string, silent bool) map[string]string {
	if !silent {
		fmt.Printf(" %s[*] SCANNING FOR SUBDOMAIN TAKEOVERS...%s\n", Gold, EndC)
	}

	results := make(map[string]string) // subdomain -> vulnerable service
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 50) // Max 50 concurrent lookups

	for _, sub := range subdomains {
		wg.Add(1)
		sem <- struct{}{}
		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			cname, err := net.LookupCNAME(s)
			if err != nil {
				return
			}
			
			// net.LookupCNAME returns the FQDN with a trailing dot
			cname = strings.TrimSuffix(cname, ".")
			
			for signature, service := range TakeoverSignatures {
				if strings.Contains(cname, signature) {
					mu.Lock()
					results[s] = fmt.Sprintf("Vulnerable (%s -> %s)", service, cname)
					mu.Unlock()
					break
				}
			}
		}(sub)
	}

	wg.Wait()

	if !silent {
		if len(results) > 0 {
			fmt.Printf(" %s[!] FOUND %d POTENTIAL TAKEOVERS!%s\n", Coral, len(results), EndC)
		} else {
			fmt.Printf(" %s[✓] NO TAKEOVERS FOUND.%s\n", Mint, EndC)
		}
	}

	return results
}
