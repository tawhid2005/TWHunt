package core

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var wafClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// Known WAF Headers and Server names
var wafSignatures = map[string]string{
	"cloudflare": "Cloudflare",
	"x-amz-cf":   "AWS CloudFront/WAF",
	"akamai":     "Akamai",
	"sucuri":     "Sucuri",
	"imperva":    "Imperva Incapsula",
	"f5":         "F5 BIG-IP",
}

// DetectWAF sends a harmless XSS payload to trigger and identify WAFs
func DetectWAF(subdomains []string, silent bool) map[string]string {
	if !silent {
		fmt.Printf(" %s[*] PROBING FOR WEB APPLICATION FIREWALLS (WAF)...%s\n", Gold, EndC)
	}

	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 40)
	
	payload := "/?id=<script>alert('twhunt')</script>"

	for _, sub := range subdomains {
		wg.Add(1)
		sem <- struct{}{}
		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			req, _ := http.NewRequest("GET", "http://"+s+payload, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)") // Standard UA to avoid immediate drop

			resp, err := wafClient.Do(req)
			if err != nil {
				req, _ = http.NewRequest("GET", "https://"+s+payload, nil)
				resp, err = wafClient.Do(req)
				if err != nil {
					return
				}
			}
			defer resp.Body.Close()

			serverHeader := strings.ToLower(resp.Header.Get("Server"))
			
			// 1. Check Server Header
			for sig, name := range wafSignatures {
				if strings.Contains(serverHeader, sig) {
					mu.Lock()
					results[s] = name
					mu.Unlock()
					return
				}
			}

			// 2. Check Custom Headers (e.g. x-amz-cf-id)
			for header := range resp.Header {
				headerLower := strings.ToLower(header)
				for sig, name := range wafSignatures {
					if strings.Contains(headerLower, sig) {
						mu.Lock()
						results[s] = name
						mu.Unlock()
						return
					}
				}
			}

			// 3. Fallback: If blocked (403, 406) but no obvious header, it's an Unknown WAF
			if resp.StatusCode == 403 || resp.StatusCode == 406 {
				mu.Lock()
				results[s] = "Unknown WAF (Blocked Payload)"
				mu.Unlock()
			}
		}(sub)
	}

	wg.Wait()

	if !silent {
		fmt.Printf(" %s[✓] WAF DETECTION COMPLETE.%s\n", Mint, EndC)
	}
	return results
}
