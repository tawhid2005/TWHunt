package core

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"twhunt/utils"
)

// FetchWaybackURLs queries the Wayback Machine CDX API for historical URLs of a domain
func FetchWaybackURLs(domain string, silent bool) []string {
	if !silent {
		fmt.Printf(" %s[*] FETCHING WAYBACK MACHINE URLS FOR %s...%s\n", Gold, domain, EndC)
	}

	url := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=txt&fl=original&collapse=urlkey", domain)
	
	body, err := utils.RequestAPI(url)
	if err != nil {
		if !silent {
			fmt.Printf(" %s[!] Error fetching Wayback URLs: %v%s\n", Coral, err, EndC)
		}
		return nil
	}

	var urls []string
	scanner := bufio.NewScanner(bytes.NewReader(body))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			urls = append(urls, line)
		}
	}

	if !silent {
		if len(urls) > 0 {
			fmt.Printf(" %s[✓] FOUND %d HISTORICAL URLS FROM WAYBACK MACHINE!%s\n", Mint, len(urls), EndC)
		} else {
			fmt.Printf(" %s[!] NO WAYBACK URLS FOUND.%s\n", Coral, EndC)
		}
	}

	return urls
}
