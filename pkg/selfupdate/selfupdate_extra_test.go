package selfupdate

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckAndApplyUpdateNoPlatformAsset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]Release{{
			TagName: "v2.0.0",
			Assets:  []ReleaseAsset{{Name: "wrong-asset.txt", URL: "http://example.com"}},
		}})
	}))
	defer server.Close()

	oldAPI := githubAPI
	githubAPI = server.URL
	defer func() { githubAPI = oldAPI }()

	_, err := CheckAndApplyUpdate("v1.0.0")
	if err == nil || !strings.Contains(err.Error(), "could not find") {
		t.Errorf("expected asset-not-found error, got %v", err)
	}
}

func TestFindAssetRejectsDuplicateNames(t *testing.T) {
	_, err := findAsset([]ReleaseAsset{{Name: "checksums.txt"}, {Name: "checksums.txt"}}, "checksums.txt")
	if err == nil || !strings.Contains(err.Error(), "multiple assets") {
		t.Fatalf("expected duplicate-asset error, got %v", err)
	}
}
