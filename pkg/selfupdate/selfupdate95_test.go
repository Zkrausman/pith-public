package selfupdate

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCheckForUpdateSilent_UpToDate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		releases := []Release{{TagName: "v1.0.0"}}
		json.NewEncoder(w).Encode(releases)
	}))
	defer server.Close()

	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	version, err := CheckForUpdateSilent("v1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if version != "" {
		t.Errorf("Expected empty version (up to date), got %s", version)
	}
}

func TestCheckForUpdateSilent_EmptyReleases(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Release{})
	}))
	defer server.Close()

	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	version, err := CheckForUpdateSilent("v1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if version != "" {
		t.Errorf("Expected empty version for empty releases, got %s", version)
	}
}

func TestCheckAndApplyUpdate_WithChangelog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suffix := "pith-" + runtime.GOOS + "-" + runtime.GOARCH
		if runtime.GOOS == "windows" {
			suffix += ".exe"
		}
		releases := []Release{
			{
				TagName: "v99.9.9",
				Body:    "## Changelog\n- Fixed something\n- Added another",
				Assets: []struct {
					ID                 int64  `json:"id"`
					Name               string `json:"name"`
					BrowserDownloadURL string `json:"browser_download_url"`
					URL                string `json:"url"`
				}{
					{
						Name: suffix,
						URL:  "http://" + r.Host + "/download",
					},
				},
			},
		}
		if r.URL.Path == "/download" {
			w.Write([]byte("fake binary"))
			return
		}
		json.NewEncoder(w).Encode(releases)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	fakeExe := filepath.Join(tmpDir, "pith.exe")
	os.WriteFile(fakeExe, []byte("old"), 0755)

	oldAPI := githubAPI
	oldExe := osExecutable
	githubAPI = server.URL
	osExecutable = func() (string, error) { return fakeExe, nil }
	defer func() {
		githubAPI = oldAPI
		osExecutable = oldExe
	}()

	updated, err := CheckAndApplyUpdate("v0.0.1")
	if err != nil {
		t.Fatalf("CheckAndApplyUpdate failed: %v", err)
	}
	if !updated {
		t.Error("Expected updated=true")
	}
}

func TestCheckAndApplyUpdate_OtherStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	_, err := CheckAndApplyUpdate("v0.0.1")
	if err == nil {
		t.Error("Expected error for 500 status")
	}
}

func TestCheckAndApplyUpdate_EmptyReleases(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Release{})
	}))
	defer server.Close()

	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	_, err := CheckAndApplyUpdate("v0.0.1")
	if err == nil {
		t.Error("Expected error for empty releases")
	}
}

func TestDownloadAndReplace_DownloadError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	fakeExe := filepath.Join(tmpDir, "pith.exe")
	os.WriteFile(fakeExe, []byte("old"), 0755)

	oldExe := osExecutable
	osExecutable = func() (string, error) { return fakeExe, nil }
	defer func() { osExecutable = oldExe }()

	err := downloadAndReplace(server.URL, "")
	if err == nil {
		t.Error("Expected error for 403 download")
	}
}

func TestDownloadAndReplace_WithToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer mytoken" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte("authorized binary"))
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	fakeExe := filepath.Join(tmpDir, "pith.exe")
	os.WriteFile(fakeExe, []byte("old"), 0755)

	oldExe := osExecutable
	osExecutable = func() (string, error) { return fakeExe, nil }
	defer func() { osExecutable = oldExe }()

	err := downloadAndReplace(server.URL, "mytoken")
	if err != nil {
		t.Fatalf("Expected success with token, got %v", err)
	}

	data, _ := os.ReadFile(fakeExe)
	if string(data) != "authorized binary" {
		t.Errorf("Expected authorized binary, got %s", string(data))
	}
}

func TestGetAuthToken_FallbackEmpty(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	// gh auth token will likely fail in test env, so token should be ""
	token := getAuthToken()
	// Just ensure no panic; token could be empty or populated by gh CLI
	_ = token
}
