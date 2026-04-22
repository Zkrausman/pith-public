package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

var SurveyAskOne = survey.AskOne
var SurveyAsk = survey.Ask

type Config struct {
	EnabledParsers  map[string]bool `json:"enabled_parsers"`
	LastUpdateCheck int64           `json:"last_update_check"`
	MaxLines        int             `json:"max_lines"`
	HeadLines       int             `json:"head_lines"`
	TailLines       int             `json:"tail_lines"`
	StoragePath     string          `json:"storage_path"`
}

func GetConfigPath() (string, error) {
	// 1. Check environment variable
	if env := os.Getenv("PITH_STORAGE"); env != "" {
		return filepath.Join(env, "config.json"), nil
	}

	// 2. Check if config exists in the new default location
	newDefault := filepath.Join("E:\\", "TheBrain", "PithBackup")
	newConfig := filepath.Join(newDefault, "config.json")
	if _, err := os.Stat(newConfig); err == nil {
		return newConfig, nil
	}

	// 3. Fallback to old default
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".pith", "config.json"), nil
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
		StoragePath:    filepath.Dir(path),
	}

	// If we are using the old default, but E:\TheBrain\PithBackup is the intended new default,
	// we should set StoragePath to the new default to trigger migration.
	newDefault := filepath.Join("E:\\", "TheBrain", "PithBackup")
	if env := os.Getenv("PITH_STORAGE"); env != "" {
		cfg.StoragePath = env
	} else if _, err := os.Stat(filepath.Join(newDefault, "config.json")); err == nil {
		cfg.StoragePath = newDefault
	} else if strings.Contains(path, ".pith") {
		// If we are currently in .pith and newDefault doesn't exist, we WANT to go to newDefault
		cfg.StoragePath = newDefault
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
	err := SurveyAskOne(prompt, &selectedParsers)
	if err != nil {
		return err
	}

	// Update enabled parsers
	for _, p := range availableParsers {
		c.EnabledParsers[p] = false
	}
	for _, p := range selectedParsers {
		c.EnabledParsers[p] = true
	}

	// Ask for other settings
	qs := []*survey.Question{
		{
			Name: "maxlines",
			Prompt: &survey.Input{
				Message: "Maximum total lines to keep in output (MaxLines):",
				Default: fmt.Sprintf("%d", c.MaxLines),
			},
		},
		{
			Name: "headlines",
			Prompt: &survey.Input{
				Message: "Lines to keep at the start (HeadLines):",
				Default: fmt.Sprintf("%d", c.HeadLines),
			},
		},
		{
			Name: "taillines",
			Prompt: &survey.Input{
				Message: "Lines to keep at the end (TailLines):",
				Default: fmt.Sprintf("%d", c.TailLines),
			},
		},
	}

	answers := struct {
		MaxLines  int `survey:"maxlines"`
		HeadLines int `survey:"headlines"`
		TailLines int `survey:"taillines"`
	}{}

	err = SurveyAsk(qs, &answers)
	if err != nil {
		return err
	}


	c.MaxLines = answers.MaxLines
	c.HeadLines = answers.HeadLines
	c.TailLines = answers.TailLines

	return c.Save()
}
