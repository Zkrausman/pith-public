package anomaly

import (
	"strings"
	"testing"
)

func TestDiagnosticSnippetRedactsAndBoundsPayload(t *testing.T) {
	value := "password=super-secret api_key: abc123 " + strings.Repeat("x", 600)
	got := diagnosticSnippet(value)
	if strings.Contains(got, "super-secret") || strings.Contains(got, "abc123") {
		t.Fatalf("secret was not redacted: %q", got)
	}
	if len(got) > 540 || !strings.Contains(got, "[truncated]") {
		t.Fatalf("diagnostic payload was not bounded: length %d", len(got))
	}
}
