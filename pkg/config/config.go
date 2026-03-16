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
	EnabledParsers  map[string]bool `json:"enabled_parsers"`
	LastUpdateCheck int64           `json:"last_update_check"`
	MaxLines        int             `json:"max_lines"`
	HeadLines       int             `json:"head_lines"`
	TailLines       int             `json:"tail_lines"`
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
		MaxLines:       500,
		HeadLines:      100,
		TailLines:      100,
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

	// Ensure defaults if missing from JSON
	if cfg.MaxLines == 0 {
		cfg.MaxLines = 500
	}
	if cfg.HeadLines == 0 {
		cfg.HeadLines = 100
	}
	if cfg.TailLines == 0 {
		cfg.TailLines = 100
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

	var selectedParsers []string
	prompt := &survey.MultiSelect{
		Message: "Enable/Disable Parsers (Space to toggle, Enter to save, Ctrl+C to exit):",
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

	// Truncation settings
	var qs = []*survey.Question{
		{
			Name: "MaxLines",
			Prompt: &survey.Input{
				Message: "Max lines before middle-out truncation:",
				Default: fmt.Sprintf("%d", c.MaxLines),
			},
		},
		{
			Name: "HeadLines",
			Prompt: &survey.Input{
				Message: "Number of lines to keep at the START:",
				Default: fmt.Sprintf("%d", c.HeadLines),
			},
		},
		{
			Name: "TailLines",
			Prompt: &survey.Input{
				Message: "Number of lines to keep at the END:",
				Default: fmt.Sprintf("%d", c.TailLines),
			},
		},
	}

	answers := struct {
		MaxLines  int
		HeadLines int
		TailLines int
	}{}

	if err := survey.Ask(qs, &answers); err != nil {
		return err
	}

	c.MaxLines = answers.MaxLines
	c.HeadLines = answers.HeadLines
	c.TailLines = answers.TailLines

	return c.Save()
}
