package install

import (
	"os"
	"testing"
)

func TestSetupHooks(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pith-install-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Change working directory to tmpDir for the test
	origCwd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origCwd)

	// Test Gemini Hook
	if err := SetupGeminiHook(false); err != nil {
		t.Errorf("SetupGeminiHook failed: %v", err)
	}
	if _, err := os.Stat(".gemini/settings.json"); os.IsNotExist(err) {
		t.Error("Gemini settings.json not created")
	}

	// Test Claude Hook
	if err := SetupClaudeHook(false); err != nil {
		t.Errorf("SetupClaudeHook failed: %v", err)
	}
	if _, err := os.Stat(".claude/settings.json"); os.IsNotExist(err) {
		t.Error("Claude settings.json not created")
	}

	// Test skip if exists
	if err := SetupGeminiHook(false); err != nil {
		t.Error(err)
	}
}
