package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Install() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)

	if runtime.GOOS == "windows" {
		return installWindows(exeDir)
	}
	return fmt.Errorf("install command not yet supported on %s", runtime.GOOS)
}

func installWindows(exeDir string) error {
	// Use setx to add to user PATH
	// Note: setx is limited to 1024 chars, which is a risk for PATH, but for a simple tool it usually works
	// A better way would be using the registry, but setx is more standard for quick CLIs
	
	path := os.Getenv("PATH")
	if strings.Contains(strings.ToLower(path), strings.ToLower(exeDir)) {
		fmt.Println("Diet directory is already in your PATH.")
		return nil
	}

	fmt.Printf("Adding %s to your PATH...\n", exeDir)
	
	// We use powershell to safely append to PATH
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`
		$oldPath = [Environment]::GetEnvironmentVariable("Path", "User");
		if ($oldPath -notlike "*%s*") {
			[Environment]::SetEnvironmentVariable("Path", $oldPath + ";%s", "User");
			Write-Host "Success";
		}
	`, exeDir, exeDir))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update PATH: %v (output: %s)", err, string(output))
	}

	fmt.Println("Successfully added Diet to your User PATH.")
	fmt.Println("Please restart your terminal for changes to take effect.")
	return nil
}

func SetupGeminiHook() error {
	configDir := ".gemini"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	settingsPath := filepath.Join(configDir, "settings.json")
	if _, err := os.Stat(settingsPath); err == nil {
		fmt.Println(".gemini/settings.json already exists. Skipping.")
		return nil
	}

	exePath, _ := os.Executable()
	// Escape backslashes for JSON
	escapedPath := strings.ReplaceAll(exePath, "\\", "\\\\")

	settings := fmt.Sprintf(`{
  "hooks": {
    "AfterTool": [
      {
        "matcher": "run_shell_command",
        "hooks": [
          {
            "name": "diet-optimizer",
            "type": "command",
            "command": "%s _hook",
            "timeout": 5000
          }
        ]
      }
    ]
  }
}`, escapedPath)

	if err := os.WriteFile(settingsPath, []byte(settings), 0644); err != nil {
		return err
	}

	fmt.Println("Successfully created .gemini/settings.json with Diet hook.")
	return nil
}

func SetupClaudeHook() error {
	configDir := ".claude"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	settingsPath := filepath.Join(configDir, "settings.json")
	if _, err := os.Stat(settingsPath); err == nil {
		fmt.Println(".claude/settings.json already exists. Skipping.")
		return nil
	}

	exePath, _ := os.Executable()
	// Escape backslashes for JSON
	escapedPath := strings.ReplaceAll(exePath, "\\", "\\\\")

	settings := fmt.Sprintf(`{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "name": "diet-optimizer",
            "type": "command",
            "command": "%s _hook",
            "timeout": 5000
          }
        ]
      }
    ]
  }
}`, escapedPath)

	if err := os.WriteFile(settingsPath, []byte(settings), 0644); err != nil {
		return err
	}

	fmt.Println("Successfully created .claude/settings.json with Diet hook.")
	return nil
}
