package runner

import (
	"os"
	"path/filepath"
	"pith/pkg/config"
	"pith/pkg/telemetry"
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
	// We allow it to fail if go is not in path or behaves weirdly in test env
	_ = run.Run([]string{"go", "version"})

	// Test a command that doesn't exist
	err = run.Run([]string{"this_command_should_not_exist_xyz123"})
	if err == nil {
		t.Error("Expected error for non-existent command")
	}

	// Test command with pipe (composite)
	// We allow it to fail if tools are missing, just ensure no crash
	if os.PathSeparator == '\\' {
		_ = run.Run([]string{"go", "version", "|", "findstr", "go"})
	} else {
		_ = run.Run([]string{"go", "version", "|", "grep", "go"})
	}

	// Test command passed as single string with spaces
	_ = run.Run([]string{"go version"})
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
