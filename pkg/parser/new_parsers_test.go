package parser

import (
	"testing"
	"strings"
)

func TestSourceParser(t *testing.T) {
	p := &SourceParser{}
	input := `package main
	
	import "fmt"
	
	// This is a comment
	func main() {
		/* Multi-line
		   comment */
		fmt.Println("Hello") // Inline comment
	}`

	output := p.Parse(input)
	if strings.Contains(output, "// This is a comment") {
		t.Errorf("Expected inline comment to be removed")
	}
	if strings.Contains(output, "Multi-line") {
		t.Errorf("Expected multi-line comment to be removed")
	}
	if !strings.Contains(output, "func main()") {
		t.Errorf("Expected function signature to be preserved")
	}
}

func TestGitHubReleaseParser(t *testing.T) {
	p := &GitHubReleaseParser{}
	input := `Tag: v0.5.4
Title: Release v0.5.4
Latest: true

Assets
NAME                    DIGEST                                                                   SIZE     
diet-linux-amd64        sha256:9f6385346102d93580eba389bf24c7ea6c95e9e60b58b0c1ebe2f989863a897a  19.31 MiB
diet-windows-amd64.exe  sha256:5edcbe01c20005a3e24f088bdc7f34535684b29d07bd7ad35d0ed2b78bf11e49  19.71 MiB

View on GitHub: https://github.com/Zkrausman/Diet/releases/tag/v0.5.4`

	output := p.Parse(input)
	if strings.Contains(output, "sha256:") {
		t.Errorf("Expected SHA digests to be removed")
	}
	if !strings.Contains(output, "- diet-linux-amd64 (19.31 MiB)") {
		t.Errorf("Expected asset name and size to be preserved")
	}
}

func TestChainParser(t *testing.T) {
	p := &ChainParser{}
	if !p.CanParse("git status; git log", []string{}) {
		t.Errorf("Expected ChainParser to handle semicolon")
	}
	if !p.CanParse("go test && go build", []string{}) {
		t.Errorf("Expected ChainParser to handle &&")
	}
	
	subcmds := p.SplitSubCommands("git status; git log && ls")
	if len(subcmds) != 3 {
		t.Errorf("Expected 3 sub-commands, got %d", len(subcmds))
	}
}
