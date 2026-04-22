package install

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetupHook_ExistingSettings(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	
	// Create existing settings with some other content
	existing := map[string]interface{}{
		"other_field": "some value",
		"hooks": map[string]interface{}{
			"AfterTool": []interface{}{
				map[string]interface{}{
					"matcher": "run_shell_command",
					"hooks": []interface{}{
						map[string]interface{}{
							"name": "existing-hook",
							"type": "command",
							"command": "echo hello",
						},
					},
				},
			},
		},
	}
	data, _ := json.Marshal(existing)
	os.WriteFile(settingsPath, data, 0644)
	
	path, err := setupHook(tmpDir, "AfterTool", "run_shell_command", false, "test")
	if err != nil {
		t.Fatalf("setupHook failed: %v", err)
	}
	if path != settingsPath {
		t.Errorf("Expected path %s, got %s", settingsPath, path)
	}
	
	// Verify backup
	if _, err := os.Stat(settingsPath + ".bak"); os.IsNotExist(err) {
		t.Error("Backup file was not created")
	}
	
	// Verify new content has both old and new hooks
	newData, _ := os.ReadFile(settingsPath)
	var settings Settings
	json.Unmarshal(newData, &settings)
	
	if settings.Other["other_field"] != "some value" {
		t.Error("Existing fields were lost")
	}
	
	hooks := settings.Hooks["AfterTool"][0].Hooks
	if len(hooks) != 2 {
		t.Errorf("Expected 2 hooks, got %d", len(hooks))
	}
}

func TestSetupHook_OverwriteExistingPith(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	
	// Create settings with an OLD pith hook
	existing := map[string]interface{}{
		"hooks": map[string]interface{}{
			"AfterTool": []interface{}{
				map[string]interface{}{
					"matcher": "run_shell_command",
					"hooks": []interface{}{
						map[string]interface{}{
							"name": "pith-optimizer",
							"type": "command",
							"command": "OLD_COMMAND",
						},
					},
				},
			},
		},
	}
	data, _ := json.Marshal(existing)
	os.WriteFile(settingsPath, data, 0644)
	
	_, err := setupHook(tmpDir, "AfterTool", "run_shell_command", false, "test")
	if err != nil {
		t.Fatalf("setupHook failed: %v", err)
	}
	
	newData, _ := os.ReadFile(settingsPath)
	var settings Settings
	json.Unmarshal(newData, &settings)
	
	hook := settings.Hooks["AfterTool"][0].Hooks[0]
	if hook.Name != "pith-optimizer" || hook.Command == "OLD_COMMAND" {
		t.Errorf("Expected pith-optimizer hook to be overwritten, got command: %s", hook.Command)
	}
}

func TestSettings_MarshalUnmarshal_OtherFields(t *testing.T) {
	data := []byte(`{"hooks":{},"foo":"bar","num":123}`)
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	
	if s.Other["foo"] != "bar" || s.Other["num"] != 123.0 {
		t.Errorf("Other fields not captured correctly: %v", s.Other)
	}
	
	marshaled, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	
	if !strings.Contains(string(marshaled), `"foo":"bar"`) || !strings.Contains(string(marshaled), `"num":123`) {
		t.Errorf("Marshaled JSON missing other fields: %s", string(marshaled))
	}
}

func TestInstall_MkdirError(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("HOME", tmpHome)

	// Create a file where .pith directory should be
	os.WriteFile(filepath.Join(tmpHome, ".pith"), []byte("blocked"), 0644)

	err := Install()
	if err == nil {
		t.Error("Expected error when MkdirAll fails, got nil")
	}
}
