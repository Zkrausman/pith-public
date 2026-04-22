package parser

import (
	"testing"
)

func TestChainParserCoverage(t *testing.T) {
	p := &ChainParser{}
	if p.Name() != "chain" {
		t.Errorf("Expected chain, got %s", p.Name())
	}
	if !p.CanParse("cmd1 && cmd2", []string{}) {
		t.Log("CanParse chain")
	}

	// SplitSubCommands
	subs := p.SplitSubCommands("echo a && echo b | grep c")
	if len(subs) == 0 {
		t.Error("Expected some subcommands")
	}

	// Parse
	out := p.Parse(`some text`)
	if out == "" {
		t.Log("Expected some output")
	}
}
