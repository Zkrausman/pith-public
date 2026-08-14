package selfupdate

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckForUpdateSilentUpToDate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"tag_name":"v1.0.0"}]`))
	}))
	defer server.Close()

	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	version, err := CheckForUpdateSilent("v1.0.0")
	if err != nil || version != "" {
		t.Fatalf("expected no update, got version %q and error %v", version, err)
	}
}

func TestCheckAndApplyUpdateRejectsReleaseAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	_, err := CheckAndApplyUpdate("v0.0.1")
	if err == nil || !strings.Contains(err.Error(), "status 500") {
		t.Fatalf("expected API status error, got %v", err)
	}
}

func TestDownloadAndReplaceExecutableError(t *testing.T) {
	oldExe := osExecutable
	osExecutable = func() (string, error) { return "", errors.New("exec error") }
	defer func() { osExecutable = oldExe }()

	err := downloadAndReplace("http://example.com", "", strings.Repeat("0", 64))
	if err == nil || err.Error() != "exec error" {
		t.Errorf("expected exec error, got %v", err)
	}
}

func TestDownloadAndReplaceRejectsDownloadError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	fakeExe := filepath.Join(t.TempDir(), "pith")
	if err := os.WriteFile(fakeExe, []byte("old"), 0755); err != nil {
		t.Fatal(err)
	}
	oldExe := osExecutable
	osExecutable = func() (string, error) { return fakeExe, nil }
	defer func() { osExecutable = oldExe }()
	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	if err := downloadAndReplace(server.URL, "", strings.Repeat("0", 64)); err == nil {
		t.Fatal("expected error for forbidden download")
	}
}
