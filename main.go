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
	domainPtr := flag.String("d", "", "TARGET DOMAIN")
	verifyPtr := flag.Bool("v", false, "VERIFY LIVE STATUS")
	outputPtr := flag.String("o", "", "OUTPUT FILE TO SAVE RESULTS")
	silentPtr := flag.Bool("silent", false, "SILENT MODE")
	jsonPtr := flag.Bool("json", false, "JSON OUTPUT")
	listPtr := flag.String("dL", "", "TARGET DOMAINS LIST FILE")
	updatePtr := flag.Bool("update", false, "UPDATE TWHUNT")
	flag.Parse()

	if *updatePtr {
		utils.AutoUpdate()
	}

	if *domainPtr == "" && *listPtr == "" {
		if !*silentPtr {
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
			fmt.Printf("%s   📞 Phone   :%s +8801711729858%s\n", core.Sky, core.Mint, core.EndC)
			fmt.Printf("%s   🌐 GitHub  :%s https://github.com/tawhid2005%s\n", core.Sky, core.Mint, core.EndC)
			fmt.Printf("%s   -------------------------------------------------------%s\n", core.Slate, core.EndC)
			fmt.Printf("\n%s   [USAGE]%s twhunt -d <domain> [-v] [-o output.txt] [-silent] [-json] [-dL list.txt] [-update]%s\n\n", core.Gold, core.Silver, core.EndC)
		}
		os.Exit(1)
	}

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
		engine := core.NewEngine(sourceList, *silentPtr)
		subdomains := engine.Run(target)

		if *verifyPtr && len(subdomains) > 0 {
			if !*silentPtr {
				fmt.Printf(" %s[*] VERIFYING LIVE STATUS FOR %s...%s\n", core.Gold, target, core.EndC)
			}
			var live []string
			for _, sub := range subdomains {
				if _, err := net.LookupHost(sub); err == nil {
					live = append(live, sub)
				}
			}
			subdomains = live
			if !*silentPtr {
				fmt.Printf(" %s[✓] TOTAL LIVE HOSTS FOUND FOR %s: %d%s\n", core.Mint, target, len(subdomains), core.EndC)
			}
		}

		allDiscoveredSubdomains = append(allDiscoveredSubdomains, subdomains...)
	}

	// Remove completely global duplicates across all domains
	uniqueSubdomains := make(map[string]bool)
	var finalSubdomains []string
	for _, sub := range allDiscoveredSubdomains {
		if !uniqueSubdomains[sub] {
			uniqueSubdomains[sub] = true
			finalSubdomains = append(finalSubdomains, sub)
		}
	}
	sort.Strings(finalSubdomains)

	if *jsonPtr {
		jsonData, _ := json.MarshalIndent(finalSubdomains, "", "  ")
		fmt.Println(string(jsonData))
	} else if !*silentPtr && len(finalSubdomains) > 0 {
		fmt.Printf("\n%s%s========================= RESULTS =========================%s\n", core.Sky, core.Bold, core.EndC)
		for _, sub := range finalSubdomains {
			fmt.Printf(" %s→%s %s%s%s\n", core.Mint, core.EndC, core.Silver, sub, core.EndC)
		}
		fmt.Printf("%s%s===========================================================%s\n", core.Sky, core.Bold, core.EndC)
	} else if *silentPtr && len(finalSubdomains) > 0 {
		// In silent mode without JSON, just print one domain per line
		for _, sub := range finalSubdomains {
			fmt.Println(sub)
		}
	}

	if !*silentPtr && !*jsonPtr {
		fmt.Printf("\n %s[!] SCAN COMPLETED. HAPPY HUNTING!%s\n", core.Sky, core.EndC)
	}

	// সেভ করার লজিক
	if *outputPtr != "" && len(finalSubdomains) > 0 {
		saveToFile(*outputPtr, finalSubdomains, *jsonPtr, *silentPtr)
	} else if len(finalSubdomains) > 0 && !*silentPtr && !*jsonPtr {
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
				saveToFile(filename, finalSubdomains, false, *silentPtr)
			} else {
				fmt.Printf(" %s[!] No filename provided. Skipping save.%s\n", core.Coral, core.EndC)
			}
		}
	}
}

func saveToFile(filename string, data []string, isJson bool, isSilent bool) {
	file, err := os.Create(filename)
	if err != nil {
		if !isSilent {
			fmt.Printf(" %s[!] ERROR SAVING FILE: %v%s\n", core.Coral, err, core.EndC)
		}
		return
	}
	defer file.Close()
	
	if isJson {
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		file.Write(jsonData)
	} else {
		for _, sub := range data {
			file.WriteString(sub + "\n")
		}
	}
	
	if !isSilent {
		fmt.Printf(" %s[✓] RESULTS SAVED TO: %s%s\n", core.Mint, filename, core.EndC)
	}
}
