package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	GitHubToken       string `json:"github_token"`
	SecurityTrailsKey string `json:"securitytrails_key"`
	ShodanKey         string `json:"shodan_key"`
	ChaosKey          string `json:"chaos_key"`
	CensysID          string `json:"censys_id"`
	CensysSecret      string `json:"censys_secret"`
}

var AppConfig Config

func LoadConfig(silent bool) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	configDir := filepath.Join(homeDir, ".config", "twhunt")
	configFile := filepath.Join(configDir, "config.json")

	// If config doesn't exist, create a skeleton
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
		defaultConfig := Config{}
		data, _ := json.MarshalIndent(defaultConfig, "", "  ")
		os.WriteFile(configFile, data, 0644)
		if !silent {
			fmt.Printf(" %s[*] Created default config at: %s%s\n", Sky, configFile, EndC)
		}
		return
	}

	// Read existing config
	data, err := os.ReadFile(configFile)
	if err != nil {
		return
	}

	json.Unmarshal(data, &AppConfig)
	if !silent {
		// Just a small notification that config is loaded if any key exists
		if AppConfig.GitHubToken != "" || AppConfig.ShodanKey != "" || AppConfig.SecurityTrailsKey != "" {
			fmt.Printf(" %s[*] Loaded API Keys from config.json%s\n", Mint, EndC)
		}
	}
}
