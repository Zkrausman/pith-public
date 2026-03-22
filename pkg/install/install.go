package install

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Install() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	pithBinDir := filepath.Join(home, ".pith", "bin")
	if err := os.MkdirAll(pithBinDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	destPath := filepath.Join(pithBinDir, filepath.Base(exePath))
	
	// Skip copy if we are already running from the destination path
	absExe, _ := filepath.Abs(exePath)
	absDest, _ := filepath.Abs(destPath)
	if strings.EqualFold(absExe, absDest) {
		fmt.Printf("Pith is already running from %s. Skipping copy.\n", destPath)
	} else {
		// Copy the executable to ~/.pith/bin
		input, err := os.ReadFile(exePath)
		if err != nil {
			return fmt.Errorf("failed to read executable: %v", err)
		}
		if err := os.WriteFile(destPath, input, 0755); err != nil {
			return fmt.Errorf("failed to copy executable: %v", err)
		}
		fmt.Printf("Copied Pith to %s\n", destPath)
	}

	if runtime.GOOS == "windows" {
		return installWindows(pithBinDir)
	}
	return fmt.Errorf("install command not yet supported on %s", runtime.GOOS)
}

func installWindows(binDir string) error {
	// Use powershell to safely append to User PATH
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`
		$oldPath = [Environment]::GetEnvironmentVariable("Path", "User");
		if ($oldPath -notlike "*%s*") {
			[Environment]::SetEnvironmentVariable("Path", $oldPath + ";%s", "User");
			Write-Host "Success";
		} else {
			Write-Host "Already in PATH";
		}
	`, binDir, binDir))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update PATH: %v (output: %s)", err, string(output))
	}

	fmt.Println("Successfully added Pith to your User PATH.")
	fmt.Println("Please restart your terminal for changes to take effect.")
	return nil
}

type HookEntry struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Command string `json:"command"`
	Timeout int    `json:"timeout"`
}

type HookGroup struct {
	Matcher string      `json:"matcher"`
	Hooks   []HookEntry `json:"hooks"`
}

type Settings struct {
	Hooks map[string][]HookGroup `json:"hooks"`
	Other map[string]interface{} `json:"-"` // Catch-all for other fields
}

// Custom Unmarshal to capture other fields
func (s *Settings) UnmarshalJSON(data []byte) error {
	type Alias Settings
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Capture everything else into s.Other
	var fullMap map[string]interface{}
	if err := json.Unmarshal(data, &fullMap); err != nil {
		return err
	}
	delete(fullMap, "hooks")
	s.Other = fullMap
	return nil
}

// Custom Marshal to combine hooks and other fields
func (s Settings) MarshalJSON() ([]byte, error) {
	fullMap := make(map[string]interface{})
	for k, v := range s.Other {
		fullMap[k] = v
	}
	fullMap["hooks"] = s.Hooks
	return json.Marshal(fullMap)
}

func setupHook(dirName, eventName, matcher string, global bool) error {
	configDir := dirName
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configDir = filepath.Join(home, dirName)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	settingsPath := filepath.Join(configDir, "settings.json")
	var settings Settings

	// Load existing settings if they exist
	if _, err := os.Stat(settingsPath); err == nil {
		data, err := os.ReadFile(settingsPath)
		if err == nil {
			// Create a backup before modifying
			backupPath := settingsPath + ".bak"
			_ = os.WriteFile(backupPath, data, 0644)
			
			_ = json.Unmarshal(data, &settings)
		}
	}

	if settings.Hooks == nil {
		settings.Hooks = make(map[string][]HookGroup)
	}

	home, _ := os.UserHomeDir()
	exePath := filepath.Join(home, ".pith", "bin", "pith.exe")
	if runtime.GOOS != "windows" {
		exePath = filepath.Join(home, ".pith", "bin", "pith")
	}

	pithHook := HookEntry{
		Name:    "pith-optimizer",
		Type:    "command",
		Command: fmt.Sprintf("%s _hook", exePath),
		Timeout: 5000,
	}

	// Find or create the group for this matcher
	eventHooks := settings.Hooks[eventName]
	foundGroup := false
	for i, group := range eventHooks {
		if group.Matcher == matcher {
			// Check if Pith hook already exists in this group
			exists := false
			for _, h := range group.Hooks {
				if h.Name == "pith-optimizer" {
					exists = true
					break
				}
			}
			if !exists {
				eventHooks[i].Hooks = append(eventHooks[i].Hooks, pithHook)
			}
			foundGroup = true
			break
		}
	}

	if !foundGroup {
		eventHooks = append(eventHooks, HookGroup{
			Matcher: matcher,
			Hooks:   []HookEntry{pithHook},
		})
	}
	settings.Hooks[eventName] = eventHooks

	// Write back
	finalData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(settingsPath, finalData, 0644); err != nil {
		return err
	}

	fmt.Printf("Successfully updated %s with Pith hook.\n", settingsPath)
	return nil
}

func SetupGeminiHook(global bool) error {
	matchers := []string{"run_shell_command", "run_command", "send_command_input"}
	for _, matcher := range matchers {
		if err := setupHook(".gemini", "AfterTool", matcher, global); err != nil {
			return err
		}
	}
	return nil
}

func SetupClaudeHook(global bool) error {
	return setupHook(".claude", "PostToolUse", "Bash", global)
}

func SetupCodexHook(global bool) error {
	return setupHook(".codex", "AfterTool", "run_shell_command", global)
}
