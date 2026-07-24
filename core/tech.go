package core

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// DetectTechStack analyzes headers to find technologies
func DetectTechStack(subdomains []string, silent bool) map[string][]string {
	if !silent {
		fmt.Printf(" %s[*] DETECTING TECHNOLOGY STACK...%s\n", Gold, EndC)
	}

	results := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 40)

	for _, sub := range subdomains {
		wg.Add(1)
		sem <- struct{}{}
		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			req, _ := http.NewRequest("HEAD", "http://"+s, nil)
			req.Header.Set("User-Agent", "TWHunt-Scanner/1.0")
			
			resp, err := jsClient.Do(req)
			if err != nil {
				req, _ = http.NewRequest("HEAD", "https://"+s, nil)
				resp, err = jsClient.Do(req)
				if err != nil {
					return
				}
			}
			defer resp.Body.Close()

			var tech []string

			// Header Checks
			if server := resp.Header.Get("Server"); server != "" {
				tech = append(tech, "Server: "+server)
			}
			if poweredBy := resp.Header.Get("X-Powered-By"); poweredBy != "" {
				tech = append(tech, "PoweredBy: "+poweredBy)
			}
			if asp := resp.Header.Get("X-AspNet-Version"); asp != "" {
				tech = append(tech, "ASP.NET")
			}
			
			// Simple Cookie Checks
			cookies := resp.Header.Get("Set-Cookie")
			if strings.Contains(cookies, "PHPSESSID") {
				tech = append(tech, "PHP")
			}
			if strings.Contains(cookies, "JSESSIONID") {
				tech = append(tech, "Java/Tomcat")
			}

			if len(tech) > 0 {
				mu.Lock()
				results[s] = tech
				mu.Unlock()
			}
		}(sub)
	}

	wg.Wait()

	if !silent {
		fmt.Printf(" %s[✓] TECH STACK DETECTION COMPLETE.%s\n", Mint, EndC)
	}
	return results
}
