package core

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

// RunBruteforce active brute-forces subdomains using a wordlist
func RunBruteforce(domain string, wordlist string, threads int, silent bool) []string {
	file, err := os.Open(wordlist)
	if err != nil {
		if !silent {
			fmt.Printf(" %s[!] Error opening wordlist: %v%s\n", Coral, err, EndC)
		}
		return nil
	}
	defer file.Close()

	if !silent {
		fmt.Printf(" %s[*] STARTING ACTIVE BRUTE-FORCE ENGINE (%d THREADS)...%s\n", Gold, threads, EndC)
	}

	words := make(chan string, threads*2)
	results := make(chan string, threads*2)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for word := range words {
				sub := word + "." + domain
				if _, err := net.LookupHost(sub); err == nil {
					results <- sub
				}
			}
		}()
	}

	// Read file and send to workers
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			word := strings.TrimSpace(scanner.Text())
			if word != "" {
				words <- strings.ToLower(word)
			}
		}
		close(words)
	}()

	// Wait for workers to finish in background
	go func() {
		wg.Wait()
		close(results)
	}()

	var found []string
	for sub := range results {
		found = append(found, sub)
	}

	if !silent {
		fmt.Printf(" %s[✓] BRUTE-FORCE FOUND %d VALID SUBDOMAINS%s\n", Mint, len(found), EndC)
	}

	return found
}
