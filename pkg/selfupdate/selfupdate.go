package selfupdate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

const repo = "Zkrausman/Diet"

type Release struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func CheckAndApplyUpdate(currentVersion string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, err
	}

	if release.TagName == currentVersion {
		return false, nil
	}

	fmt.Printf("New version found: %s (current: %s)\n", release.TagName, currentVersion)
	if release.Body != "" {
		fmt.Printf("\n--- Changelog ---\n%s\n-----------------\n\n", release.Body)
	}
	
	var downloadURL string
	expectedName := "diet"
	if runtime.GOOS == "windows" {
		expectedName = "diet.exe"
	}

	for _, asset := range release.Assets {
		if asset.Name == expectedName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return false, fmt.Errorf("could not find binary for %s in release %s", runtime.GOOS, release.TagName)
	}

	fmt.Println("Downloading update...")
	if err := downloadAndReplace(downloadURL); err != nil {
		return false, err
	}

	return true, nil
}

func downloadAndReplace(url string) error {
	executablePath, err := os.Executable()
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create a temporary file for the new binary
	tmpPath := executablePath + ".tmp"
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return err
	}
	tmpFile.Close()

	// On Windows, we can't replace the running executable directly easily
	// We rename the current one and move the new one in
	oldPath := executablePath + ".old"
	_ = os.Remove(oldPath) // Clean up any previous old file
	
	if err := os.Rename(executablePath, oldPath); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, executablePath); err != nil {
		// Try to rollback if possible
		_ = os.Rename(oldPath, executablePath)
		return err
	}

	fmt.Println("Update successful! Please restart the application.")
	return nil
}
