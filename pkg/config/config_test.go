package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
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

func TestMigrateStorage(t *testing.T) {
	tmpHome, err := os.MkdirTemp("", "pith-home-mock-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	oldPath := filepath.Join(tmpHome, ".pith")
	os.MkdirAll(oldPath, 0755)
	
	// Create old config and db
	os.WriteFile(filepath.Join(oldPath, "config.json"), []byte(`{"max_lines": 999}`), 0644)
	os.WriteFile(filepath.Join(oldPath, "pith.db"), []byte("mock db content"), 0644)

	newPath := filepath.Join(tmpHome, "new-storage")

	if err := MigrateStorage(newPath); err != nil {
		t.Fatalf("MigrateStorage failed: %v", err)
	}

	// Verify files moved
	if _, err := os.Stat(filepath.Join(newPath, "config.json")); err != nil {
		t.Error("config.json not migrated")
	}
	if _, err := os.Stat(filepath.Join(newPath, "pith.db")); err != nil {
		t.Error("pith.db not migrated")
	}

	// Verify old files renamed to .bak
	if _, err := os.Stat(filepath.Join(oldPath, "config.json.bak")); err != nil {
		t.Error("config.json.bak not created")
	}
}

