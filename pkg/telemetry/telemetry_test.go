package telemetry

import (
	"os"
	"path/filepath"
	"strings"
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
	if details.OriginalContent != "" || details.CompressedContent != "" {
		t.Errorf("Telemetry must not retain command output: %#v", details)
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

func TestNewTelemetry_Default(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("HOME", tmpHome)

	tel, err := NewTelemetry("")
	if err != nil {
		t.Fatalf("NewTelemetry default failed: %v", err)
	}
	defer tel.Close()

	expectedDB := filepath.Join(tmpHome, ".pith", "pith.db")
	if _, err := os.Stat(expectedDB); os.IsNotExist(err) {
		t.Errorf("DB not created at expected default path: %s", expectedDB)
	}
}

func TestTelemetryFiltering(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := NewTelemetry(tmpDir)
	defer tel.Close()

	tel.Record(ExecutionRecord{Command: "cmd1", Source: "src1", OriginalTokens: 10, CompressedTokens: 5})
	tel.Record(ExecutionRecord{Command: "cmd2", Source: "src2", OriginalTokens: 20, CompressedTokens: 10})
	tel.Record(ExecutionRecord{Command: "cmd3", Source: "", OriginalTokens: 30, CompressedTokens: 15})

	// Test GetStats with empty source
	orig, _, _ := tel.GetStats("")
	if orig != 60 {
		t.Errorf("Expected 60 total tokens, got %d", orig)
	}

	// Test GetStatsByCommand with specific source
	stats, _ := tel.GetStatsByCommand("src1")
	if len(stats) != 1 || stats[0].Command != "cmd1" {
		t.Errorf("Expected 1 stats for src1, got %v", stats)
	}

	// Test GetUnparsedCommands with specific source
	tel.Record(ExecutionRecord{Command: "unparsed1", Source: "src1", IsPassthrough: true})
	unparsed, _ := tel.GetUnparsedCommands("src1")
	if len(unparsed) != 1 || unparsed[0].Pattern != "unparsed1" {
		t.Errorf("Expected 1 unparsed for src1, got %v", unparsed)
	}
}

func TestTelemetryFTSSearch(t *testing.T) {
	tmpDir := t.TempDir()
	tel, err := NewTelemetry(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer tel.Close()

	tel.Record(ExecutionRecord{Command: "npm run build", OriginalContent: "webpack compilation success", CompressedContent: "compilation ok", Source: "term"})
	tel.Record(ExecutionRecord{Command: "git commit -m 'fix bug'", OriginalContent: "git commit output success", CompressedContent: "success", Source: "term"})
	tel.Record(ExecutionRecord{Command: "ls -la", OriginalContent: "list directory", CompressedContent: "dir", Source: "term"})

	// Search stays useful for command discovery without retaining output.
	res, err := tel.SearchExecutions("build", "all", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(res) != 1 || res[0].Command != "npm run build" {
		t.Errorf("Expected npm run build, got %v", res)
	}

	// Test simple command-word match
	res, err = tel.SearchExecutions("npm", "all", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(res) != 1 || res[0].Command != "npm run build" {
		t.Errorf("Expected npm run build, got %v", res)
	}

	// Test boolean operator OR over command names.
	res, err = tel.SearchExecutions("npm OR ls", "all", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(res) != 2 {
		t.Errorf("Expected 2 matches for OR search, got %d", len(res))
	}

	// Test command search with punctuation.
	tel.Record(ExecutionRecord{Command: "git commit -m 'fix bug (error)'", OriginalContent: "special (log)", CompressedContent: "special", Source: "term"})
	res, err = tel.SearchExecutions("fix bug", "all", 10)
	if err != nil {
		t.Fatalf("Fallback search failed: %v", err)
	}
	if len(res) < 1 || !strings.Contains(res[0].Command, "fix bug") {
		t.Errorf("Expected command search to find a matching record, got %v", res)
	}
}

func TestTelemetryJSONL(t *testing.T) {
	tmpDir := t.TempDir()
	tel, err := NewTelemetry(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer tel.Close()

	rec := ExecutionRecord{
		Command:           "grep pattern",
		OriginalTokens:    150,
		CompressedTokens:  75,
		OriginalContent:   "grep output",
		CompressedContent: "grep short",
		DurationMs:        30,
		ParserUsed:        "grep",
		IsPassthrough:     false,
		Source:            "terminal",
	}

	if err := tel.Record(rec); err != nil {
		t.Fatalf("Failed to record: %v", err)
	}

	var buf strings.Builder
	if err := tel.ExportJSONL(&buf); err != nil {
		t.Fatalf("ExportJSONL failed: %v", err)
	}

	tmpDir2 := t.TempDir()
	tel2, err := NewTelemetry(tmpDir2)
	if err != nil {
		t.Fatal(err)
	}
	defer tel2.Close()

	inReader := strings.NewReader(buf.String())
	if err := tel2.ImportJSONL(inReader); err != nil {
		t.Fatalf("ImportJSONL failed: %v", err)
	}

	stats, err := tel2.GetStatsByCommand("terminal")
	if err != nil {
		t.Fatalf("GetStatsByCommand failed: %v", err)
	}
	if len(stats) != 1 || stats[0].Command != "grep pattern" {
		t.Errorf("Expected grep pattern in imported telemetry, got %v", stats)
	}

	inReader2 := strings.NewReader(buf.String())
	if err := tel2.ImportJSONL(inReader2); err != nil {
		t.Fatalf("Re-importing same data failed: %v", err)
	}

	stats2, _ := tel2.GetStatsByCommand("terminal")
	if len(stats2) != 1 {
		t.Errorf("Expected deduplicated count to be 1, got %d", len(stats2))
	}
}
