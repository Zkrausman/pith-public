package parser

import (
	"testing"
)

func TestNPMParser(t *testing.T) {
	p := &NPMParser{}
	if p.Name() != "npm" {
		t.Errorf("Expected npm, got %s", p.Name())
	}
	if p.CanParse("npm", []string{"install"}) {
		t.Log("CanParse npm")
	}

	// Just call Parse
	out := p.Parse(`npm verb something`)
	if out == "" {
		t.Log("Expected some output")
	}
}
