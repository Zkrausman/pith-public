package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/AlecAivazis/survey/v2"
)

func TestGetConfigPath_NewDefault(t *testing.T) {
	tmpDir := t.TempDir()
	
	oldBase := TheBrainBase
	TheBrainBase = tmpDir
	defer func() { TheBrainBase = oldBase }()

	newDefault := filepath.Join(tmpDir, "TheBrain", "PithBackup")
	os.MkdirAll(newDefault, 0755)
	cfgFile := filepath.Join(newDefault, "config.json")
	os.WriteFile(cfgFile, []byte("{}"), 0644)

	t.Setenv("PITH_STORAGE", "")

	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}
	if path != cfgFile {
		t.Errorf("Expected path %s, got %s", cfgFile, path)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.StoragePath != newDefault {
		t.Errorf("Expected StoragePath %s, got %s", newDefault, cfg.StoragePath)
	}
}

func TestLoadConfig_WithDefaultsApplied(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	// Write config with zeros (should get defaults applied)
	cfgPath := filepath.Join(tmpDir, "config.json")
	os.WriteFile(cfgPath, []byte(`{"max_lines":0,"head_lines":0,"tail_lines":0}`), 0644)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
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

func TestInteractiveConfig_SurveyAskOneError(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	oldAskOne := SurveyAskOne
	defer func() { SurveyAskOne = oldAskOne }()

	SurveyAskOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
		return errors.New("user cancelled")
	}

	cfg := &Config{EnabledParsers: make(map[string]bool)}
	err := cfg.InteractiveConfig([]string{"git"})
	if err == nil || err.Error() != "user cancelled" {
		t.Errorf("Expected 'user cancelled' error, got %v", err)
	}
}

func TestInteractiveConfig_SurveyAskError(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	oldAskOne := SurveyAskOne
	oldAsk := SurveyAsk
	defer func() {
		SurveyAskOne = oldAskOne
		SurveyAsk = oldAsk
	}()

	SurveyAskOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
		resp := response.(*[]string)
		*resp = []string{"git"}
		return nil
	}
	SurveyAsk = func(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
		return errors.New("survey ask failed")
	}

	cfg := &Config{EnabledParsers: make(map[string]bool)}
	err := cfg.InteractiveConfig([]string{"git"})
	if err == nil || err.Error() != "survey ask failed" {
		t.Errorf("Expected 'survey ask failed' error, got %v", err)
	}
}

func TestInteractiveConfig_EnablesExistingParsers(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	oldAskOne := SurveyAskOne
	oldAsk := SurveyAsk
	defer func() {
		SurveyAskOne = oldAskOne
		SurveyAsk = oldAsk
	}()

	SurveyAskOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
		resp := response.(*[]string)
		*resp = []string{"git", "docker"}
		return nil
	}
	SurveyAsk = func(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
		resp := response.(*struct {
			MaxLines  int `survey:"maxlines"`
			HeadLines int `survey:"headlines"`
			TailLines int `survey:"taillines"`
		})
		resp.MaxLines = 750
		resp.HeadLines = 50
		resp.TailLines = 50
		return nil
	}

	// Config already has docker enabled and git disabled
	cfg := &Config{
		EnabledParsers: map[string]bool{"git": false, "docker": true},
		MaxLines:       500,
		HeadLines:      100,
		TailLines:      100,
	}

	err := cfg.InteractiveConfig([]string{"git", "docker"})
	if err != nil {
		t.Fatalf("InteractiveConfig failed: %v", err)
	}
	if !cfg.EnabledParsers["git"] {
		t.Error("Expected git to be enabled")
	}
	if !cfg.EnabledParsers["docker"] {
		t.Error("Expected docker to remain enabled")
	}
	if cfg.MaxLines != 750 {
		t.Errorf("Expected MaxLines 750, got %d", cfg.MaxLines)
	}
}

func TestMigrateStorage_NoOldPath2(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir)

	// ~/.pith doesn't exist in tmpDir, so migration should return nil
	err := MigrateStorage(filepath.Join(tmpDir, "new_storage"))
	if err != nil {
		t.Errorf("Expected no error when old path doesn't exist, got %v", err)
	}
}

