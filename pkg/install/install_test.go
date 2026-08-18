package install

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSettingsMarshaling(t *testing.T) {
	data := []byte(`{
		"hooks": {
			"Event": [
				{
					"matcher": "Match",
					"hooks": [
						{"name": "pith", "type": "command", "command": "cmd", "timeout": 1000}
					]
				}
			]
		},
		"other_field": "value",
		"nested": {"key": 123}
	}`)

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if s.Other["other_field"] != "value" {
		t.Errorf("Expected other_field value, got %v", s.Other["other_field"])
	}

	// Marshal back
	marshaled, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var roundtrip map[string]interface{}
	if err := json.Unmarshal(marshaled, &roundtrip); err != nil {
		t.Fatal(err)
	}

	if roundtrip["other_field"] != "value" {
		t.Errorf("Roundtrip failed for other_field")
	}
	if roundtrip["hooks"] == nil {
		t.Error("Roundtrip missing hooks")
	}
}

func TestSetupHooksComplex(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pith-install-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	origCwd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origCwd)

	// Create existing settings.json with other content
	os.MkdirAll(".gemini", 0755)
	existing := []byte(`{"other": "data"}`)
	os.WriteFile(".gemini/settings.json", existing, 0644)

	// Test Gemini Hook
	if err := SetupGeminiHook(false); err != nil {
		t.Errorf("SetupGeminiHook failed: %v", err)
	}

	// Verify backup
	if _, err := os.Stat(".gemini/settings.json.bak"); os.IsNotExist(err) {
		t.Error("Backup settings.json.bak not created")
	}

	// Verify content
	data, _ := os.ReadFile(".gemini/settings.json")
	var settings map[string]interface{}
	json.Unmarshal(data, &settings)
	if settings["other"] != "data" {
		t.Error("Original data lost")
	}

	// Test Claude and Codex
	SetupClaudeHook(false)
	SetupCodexHook(false)

	if _, err := os.Stat(".claude/settings.json"); os.IsNotExist(err) {
		t.Error("Claude settings not created")
	}
	if _, err := os.Stat(".codex/settings.json"); os.IsNotExist(err) {
		t.Error("Codex settings not created")
	}
}

func TestSetupHookGlobal(t *testing.T) {
	tmpHome, err := os.MkdirTemp("", "pith-home-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("HOME", tmpHome)

	if err := SetupGeminiHook(true); err != nil {
		t.Fatal(err)
	}

	expectedPath := filepath.Join(tmpHome, ".gemini", "settings.json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Global hook not created at %s", expectedPath)
	}
}
func TestRenamedGoTestBinaryCannotInstall(t *testing.T) {
	binaryName := "pith"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(t.TempDir(), binaryName)

	build := exec.Command("go", "test", "-c", "-o", binaryPath, ".")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build renamed test binary: %v\n%s", err, output)
	}

	run := exec.Command(binaryPath, "-test.run", "^TestRenamedGoTestBinaryInstallHelper$")
	run.Env = append(os.Environ(), "PITH_RENAMED_TEST_BINARY_HELPER=1")
	if output, err := run.CombinedOutput(); err != nil {
		t.Fatalf("renamed test binary installed successfully or failed unexpectedly: %v\n%s", err, output)
	}
}

func TestRenamedGoTestBinaryInstallHelper(t *testing.T) {
	if os.Getenv("PITH_RENAMED_TEST_BINARY_HELPER") != "1" {
		return
	}
	if err := Install(); err == nil || !strings.Contains(err.Error(), "Go test executable") {
		t.Fatalf("expected renamed test binary to be rejected, got %v", err)
	}
}

func TestInstall(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("HOME", tmpHome)

	err := Install()
	if err == nil || !strings.Contains(err.Error(), "Go test executable") {
		t.Fatalf("expected test install to be rejected, got %v", err)
	}

	pithBinDir := filepath.Join(tmpHome, ".pith", "bin")
	if _, err := os.Stat(filepath.Join(pithBinDir, installedBinaryName())); !os.IsNotExist(err) {
		t.Errorf("test install created binary: %v", err)
	}
}

func TestSettingsEmpty(t *testing.T) {
	var s Settings
	data, _ := json.Marshal(s)
	if !strings.Contains(string(data), `"hooks":null`) {
		t.Errorf("Expected null hooks in empty settings, got %s", string(data))
	}
}

func TestSetupHook_Malformed(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, ".gemini"), 0755)
	os.WriteFile(filepath.Join(tmpDir, ".gemini", "settings.json"), []byte("{invalid}"), 0644)

	t.Setenv("USERPROFILE", tmpDir)
	t.Setenv("HOME", tmpDir)

	// Should still work (overwrite or backup)
	_, err := setupHook(".gemini", "AfterTool", "run_shell_command", true, "gemini")
	if err != nil {
		t.Errorf("setupHook failed with malformed json: %v", err)
	}
}

func TestSetupHook_Existing(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("USERPROFILE", tmpDir)
	t.Setenv("HOME", tmpDir)

	// First install
	setupHook(".gemini", "AfterTool", "run_shell_command", true, "gemini")

	// Second install (should hit the "exists" branch)
	_, err := setupHook(".gemini", "AfterTool", "run_shell_command", true, "gemini")
	if err != nil {
		t.Errorf("setupHook failed on existing: %v", err)
	}
}

func TestSetupAntigravityHook_GlobalAndLocal(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("HOME", tmpHome)

	origCwd, _ := os.Getwd()
	tmpWorkspace := t.TempDir()
	os.Chdir(tmpWorkspace)
	defer os.Chdir(origCwd)

	// 1. Test local workspace setup with existing custom hook
	os.MkdirAll(".agents", 0755)
	existingContent := []byte(`{"other-hook": {"enabled": true}}`)
	os.WriteFile(filepath.Join(".agents", "hooks.json"), existingContent, 0644)

	if err := SetupAntigravityHook(false); err != nil {
		t.Fatalf("SetupAntigravityHook(false) failed: %v", err)
	}

	// Verify backup
	if _, err := os.Stat(filepath.Join(".agents", "hooks.json.bak")); os.IsNotExist(err) {
		t.Error("Backup .agents/hooks.json.bak was not created")
	}

	// Verify content preserved and pith-optimizer added
	localData, err := os.ReadFile(filepath.Join(".agents", "hooks.json"))
	if err != nil {
		t.Fatal(err)
	}
	var localMap map[string]interface{}
	if err := json.Unmarshal(localData, &localMap); err != nil {
		t.Fatal(err)
	}
	if localMap["other-hook"] == nil {
		t.Error("Existing other-hook was lost in .agents/hooks.json")
	}
	if localMap["pith-optimizer"] == nil {
		t.Error("pith-optimizer was not added to .agents/hooks.json")
	}

	// 2. Test global setup
	if err := SetupAntigravityHook(true); err != nil {
		t.Fatalf("SetupAntigravityHook(true) failed: %v", err)
	}
	globalPath := filepath.Join(tmpHome, ".gemini", "config", "hooks.json")
	globalData, err := os.ReadFile(globalPath)
	if err != nil {
		t.Fatalf("Failed to read global hooks.json: %v", err)
	}
	var globalMap map[string]interface{}
	if err := json.Unmarshal(globalData, &globalMap); err != nil {
		t.Fatal(err)
	}
	if globalMap["pith-optimizer"] == nil {
		t.Error("pith-optimizer was not added to global hooks.json")
	}
}

