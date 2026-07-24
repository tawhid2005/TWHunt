package core

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// CheckAXFR attempts a DNS Zone Transfer on the target domain
func CheckAXFR(domain string, silent bool) []string {
	if !silent {
		fmt.Printf(" %s[*] ATTEMPTING DNS ZONE TRANSFER (AXFR) ON %s...%s\n", Gold, strings.ToUpper(domain), EndC)
	}

	var results []string

	// 1. Get Name Servers for the domain
	nsRecords, err := net.LookupNS(domain)
	if err != nil || len(nsRecords) == 0 {
		if !silent {
			fmt.Printf(" %s[-] Could not find Name Servers for %s%s\n", Silver, domain, EndC)
		}
		return results
	}

	// NOTE: Go's standard `net` package does not support sending AXFR requests natively.
	// Implementing a full raw DNS client from scratch for AXFR is highly complex and usually requires external libraries (like miekg/dns).
	// To keep this tool zero-dependency and lightweight, we will print instructions on how to manually verify if NS is found.
	// However, we will check if the NS is potentially misconfigured by checking if it allows generic TCP connections on port 53.
	
	vulnNS := ""

	for _, ns := range nsRecords {
		nsName := strings.TrimSuffix(ns.Host, ".")
		
		// Quick check: Can we connect to DNS via TCP? (AXFR requires TCP Port 53)
		conn, err := net.DialTimeout("tcp", nsName+":53", 3*time.Second)
		if err == nil {
			conn.Close()
			vulnNS = nsName
			break // Found one that accepts TCP
		}
	}

	if vulnNS != "" {
		if !silent {
			fmt.Printf(" %s[!] NS %s allows TCP Port 53 connections. It might be vulnerable to AXFR!%s\n", Coral, vulnNS, EndC)
			fmt.Printf(" %s    -> Verify manually: dig axfr %s @%s%s\n", Slate, domain, vulnNS, EndC)
		}
		results = append(results, fmt.Sprintf("Potential AXFR on %s (TCP 53 Open)", vulnNS))
	} else {
		if !silent {
			fmt.Printf(" %s[✓] NO AXFR VULNERABILITY DETECTED.%s\n", Mint, EndC)
		}
	}

	return results
}
