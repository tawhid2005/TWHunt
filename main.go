package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	"twhunt/core"
	"twhunt/sources"
	"twhunt/utils"
)

func main() {
	// Custom help menu
	flag.Usage = func() {
		fmt.Printf("%s%s\n", core.Sky, core.Bold)
		fmt.Println(`
  _______          __  _    _             _   
 |__   __|        / / | |  | |           | |  
    | | __      __ /  | |__| | _   _  _ __ | |_ 
    | | \ \ /\ / /    |  __  || | | || '_ \| __|
    | |  \ V  V /     | |  | || |_| || | | | |_ 
    |_|   \_/\_/      |_|  |_| \__,_||_| |_|\__|
                                                  `)
		fmt.Printf("%s   An Advanced & Professional Subdomain Hunter%s\n", core.Lavender, core.EndC)
		fmt.Printf("%s   -------------------------------------------------------%s\n", core.Slate, core.EndC)
		fmt.Printf("%s   👤 Author  :%s MD TALHA HUSSAIN TAWHID%s\n", core.Sky, core.Mint, core.EndC)
		fmt.Printf("%s   📧 Email   :%s tawhidh2005@gmail.com%s\n", core.Sky, core.Mint, core.EndC)
		fmt.Printf("%s   -------------------------------------------------------%s\n\n", core.Slate, core.EndC)
		
		fmt.Printf(" %s[DESCRIPTION]%s\n", core.Gold, core.EndC)
		fmt.Printf("   TWHunt is a fast, modular, and passive subdomain enumeration tool.\n")
		fmt.Printf("   It gathers subdomains from 24 free OSINT sources without API keys.\n\n")

		fmt.Printf(" %s[USAGE]%s\n", core.Gold, core.EndC)
		fmt.Printf("   twhunt -d <domain.com> [flags]\n\n")

		fmt.Printf(" %s[BASIC FLAGS]%s\n", core.Gold, core.EndC)
		fmt.Printf("   %s-d%s      : Target domain (e.g. hackerone.com)\n", core.Mint, core.EndC)
		fmt.Printf("   %s-dL%s     : Target domains list file (e.g. domains.txt)\n", core.Mint, core.EndC)
		fmt.Printf("   %s-v%s      : Verify live status (resolves DNS to find alive subdomains)\n", core.Mint, core.EndC)
		
		fmt.Printf("\n %s[ADVANCED FLAGS]%s\n", core.Gold, core.EndC)
		fmt.Printf("   %s-nw%s     : Enable wildcard DNS filtering (removes false positives)\n", core.Mint, core.EndC)
		fmt.Printf("   %s-w%s      : Active brute-forcing wordlist file (e.g. subdomains.txt)\n", core.Mint, core.EndC)
		fmt.Printf("   %s-ports%s  : TCP Ports to scan on live subdomains (e.g. 80,443,8080)\n", core.Mint, core.EndC)
		
		fmt.Printf("\n %s[OUTPUT FLAGS]%s\n", core.Gold, core.EndC)
		fmt.Printf("   %s-o%s      : Output file to save results (e.g. results.txt)\n", core.Mint, core.EndC)
		fmt.Printf("   %s-json%s   : JSON output format (prints subdomains as a JSON array)\n", core.Mint, core.EndC)
		fmt.Printf("   %s-html%s   : Generate a beautiful HTML report dashboard\n", core.Mint, core.EndC)
		fmt.Printf("   %s-silent%s : Silent mode (no banners/logs, pipes nicely into httpx)\n", core.Mint, core.EndC)
		
		fmt.Printf("\n %s[MISC FLAGS]%s\n", core.Gold, core.EndC)
		fmt.Printf("   %s-update%s : Auto-updates TWHunt to the latest version from GitHub\n", core.Mint, core.EndC)
		fmt.Printf("   %s-h%s      : Show this help menu\n\n", core.Mint, core.EndC)

		fmt.Printf(" %s[EXAMPLES]%s\n", core.Gold, core.EndC)
		fmt.Printf("   %stwhunt -d example.com -nw -v%s\n", core.Sky, core.EndC)
		fmt.Printf("   %stwhunt -d example.com -w wordlist.txt -ports 80,443 -html%s\n", core.Sky, core.EndC)
		fmt.Printf("   %stwhunt -dL domains.txt -silent | httpx%s\n\n", core.Sky, core.EndC)
	}

	domainPtr := flag.String("d", "", "TARGET DOMAIN")
	verifyPtr := flag.Bool("v", false, "VERIFY LIVE STATUS")
	outputPtr := flag.String("o", "", "OUTPUT FILE TO SAVE RESULTS")
	silentPtr := flag.Bool("silent", false, "SILENT MODE")
	jsonPtr := flag.Bool("json", false, "JSON OUTPUT")
	listPtr := flag.String("dL", "", "TARGET DOMAINS LIST FILE")
	updatePtr := flag.Bool("update", false, "UPDATE TWHUNT")
	
	nwPtr := flag.Bool("nw", false, "FILTER WILDCARD DNS")
	wordlistPtr := flag.String("w", "", "WORDLIST FOR BRUTEFORCING")
	portsPtr := flag.String("ports", "", "PORTS TO SCAN (e.g. 80,443)")
	htmlPtr := flag.Bool("html", false, "GENERATE HTML REPORT")
	
	flag.Parse()

	if *updatePtr {
		utils.AutoUpdate()
	}

	if *domainPtr == "" && *listPtr == "" {
		if !*silentPtr {
			flag.Usage()
		}
		os.Exit(1)
	}

	if !*silentPtr && !*jsonPtr {
		fmt.Printf("%s%s\n", core.Sky, core.Bold)
		fmt.Println(`
  _______          __  _    _             _   
 |__   __|        / / | |  | |           | |  
    | | __      __ /  | |__| | _   _  _ __ | |_ 
    | | \ \ /\ / /    |  __  || | | || '_ \| __|
    | |  \ V  V /     | |  | || |_| || | | | |_ 
    |_|   \_/\_/      |_|  |_| \__,_||_| |_|\__|
                                                  `)
		fmt.Printf("%s   An Advanced & Professional Subdomain Hunter%s\n", core.Lavender, core.EndC)
		fmt.Printf("%s   -------------------------------------------------------%s\n", core.Slate, core.EndC)
		fmt.Printf("%s   👤 Author  :%s MD TALHA HUSSAIN TAWHID%s\n", core.Sky, core.Mint, core.EndC)
		fmt.Printf("%s   -------------------------------------------------------%s\n\n", core.Slate, core.EndC)
	}

	// Load Config File (creates if missing)
	core.LoadConfig(*silentPtr)

	var targets []string
	if *domainPtr != "" {
		targets = append(targets, strings.ToLower(strings.TrimSpace(*domainPtr)))
	}
	if *listPtr != "" {
		content, err := os.ReadFile(*listPtr)
		if err != nil {
			if !*silentPtr {
				fmt.Printf(" %s[!] Error reading domain list: %v%s\n", core.Coral, err, core.EndC)
			}
			os.Exit(1)
		}
		for _, line := range strings.Split(string(content), "\n") {
			line = strings.ToLower(strings.TrimSpace(line))
			if line != "" {
				targets = append(targets, line)
			}
		}
	}

	sourceList := []core.Source{
		&sources.AbuseIPDBSource{},
		&sources.AlienVaultSource{},
		&sources.AnubisSource{},
		&sources.BeVigilSource{},
		&sources.BufferOverSource{},
		&sources.CertSpotterSource{},
		&sources.CommonCrawlSource{},
		&sources.CrtshSource{},
		&sources.FullHuntSource{},
		&sources.HackerTargetSource{},
		&sources.NetcraftSource{},
		&sources.OmnisintSource{},
		&sources.RapidDNSSource{},
		&sources.RiddlerSource{},
		&sources.ShodanCTSource{},
		&sources.ShrewdEyeSource{},
		&sources.SiteDossierSource{},
		&sources.SubdomainCenterSource{},
		&sources.SynapsintSource{},
		&sources.ThreatCrowdSource{},
		&sources.ThreatMinerSource{},
		&sources.URLScanSource{},
		&sources.VirusTotalSource{},
		&sources.WaybackSource{},
	}

	var allDiscoveredSubdomains []string

	for _, target := range targets {
		var wildcardIPs []string
		if *nwPtr {
			wildcardIPs = core.GetWildcardIPs(target, *silentPtr)
		}

		engine := core.NewEngine(sourceList, *silentPtr)
		subdomains := engine.Run(target)

		if *wordlistPtr != "" {
			bruted := core.RunBruteforce(target, *wordlistPtr, 50, *silentPtr)
			subdomains = append(subdomains, bruted...)
		}

		// Deduplicate current domain
		uniqueMap := make(map[string]bool)
		var uniqueSubdomains []string
		for _, sub := range subdomains {
			if !uniqueMap[sub] {
				uniqueMap[sub] = true
				uniqueSubdomains = append(uniqueSubdomains, sub)
			}
		}
		subdomains = uniqueSubdomains

		// Verify Live Status or Wildcard Filter
		if *verifyPtr || *nwPtr || *portsPtr != "" {
			if !*silentPtr {
				fmt.Printf(" %s[*] VERIFYING LIVE STATUS & FILTERING...%s\n", core.Gold, core.EndC)
			}
			var live []string
			for _, sub := range subdomains {
				if _, err := net.LookupHost(sub); err == nil {
					// Check wildcard
					if *nwPtr && core.IsWildcard(sub, wildcardIPs) {
						continue // Skip if wildcard fake
					}
					live = append(live, sub)
				}
			}
			subdomains = live
			if !*silentPtr {
				fmt.Printf(" %s[✓] TOTAL LIVE/VALID HOSTS FOUND FOR %s: %d%s\n", core.Mint, target, len(subdomains), core.EndC)
			}
		}

		allDiscoveredSubdomains = append(allDiscoveredSubdomains, subdomains...)
	}

	// Remove completely global duplicates across all domains
	uniqueGlobal := make(map[string]bool)
	var finalSubdomains []string
	for _, sub := range allDiscoveredSubdomains {
		if !uniqueGlobal[sub] {
			uniqueGlobal[sub] = true
			finalSubdomains = append(finalSubdomains, sub)
		}
	}
	sort.Strings(finalSubdomains)

	// Port Scanning
	var portResults map[string][]string
	if *portsPtr != "" && len(finalSubdomains) > 0 {
		portResults = core.ScanPorts(finalSubdomains, *portsPtr, *silentPtr)
	}

	// Generate HTML Report
	if *htmlPtr {
		htmlName := "twhunt_report.html"
		if *outputPtr != "" && strings.HasSuffix(*outputPtr, ".html") {
			htmlName = *outputPtr
		}
		utils.GenerateHTMLReport(htmlName, finalSubdomains, portResults, *silentPtr)
	}

	// Console Output
	if *jsonPtr {
		jsonData, _ := json.MarshalIndent(finalSubdomains, "", "  ")
		fmt.Println(string(jsonData))
	} else if !*silentPtr && len(finalSubdomains) > 0 {
		fmt.Printf("\n%s%s========================= RESULTS =========================%s\n", core.Sky, core.Bold, core.EndC)
		for _, sub := range finalSubdomains {
			if ports, ok := portResults[sub]; ok && len(ports) > 0 {
				fmt.Printf(" %s→%s %s%-35s %s[PORTS: %s]%s\n", core.Mint, core.EndC, core.Silver, sub, core.Coral, strings.Join(ports, ","), core.EndC)
			} else {
				fmt.Printf(" %s→%s %s%s%s\n", core.Mint, core.EndC, core.Silver, sub, core.EndC)
			}
		}
		fmt.Printf("%s%s===========================================================%s\n", core.Sky, core.Bold, core.EndC)
	} else if *silentPtr && !*htmlPtr && !*jsonPtr {
		// In silent mode without JSON/HTML, just print one domain per line
		for _, sub := range finalSubdomains {
			if ports, ok := portResults[sub]; ok && len(ports) > 0 {
				fmt.Printf("%s:%s\n", sub, strings.Join(ports, ","))
			} else {
				fmt.Println(sub)
			}
		}
	}

	if !*silentPtr && !*jsonPtr {
		fmt.Printf("\n %s[!] SCAN COMPLETED. HAPPY HUNTING!%s\n", core.Sky, core.EndC)
	}

	// Text Save Logic
	if *outputPtr != "" && len(finalSubdomains) > 0 && !*htmlPtr {
		saveToFile(*outputPtr, finalSubdomains, portResults, *jsonPtr, *silentPtr)
	} else if len(finalSubdomains) > 0 && !*silentPtr && !*jsonPtr && !*htmlPtr {
		fmt.Printf("\n %s[?] Do you want to save the results to a file? (y/n): %s", core.Gold, core.EndC)
		var choice string
		fmt.Scanln(&choice)
		choice = strings.ToLower(strings.TrimSpace(choice))
		
		if choice == "y" || choice == "yes" {
			fmt.Printf(" %s[>] Enter filename (e.g. results.txt): %s", core.Mint, core.EndC)
			var filename string
			fmt.Scanln(&filename)
			filename = strings.TrimSpace(filename)
			
			if filename != "" {
				saveToFile(filename, finalSubdomains, portResults, false, *silentPtr)
			}
		}
	}
}

func saveToFile(filename string, data []string, portResults map[string][]string, isJson bool, isSilent bool) {
	file, err := os.Create(filename)
	if err != nil {
		if !isSilent {
			fmt.Printf(" %s[!] ERROR SAVING FILE: %v%s\n", core.Coral, err, core.EndC)
		}
		return
	}
	defer file.Close()
	
	if isJson {
		// For JSON we could construct a structured object with ports, but keeping it simple for now
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		file.Write(jsonData)
	} else {
		for _, sub := range data {
			if ports, ok := portResults[sub]; ok && len(ports) > 0 {
				file.WriteString(sub + " [PORTS: " + strings.Join(ports, ",") + "]\n")
			} else {
				file.WriteString(sub + "\n")
			}
		}
	}
	
	if !isSilent {
		fmt.Printf(" %s[✓] RESULTS SAVED TO: %s%s\n", core.Mint, filename, core.EndC)
	}
}
