package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type StateFile struct {
	KnownSubdomains []string `json:"known_subdomains"`
}

func SendDiscordNotification(webhookURL, msg string) {
	payload := map[string]string{"content": msg}
	jsonPayload, _ := json.Marshal(payload)
	
	http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
}

// CompareAndNotify compares found subdomains with previous state and sends discord alerts for new ones
func CompareAndNotify(domain string, newSubdomains []string, webhookURL string, silent bool) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	
	statePath := filepath.Join(homeDir, ".config", "twhunt", fmt.Sprintf("%s_state.json", domain))
	
	var state StateFile
	
	if _, err := os.Stat(statePath); err == nil {
		data, _ := os.ReadFile(statePath)
		json.Unmarshal(data, &state)
	}

	knownMap := make(map[string]bool)
	for _, k := range state.KnownSubdomains {
		knownMap[k] = true
	}

	var brandNew []string
	for _, n := range newSubdomains {
		if !knownMap[n] {
			brandNew = append(brandNew, n)
			state.KnownSubdomains = append(state.KnownSubdomains, n)
		}
	}

	// Save new state
	data, _ := json.MarshalIndent(state, "", "  ")
	os.WriteFile(statePath, data, 0644)

	// Send notification if there are new subdomains and a webhook is provided
	if len(brandNew) > 0 {
		if !silent {
			fmt.Printf(" %s[!] FOUND %d BRAND NEW SUBDOMAINS SINCE LAST SCAN!%s\n", Mint, len(brandNew), EndC)
		}
		
		if webhookURL != "" {
			if !silent {
				fmt.Printf(" %s[*] SENDING ALERTS TO DISCORD...%s\n", Gold, EndC)
			}
			msg := fmt.Sprintf("🚨 **TWHunt Alert** 🚨\nFound **%d** new subdomains for `%s`:\n```\n", len(brandNew), domain)
			for i, sub := range brandNew {
				if i >= 15 {
					msg += fmt.Sprintf("...and %d more\n", len(brandNew)-15)
					break
				}
				msg += sub + "\n"
			}
			msg += "```"
			SendDiscordNotification(webhookURL, msg)
		}
	}
}
