package selfupdate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const repo = "Zkrausman/Pith"

var githubAPI = "https://api.github.com"
var osExecutable = os.Executable


type Release struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
	Assets  []struct {
		ID                 int64  `json:"id"`
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		URL                string `json:"url"` // This is the API URL for the asset
	} `json:"assets"`
}

func getAuthToken() string {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token
	}

	// Fallback: Try to get token from GitHub CLI
	cmd := exec.Command("gh", "auth", "token")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		return strings.TrimSpace(out.String())
	}

	return ""
}

func CheckAndApplyUpdate(currentVersion string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repos/%s/releases", githubAPI, repo), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "Pith-Updater")

	// Support private repos via GITHUB_TOKEN or gh CLI
	token := getAuthToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, fmt.Errorf("repository not found (if it is private, set GITHUB_TOKEN or run 'gh auth login')")
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return false, err
	}

	if len(releases) == 0 {
		return false, fmt.Errorf("no releases found")
	}

	release := releases[0] // Latest is first

	if release.TagName == currentVersion {
		return false, nil
	}

	fmt.Printf("New version found: %s (current: %s)\n", release.TagName, currentVersion)
	if release.Body != "" {
		fmt.Printf("\n--- Changelog ---\n%s\n-----------------\n\n", release.Body)
	}
	
	var assetURL string
	expectedSuffix := fmt.Sprintf("-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		expectedSuffix += ".exe"
	}

	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, expectedSuffix) {
			assetURL = asset.URL
			break
		}
	}

	if assetURL == "" {
		return false, fmt.Errorf("could not find binary for %s in release %s", runtime.GOOS, release.TagName)
	}

	fmt.Println("Downloading update...")
	if err := downloadAndReplace(assetURL, token); err != nil {
		return false, err
	}

	return true, nil
}

func CheckForUpdateSilent(currentVersion string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repos/%s/releases", githubAPI, repo), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Pith-Updater")

	if token := getAuthToken(); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", err
	}

	if len(releases) > 0 && releases[0].TagName != currentVersion {
		return releases[0].TagName, nil
	}

	return "", nil
}

func downloadAndReplace(url string, token string) error {
	executablePath, err := osExecutable()
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Pith-Updater")
	
	// CRITICAL for private asset downloads via API
	req.Header.Set("Accept", "application/octet-stream")
	
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub returned status %d while downloading asset", resp.StatusCode)
	}

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
