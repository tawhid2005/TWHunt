package core

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var screenClient = &http.Client{
	Timeout: 10 * time.Second, // Screenshots can take time
}

// TakeScreenshots uses a public API (thum.io) to fetch screenshots of live subdomains without needing Headless Chrome
func TakeScreenshots(subdomains []string, silent bool) map[string]string {
	if !silent {
		fmt.Printf(" %s[*] CAPTURING SCREENSHOTS OF LIVE SUBDOMAINS...%s\n", Gold, EndC)
	}

	results := make(map[string]string)
	
	// Create output directory
	outDir := "twhunt_screenshots"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		if !silent {
			fmt.Printf(" %s[-] Could not create screenshot directory%s\n", Silver, EndC)
		}
		return results
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 15) // Lower concurrency to avoid api limits

	for _, sub := range subdomains {
		wg.Add(1)
		sem <- struct{}{}
		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			// Ensure we are screenshotting HTTPs first
			targetUrl := "https://" + s
			// Using thum.io public API for fast zero-dependency screenshots
			apiUrl := fmt.Sprintf("https://image.thum.io/get/width/1024/crop/800/%s", targetUrl)
			
			req, _ := http.NewRequest("GET", apiUrl, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0")
			
			resp, err := screenClient.Do(req)
			if err != nil || resp.StatusCode != 200 {
				// Fallback to HTTP
				targetUrl = "http://" + s
				apiUrl = fmt.Sprintf("https://image.thum.io/get/width/1024/crop/800/%s", targetUrl)
				req, _ = http.NewRequest("GET", apiUrl, nil)
				resp, err = screenClient.Do(req)
				if err != nil || resp.StatusCode != 200 {
					return
				}
			}
			defer resp.Body.Close()

			// Check if it's an image
			if !strings.Contains(resp.Header.Get("Content-Type"), "image") {
				return
			}

			// Generate a unique filename based on the subdomain
			hash := md5.Sum([]byte(s))
			filename := filepath.Join(outDir, fmt.Sprintf("%s_%x.jpg", strings.ReplaceAll(s, ".", "_"), hash[:4]))
			
			file, err := os.Create(filename)
			if err != nil {
				return
			}
			defer file.Close()
			
			_, err = io.Copy(file, resp.Body)
			if err == nil {
				mu.Lock()
				results[s] = filename
				mu.Unlock()
			}
		}(sub)
	}

	wg.Wait()

	if !silent {
		if len(results) > 0 {
			fmt.Printf(" %s[!] SAVED %d SCREENSHOTS IN '%s'!%s\n", Coral, len(results), outDir, EndC)
		} else {
			fmt.Printf(" %s[✓] COULD NOT CAPTURE ANY SCREENSHOTS.%s\n", Mint, EndC)
		}
	}
	return results
}
