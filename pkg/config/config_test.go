package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/AlecAivazis/survey/v2"
)

func TestGetConfigPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pith-storage-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Setenv("PITH_STORAGE", tmpDir)

	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := filepath.Join(tmpDir, "config.json")
	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}
}

func TestLoadConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pith-load-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Setenv("PITH_STORAGE", tmpDir)

	// Case 1: Config does not exist
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if cfg.MaxLines != 500 {
		t.Errorf("Expected MaxLines 500, got %d", cfg.MaxLines)
	}

	// Case 2: Config exists
	cfg.MaxLines = 1000
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	cfg2, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if cfg2.MaxLines != 1000 {
		t.Errorf("Expected MaxLines 1000, got %d", cfg2.MaxLines)
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pith-save-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Setenv("PITH_STORAGE", tmpDir)

	cfg := &Config{
		EnabledParsers: map[string]bool{"test": true},
		MaxLines:       123,
	}

	if err := cfg.Save(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify file exists and content is correct
	path := filepath.Join(tmpDir, "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Expected file to exist, got error %v", err)
	}

	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Failed to unmarshal saved config: %v", err)
	}

	if loaded.MaxLines != 123 || !loaded.EnabledParsers["test"] {
		t.Errorf("Loaded config does not match saved config: %+v", loaded)
	}
}

func TestConfigDefaults(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pith-defaults-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Setenv("PITH_STORAGE", tmpDir)

	// Create a config file with missing fields
	path := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(path, []byte(`{"enabled_parsers": {"git": true}}`), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.MaxLines != 500 {
		t.Errorf("Expected default MaxLines 500, got %d", cfg.MaxLines)
	}
	if cfg.HeadLines != 100 {
		t.Errorf("Expected default HeadLines 100, got %d", cfg.HeadLines)
	}
	if cfg.TailLines != 100 {
		t.Errorf("Expected default TailLines 100, got %d", cfg.TailLines)
	}
}
func TestLoadConfig_Malformed(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	path := filepath.Join(tmpDir, "config.json")
	os.WriteFile(path, []byte(`{invalid-json}`), 0644)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error (falling back to defaults), got %v", err)
	}
	if cfg.MaxLines != 500 {
		t.Error("Expected default MaxLines for malformed config")
	}
}

func TestGetConfigPath_Default(t *testing.T) {
	// Temporarily unset PITH_STORAGE
	oldEnv := os.Getenv("PITH_STORAGE")
	t.Setenv("PITH_STORAGE", "")
	defer os.Setenv("PITH_STORAGE", oldEnv)

	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}
	if path == "" {
		t.Error("Expected non-empty path")
	}
}

func TestInteractiveConfig(t *testing.T) {
	// Mock survey functions
	oldAskOne := SurveyAskOne
	oldAsk := SurveyAsk
	defer func() {
		SurveyAskOne = oldAskOne
		SurveyAsk = oldAsk
	}()

	SurveyAskOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
		// Simulate selecting "git" parser
		resp := response.(*[]string)
		*resp = []string{"git"}
		return nil
	}

	SurveyAsk = func(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {

		// Simulate answers
		resp := response.(*struct {
			MaxLines  int     `survey:"maxlines"`
			HeadLines int     `survey:"headlines"`
			TailLines int     `survey:"taillines"`
			USDRate   float64 `survey:"usdrate"`
			Heuristic float64 `survey:"heuristic"`
		})
		resp.MaxLines = 1000
		resp.HeadLines = 200
		resp.TailLines = 200
		resp.USDRate = 5.0
		resp.Heuristic = 3.5
		return nil
	}

	cfg := &Config{
		EnabledParsers: make(map[string]bool),
		MaxLines:       500,
	}

	err := cfg.InteractiveConfig([]string{"git", "docker"})
	if err != nil {
		t.Fatalf("InteractiveConfig failed: %v", err)
	}

	if !cfg.EnabledParsers["git"] {
		t.Error("Expected git parser to be enabled")
	}
	if cfg.EnabledParsers["docker"] {
		t.Error("Expected docker parser to be disabled")
	}
	if cfg.MaxLines != 1000 {
		t.Errorf("Expected MaxLines 1000, got %d", cfg.MaxLines)
	}
}
