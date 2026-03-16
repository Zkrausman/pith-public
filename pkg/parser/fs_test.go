package parser

import (
	"strings"
	"testing"
)

func TestLsParser(t *testing.T) {
	p := &LsParser{}
	input := `total 8
-rw-r--r--  1 user  group  123 Mar 15 22:30 main.go
-rw-r--r--  1 user  group  456 Mar 15 22:31 README.md
`
	output := p.Parse(input)
	if !strings.Contains(output, "main.go") || !strings.Contains(output, "README.md") {
		t.Errorf("Expected filenames in output, got: %s", output)
	}
	if strings.Contains(output, "total") || strings.Contains(output, "user") {
		t.Errorf("Output contains metadata: %s", output)
	}
}

func TestFindParser(t *testing.T) {
	p := &FindParser{}
	input := ""
	for i := 0; i < 60; i++ {
		input += "file" + strings.Repeat("a", i) + ".go\n"
	}
	output := p.Parse(input)
	if !strings.Contains(output, "(+ 10 more results)") {
		t.Errorf("Expected truncation message, got output of length %d", len(strings.Split(output, "\n")))
	}
}

func TestTreeParser(t *testing.T) {
	p := &TreeParser{}
	input := `
.
├── pkg
│   └── parser
│       └── fs.go
└── main.go
`
	output := p.Parse(input)
	if strings.Contains(output, "├──") || strings.Contains(output, "│") {
		t.Errorf("Output still contains tree characters: %s", output)
	}
	if !strings.Contains(output, "pkg") || !strings.Contains(output, "fs.go") {
		t.Errorf("Output missing directory structure: %s", output)
	}
}

func TestDuParser(t *testing.T) {
	p := &DuParser{}
	input := `
123	./pkg/parser
456	./pkg/config
789	.
`
	output := p.Parse(input)
	if !strings.Contains(output, "123	./pkg/parser") {
		t.Errorf("Unexpected du output format: %s", output)
	}
}
