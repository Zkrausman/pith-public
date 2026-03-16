package parser

import (
	"strings"
	"testing"
)

func TestGrepParser(t *testing.T) {
	p := &GrepParser{}
	input := `pkg/parser/git.go:10:func TestSomething() {
pkg/parser/git.go:15:	fmt.Println("hi")
pkg/runner/runner.go:20:func NewRunner() {
`
	output := p.Parse(input)
	if !strings.Contains(output, "pkg/parser/git.go:") {
		t.Errorf("Output missing file header: %s", output)
	}
	if strings.Count(output, "pkg/parser/git.go:") > 1 {
		t.Error("File header should only appear once per file group")
	}
}

func TestMinifyParser(t *testing.T) {
	p := &MinifyParser{}
	input := `{
  "name": "diet",
  // This is a comment
  "version": "0.3.0",
  "enabled_parsers": {
    "git_status": true
  }
}
`
	// Test CanParse first
	if !p.CanParse("cat", []string{"config.json"}) {
		t.Error("MinifyParser should handle .json files via cat")
	}

	output := p.Parse(input)
	if strings.Contains(output, "//") {
		t.Error("Minified output should not contain comments")
	}
	if !strings.Contains(output, "\n") {
		t.Error("Balanced minification should preserve some newlines for readability")
	}
	if strings.Contains(output, "  ") {
		t.Error("Minified output should have collapsed indentation")
	}
}
