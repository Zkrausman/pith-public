package runner

import (
	"os"
	"path/filepath"
	"pith/pkg/config"
	"pith/pkg/telemetry"
	"strings"
	"testing"
)

func TestRunner_RunWithOptions_CombinedArgs(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{StoragePath: tmpDir, MaxLines: 100, HeadLines: 10, TailLines: 10}
	tel, _ := telemetry.NewTelemetry(tmpDir)
	defer tel.Close()

	r := NewRunner(cfg, tel)
	err := r.RunWithOptions([]string{"echo hello"}, true)
	if err != nil {
		t.Fatalf("RunWithOptions failed: %v", err)
	}
}

func TestRunner_RunWithOptions_Error(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{StoragePath: tmpDir, MaxLines: 100}
	tel, _ := telemetry.NewTelemetry(tmpDir)
	defer tel.Close()

	r := NewRunner(cfg, tel)
	err := r.RunWithOptions([]string{"non-existent-command-xyz"}, true)
	if err == nil {
		t.Error("Expected error for non-existent command, got nil")
	}
}

func TestRunner_ApplyMiddleOutTruncation_ExactLines(t *testing.T) {
	r := &Runner{cfg: &config.Config{MaxLines: 10, HeadLines: 5, TailLines: 5}}
	output := "1\n2\n3\n4\n5\n6\n7\n8\n9\n10"
	result := r.ApplyMiddleOutTruncation(output)
	if result != output {
		t.Errorf("Should not have truncated when lines == MaxLines. Got:\n%s", result)
	}

	r.cfg.MaxLines = 5
	r.cfg.HeadLines = 3
	r.cfg.TailLines = 3
	output5 := "1\n2\n3\n4\n5"
	result5 := r.ApplyMiddleOutTruncation(output5)
	if result5 != output5 {
		t.Error("Should not have truncated when head+tail >= len(lines)")
	}
}

func TestRunner_LogForSnag_EmptyPath(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("HOME", tmpHome)

	r := &Runner{cfg: &config.Config{StoragePath: "", SnagLogging: true}}
	r.LogForSnag("test cmd", "test output", 0)

	logPath := filepath.Join(tmpHome, ".pith", "pith.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("Log file not created at default path: %s", logPath)
	}
}

func TestRunner_LogForSnag_NoNewline(t *testing.T) {
	tmpDir := t.TempDir()
	r := &Runner{cfg: &config.Config{StoragePath: tmpDir, SnagLogging: true}}
	r.LogForSnag("test cmd", "output without newline", 0)

	data, _ := os.ReadFile(filepath.Join(tmpDir, "pith.log"))
	if !strings.Contains(string(data), "output without newline\n[EXIT] 0") {
		t.Error("Log entry should ensure newline before [EXIT]")
	}
}

func TestRunner_LogForSnag_OpenError(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "pith.log")
	os.MkdirAll(logPath, 0755)

	r := &Runner{cfg: &config.Config{StoragePath: tmpDir, SnagLogging: true}}
	r.LogForSnag("test cmd", "output", 0)
}

func TestRunner_RunWithOptions_DisabledParser(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := telemetry.NewTelemetry(tmpDir)
	defer tel.Close()
	cfg := &config.Config{
		StoragePath:    tmpDir,
		EnabledParsers: map[string]bool{"git": false},
	}
	r := NewRunner(cfg, tel)

	err := r.RunWithOptions([]string{"echo git status"}, false)
	if err != nil {
		t.Fatalf("RunWithOptions failed: %v", err)
	}
}

func TestRunner_RunWithOptions_GenericError(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := telemetry.NewTelemetry(tmpDir)
	defer tel.Close()
	cfg := &config.Config{StoragePath: tmpDir, SnagLogging: true}
	r := NewRunner(cfg, tel)

	err := r.RunWithOptions([]string{tmpDir}, false)
	if err == nil {
		t.Error("Expected error for executing directory, got nil")
	}
}

func TestRunner_LogForSnag_EmptyOutput(t *testing.T) {
	tmpDir := t.TempDir()
	r := &Runner{cfg: &config.Config{StoragePath: tmpDir, SnagLogging: true}}
	r.LogForSnag("test cmd", "", 0)

	data, _ := os.ReadFile(filepath.Join(tmpDir, "pith.log"))
	if strings.Contains(string(data), "output") {
		t.Error("Log entry should not contain output when it is empty")
	}
}

func TestRunner_RunWithOptions_ChainMatch(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := telemetry.NewTelemetry(tmpDir)
	defer tel.Close()
	cfg := &config.Config{StoragePath: tmpDir, SnagLogging: true}
	r := NewRunner(cfg, tel)

	// Command that should trigger ChainParser and find a sub-parser (git)
	// We allow it to fail if git is not in a repo or returns 255 in test env, just ensure no crash
	_ = r.RunWithOptions([]string{"git version | echo"}, false)
}

func TestRunner_RunWithOptions_EmptyArgs(t *testing.T) {
	r := &Runner{}
	err := r.RunWithOptions([]string{}, false)
	if err == nil {
		t.Error("Expected error for empty args, got nil")
	}
}

func TestRunner_RunWithOptions_NoParts(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := telemetry.NewTelemetry(tmpDir)
	defer tel.Close()
	r := NewRunner(&config.Config{StoragePath: tmpDir, SnagLogging: true}, tel)

	err := r.RunWithOptions([]string{"   "}, false)
	if err == nil {
		t.Error("Expected error for command with only spaces")
	}
}
