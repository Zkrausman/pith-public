package parser

import (
	"strings"
	"testing"
)

func TestThneedParser(t *testing.T) {
	p := &ThneedParser{}

	// Test CanParse
	if !p.CanParse("thneed", []string{"query", "--text", "pith"}) {
		t.Error("ThneedParser should handle thneed command")
	}
	if !p.CanParse(".\\thneed.exe", []string{"query"}) {
		t.Error("ThneedParser should handle thneed.exe with path")
	}

	// Test Parse with valid JSON results
	input := `[
		{
			"id": "node-1",
			"path": "pkg\\parser\\infra.go",
			"content": "func (t *TestParser) Parse(output string) string {\n\t// Very long content here..."
		},
		{
			"id": "node-2",
			"path": "main.go",
			"content": "package main\n\nimport \"fmt\""
		}
	]`

	output := p.Parse(input)

	if !strings.Contains(output, "Thneed found 2 nodes:") {
		t.Errorf("Expected summary line, got: %s", output)
	}
	if !strings.Contains(output, "[infra.go] (node-1)") {
		t.Error("Expected infra.go details to be preserved")
	}
	if !strings.Contains(output, "package main") {
		t.Error("Expected main.go snippet to be preserved")
	}

	// Test Parse with Thneed error JSON
	errInput := `{"error": "failed to embed query"}`
	errOutput := p.Parse(errInput)
	if !strings.Contains(errOutput, "Thneed Error: failed to embed query") {
		t.Errorf("Expected error message, got: %s", errOutput)
	}
}
