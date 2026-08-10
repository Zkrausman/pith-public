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
