package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
)

type Config struct {
	MaxTokens       int             `json:"max_tokens"`
	EnabledParsers  map[string]bool `json:"enabled_parsers"`
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
		MaxTokens: 4000,
		EnabledParsers: map[string]bool{
			"git_status": true,
		},
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

func (c *Config) InteractiveConfig() error {
	var qs = []*survey.Question{
		{
			Name: "max_tokens",
			Prompt: &survey.Input{
				Message: "Maximum output token limit:",
				Default: "4000",
			},
		},
	}

	answers := struct {
		MaxTokens int
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return err
	}

	c.MaxTokens = answers.MaxTokens

	// Toggle parsers
	var parserOpts []string
	var defaultSelected []string
	for k, v := range c.EnabledParsers {
		parserOpts = append(parserOpts, k)
		if v {
			defaultSelected = append(defaultSelected, k)
		}
	}

	var selectedParsers []string
	prompt := &survey.MultiSelect{
		Message: "Enable/Disable Parsers:",
		Options: parserOpts,
		Default: defaultSelected,
	}
	err = survey.AskOne(prompt, &selectedParsers)
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
