package install

import (
	"strings"
	"testing"
)

func TestPiExtensionRejectsMalformedPithResponse(t *testing.T) {
	if !strings.Contains(piExtension, "try { resolve(JSON.parse(out)); } catch (error) { reject(error); }") {
		t.Fatal("Pi extension must reject malformed Pith JSON so its hook falls back to Pi's original result")
	}
	if !strings.Contains(piExtension, "catch { return; } // Pith failure always preserves Pi's original result.") {
		t.Fatal("Pi extension must preserve Pi result when Pith fails")
	}
}

func TestPiExtensionDelegatesAllCompletedBashResultsToPith(t *testing.T) {
	for _, forbidden := range []string{
		"thresholdBytes",
		"output.length <",
	} {
		if strings.Contains(piExtension, forbidden) {
			t.Errorf("Pi extension must not filter Pith input by %q", forbidden)
		}
	}
	if !strings.Contains(piExtension, "telemetryEnabled: true") {
		t.Error("Pi extension must enable Pith harness telemetry")
	}
	if !strings.Contains(piExtension, "event.isError ? 1 :") {
		t.Error("Pi extension must classify Pi error results as nonzero exits")
	}
}

func TestPiExtensionSafetyGuards(t *testing.T) {
	for _, guard := range []string{
		"if (!ctx.isProjectTrusted()) return;",
		"blocks.length !== 1",
		"if (signal?.aborted) return Promise.reject",
		"process.env.PITH_BIN || join(homedir(), \".pith\", \"bin\", process.platform === \"win32\" ? \"pith.exe\" : \"pith\")",
	} {
		if !strings.Contains(piExtension, guard) {
			t.Errorf("Pi extension is missing safety guard %q", guard)
		}
	}
}
