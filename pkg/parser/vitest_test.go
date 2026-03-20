package parser

import (
	"strings"
	"testing"
)

func TestVitestParser(t *testing.T) {
	p := &VitestParser{}

	// Test CanParse
	if !p.CanParse("npx", []string{"vitest", "run"}) {
		t.Errorf("Expected CanParse to handle npx vitest run")
	}
	if !p.CanParse("vitest", []string{"run"}) {
		t.Errorf("Expected CanParse to handle vitest run")
	}
	if p.CanParse("ls", []string{}) {
		t.Errorf("Expected CanParse to reject ls")
	}

	// Mock vitest output
	input := `
 ❯ server.test.js (7 tests | 3 failed) 118ms
     × POST /api/issues/:repoId creates an issue 9ms
     × PATCH /api/issues/:repoId/:issueId updates an issue 8ms
     × DELETE /api/issues/:repoId/:issueId deletes an issue 11ms

⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯ Failed Tests 3
⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯

 FAIL  server.test.js > Backend API > POST /api/issues/:repoId creates an issue
TypeError: exec.mockImplementation is not a function
 ❯ server.test.js:84:10
     82|     const repoId = repoRes.body.id;
     83|
     84|     exec.mockImplementation((cmd, options, callback) => {
       |          ^
     85|       callback(null, { stdout: 'NEW-ISSUE-ID\n' });
     86|     });

⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯
⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯[1/3]⎯

 Test Files  1 failed (1)
      Tests  3 failed | 4 passed (7)
   Start at  01:10:41
   Duration  2.30s (transform 34ms, setup 0ms, import 1.32s, tests 118ms, environment 0ms)
`

	output := p.Parse(input)

	if strings.Contains(output, "⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯") {
		t.Errorf("Expected decorative lines to be removed")
	}
	if !strings.Contains(output, "Test Files  1 failed (1)") {
		t.Errorf("Expected summary stats to be preserved")
	}
	if !strings.Contains(output, "FAIL  server.test.js") {
		t.Errorf("Expected failure header to be preserved")
	}
	if !strings.Contains(output, "TypeError: exec.mockImplementation is not a function") {
		t.Errorf("Expected error message to be preserved")
	}
	if !strings.Contains(output, "❯ server.test.js (7 tests | 3 failed) 118ms") {
		t.Errorf("Expected individual test file summary to be preserved")
	}
}

func TestBDParser(t *testing.T) {
	p := &BDParser{}
	
	if !p.CanParse("bd", []string{"ready"}) {
		t.Errorf("Expected CanParse to handle bd ready")
	}

	helpInput := `Issues chained together like beads. A lightweight issue tracker with first-class dependency support.

Usage:
  bd [flags]
  bd [command]

Working With Issues:
  children        List child beads of a parent
  close           Close one or more issues
  comments        View or manage comments on an issue
  create          Create a new issue (or multiple issues from markdown file)
  delete          Delete one or more issues
  edit            Edit an issue
  list            List issues
  show            Show issue details
  update          Update an issue
  status          Show status
  ready           Show ready work
  dep             Manage dependencies
  graph           Show graph
  search          Search issues
  label           Manage labels
  config          Manage config
  version         Show version
  help            Show help
  bootstrap       Bootstrap database
  init            Initialize bd
  doctor          Check health
  migrate         Migrate database
  sql             Execute SQL
  backup          Backup database
  restore         Restore database
... many many lines ...
`
	output := p.Parse(helpInput)
	if !strings.Contains(output, "(truncated bd help)") {
		t.Errorf("Expected bd help to be truncated")
	}
}

func TestPromptfooParser(t *testing.T) {
	p := &PromptfooParser{}

	if !p.CanParse("npx", []string{"promptfoo", "eval"}) {
		t.Errorf("Expected CanParse to handle npx promptfoo eval")
	}

	input := `
Evaluating [████████████████████████████████████████] 100% | 1/1 (errors: 1)
┌────────────────┬────────────────┬────────────────┐
│ command        │ output         │ task           │
├────────────────┼────────────────┼────────────────┤
│ git status     │ On branch main │ List files     │
└────────────────┴────────────────┴────────────────┘
`
	output := p.Parse(input)
	if strings.Contains(output, "┌──") {
		t.Errorf("Expected table borders to be removed")
	}
	if !strings.Contains(output, "git status | On branch main | List files") {
		t.Errorf("Expected table content to be preserved and cleaned")
	}
	if !strings.Contains(output, "100%") {
		t.Errorf("Expected progress bar to be preserved")
	}
}
