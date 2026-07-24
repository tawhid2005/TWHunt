package core

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// ScanPorts scans the provided ports on all live subdomains concurrently
func ScanPorts(subdomains []string, ports string, silent bool) map[string][]string {
	if !silent {
		fmt.Printf(" %s[*] STARTING TCP PORT SCANNING (%s)...%s\n", Gold, ports, EndC)
	}
	
	portList := strings.Split(ports, ",")
	results := make(map[string][]string) // subdomain -> []openPorts
	var mu sync.Mutex

	var wg sync.WaitGroup
	// Limit concurrency using a semaphore (e.g., max 100 concurrent dials)
	sem := make(chan struct{}, 100)

	for _, sub := range subdomains {
		for _, p := range portList {
			port := strings.TrimSpace(p)
			if port == "" {
				continue
			}

			wg.Add(1)
			sem <- struct{}{}
			go func(s, p string) {
				defer wg.Done()
				defer func() { <-sem }()
				
				target := net.JoinHostPort(s, p)
				conn, err := net.DialTimeout("tcp", target, 2*time.Second)
				if err == nil {
					conn.Close()
					mu.Lock()
					results[s] = append(results[s], p)
					mu.Unlock()
				}
			}(sub, port)
		}
	}

	wg.Wait()

	if !silent {
		openCount := 0
		for _, p := range results {
			if len(p) > 0 {
				openCount++
			}
		}
		fmt.Printf(" %s[✓] FOUND OPEN PORTS ON %d SUBDOMAINS%s\n", Mint, openCount, EndC)
	}

	return results
}
