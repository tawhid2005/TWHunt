package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
)

// GetWildcardIPs generates a random subdomain and resolves it to detect wildcard IPs
func GetWildcardIPs(domain string, silent bool) []string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	randomSub := hex.EncodeToString(bytes) + "." + domain

	ips, err := net.LookupHost(randomSub)
	if err == nil && len(ips) > 0 {
		if !silent {
			fmt.Printf(" %s[!] WILDCARD DNS DETECTED! Filtering out fake subdomains pointing to: %v%s\n", Coral, ips, EndC)
		}
		return ips
	}
	return nil
}

// IsWildcard checks if a subdomain's IPs match any of the wildcard IPs
func IsWildcard(subdomain string, wildcardIPs []string) bool {
	if len(wildcardIPs) == 0 {
		return false
	}
	
	ips, err := net.LookupHost(subdomain)
	if err != nil {
		return false
	}

	for _, ip := range ips {
		for _, wip := range wildcardIPs {
			if ip == wip {
				return true
			}
		}
	}
	return false
}
