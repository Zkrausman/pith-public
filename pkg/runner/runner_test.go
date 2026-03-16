package runner

import (
	"diet/pkg/config"
	"diet/pkg/telemetry"
	"os"
	"path/filepath"
	"strings"
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

	if !strings.Contains(output, "removed by Diet") {
		t.Error("Expected output to be truncated")
	}
	if !strings.HasPrefix(output, "1\n2") {
		t.Error("Expected head to be preserved")
	}
	if !strings.HasSuffix(output, "14\n15") {
		t.Error("Expected tail to be preserved")
	}
}
