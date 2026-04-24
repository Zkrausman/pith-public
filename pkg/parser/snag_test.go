package parser

import (
	"strings"
	"testing"
)

func TestSnagParser(t *testing.T) {
	p := &SnagParser{}

	// Test CanParse
	if !p.CanParse("snag", []string{"list", "--json"}) {
		t.Error("SnagParser should handle snag command")
	}
	if !p.CanParse("./snag.exe", []string{"list"}) {
		t.Error("SnagParser should handle snag.exe with path")
	}

	// Test Parse with valid JSON
	input := `[
		{
			"id": "snag-1",
			"status": "pending",
			"failed_command": "git commit",
			"error_message": "fatal: pathspec '...' did not match any files",
			"advice": "Check your path and try again",
			"context_snippet": "THIS IS A VERY LARGE CONTEXT SNIPPET THAT SHOULD BE STRIPPED"
		},
		{
			"id": "snag-2",
			"status": "resolved",
			"failed_command": "go test ./pkg/...",
			"error_message": "FAIL: TestOne",
			"advice": "Fix the test logic",
			"context_snippet": "ANOTHER LARGE SNIPPET"
		}
	]`

	output := p.Parse(input)
	
	if !strings.Contains(output, "Snag found 2 issues:") {
		t.Errorf("Expected summary line, got: %s", output)
	}
	if !strings.Contains(output, "[pending] snag-1: git commit") {
		t.Error("Expected snag-1 details to be preserved")
	}
	if !strings.Contains(output, "Advice: Check your path and try again") {
		t.Error("Expected advice to be preserved")
	}
	if strings.Contains(output, "THIS IS A VERY LARGE CONTEXT SNIPPET") {
		t.Error("Expected context_snippet to be stripped")
	}

	// Test Parse with empty list
	if p.Parse("[]") != "Snag: No snags found." {
		t.Error("Expected empty message for []")
	}

	// Test Parse with invalid JSON
	if !strings.Contains(p.Parse("invalid json"), "invalid json") {
		t.Error("Expected plain text fallback for invalid JSON")
	}
}
