package selfupdate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetAuthTokenEnv(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "test-token")
	if token := getAuthToken(); token != "test-token" {
		t.Errorf("expected test-token, got %s", token)
	}
}

func TestCheckForUpdateSilent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Release{{TagName: "v9.9.9"}})
	}))
	defer server.Close()

	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	version, err := CheckForUpdateSilent("v0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if version != "v9.9.9" {
		t.Errorf("expected v9.9.9, got %s", version)
	}
}

func TestCheckAndApplyUpdateVerifiesChecksumBeforeReplacement(t *testing.T) {
	const newBinary = "verified new binary"
	checksum := sha256.Sum256([]byte(newBinary))
	binaryName := platformAssetName()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/Zkrausman/pith-public/releases":
			json.NewEncoder(w).Encode([]Release{{
				TagName: "v9.9.9",
				Assets: []ReleaseAsset{
					{Name: binaryName, URL: serverURL(r, "/assets/binary")},
					{Name: "checksums.txt", URL: serverURL(r, "/assets/checksums")},
				},
			}})
		case "/assets/checksums":
			_, _ = w.Write([]byte(hex.EncodeToString(checksum[:]) + "  " + binaryName + "\n"))
		case "/assets/binary":
			_, _ = w.Write([]byte(newBinary))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	fakeExe := filepath.Join(t.TempDir(), "pith")
	if err := os.WriteFile(fakeExe, []byte("old binary"), 0755); err != nil {
		t.Fatal(err)
	}
	oldAPI, oldExecutable := githubAPI, osExecutable
	githubAPI = server.URL
	osExecutable = func() (string, error) { return fakeExe, nil }
	defer func() {
		githubAPI = oldAPI
		osExecutable = oldExecutable
	}()

	updated, err := CheckAndApplyUpdate("v0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if !updated {
		t.Fatal("expected update")
	}
	data, err := os.ReadFile(fakeExe)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != newBinary {
		t.Errorf("expected verified binary, got %q", data)
	}
}

func TestDownloadAndReplaceRejectsTamperedBinary(t *testing.T) {
	fakeExe := filepath.Join(t.TempDir(), "pith")
	if err := os.WriteFile(fakeExe, []byte("old binary"), 0755); err != nil {
		t.Fatal(err)
	}
	oldExecutable := osExecutable
	osExecutable = func() (string, error) { return fakeExe, nil }
	defer func() { osExecutable = oldExecutable }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("tampered binary"))
	}))
	defer server.Close()
	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	err := downloadAndReplace(server.URL, "", strings.Repeat("0", sha256.Size*2))
	if err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("expected checksum mismatch, got %v", err)
	}
	data, err := os.ReadFile(fakeExe)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "old binary" {
		t.Errorf("tampered update replaced executable: %q", data)
	}
	if _, err := os.Stat(fakeExe + ".tmp"); !os.IsNotExist(err) {
		t.Error("temporary file was not removed after checksum mismatch")
	}
}

func TestCheckAndApplyUpdateRequiresChecksums(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Release{{
			TagName: "v9.9.9",
			Assets:  []ReleaseAsset{{Name: platformAssetName(), URL: serverURL(r, "/assets/binary")}},
		}})
	}))
	defer server.Close()

	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	_, err := CheckAndApplyUpdate("v0.0.1")
	if err == nil || !strings.Contains(err.Error(), "no checksums.txt") {
		t.Fatalf("expected missing checksums error, got %v", err)
	}
}

func TestDownloadAssetRejectsForeignURL(t *testing.T) {
	trusted := httptest.NewServer(http.NotFoundHandler())
	defer trusted.Close()
	foreign := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("foreign asset host was contacted")
	}))
	defer foreign.Close()

	oldAPI := githubAPI
	githubAPI = trusted.URL
	defer func() { githubAPI = oldAPI }()

	_, err := downloadAsset(foreign.URL, "sensitive-token", maxChecksumFileSize)
	if err == nil || !strings.Contains(err.Error(), "not hosted") {
		t.Fatalf("expected foreign-host rejection, got %v", err)
	}
}

func TestChecksumForAsset(t *testing.T) {
	name := "pith-" + runtime.GOOS + "-" + runtime.GOARCH
	digest := strings.Repeat("a", sha256.Size*2)
	got, err := checksumForAsset([]byte(digest+"  "+name+"\n"), name)
	if err != nil || got != digest {
		t.Fatalf("expected checksum %q, got %q, err %v", digest, got, err)
	}
	_, err = checksumForAsset([]byte(digest+"  other\n"), name)
	if err == nil {
		t.Fatal("expected error for a missing platform checksum")
	}
	_, err = checksumForAsset([]byte(digest+"  "+name+"\n"+digest+"  "+name+"\n"), name)
	if err == nil {
		t.Fatal("expected error for duplicate platform checksums")
	}
}

func TestCheckAndApplyUpdateSameVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Release{{TagName: "v0.0.1"}})
	}))
	defer server.Close()

	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	updated, err := CheckAndApplyUpdate("v0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if updated {
		t.Error("should not update at the same version")
	}
}

func serverURL(r *http.Request, path string) string {
	return "http://" + r.Host + path
}
