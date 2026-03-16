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
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dietBinDir := filepath.Join(home, ".diet", "bin")
	if err := os.MkdirAll(dietBinDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	destPath := filepath.Join(dietBinDir, filepath.Base(exePath))
	
	// Copy the executable to ~/.diet/bin
	input, err := os.ReadFile(exePath)
	if err != nil {
		return fmt.Errorf("failed to read executable: %v", err)
	}
	if err := os.WriteFile(destPath, input, 0755); err != nil {
		return fmt.Errorf("failed to copy executable: %v", err)
	}

	fmt.Printf("Copied Diet to %s\n", destPath)

	if runtime.GOOS == "windows" {
		return installWindows(dietBinDir)
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

	fmt.Println("Successfully added Diet to your User PATH.")
	fmt.Println("Please restart your terminal for changes to take effect.")
	return nil
}

func SetupGeminiHook(global bool) error {
	configDir := ".gemini"
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configDir = filepath.Join(home, ".gemini")
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	settingsPath := filepath.Join(configDir, "settings.json")
	if _, err := os.Stat(settingsPath); err == nil {
		fmt.Printf("%s already exists. Skipping.\n", settingsPath)
		return nil
	}

	home, _ := os.UserHomeDir()
	exePath := filepath.Join(home, ".diet", "bin", "diet.exe")
	if runtime.GOOS != "windows" {
		exePath = filepath.Join(home, ".diet", "bin", "diet")
	}
	
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

	fmt.Printf("Successfully created %s with Diet hook.\n", settingsPath)
	return nil
}

func SetupClaudeHook(global bool) error {
	configDir := ".claude"
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configDir = filepath.Join(home, ".claude")
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	settingsPath := filepath.Join(configDir, "settings.json")
	if _, err := os.Stat(settingsPath); err == nil {
		fmt.Printf("%s already exists. Skipping.\n", settingsPath)
		return nil
	}

	home, _ := os.UserHomeDir()
	exePath := filepath.Join(home, ".diet", "bin", "diet.exe")
	if runtime.GOOS != "windows" {
		exePath = filepath.Join(home, ".diet", "bin", "diet")
	}
	
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

	fmt.Printf("Successfully created %s with Diet hook.\n", settingsPath)
	return nil
}
