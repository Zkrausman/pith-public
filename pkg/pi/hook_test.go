package pi

import (
	"strings"
	"testing"
)

func TestOptimizeHookUsesPithParser(t *testing.T) {
	out := "On branch main\n\nChanges not staged for commit:\n\tmodified: pkg/pi/hook.go\n"
	got := OptimizeHook(HookRequest{Command: "git status", Output: strings.Repeat(out, 200)})
	if got.Parser != "git_status" || got.Passthrough {
		t.Fatalf("expected git parser, got %#v", got)
	}
}
func TestOptimizeHookPreservesFailuresAndRaw(t *testing.T) {
	failure := "token=secret\nERROR: boom\n" + strings.Repeat("x\n", 5000)
	got := OptimizeHook(HookRequest{Command: "go test ./...", Output: failure})
	if !got.Passthrough || !strings.Contains(got.Output, "[REDACTED]") { t.Fatalf("failure must be redacted passthrough: %#v", got) }
	raw := OptimizeHook(HookRequest{Command: "git status", Output: strings.Repeat("x\n", 5000), RawBypass: true})
	if !raw.Passthrough { t.Fatal("raw bypass must pass through") }
}
