package core

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var paramFuzzClient = &http.Client{
	Timeout: 4 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

var sensitiveParams = []string{
	"?admin=1", "?admin=true", "?debug=1", "?debug=true", 
	"?test=1", "?test=true", "?dev=1", "?dev=true",
}

// FuzzHiddenParams tests live subdomains for hidden debug/admin parameters
func FuzzHiddenParams(subdomains []string, silent bool) map[string]string {
	if !silent {
		fmt.Printf(" %s[*] FUZZING FOR HIDDEN SENSITIVE PARAMETERS...%s\n", Gold, EndC)
	}

	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 30)

	for _, sub := range subdomains {
		wg.Add(1)
		sem <- struct{}{}
		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			// 1. Get Baseline response size
			baseReq, _ := http.NewRequest("GET", "http://"+s, nil)
			baseReq.Header.Set("User-Agent", "Mozilla/5.0")
			baseResp, err := paramFuzzClient.Do(baseReq)
			
			if err != nil {
				baseReq, _ = http.NewRequest("GET", "https://"+s, nil)
				baseResp, err = paramFuzzClient.Do(baseReq)
				if err != nil {
					return
				}
			}
			
			baseBody, _ := ioutil.ReadAll(baseResp.Body)
			baseResp.Body.Close()
			baseLen := len(baseBody)
			baseStatus := baseResp.StatusCode

			// 2. Test parameters
			scheme := "http://"
			if strings.HasPrefix(baseReq.URL.String(), "https") {
				scheme = "https://"
			}

			for _, param := range sensitiveParams {
				req, _ := http.NewRequest("GET", scheme+s+param, nil)
				req.Header.Set("User-Agent", "Mozilla/5.0")
				
				resp, err := paramFuzzClient.Do(req)
				if err != nil {
					continue
				}

				body, _ := ioutil.ReadAll(resp.Body)
				resp.Body.Close()
				paramLen := len(body)
				paramStatus := resp.StatusCode

				// If response size changes significantly (> 100 bytes difference) or status code changes
				diff := baseLen - paramLen
				if diff < 0 { diff = -diff }

				if (paramStatus != baseStatus && paramStatus == 200) || (paramStatus == 200 && diff > 150) {
					finding := fmt.Sprintf("Response changed for %s (Status: %d, Size diff: %d bytes)", param, paramStatus, diff)
					mu.Lock()
					results[s] = finding
					mu.Unlock()
					return // Only need to flag it once per subdomain
				}
			}
		}(sub)
	}

	wg.Wait()

	if !silent {
		if len(results) > 0 {
			fmt.Printf(" %s[!] FOUND %d HIDDEN PARAMETER ANOMALIES!%s\n", Coral, len(results), EndC)
		} else {
			fmt.Printf(" %s[✓] NO HIDDEN PARAMETERS FOUND.%s\n", Mint, EndC)
		}
	}
	return results
}
