package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/AlecAivazis/survey/v2"
)

type Config struct {
	EnabledParsers map[string]bool `json:"enabled_parsers"`
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".diet", "config.json"), nil
}

func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		EnabledParsers: make(map[string]bool),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return cfg, nil
	}

	return cfg, nil
}

func (c *Config) Save() error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (c *Config) InteractiveConfig(availableParsers []string) error {
	// Sync available parsers into map (enable new ones by default)
	for _, p := range availableParsers {
		if _, ok := c.EnabledParsers[p]; !ok {
			c.EnabledParsers[p] = true
		}
	}

	// Toggle parsers
	var parserOpts []string
	var defaultSelected []string
	
	// Sort for consistent UI
	sort.Strings(availableParsers)
	
	for _, p := range availableParsers {
		parserOpts = append(parserOpts, p)
		if c.EnabledParsers[p] {
			defaultSelected = append(defaultSelected, p)
		}
	}

	fmt.Println("\n--- Diet Configuration ---")
	fmt.Println("Use [Space] to toggle, [Enter] to save, or [Ctrl+C] to cancel and exit.")
	fmt.Println("--------------------------")

	var selectedParsers []string
	prompt := &survey.MultiSelect{
		Message: "Enable/Disable Parsers (unchecked will passthrough):",
		Options: parserOpts,
		Default: defaultSelected,
	}
	err := survey.AskOne(prompt, &selectedParsers)
	if err != nil {
		return err
	}

	newEnabled := make(map[string]bool)
	for _, p := range parserOpts {
		newEnabled[p] = false
	}
	for _, p := range selectedParsers {
		newEnabled[p] = true
	}
	c.EnabledParsers = newEnabled

	return c.Save()
}
