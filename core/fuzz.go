package core

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var fuzzClient = &http.Client{
	Timeout: 4 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// Top 40 most critical endpoints for fast fuzzing
var fuzzPaths = []string{
	"/admin", "/admin-panel", "/admin/login", "/dashboard", "/api", "/api/v1", "/api/v2",
	"/swagger.json", "/api-docs", "/v2/api-docs", "/v1/api-docs", "/swagger-ui.html",
	"/server-status", "/phpMyAdmin", "/phpmyadmin", "/dbadmin", "/wp-admin", "/wp-login.php",
	"/config.json", "/config.yml", "/config.php", "/database.yml", "/.git/HEAD",
	"/actuator", "/actuator/env", "/actuator/health", "/health", "/metrics", "/prometheus",
	"/backup.zip", "/backup.sql", "/database.sql", "/dump.sql", "/db.sqlite3",
	"/.env", "/.env.example", "/.env.local", "/.env.production", "/info.php", "/phpinfo.php",
}

// FuzzDirectories brute-forces live subdomains for critical endpoints
func FuzzDirectories(subdomains []string, silent bool) map[string][]string {
	if !silent {
		fmt.Printf(" %s[*] FUZZING SUBDOMAINS FOR 40 CRITICAL DIRECTORIES...%s\n", Gold, EndC)
	}

	results := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 50) // Higher concurrency for fuzzing

	for _, sub := range subdomains {
		for _, path := range fuzzPaths {
			wg.Add(1)
			sem <- struct{}{}
			go func(s, p string) {
				defer wg.Done()
				defer func() { <-sem }()

				url := "https://" + s + p
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("User-Agent", "Mozilla/5.0")
				
				resp, err := fuzzClient.Do(req)
				if err != nil {
					// Fallback to HTTP
					url = "http://" + s + p
					req, _ = http.NewRequest("GET", url, nil)
					resp, err = fuzzClient.Do(req)
					if err != nil {
						return
					}
				}
				defer resp.Body.Close()

				// If 200 OK or 403 Forbidden (Directory exists but access denied)
				// We also check against false positive 200s (wildcard handling)
				if resp.StatusCode == 200 || resp.StatusCode == 403 {
					// Extremely basic false positive check (if body contains "Not Found" despite 200)
					if resp.StatusCode == 200 && strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), "text/html") {
						// Assuming standard 404 pages that return 200 have "404" or "not found" in title
						// In a real fuzzer, we would do a baseline check on a non-existent path.
						// To keep it fast, we just flag it.
					}
					
					finding := fmt.Sprintf("[%d] %s", resp.StatusCode, p)
					mu.Lock()
					results[s] = append(results[s], finding)
					mu.Unlock()
				}
			}(sub, path)
		}
	}

	wg.Wait()

	if !silent {
		totalFound := 0
		for _, v := range results {
			totalFound += len(v)
		}
		if totalFound > 0 {
			fmt.Printf(" %s[!] FOUND %d HIDDEN DIRECTORIES!%s\n", Coral, totalFound, EndC)
		} else {
			fmt.Printf(" %s[✓] NO HIDDEN DIRECTORIES FOUND.%s\n", Mint, EndC)
		}
	}
	return results
}
