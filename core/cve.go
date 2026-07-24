package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var cveClient = &http.Client{
	Timeout: 5 * time.Second,
}

// CheckCVE uses a public vulnerability API to find CVEs for detected technologies
func CheckCVE(techResults map[string][]string, silent bool) map[string][]string {
	if !silent {
		fmt.Printf(" %s[*] CROSS-REFERENCING TECHNOLOGIES WITH CVE DATABASES...%s\n", Gold, EndC)
	}

	results := make(map[string][]string)
	
	// We only want to check a few unique technologies to avoid rate limits
	checkedTechs := make(map[string][]string)

	for sub, techs := range techResults {
		for _, tech := range techs {
			// Clean up tech name for API (e.g., "WordPress" -> "wordpress")
			cleanTech := strings.ToLower(strings.Split(tech, " ")[0])
			
			// Skip generic or extremely common ones that will return too many results
			if cleanTech == "html" || cleanTech == "css" || cleanTech == "utf-8" || cleanTech == "php" {
				continue
			}

			// If we haven't checked this technology yet
			if _, exists := checkedTechs[cleanTech]; !exists {
				// Using CIRCL public CVE API or similar basic search
				apiUrl := "https://cve.circl.lu/api/search/" + url.PathEscape(cleanTech)
				
				req, _ := http.NewRequest("GET", apiUrl, nil)
				req.Header.Set("User-Agent", "Mozilla/5.0")
				
				resp, err := cveClient.Do(req)
				if err == nil && resp.StatusCode == 200 {
					body, _ := ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					
					var cves []struct {
						ID      string  `json:"id"`
						Summary string  `json:"summary"`
						CVSS    float64 `json:"cvss"`
					}
					
					json.Unmarshal(body, &cves)
					
					// Store top 3 critical CVEs
					var foundCves []string
					count := 0
					for _, cve := range cves {
						if cve.CVSS >= 7.0 && count < 3 { // High or Critical
							foundCves = append(foundCves, fmt.Sprintf("%s (CVSS: %.1f)", cve.ID, cve.CVSS))
							count++
						}
					}
					checkedTechs[cleanTech] = foundCves
				} else {
					checkedTechs[cleanTech] = []string{}
				}
				
				// Small delay to prevent API bans
				time.Sleep(500 * time.Millisecond)
			}
			
			// If we found CVEs for this tech, add it to the subdomain results
			if len(checkedTechs[cleanTech]) > 0 {
				results[sub] = append(results[sub], fmt.Sprintf("%s -> %s", tech, strings.Join(checkedTechs[cleanTech], ", ")))
			}
		}
	}

	if !silent {
		totalFound := 0
		for _, v := range results {
			totalFound += len(v)
		}
		if totalFound > 0 {
			fmt.Printf(" %s[!] FOUND %d POTENTIAL CVE VULNERABILITIES!%s\n", Coral, totalFound, EndC)
		} else {
			fmt.Printf(" %s[✓] NO CRITICAL CVEs FOUND FOR DETECTED TECH.%s\n", Mint, EndC)
		}
	}
	return results
}
