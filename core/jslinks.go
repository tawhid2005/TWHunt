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

var jsLinkClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// ExtractJSLinks acts as a LinkFinder to extract hidden API endpoints from JS files
func ExtractJSLinks(subdomains []string, silent bool) map[string][]string {
	if !silent {
		fmt.Printf(" %s[*] SEARCHING FOR HIDDEN API ENDPOINTS IN JS FILES...%s\n", Gold, EndC)
	}

	results := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20)

	// Regex to find endpoints like /api/v1/users, api/data, etc. in JS
	linkRegex := regexp.MustCompile(`(?i)(?:\"|\')(((?:[a-zA-Z]{1,10}://|/)[^\"\'\s]+|([a-zA-Z0-9_\-]+/)+[a-zA-Z0-9_\-]+))(?:\?|/|\"|\')`)

	for _, sub := range subdomains {
		wg.Add(1)
		sem <- struct{}{}
		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			// Simple check on the root page for JS links, a full implementation would crawl first
			url := "https://" + s
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0")
			
			resp, err := jsLinkClient.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			
			body, _ := ioutil.ReadAll(resp.Body)
			
			// Find JS files in HTML
			jsRegex := regexp.MustCompile(`src=["']([^"']+\.js)["']`)
			jsFiles := jsRegex.FindAllStringSubmatch(string(body), -1)

			for _, match := range jsFiles {
				if len(match) > 1 {
					jsUrl := match[1]
					if !strings.HasPrefix(jsUrl, "http") {
						if strings.HasPrefix(jsUrl, "/") {
							jsUrl = "https://" + s + jsUrl
						} else {
							jsUrl = "https://" + s + "/" + jsUrl
						}
					}

					// Download JS file
					jsReq, _ := http.NewRequest("GET", jsUrl, nil)
					jsReq.Header.Set("User-Agent", "Mozilla/5.0")
					jsResp, jsErr := jsLinkClient.Do(jsReq)
					if jsErr == nil {
						jsBody, _ := ioutil.ReadAll(jsResp.Body)
						jsResp.Body.Close()

						links := linkRegex.FindAllString(string(jsBody), -1)
						
						uniqueLinks := make(map[string]bool)
						for _, l := range links {
							cleanLink := strings.Trim(l, "\"'")
							if strings.HasPrefix(cleanLink, "/api") || strings.HasPrefix(cleanLink, "api/") || strings.HasPrefix(cleanLink, "/v1") {
								uniqueLinks[cleanLink] = true
							}
						}

						if len(uniqueLinks) > 0 {
							var foundLinks []string
							for l := range uniqueLinks {
								foundLinks = append(foundLinks, l)
							}
							mu.Lock()
							results[s] = append(results[s], fmt.Sprintf("JS: %s -> %v", match[1], foundLinks))
							mu.Unlock()
						}
					}
				}
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
			fmt.Printf(" %s[!] FOUND %d HIDDEN ENDPOINTS IN JS!%s\n", Coral, totalFound, EndC)
		} else {
			fmt.Printf(" %s[✓] NO HIDDEN ENDPOINTS FOUND IN JS.%s\n", Mint, EndC)
		}
	}
	return results
}
