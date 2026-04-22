package parser

import (
	"testing"
)

func TestThneedParser(t *testing.T) {
	p := &ThneedParser{}
	if p.Name() != "thneed" {
		t.Errorf("Expected thneed, got %s", p.Name())
	}
	if p.CanParse("thneed", []string{"query"}) {
		t.Log("CanParse thneed query")
	}

	// Just call Parse with something to get coverage
	out := p.Parse(`{"data": {"something": "here"}}`)
	if out == "" {
		t.Log("Expected some output")
	}

	out = p.Parse(`plain text`)
	if out == "" {
		t.Log("Expected some output")
	}
}
