package parser

import "testing"

func TestGetAllParsers(t *testing.T) {
	parsers := GetAllParsers()
	if len(parsers) == 0 {
		t.Error("Expected to get some parsers")
	}

	for _, p := range parsers {
		if p.Name() == "" {
			t.Error("Parser should have a non-empty name")
		}
		// Some commands have CanParse = true or false depending on input
		// Just to cover Name() we call it.
	}
}
