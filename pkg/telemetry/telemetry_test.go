package telemetry

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTelemetryFull(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pith-telemetry-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "pith.db")
	tel, err := NewTelemetryWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create telemetry: %v", err)
	}
	defer tel.Close()

	rec1 := ExecutionRecord{
		Command:           "ls -la",
		OriginalTokens:    200,
		CompressedTokens:  100,
		OriginalContent:   "long list",
		CompressedContent: "short list",
		DurationMs:        50,
		ParserUsed:        "ls",
		IsPassthrough:     false,
		Source:            "terminal",
	}

	rec2 := ExecutionRecord{
		Command:           "git status",
		OriginalTokens:    500,
		CompressedTokens:  250,
		OriginalContent:   "git status output",
		CompressedContent: "git status summary",
		DurationMs:        100,
		ParserUsed:        "git",
		IsPassthrough:     false,
		Source:            "gemini",
	}

	if err := tel.Record(rec1); err != nil {
		t.Fatalf("Failed to record rec1: %v", err)
	}
	// Sleep to ensure distinct timestamps for reliable DESC order
	time.Sleep(1100 * time.Millisecond) 
	if err := tel.Record(rec2); err != nil {
		t.Fatalf("Failed to record rec2: %v", err)
	}

	// Test GetStats
	orig, comp, err := tel.GetStats("terminal")
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	if orig != 200 || comp != 100 {
		t.Errorf("Expected 200/100 for terminal, got %d/%d", orig, comp)
	}

	origAll, compAll, _ := tel.GetStats("all")
	if origAll != 700 || compAll != 350 {
		t.Errorf("Expected 700/350 for all, got %d/%d", origAll, compAll)
	}

	// Test GetStatsByDay
	statsByDay, err := tel.GetStatsByDay("all")
	if err != nil {
		t.Fatalf("GetStatsByDay failed: %v", err)
	}
	if len(statsByDay) == 0 {
		t.Fatal("Expected at least one day of stats")
	}
	// Just check if date is in YYYY-MM-DD format
	if len(statsByDay[0].Date) != 10 {
		t.Errorf("Unexpected date format: %s", statsByDay[0].Date)
	}

	// Test GetStatsByCommand
	statsByCmd, err := tel.GetStatsByCommand("all")
	if err != nil {
		t.Fatalf("GetStatsByCommand failed: %v", err)
	}
	if len(statsByCmd) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(statsByCmd))
	}

	// Test GetSources
	sources, err := tel.GetSources()
	if err != nil {
		t.Fatalf("GetSources failed: %v", err)
	}
	if len(sources) != 2 {
		t.Errorf("Expected 2 sources, got %d: %v", len(sources), sources)
	}

	// Test GetRecentExecutions
	recent, err := tel.GetRecentExecutions(10, "all")
	if err != nil {
		t.Fatalf("GetRecentExecutions failed: %v", err)
	}
	if len(recent) != 2 {
		t.Errorf("Expected 2 recent executions, got %d", len(recent))
	}
	if recent[0].Command != "git status" { // Should be most recent
		t.Errorf("Expected most recent command to be git status, got %s", recent[0].Command)
	}

	// Test GetExecutionDetails
	details, err := tel.GetExecutionDetails(recent[0].ID)
	if err != nil {
		t.Fatalf("GetExecutionDetails failed: %v", err)
	}
	if details.Command != "git status" {
		t.Errorf("Expected command 'git status', got '%s'", details.Command)
	}

	// Test ResetPassthrough with a passthrough record
	passthroughRec := ExecutionRecord{
		Command:       "custom-command",
		IsPassthrough: true,
		Source:        "unknown",
	}
	tel.Record(passthroughRec)
	unparsed, _ := tel.GetUnparsedCommands("all")
	if len(unparsed) != 1 {
		t.Errorf("Expected 1 unparsed command, got %d", len(unparsed))
	}
	tel.ResetPassthrough()
	unparsed, _ = tel.GetUnparsedCommands("all")
	if len(unparsed) != 0 {
		t.Error("ResetPassthrough did not clear unparsed commands")
	}

	// Test ResetAll
	tel.ResetAll()
	origReset, _, _ := tel.GetStats("all")
	if origReset != 0 {
		t.Error("ResetAll did not clear all records")
	}
}

func TestNewTelemetry(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pith-new-tel-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tel, err := NewTelemetry(tmpDir)
	if err != nil {
		t.Fatalf("NewTelemetry failed: %v", err)
	}
	defer tel.Close()

	if _, err := os.Stat(filepath.Join(tmpDir, "pith.db")); os.IsNotExist(err) {
		t.Error("pith.db was not created")
	}
}
