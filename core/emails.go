package core

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var emailClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// ExtractEmails uses Regex to scrape employee and company emails from subdomains
func ExtractEmails(subdomains []string, domain string, silent bool) map[string][]string {
	if !silent {
		fmt.Printf(" %s[*] OSINT: SCRAPING EXPOSED EMAILS...%s\n", Gold, EndC)
	}

	results := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20)

	// Regex for standard emails
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	for _, sub := range subdomains {
		wg.Add(1)
		sem <- struct{}{}
		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			url := "https://" + s
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0")
			
			resp, err := emailClient.Do(req)
			if err != nil {
				url = "http://" + s
				req, _ = http.NewRequest("GET", url, nil)
				resp, err = emailClient.Do(req)
				if err != nil {
					return
				}
			}
			
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()

			foundEmails := emailRegex.FindAllString(string(body), -1)
			
			// Deduplicate and filter emails (Prefer emails ending with the target domain)
			uniqueEmails := make(map[string]bool)
			for _, e := range foundEmails {
				e = strings.ToLower(e)
				// Basic filter to avoid false positives (like png@2x or example@example.com)
				if !strings.Contains(e, "example") && !strings.Contains(e, "test") {
					uniqueEmails[e] = true
				}
			}

			if len(uniqueEmails) > 0 {
				var finalEmails []string
				for e := range uniqueEmails {
					// Highlight if the email matches the target domain
					if strings.HasSuffix(e, "@"+domain) || strings.HasSuffix(e, "@"+strings.Split(domain, ".")[0]+".") {
						finalEmails = append(finalEmails, "Target: "+e)
					} else {
						finalEmails = append(finalEmails, e)
					}
				}
				
				mu.Lock()
				results[s] = finalEmails
				mu.Unlock()
			}

		}(sub)
	}

	wg.Wait()

	if !silent {
		totalFound := 0
		for _, v := range results {
			totalFound += len(v)
		}
		if totalFound > 0 {
			fmt.Printf(" %s[!] FOUND %d EXPOSED EMAILS!%s\n", Coral, totalFound, EndC)
		} else {
			fmt.Printf(" %s[✓] NO EXPOSED EMAILS FOUND.%s\n", Mint, EndC)
		}
	}
	return results
}
