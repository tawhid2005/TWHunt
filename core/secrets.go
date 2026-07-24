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

// Common regexes for high-value secrets
var secretPatterns = map[string]*regexp.Regexp{
	"AWS Access Key":  regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
	"Stripe API Key":  regexp.MustCompile(`(?i)(sk_live_[0-9a-zA-Z]{24})`),
	"Google API Key":  regexp.MustCompile(`AIza[0-9A-Za-z-_]{35}`),
	"Slack Token":     regexp.MustCompile(`xox[baprs]-[0-9a-zA-Z]{10,48}`),
	"RSA Private Key": regexp.MustCompile(`-----BEGIN RSA PRIVATE KEY-----`),
	"GitHub Token":    regexp.MustCompile(`gh[pso]_[0-9a-zA-Z]{36}`),
}

// jsRegex to find script src attributes
var jsRegex = regexp.MustCompile(`(?i)<script[^>]+src=["'](.*?\.js)["']`)

var jsClient = &http.Client{
	Timeout: 7 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// FindJSSecrets fetches HTML, extracts JS links, and scans them for secrets
func FindJSSecrets(subdomains []string, silent bool) map[string][]string {
	if !silent {
		fmt.Printf(" %s[*] SCANNING JS FILES FOR HARDCODED SECRETS...%s\n", Gold, EndC)
	}

	results := make(map[string][]string) // subdomain -> list of found secrets
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 30) // Concurrent limit

	for _, sub := range subdomains {
		wg.Add(1)
		sem <- struct{}{}
		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			url := "http://" + s
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", "TWHunt-Scanner/1.0")
			
			resp, err := jsClient.Do(req)
			if err != nil {
				// try https
				url = "https://" + s
				req, _ = http.NewRequest("GET", url, nil)
				resp, err = jsClient.Do(req)
				if err != nil {
					return
				}
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			bodyStr := string(body)

			// 1. Scan inline HTML
			var localFound []string
			for name, regex := range secretPatterns {
				if regex.MatchString(bodyStr) {
					localFound = append(localFound, fmt.Sprintf("%s (in HTML)", name))
				}
			}

			// 2. Extract JS files and scan them
			matches := jsRegex.FindAllStringSubmatch(bodyStr, -1)
			for _, match := range matches {
				if len(match) > 1 {
					jsUrl := match[1]
					if !strings.HasPrefix(jsUrl, "http") {
						if strings.HasPrefix(jsUrl, "/") {
							jsUrl = url + jsUrl
						} else {
							jsUrl = url + "/" + jsUrl
						}
					}
					
					// Fetch JS
					jsReq, _ := http.NewRequest("GET", jsUrl, nil)
					jsResp, err := jsClient.Do(jsReq)
					if err == nil {
						jsBody, _ := ioutil.ReadAll(jsResp.Body)
						jsResp.Body.Close()
						jsStr := string(jsBody)
						
						for name, regex := range secretPatterns {
							if regex.MatchString(jsStr) {
								localFound = append(localFound, fmt.Sprintf("%s (in %s)", name, jsUrl))
							}
						}
					}
				}
			}

			if len(localFound) > 0 {
				mu.Lock()
				results[s] = localFound
				mu.Unlock()
			}
		}(sub)
	}

	wg.Wait()

	if !silent {
		if len(results) > 0 {
			fmt.Printf(" %s[!] CRITICAL: FOUND SECRETS ON %d SUBDOMAINS!%s\n", Coral, len(results), EndC)
		} else {
			fmt.Printf(" %s[✓] NO LEAKED SECRETS FOUND.%s\n", Mint, EndC)
		}
	}

	return results
}
