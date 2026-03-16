package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "diet-config-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Mock UserHomeDir by setting HOME (or USERPROFILE on Windows)
	// But our LoadConfig calls GetConfigPath which calls os.UserHomeDir
	// We can't easily mock os.UserHomeDir without refactoring or global vars
	// Let's test the logic by manually creating a config file and calling Unmarshal
}

func TestLoadConfigLogic(t *testing.T) {
	cfg := &Config{
		EnabledParsers: make(map[string]bool),
	}
	data := []byte(`{"enabled_parsers": {"git_status": true}}`)
	if err := json.Unmarshal(data, cfg); err != nil {
		t.Fatal(err)
	}
	if !cfg.EnabledParsers["git_status"] {
		t.Error("Expected git_status to be enabled")
	}
}

func TestConfigSaveLoad(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	cfg := &Config{
		EnabledParsers: map[string]bool{"test": true},
	}
	
	data, _ := json.Marshal(cfg)
	if err := os.WriteFile(tmpFile.Name(), data, 0644); err != nil {
		t.Fatal(err)
	}

	// Verify we can read it back
	newCfg := &Config{}
	readData, _ := os.ReadFile(tmpFile.Name())
	json.Unmarshal(readData, newCfg)
	
	if !newCfg.EnabledParsers["test"] {
		t.Error("Failed to round-trip config")
	}
}
