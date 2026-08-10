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
	OptimizeHook(HookRequest{Command: "git status --token=raw-secret", Output: strings.Repeat("secret-output\\n", 1000), StoragePath: dir})
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
	if detail.Command != "git status --token=[REDACTED]" || strings.Contains(detail.Command, "raw-secret") {
		t.Fatalf("Pi telemetry did not retain a redacted command: %#v", detail)
	}
	if detail.OriginalContent != "" || detail.CompressedContent != "" {
		t.Fatalf("Pi telemetry retained content: %#v", detail)
	}
}

func TestOptimizeHookPreservesStructuredJSON(t *testing.T) {
	output := `{"state":"OPEN","title":"fix: restore reliable hook telemetry"}`
	got := OptimizeHook(HookRequest{Command: "gh pr view --json state,title", Output: output})
	if !got.Passthrough || got.Parser != "" || got.Output != output {
		t.Fatalf("structured output must remain lossless, got %#v", got)
	}
}

func TestOptimizeHookRecordsUnparsedCommandsForDiscovery(t *testing.T) {
	dir := t.TempDir()
	OptimizeHook(HookRequest{Command: "custom-tool --token=secret", Output: "verbose output", StoragePath: dir})
	tel, err := telemetry.NewTelemetry(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer tel.Close()
	unparsed, err := tel.GetUnparsedCommands("")
	if err != nil || len(unparsed) != 1 {
		t.Fatalf("unparsed telemetry: %v %#v", err, unparsed)
	}
	if unparsed[0].Pattern != "custom-tool --token=[REDACTED]" || unparsed[0].Harness != HarnessPi {
		t.Fatalf("unexpected discovery record: %#v", unparsed[0])
	}
}

func TestOptimizeHookHonorsEnabledParsers(t *testing.T) {
	out := strings.Repeat("On branch main\n\nChanges not staged for commit:\n\tmodified: pkg/pi/hook.go\n", 200)
	got := OptimizeHook(HookRequest{
		Command:        "git status",
		Output:         out,
		EnabledParsers: map[string]bool{"git_status": false},
	})
	if !got.Passthrough || got.Parser != "" || got.Output != out {
		t.Fatalf("disabled parser must preserve output, got %#v", got)
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
