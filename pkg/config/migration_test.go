package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrateStorage(t *testing.T) {
	var notices bytes.Buffer
	previousOutput := migrationOutput
	migrationOutput = &notices
	t.Cleanup(func() { migrationOutput = previousOutput })

	// Create a temporary home directory
	tempHome, err := os.MkdirTemp("", "pith-home-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempHome)

	// Set home directory for testing
	t.Setenv("USERPROFILE", tempHome) // Windows
	t.Setenv("HOME", tempHome)        // Unix

	oldPath := filepath.Join(tempHome, ".pith")
	if err := os.MkdirAll(oldPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Create some old files
	dbContent := "fake db"
	configContent := `{"max_lines": 123}`
	if err := os.WriteFile(filepath.Join(oldPath, "pith.db"), []byte(dbContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(oldPath, "config.json"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Target path
	targetPath, err := os.MkdirTemp("", "pith-target-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(targetPath)

	// Migrate
	if err := MigrateStorage(targetPath); err != nil {
		t.Fatalf("MigrateStorage failed: %v", err)
	}
	if !strings.Contains(notices.String(), "Migrating pith.db") || !strings.Contains(notices.String(), "Migrating config.json") {
		t.Fatalf("expected migration notices on stderr writer, got %q", notices.String())
	}

	// Verify files in target
	targetDB := filepath.Join(targetPath, "pith.db")
	targetConfig := filepath.Join(targetPath, "config.json")

	if data, err := os.ReadFile(targetDB); err != nil || string(data) != dbContent {
		t.Errorf("DB migration failed: %v, content: %s", err, string(data))
	}
	if data, err := os.ReadFile(targetConfig); err != nil || string(data) != configContent {
		t.Errorf("Config migration failed: %v, content: %s", err, string(data))
	}

	// Verify old files are renamed to .bak
	if _, err := os.Stat(filepath.Join(oldPath, "pith.db.bak")); os.IsNotExist(err) {
		t.Error("Old DB was not renamed to .bak")
	}
	if _, err := os.Stat(filepath.Join(oldPath, "config.json.bak")); os.IsNotExist(err) {
		t.Error("Old Config was not renamed to .bak")
	}
}

func TestMigrateStorage_NoOldPath(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "pith-home-empty-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempHome)

	t.Setenv("USERPROFILE", tempHome)
	t.Setenv("HOME", tempHome)

	targetPath, _ := os.MkdirTemp("", "pith-target-*")
	defer os.RemoveAll(targetPath)

	if err := MigrateStorage(targetPath); err != nil {
		t.Fatalf("MigrateStorage should not fail when old path doesn't exist: %v", err)
	}
}

func TestMigrateStorage_SamePath(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "pith-home-same-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempHome)

	t.Setenv("USERPROFILE", tempHome)
	t.Setenv("HOME", tempHome)

	oldPath := filepath.Join(tempHome, ".pith")
	if err := os.MkdirAll(oldPath, 0755); err != nil {
		t.Fatal(err)
	}

	if err := MigrateStorage(oldPath); err != nil {
		t.Fatalf("MigrateStorage should not fail when target is same as old: %v", err)
	}
}
