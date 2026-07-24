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

type ProbeResult struct {
	Subdomain  string
	StatusCode int
	Title      string
	URL        string
}

// HTTPClient for probing (short timeout, ignores SSL errors)
var probeClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// titleRegex extracts the title tag content
var titleRegex = regexp.MustCompile(`(?i)<title>(.*?)</title>`)

// ProbeHTTP checks HTTP and HTTPS for subdomains and extracts status codes and titles
func ProbeHTTP(subdomains []string, silent bool) map[string]ProbeResult {
	if !silent {
		fmt.Printf(" %s[*] PROBING HTTP/HTTPS STATUS CODES & TITLES...%s\n", Gold, EndC)
	}

	results := make(map[string]ProbeResult)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 40) // 40 concurrent HTTP requests

	for _, sub := range subdomains {
		wg.Add(1)
		sem <- struct{}{}
		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			// Try HTTPS first, then fallback to HTTP
			urlsToTry := []string{"https://" + s, "http://" + s}
			var success bool

			for _, u := range urlsToTry {
				req, err := http.NewRequest("GET", u, nil)
				if err != nil {
					continue
				}
				req.Header.Set("User-Agent", "TWHunt-Prober/1.0")

				resp, err := probeClient.Do(req)
				if err != nil {
					continue
				}
				
				bodyBytes, _ := ioutil.ReadAll(resp.Body)
				resp.Body.Close()

				title := "No Title"
				matches := titleRegex.FindStringSubmatch(string(bodyBytes))
				if len(matches) > 1 {
					title = strings.TrimSpace(matches[1])
				}

				mu.Lock()
				results[s] = ProbeResult{
					Subdomain:  s,
					StatusCode: resp.StatusCode,
					Title:      title,
					URL:        u,
				}
				mu.Unlock()
				success = true
				break // Stop trying if one protocol succeeds
			}

			// If both fail, we could log it as dead, but usually we just skip adding to results
			if !success {
				mu.Lock()
				results[s] = ProbeResult{
					Subdomain:  s,
					StatusCode: 0,
					Title:      "Dead/Unreachable",
					URL:        "",
				}
				mu.Unlock()
			}
		}(sub)
	}

	wg.Wait()

	if !silent {
		fmt.Printf(" %s[✓] FINISHED HTTP PROBING FOR %d HOSTS.%s\n", Mint, len(subdomains), EndC)
	}

	return results
}
