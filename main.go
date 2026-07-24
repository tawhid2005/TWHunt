package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	"twhunt/core"
	"twhunt/sources"
)

func main() {
	domainPtr := flag.String("d", "", "TARGET DOMAIN")
	verifyPtr := flag.Bool("v", false, "VERIFY LIVE STATUS")
	outputPtr := flag.String("o", "", "OUTPUT FILE TO SAVE RESULTS")
	flag.Parse()

	if *domainPtr == "" {
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
		fmt.Printf("%s   Created by: %sMD TALHA HUSSAIN TAWHID%s\n", core.Slate, core.Mint, core.EndC)
		fmt.Printf("\n%s   [USAGE]%s ./twhunt -d <domain> [-v] [-o output.txt]%s\n\n", core.Gold, core.Silver, core.EndC)
		os.Exit(1)
	}

	target := strings.ToLower(*domainPtr)

	// ২৪টি সোর্স যোগ করা হচ্ছে
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

	engine := core.NewEngine(sourceList)
	subdomains := engine.Run(target)

	// ডোমেইনগুলো সর্ট করা
	sort.Strings(subdomains)

	// লাইভ স্ট্যাটাস চেক করা (যদি -v ফ্ল্যাগ দেওয়া থাকে)
	if *verifyPtr && len(subdomains) > 0 {
		fmt.Printf(" %s[*] VERIFYING LIVE STATUS...%s\n", core.Gold, core.EndC)
		var live []string
		for _, sub := range subdomains {
			if _, err := net.LookupHost(sub); err == nil {
				live = append(live, sub)
			}
		}
		subdomains = live
		fmt.Printf(" %s[✓] TOTAL LIVE HOSTS FOUND: %d%s\n", core.Mint, len(subdomains), core.EndC)
	}

	// রেজাল্ট প্রিন্ট করা
	if len(subdomains) > 0 {
		fmt.Printf("\n%s%s========================= RESULTS =========================%s\n", core.Sky, core.Bold, core.EndC)
		for _, sub := range subdomains {
			fmt.Printf(" %s→%s %s%s%s\n", core.Mint, core.EndC, core.Silver, sub, core.EndC)
		}
		fmt.Printf("%s%s===========================================================%s\n", core.Sky, core.Bold, core.EndC)
	}

	fmt.Printf("\n %s[!] SCAN COMPLETED. HAPPY HUNTING!%s\n", core.Sky, core.EndC)

	// সেভ করার লজিক
	if *outputPtr != "" && len(subdomains) > 0 {
		file, err := os.Create(*outputPtr)
		if err != nil {
			fmt.Printf(" %s[!] ERROR SAVING FILE: %v%s\n", core.Coral, err, core.EndC)
			return
		}
		defer file.Close()
		
		for _, sub := range subdomains {
			file.WriteString(sub + "\n")
		}
		fmt.Printf(" %s[✓] RESULTS SAVED TO: %s%s\n", core.Mint, *outputPtr, core.EndC)
	}
}
