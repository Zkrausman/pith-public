package runner

import (
	"pith/pkg/config"
	"pith/pkg/telemetry"
	"os"
	"path/filepath"
	"testing"
)

func TestRunWithOptions(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "pith-run-test-*")
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	tel, _ := telemetry.NewTelemetryWithPath(dbPath)
	defer tel.Close()

	cfg := &config.Config{
		EnabledParsers: map[string]bool{},
		StoragePath:    tmpDir,
		MaxLines:       10,
	}

	run := NewRunner(cfg, tel)

	// Test empty args
	err := run.RunWithOptions([]string{}, false)
	if err == nil {
		t.Error("Expected error for empty args")
	}

	// Test basic command (go version)
	err = run.Run([]string{"go", "version"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test a command that doesn't exist
	err = run.Run([]string{"this_command_should_not_exist_xyz123"})
	if err == nil {
		t.Error("Expected error for non-existent command")
	}

	// Test command with pipe (composite)
	if os.PathSeparator == '\\' {
		err = run.Run([]string{"go", "version", "|", "findstr", "go"})
	} else {
		err = run.Run([]string{"go", "version", "|", "grep", "go"})
	}
	if err != nil {
		t.Errorf("Unexpected error for pipe command: %v", err)
	}

	// Test command passed as single string with spaces
	err = run.Run([]string{"go version"})
	if err != nil {
		t.Errorf("Unexpected error for single string command: %v", err)
	}
}

func TestRunWithOptionsParsing(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "pith-run-test-parsing-*")
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	tel, _ := telemetry.NewTelemetryWithPath(dbPath)
	defer tel.Close()

	cfg := &config.Config{
		EnabledParsers: map[string]bool{"git": true},
		StoragePath:    tmpDir,
		MaxLines:       100,
	}

	run := NewRunner(cfg, tel)

	// We use "echo" command to simulate something, skipParsing = false
	// We want to hit the parser loops
	err := run.RunWithOptions([]string{"git", "status"}, false)
	if err != nil {
		t.Logf("git status failed, but that's ok if not in git repo: %v", err)
	}
}
