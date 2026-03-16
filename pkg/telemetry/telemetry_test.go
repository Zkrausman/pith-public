package telemetry

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTelemetry(t *testing.T) {
	// Create a temporary DB for testing
	tmpDir, err := os.MkdirTemp("", "diet-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	
	// We need to manually initialize if we don't want to use NewTelemetry which uses HomeDir
	tel, err := NewTelemetryWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create telemetry: %v", err)
	}
	defer tel.Close()

	rec := ExecutionRecord{
		Command:          "git status",
		OriginalTokens:   100,
		CompressedTokens: 50,
		DurationMs:       10,
		ParserUsed:       "git_status",
		IsPassthrough:    false,
	}

	if err := tel.Record(rec); err != nil {
		t.Fatalf("Failed to record: %v", err)
	}

	orig, comp, err := tel.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}
	if orig != 100 || comp != 50 {
		t.Errorf("Expected 100/50, got %d/%d", orig, comp)
	}

	byCmd, err := tel.GetStatsByCommand()
	if err != nil || len(byCmd) == 0 {
		t.Fatalf("Failed to get stats by command")
	}
	if byCmd[0].Command != "git status" {
		t.Errorf("Expected git status, got %s", byCmd[0].Command)
	}

	unparsed, err := tel.GetUnparsedCommands()
	if err != nil {
		t.Fatal(err)
	}
	if len(unparsed) != 0 {
		t.Error("Expected 0 unparsed commands")
	}
}
