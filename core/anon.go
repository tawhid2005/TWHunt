package core

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// CheckAnonymousLogin attempts anonymous logins on critical open ports (21 FTP, 6379 Redis)
func CheckAnonymousLogin(portResults map[string][]string, silent bool) map[string][]string {
	if !silent {
		fmt.Printf(" %s[*] CHECKING FOR ANONYMOUS/GUEST LOGIN ON OPEN PORTS...%s\n", Gold, EndC)
	}

	results := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 15)

	for sub, ports := range portResults {
		for _, port := range ports {
			if port == "21" || port == "6379" {
				wg.Add(1)
				sem <- struct{}{}
				go func(s, p string) {
					defer wg.Done()
					defer func() { <-sem }()

					target := s + ":" + p
					conn, err := net.DialTimeout("tcp", target, 5*time.Second)
					if err != nil {
						return
					}
					defer conn.Close()

					if p == "21" {
						// FTP Anonymous Check
						buf := make([]byte, 1024)
						conn.Read(buf) // Read banner
						
						fmt.Fprintf(conn, "USER anonymous\r\n")
						conn.Read(buf)
						
						fmt.Fprintf(conn, "PASS anonymous@example.com\r\n")
						n, _ := conn.Read(buf)
						response := string(buf[:n])
						
						if strings.Contains(response, "230") { // 230 Login successful
							mu.Lock()
							results[s] = append(results[s], "FTP (21): Anonymous Login ALLOWED!")
							mu.Unlock()
						}
					} else if p == "6379" {
						// Redis Unauthorized Check
						fmt.Fprintf(conn, "PING\r\n")
						buf := make([]byte, 1024)
						n, _ := conn.Read(buf)
						response := string(buf[:n])
						
						if strings.Contains(response, "PONG") {
							mu.Lock()
							results[s] = append(results[s], "Redis (6379): NO AUTHENTICATION REQUIRED!")
							mu.Unlock()
						}
					}

				}(sub, port)
			}
		}
	}

	wg.Wait()

	if !silent {
		totalFound := 0
		for _, v := range results {
			totalFound += len(v)
		}
		if totalFound > 0 {
			fmt.Printf(" %s[!] FOUND %d OPEN/UNAUTHORIZED DATABASES!%s\n", Coral, totalFound, EndC)
		} else {
			fmt.Printf(" %s[✓] NO UNAUTHORIZED DATABASES FOUND.%s\n", Mint, EndC)
		}
	}
	return results
}