func TestMigrateStorage_SamePath2(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir)

	// Create ~/.pith in tmpDir
	oldPath := filepath.Join(tmpDir, ".pith")
	os.MkdirAll(oldPath, 0755)

	// Migrate to same path -> should return nil immediately
	err := MigrateStorage(oldPath)
	if err != nil {
		t.Errorf("Expected no error for same-path migration, got %v", err)
	}
}

func TestMigrateStorage_WithFiles(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir)

	// Create ~/.pith with files
	oldPath := filepath.Join(tmpDir, ".pith")
	os.MkdirAll(oldPath, 0755)
	os.WriteFile(filepath.Join(oldPath, "pith.db"), []byte("db content"), 0644)
	os.WriteFile(filepath.Join(oldPath, "config.json"), []byte(`{"max_lines":100}`), 0644)

	newPath := filepath.Join(tmpDir, "new_storage")

	err := MigrateStorage(newPath)
	if err != nil {
		t.Fatalf("MigrateStorage failed: %v", err)
	}

	// Verify files were copied
	if _, err := os.Stat(filepath.Join(newPath, "pith.db")); os.IsNotExist(err) {
		t.Error("Expected pith.db to be migrated")
	}
	if _, err := os.Stat(filepath.Join(newPath, "config.json")); os.IsNotExist(err) {
		t.Error("Expected config.json to be migrated")
	}
}

func TestMigrateStorage_SkipsExistingDest(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir)

	oldPath := filepath.Join(tmpDir, ".pith")
	os.MkdirAll(oldPath, 0755)
	os.WriteFile(filepath.Join(oldPath, "pith.db"), []byte("old db"), 0644)

	newPath := filepath.Join(tmpDir, "new_storage")
	os.MkdirAll(newPath, 0755)
	// Pre-create dest file so migration skips it
	os.WriteFile(filepath.Join(newPath, "pith.db"), []byte("existing db"), 0644)

	err := MigrateStorage(newPath)
	if err != nil {
		t.Fatalf("MigrateStorage failed: %v", err)
	}

	// Verify the original dest file was NOT overwritten
	data, _ := os.ReadFile(filepath.Join(newPath, "pith.db"))
	if string(data) != "existing db" {
		t.Errorf("Expected existing db content to be preserved, got %s", string(data))
	}
}

func TestSave_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	cfg := &Config{
		EnabledParsers: map[string]bool{"git": true},
		MaxLines:       999,
		HeadLines:      77,
		TailLines:      77,
	}
	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load it back
	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig after save failed: %v", err)
	}
	if loaded.MaxLines != 999 {
		t.Errorf("Expected MaxLines 999, got %d", loaded.MaxLines)
	}
}

func TestCopyFile_Error(t *testing.T) {
	// Try to copy a non-existent file
	err := copyFile("/nonexistent/source.txt", "/tmp/dest.txt")
	if err == nil {
		t.Error("Expected error when source file doesn't exist")
	}
}

func TestMigrateStorage_WithExtraFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir)

	// Create ~/.pith with only config.json (pith.db missing)
	oldPath := filepath.Join(tmpDir, ".pith")
	os.MkdirAll(oldPath, 0755)
	os.WriteFile(filepath.Join(oldPath, "config.json"), []byte(`{"max_lines":200}`), 0644)
	// pith.db intentionally missing

	newPath := filepath.Join(tmpDir, "new_storage2")

	err := MigrateStorage(newPath)
	if err != nil {
		t.Fatalf("MigrateStorage failed: %v", err)
	}

	// config.json should be copied but pith.db should not exist
	if _, err := os.Stat(filepath.Join(newPath, "config.json")); os.IsNotExist(err) {
		t.Error("Expected config.json to be migrated")
	}
	if _, err := os.Stat(filepath.Join(newPath, "pith.db")); err == nil {
		t.Error("Expected pith.db NOT to be migrated (it didn't exist)")
	}
}

