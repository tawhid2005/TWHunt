package core

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var corsClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// CheckCORS sends a fake origin and checks if the server reflects it
func CheckCORS(subdomains []string, silent bool) map[string]string {
	if !silent {
		fmt.Printf(" %s[*] CHECKING FOR CORS MISCONFIGURATIONS...%s\n", Gold, EndC)
	}

	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 40)
	
	evilOrigin := "https://evil.com"

	for _, sub := range subdomains {
		wg.Add(1)
		sem <- struct{}{}
		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			req, _ := http.NewRequest("GET", "http://"+s, nil)
			req.Header.Set("Origin", evilOrigin)
			req.Header.Set("User-Agent", "TWHunt-Scanner/1.0")

			resp, err := corsClient.Do(req)
			if err != nil {
				req, _ = http.NewRequest("GET", "https://"+s, nil)
				req.Header.Set("Origin", evilOrigin)
				resp, err = corsClient.Do(req)
				if err != nil {
					return
				}
			}
			defer resp.Body.Close()

			allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
			allowCreds := resp.Header.Get("Access-Control-Allow-Credentials")

			if allowOrigin == evilOrigin {
				finding := "Reflects arbitrary Origin"
				if allowCreds == "true" {
					finding = "CRITICAL: Reflects Origin + Allows Credentials"
				}
				mu.Lock()
				results[s] = finding
				mu.Unlock()
			} else if allowOrigin == "*" && allowCreds == "true" {
				// While browsers block this, it's still a misconfig flag
				mu.Lock()
				results[s] = "Wildcard (*) with Credentials"
				mu.Unlock()
			}
		}(sub)
	}

	wg.Wait()

	if !silent {
		if len(results) > 0 {
			fmt.Printf(" %s[!] FOUND %d CORS MISCONFIGURATIONS!%s\n", Coral, len(results), EndC)
		} else {
			fmt.Printf(" %s[✓] NO CORS MISCONFIGURATIONS FOUND.%s\n", Mint, EndC)
		}
	}
	return results
}
