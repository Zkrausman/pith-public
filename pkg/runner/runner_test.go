package runner

import (
	"diet/pkg/config"
	"diet/pkg/telemetry"
	"os"
	"path/filepath"
	"testing"
)

func TestEstimateTokens(t *testing.T) {
	s := "12345678"
	if estimateTokens(s) != 2 {
		t.Errorf("Expected 2, got %d", estimateTokens(s))
	}
}

func TestRunner(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "diet-runner-test-*")
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

	// We can't easily test Run() because it executes real commands
	// but we can test the logic around it if we refactor or mock exec.Command
}
