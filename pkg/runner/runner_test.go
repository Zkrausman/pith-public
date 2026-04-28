package runner

import (
	"os"
	"path/filepath"
	"pith/pkg/config"
	"pith/pkg/telemetry"
	"strings"
	"testing"
)

func TestEstimateTokens(t *testing.T) {
	s := "12345678"
	if EstimateTokens(s) != 2 {
		t.Errorf("Expected 2, got %d", EstimateTokens(s))
	}
}

func TestRunner(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "pith-runner-test-*")
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	tel, _ := telemetry.NewTelemetryWithPath(dbPath)
	defer tel.Close()

	cfg := &config.Config{
		EnabledParsers: map[string]bool{"git_status": true},
	}

	run := NewRunner(cfg, tel)
	if len(run.parsers) == 0 {
		t.Error("Expected parsers to be loaded")
	}
}

func TestMiddleOutTruncation(t *testing.T) {
	cfg := &config.Config{
		MaxLines:  10,
		HeadLines: 2,
		TailLines: 2,
	}
	run := NewRunner(cfg, nil)

	input := "1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n11\n12\n13\n14\n15"
	output := run.ApplyMiddleOutTruncation(input)

	if !strings.Contains(output, "removed by Pith") {
		t.Error("Expected output to be truncated")
	}
	if !strings.HasPrefix(output, "1\n2") {
		t.Error("Expected head to be preserved")
	}
	if !strings.HasSuffix(output, "14\n15") {
		t.Error("Expected tail to be preserved")
	}
}

func TestDetectSource(t *testing.T) {
	// Gemini
	t.Setenv("GEMINI_CLI", "1")
	if DetectSource() != "gemini" {
		t.Error("Expected gemini source")
	}
	t.Setenv("GEMINI_CLI", "")

	// Claude
	t.Setenv("CLAUDE_CODE", "1")
	if DetectSource() != "claude" {
		t.Error("Expected claude source")
	}
	t.Setenv("CLAUDE_CODE", "")

	// Unknown
	t.Setenv("GEMINI_CLI", "")
	t.Setenv("GOOGLE_API_KEY", "")
	t.Setenv("CLAUDE_CODE", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	if DetectSource() != "unknown" {
		t.Error("Expected unknown source")
	}
}

func TestRunWithOptions_NoCmd(t *testing.T) {
	run := NewRunner(&config.Config{}, nil)
	err := run.RunWithOptions([]string{}, false)
	if err == nil {
		t.Error("Expected error for empty args")
	}
}

func TestRunWithOptions_Simple(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := telemetry.NewTelemetry(tmpDir)
	defer tel.Close()
	cfg := &config.Config{StoragePath: tmpDir}
	run := NewRunner(cfg, tel)

	// Run 'cmd /c echo hello' on Windows
	err := run.RunWithOptions([]string{"cmd", "/c", "echo", "hello"}, false)
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}
}

func TestRunWithOptions_Shell(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := telemetry.NewTelemetry(tmpDir)
	defer tel.Close()
	cfg := &config.Config{StoragePath: tmpDir}
	run := NewRunner(cfg, tel)

	// Run composite command
	err := run.RunWithOptions([]string{"echo hello | echo world"}, false)
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}
}

func TestRunWithOptions_Chain(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{StoragePath: tmpDir}
	tel, _ := telemetry.NewTelemetry(tmpDir)
	defer tel.Close()
	r := NewRunner(cfg, tel)

	// Test composite command
	err := r.RunWithOptions([]string{"cmd", "/c", "echo line1 & echo line2"}, false)
	if err != nil {
		t.Fatalf("RunWithOptions chain failed: %v", err)
	}
}

func TestRunWithOptions_Error(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{StoragePath: tmpDir}
	tel, _ := telemetry.NewTelemetry(tmpDir)
	defer tel.Close()
	r := NewRunner(cfg, tel)

	// Command that fails
	err := r.RunWithOptions([]string{"cmd", "/c", "exit 1"}, false)
	if err == nil {
		t.Error("Expected error for exit 1")
	}
}

func TestApplyMiddleOutTruncation_Short(t *testing.T) {
	cfg := &config.Config{MaxLines: 10, HeadLines: 5, TailLines: 5}
	r := &Runner{cfg: cfg}

	input := "1\n2\n3\n4\n5\n6\n7\n8\n9\n10"
	output := r.ApplyMiddleOutTruncation(input)
	if output != input {
		t.Errorf("Expected no truncation for 10 lines, got %s", output)
	}
}

func TestLogForSnag_Default(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("HOME", tmpHome)

	r := &Runner{cfg: &config.Config{}} // Empty storage path
	r.LogForSnag("test cmd", "test output", 0)

	logPath := filepath.Join(tmpHome, ".pith", "pith.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("Log not created at default path: %s", logPath)
	}
}

func TestLogForSnag_Truncate(t *testing.T) {
	tmpDir := t.TempDir()
	r := &Runner{cfg: &config.Config{StoragePath: tmpDir}}

	output := ""
	for i := 0; i < 100; i++ {
		output += "line\n"
	}

	r.LogForSnag("long cmd", output, 0)

	logPath := filepath.Join(tmpDir, "pith.log")
	data, _ := os.ReadFile(logPath)
	if !strings.Contains(string(data), "truncated by Pith for Snag log") {
		t.Error("Expected log truncation message")
	}
}
