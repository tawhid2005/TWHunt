package core

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var blhClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// Regex to find social media links in HTML
var socialRegex = regexp.MustCompile(`https?://(?:www\.)?(twitter\.com|x\.com|instagram\.com|youtube\.com|linkedin\.com|github\.com|facebook\.com|t\.me)/[a-zA-Z0-9_-]+`)

// DetectBLH scans HTML for broken social media links (Broken Link Hijacking)
func DetectBLH(subdomains []string, silent bool) map[string][]string {
	if !silent {
		fmt.Printf(" %s[*] SCANNING FOR BROKEN LINK HIJACKING (BLH)...%s\n", Gold, EndC)
	}

	results := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20)

	for _, sub := range subdomains {
		wg.Add(1)
		sem <- struct{}{}
		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			url := "https://" + s
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0")
			
			resp, err := blhClient.Do(req)
			if err != nil {
				url = "http://" + s
				req, _ = http.NewRequest("GET", url, nil)
				resp, err = blhClient.Do(req)
				if err != nil {
					return
				}
			}
			
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()

			links := socialRegex.FindAllString(string(body), -1)
			
			// Deduplicate links
			uniqueLinks := make(map[string]bool)
			for _, link := range links {
				uniqueLinks[link] = true
			}

			// Check each link to see if it's broken (404)
			for link := range uniqueLinks {
				checkReq, _ := http.NewRequest("GET", link, nil)
				checkReq.Header.Set("User-Agent", "Mozilla/5.0")
				checkResp, checkErr := blhClient.Do(checkReq)
				
				if checkErr == nil {
					if checkResp.StatusCode == 404 {
						finding := fmt.Sprintf("Broken Link: %s", link)
						mu.Lock()
						results[s] = append(results[s], finding)
						mu.Unlock()
					}
					checkResp.Body.Close()
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
			fmt.Printf(" %s[!] FOUND %d BROKEN SOCIAL LINKS!%s\n", Coral, totalFound, EndC)
		} else {
			fmt.Printf(" %s[✓] NO BROKEN LINKS FOUND.%s\n", Mint, EndC)
		}
	}
	return results
}
