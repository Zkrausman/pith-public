package telemetry

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewTelemetry_MkdirError(t *testing.T) {
	// Create a file where we want to create a directory
	tmpDir := t.TempDir()
	badPath := filepath.Join(tmpDir, "blocked")
	os.WriteFile(badPath, []byte("i am a file"), 0644)

	_, err := NewTelemetry(badPath)
	if err == nil {
		t.Error("Expected error when MkdirAll fails on a file path, got nil")
	}
}

func TestTelemetry_GetExecutionDetails_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := NewTelemetry(tmpDir)
	defer tel.Close()

	_, err := tel.GetExecutionDetails(9999)
	if err == nil {
		t.Error("Expected error for non-existent ID, got nil")
	}
}

func TestTelemetry_GetSources_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := NewTelemetry(tmpDir)
	defer tel.Close()

	sources, err := tel.GetSources()
	if err != nil {
		t.Fatalf("GetSources failed: %v", err)
	}
	if len(sources) != 0 {
		t.Errorf("Expected 0 sources, got %d", len(sources))
	}
}

func TestTelemetry_GetStats_NoRows(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := NewTelemetry(tmpDir)
	defer tel.Close()

	orig, comp, err := tel.GetStats("nonexistent")
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	if orig != 0 || comp != 0 {
		t.Errorf("Expected 0/0, got %d/%d", orig, comp)
	}
}

func TestTelemetry_GetStatsByDay_WithSource(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := NewTelemetry(tmpDir)
	defer tel.Close()

	tel.Record(ExecutionRecord{Command: "cmd1", Source: "src1", OriginalTokens: 10})
	tel.Record(ExecutionRecord{Command: "cmd2", Source: "src2", OriginalTokens: 20})

	stats, err := tel.GetStatsByDay("src1")
	if err != nil {
		t.Fatalf("GetStatsByDay failed: %v", err)
	}
	if len(stats) != 1 {
		t.Errorf("Expected 1 day, got %d", len(stats))
	}
	if stats[0].Original != 10 {
		t.Errorf("Expected 10 tokens, got %d", stats[0].Original)
	}
}

func TestTelemetry_GetRecentExecutions_WithSource(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := NewTelemetry(tmpDir)
	defer tel.Close()

	tel.Record(ExecutionRecord{Command: "cmd1", Source: "src1"})
	tel.Record(ExecutionRecord{Command: "cmd2", Source: "src2"})

	recent, err := tel.GetRecentExecutions(10, "src1")
	if err != nil {
		t.Fatalf("GetRecentExecutions failed: %v", err)
	}
	if len(recent) != 1 {
		t.Errorf("Expected 1 recent execution, got %d", len(recent))
	}
	if recent[0].Source != "src1" {
		t.Errorf("Expected source src1, got %s", recent[0].Source)
	}
}

func TestSearchExecutions(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := NewTelemetry(tmpDir)
	defer tel.Close()

	// Seed data
	tel.Record(ExecutionRecord{Command: "git status", OriginalContent: "on branch main", Source: "gemini"})
	tel.Record(ExecutionRecord{Command: "ls -l", OriginalContent: "file.txt", Source: "claude"})
	tel.Record(ExecutionRecord{Command: "cat main.go", OriginalContent: "package main", Source: "gemini"})

	// Search by command
	results, _ := tel.SearchExecutions("git", "all", 10)
	if len(results) != 1 || results[0].Command != "git status" {
		t.Errorf("Expected 1 result for 'git', got %d", len(results))
	}

	// Search by content
	results, _ = tel.SearchExecutions("package", "all", 10)
	if len(results) != 1 || results[0].Command != "cat main.go" {
		t.Errorf("Expected 1 result for 'package', got %d", len(results))
	}

	// Search with source filter
	results, _ = tel.SearchExecutions("l", "claude", 10)
	if len(results) != 1 || results[0].Command != "ls -l" {
		t.Errorf("Expected 1 result for 'l' with claude, got %d", len(results))
	}

	results, _ = tel.SearchExecutions("l", "gemini", 10)
	if len(results) != 0 {
		t.Errorf("Expected 0 results for 'l' with gemini, got %d", len(results))
	}
}
