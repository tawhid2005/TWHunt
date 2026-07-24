package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	"subfinder_clone/core"
	"subfinder_clone/sources"
)

func main() {
	domainPtr := flag.String("d", "", "TARGET DOMAIN")
	verifyPtr := flag.Bool("v", false, "VERIFY LIVE STATUS")
	flag.Parse()

	if *domainPtr == "" {
		fmt.Printf("%s%s\n", core.Sky, core.Bold)
		fmt.Println("░█▀▀░█▀▀░█▀▄░▀█▀░█░█░█░█░█▀█░▀█▀")
		fmt.Println("░█░░░█▀▀░█▀▄░░█░░█▀█░█░█░█░█░░█░")
		fmt.Println("░▀▀▀░▀▀▀░▀░▀░░▀░░▀░▀░▀▀▀░▀░▀░░▀░")
		fmt.Printf("\n%s         USAGE: go run main.go -d <domain>%s\n", core.Gold, core.EndC)
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
}
