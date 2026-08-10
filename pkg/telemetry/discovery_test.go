package telemetry

import (
	"strings"
	"testing"
)

func TestGetDiscoveryFamiliesNormalizesAndExcludesNoise(t *testing.T) {
	tel, err := NewTelemetry(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer tel.Close()

	records := []ExecutionRecord{
		{Command: "git status --short", OriginalTokens: 100, IsPassthrough: true, Harness: "pi"},
		{Command: "git status --porcelain", OriginalTokens: 150, IsPassthrough: true, Harness: "pi"},
		{Command: "npm pack package-a", OriginalTokens: 200, IsPassthrough: true, Harness: "pi"},
		{Command: "set -e\nprintf '%s' secret\nnpm pack package-b", OriginalTokens: 300, IsPassthrough: true, Harness: "pi"},
		{Command: "tmp=$(mktemp); git status", OriginalTokens: 50, IsPassthrough: true, Harness: "pi"},
		{Command: "pith discover", OriginalTokens: 400, IsPassthrough: true, Harness: "pi"},
		{Command: "go test ./...", OriginalTokens: 500, IsPassthrough: true, Harness: "pi"},
		{Command: "go run main.go discover", OriginalTokens: 600, IsPassthrough: true, Harness: "pi"},
	}
	for _, record := range records {
		if err := tel.Record(record); err != nil {
			t.Fatal(err)
		}
	}

	families, err := tel.GetDiscoveryFamilies("")
	if err != nil {
		t.Fatal(err)
	}
	if len(families) != 3 {
		t.Fatalf("expected three command families, got %#v", families)
	}
	if families[0].Family != "shell script" || families[0].InvocationCount != 2 || families[0].TotalRawTokens != 350 {
		t.Fatalf("expected shell script family first, got %#v", families[0])
	}
	if families[1].Family != "git status" || families[1].InvocationCount != 2 || families[1].TotalRawTokens != 250 {
		t.Fatalf("expected aggregated git status family, got %#v", families[1])
	}
	if families[2].Family != "npm pack" || families[2].TotalRawTokens != 200 {
		t.Fatalf("expected npm pack family, got %#v", families[2])
	}
}

func TestGetDiscoveryDetailsBoundsAndSanitizesCommands(t *testing.T) {
	tel, err := NewTelemetry(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer tel.Close()
	long := "custom-tool " + strings.Repeat("x", 200)
	if err := tel.Record(ExecutionRecord{Command: "custom-tool first\nsecond", OriginalTokens: 10, IsPassthrough: true, Harness: "pi"}); err != nil {
		t.Fatal(err)
	}
	if err := tel.Record(ExecutionRecord{Command: long, OriginalTokens: 20, IsPassthrough: true, Harness: "pi"}); err != nil {
		t.Fatal(err)
	}
	details, err := tel.GetDiscoveryDetails("custom-tool", "", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(details) != 1 || len([]rune(details[0].Command)) > 160 {
		t.Fatalf("expected one bounded detail, got %#v", details)
	}
}
