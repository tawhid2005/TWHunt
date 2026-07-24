package core

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var bucketPermutations = []string{
	"%s", "%s-dev", "%s-prod", "%s-staging", "%s-backup", "%s-assets", 
	"%s-public", "%s-private", "%s-static", "%s-media", "%s-data",
}

var bucketProviders = map[string]string{
	"AWS S3":      "http://%s.s3.amazonaws.com",
	"GCP Bucket":  "http://storage.googleapis.com/%s",
	"Azure Blob":  "http://%s.blob.core.windows.net",
}

var bucketClient = &http.Client{
	Timeout: 5 * time.Second,
}

// EnumerateBuckets generates potential bucket names and checks if they are open
func EnumerateBuckets(domain string, silent bool) []string {
	if !silent {
		fmt.Printf(" %s[*] ENUMERATING CLOUD BUCKETS FOR %s...%s\n", Gold, strings.ToUpper(domain), EndC)
	}

	baseName := strings.Split(domain, ".")[0] // e.g. hackerone
	var results []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 30)

	for _, perm := range bucketPermutations {
		bucketName := fmt.Sprintf(perm, baseName)
		
		for providerName, providerUrlFormat := range bucketProviders {
			url := fmt.Sprintf(providerUrlFormat, bucketName)
			
			wg.Add(1)
			sem <- struct{}{}
			go func(u, pName, bName string) {
				defer wg.Done()
				defer func() { <-sem }()

				req, _ := http.NewRequest("GET", u, nil)
				req.Header.Set("User-Agent", "TWHunt-Scanner/1.0")
				
				resp, err := bucketClient.Do(req)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				// 200 OK means bucket is completely public and listable!
				if resp.StatusCode == 200 {
					// Extra verification for AWS to avoid false positive generic pages
					if pName == "AWS S3" && !strings.Contains(resp.Header.Get("Server"), "AmazonS3") {
						return
					}
					
					finding := fmt.Sprintf("[%s] OPEN BUCKET: %s", pName, u)
					mu.Lock()
					results = append(results, finding)
					mu.Unlock()
				}
			}(url, providerName, bucketName)
		}
	}

	wg.Wait()

	if !silent {
		if len(results) > 0 {
			fmt.Printf(" %s[!] FOUND %d OPEN CLOUD BUCKETS!%s\n", Coral, len(results), EndC)
		} else {
			fmt.Printf(" %s[✓] NO OPEN BUCKETS FOUND.%s\n", Mint, EndC)
		}
	}
	
	return results
}
