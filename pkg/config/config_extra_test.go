package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigPath_EnvEmpty(t *testing.T) {
	t.Setenv("PITH_STORAGE", "")
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}
	if path == "" {
		t.Error("Expected non-empty path")
	}
}

func TestLoadConfig_MalformedJSON(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	path := filepath.Join(tmpDir, "config.json")
	os.WriteFile(path, []byte("{ malformed json }"), 0644)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig should not return error on malformed JSON, got %v", err)
	}
	if cfg.MaxLines != 500 {
		t.Errorf("Expected default MaxLines 500, got %d", cfg.MaxLines)
	}
}

func TestLoadConfig_PartialJSON(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	path := filepath.Join(tmpDir, "config.json")
	os.WriteFile(path, []byte(`{"enabled_parsers": {"git": true}}`), 0644)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.MaxLines != 500 {
		t.Errorf("Expected default MaxLines 500, got %d", cfg.MaxLines)
	}
}

func TestSave_WriteError(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	// Create a directory where the config file should be
	path := filepath.Join(tmpDir, "config.json")
	os.MkdirAll(path, 0755)

	cfg := &Config{}
	err := cfg.Save()
	if err == nil {
		t.Error("Expected error when Save fails to write to a directory path, got nil")
	}
}
