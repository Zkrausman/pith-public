package runner

import (
	"pith/pkg/config"
	"strings"
	"testing"
)

func TestApplyMiddleOutTruncation_Semantic(t *testing.T) {
	cfg := &config.Config{MaxLines: 5, HeadLines: 2, TailLines: 2}
	r := &Runner{cfg: cfg}

	output := "head1\nhead2\nmiddle1\nmiddle2\nmiddle3\nCRITICAL ERROR HERE\nmiddle5\nmiddle6\nmiddle7\ntail1\ntail2"
	res := r.ApplyMiddleOutTruncation(output)

	if !strings.Contains(res, "CRITICAL ERROR HERE") {
		t.Errorf("expected preserved hot zone, got:\n%s", res)
	}
	if !strings.Contains(res, "non-critical output removed") {
		t.Errorf("expected semantic summary, got:\n%s", res)
	}
}
