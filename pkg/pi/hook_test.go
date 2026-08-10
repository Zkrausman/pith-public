package pi

import (
	"strings"
	"testing"

	"pith/pkg/telemetry"
)

func TestOptimizeHookUsesPithParser(t *testing.T) {
	out := "On branch main\n\nChanges not staged for commit:\n\tmodified: pkg/pi/hook.go\n"
	got := OptimizeHook(HookRequest{Command: "git status", Output: strings.Repeat(out, 200)})
	if got.Parser != "git_status" || got.Passthrough {
		t.Fatalf("expected git parser, got %#v", got)
	}
}
func TestOptimizeHookTelemetryIsAggregateOnly(t *testing.T) {
	dir := t.TempDir()
	OptimizeHook(HookRequest{Command: "git status", Output: strings.Repeat("secret-output\\n", 1000), TelemetryEnabled: true, StoragePath: dir})
	tel, err := telemetry.NewTelemetry(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer tel.Close()
	records, err := tel.GetRecentExecutions(1, "")
	if err != nil || len(records) != 1 {
		t.Fatalf("records: %v %#v", err, records)
	}
	detail, err := tel.GetExecutionDetails(records[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if detail.OriginalContent != "" || detail.CompressedContent != "" {
		t.Fatalf("Pi telemetry retained content: %#v", detail)
	}
}

func TestOptimizeHookPreservesFailuresAndRaw(t *testing.T) {
	failure := "token=secret\nERROR: boom\n" + strings.Repeat("x\n", 5000)
	got := OptimizeHook(HookRequest{Command: "go test ./...", Output: failure})
	if !got.Passthrough || !strings.Contains(got.Output, "[REDACTED]") {
		t.Fatalf("failure must be redacted passthrough: %#v", got)
	}
	raw := OptimizeHook(HookRequest{Command: "git status", Output: strings.Repeat("x\n", 5000), RawBypass: true})
	if !raw.Passthrough {
		t.Fatal("raw bypass must pass through")
	}
}
