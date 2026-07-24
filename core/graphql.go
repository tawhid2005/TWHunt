package core

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var graphqlClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

var graphqlEndpoints = []string{
	"/graphql", "/api/graphql", "/v1/graphql", "/v2/graphql", "/graphql/v1",
}

// DetectGraphQL checks for GraphQL endpoints and attempts Introspection
func DetectGraphQL(subdomains []string, silent bool) map[string]string {
	if !silent {
		fmt.Printf(" %s[*] SEARCHING FOR GRAPHQL ENDPOINTS & INTROSPECTION...%s\n", Gold, EndC)
	}

	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 30)

	// Introspection Query to dump Schema
	introspectionQuery := `{"query":"query IntrospectionQuery { __schema { queryType { name } mutationType { name } } }"}`

	for _, sub := range subdomains {
		wg.Add(1)
		sem <- struct{}{}
		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			for _, ep := range graphqlEndpoints {
				url := "https://" + s + ep
				
				req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(introspectionQuery)))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("User-Agent", "Mozilla/5.0")
				
				resp, err := graphqlClient.Do(req)
				if err != nil {
					// Fallback to HTTP
					url = "http://" + s + ep
					req, _ = http.NewRequest("POST", url, bytes.NewBuffer([]byte(introspectionQuery)))
					req.Header.Set("Content-Type", "application/json")
					resp, err = graphqlClient.Do(req)
					if err != nil {
						continue
					}
				}

				body, _ := ioutil.ReadAll(resp.Body)
				resp.Body.Close()

				if resp.StatusCode == 200 && strings.Contains(string(body), "__schema") {
					finding := fmt.Sprintf("Introspection Enabled! Exposed Schema at: %s", ep)
					mu.Lock()
					results[s] = finding
					mu.Unlock()
					return // Stop after finding one vulnerable endpoint
				} else if resp.StatusCode == 200 && strings.Contains(string(body), "errors") && strings.Contains(string(body), "graphql") {
					finding := fmt.Sprintf("GraphQL Endpoint Found (Introspection Disabled): %s", ep)
					mu.Lock()
					results[s] = finding
					mu.Unlock()
					return
				}
			}
		}(sub)
	}

	wg.Wait()

	if !silent {
		if len(results) > 0 {
			fmt.Printf(" %s[!] FOUND %d GRAPHQL ENDPOINTS!%s\n", Coral, len(results), EndC)
		} else {
			fmt.Printf(" %s[✓] NO GRAPHQL ENDPOINTS FOUND.%s\n", Mint, EndC)
		}
	}
	return results
}
