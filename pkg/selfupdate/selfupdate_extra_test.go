package selfupdate

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)



func TestCheckAndApplyUpdate_NoAsset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		releases := []Release{
			{
				TagName: "v2.0.0",
				Assets: []struct {
					ID                 int64  `json:"id"`
					Name               string `json:"name"`
					BrowserDownloadURL string `json:"browser_download_url"`
					URL                string `json:"url"`
				}{
					{Name: "wrong-asset.txt", URL: "http://example.com"},
				},
			},
		}
		json.NewEncoder(w).Encode(releases)
	}))
	defer server.Close()

	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	_, err := CheckAndApplyUpdate("v1.0.0")
	if err == nil || !strings.Contains(err.Error(), "could not find binary") {
		t.Errorf("Expected asset not found error, got %v", err)
	}
}

func TestDownloadAndReplace_ExecutableError(t *testing.T) {
	oldExe := osExecutable
	osExecutable = func() (string, error) { return "", errors.New("exec error") }
	defer func() { osExecutable = oldExe }()

	err := downloadAndReplace("http://example.com", "")
	if err == nil || err.Error() != "exec error" {
		t.Errorf("Expected exec error, got %v", err)
	}
}

func TestDownloadAndReplace_RenameError(t *testing.T) {
	tmpDir := t.TempDir()
	fakeExe := filepath.Join(tmpDir, "pith.exe")
	os.WriteFile(fakeExe, []byte("old"), 0755)

	// Create a directory where the .old file should be to block Rename
	os.MkdirAll(fakeExe+".old", 0755)
	os.WriteFile(filepath.Join(fakeExe+".old", "blocker"), []byte("prevent remove"), 0644)

	oldExe := osExecutable
	osExecutable = func() (string, error) { return fakeExe, nil }
	defer func() { osExecutable = oldExe }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("new binary"))
	}))
	defer server.Close()

	err := downloadAndReplace(server.URL, "")
	if err == nil {
		t.Error("Expected rename error, got nil")
	}
}
