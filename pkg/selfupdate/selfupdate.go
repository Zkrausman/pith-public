package selfupdate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const repo = "Zkrausman/pith-public"

const (
	maxChecksumFileSize = 1 << 20
	maxBinarySize       = 256 << 20
	requestTimeout      = 30 * time.Second
)

var githubAPI = "https://api.github.com"
var osExecutable = os.Executable

func isNewerVersion(candidate, current string) bool {
	parse := func(v string) ([3]int, bool) {
		var result [3]int
		v = strings.TrimPrefix(strings.TrimSpace(v), "v")
		if strings.ContainsAny(v, "-+") {
			return result, false
		}
		parts := strings.Split(v, ".")
		if len(parts) != 3 {
			return result, false
		}
		for i, part := range parts {
			n, err := strconv.Atoi(part)
			if err != nil || n < 0 {
				return result, false
			}
			result[i] = n
		}
		return result, true
	}
	candidateParts, ok := parse(candidate)
	if !ok {
		return false
	}
	currentParts, ok := parse(current)
	if !ok {
		return false
	}
	for i := range candidateParts {
		if candidateParts[i] != currentParts[i] {
			return candidateParts[i] > currentParts[i]
		}
	}
	return false
}

type ReleaseAsset struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	URL                string `json:"url"`
}

type Release struct {
	TagName string         `json:"tag_name"`
	Body    string         `json:"body"`
	Assets  []ReleaseAsset `json:"assets"`
}

func newRequest(url, token string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Pith-Updater")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req, nil
}

func newHTTPClient() *http.Client {
	trustedAPI, _ := url.Parse(githubAPI)
	return &http.Client{
		Timeout: requestTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// GitHub release assets redirect to object storage. Never forward an
			// API token to a different host.
			if trustedAPI == nil || !strings.EqualFold(req.URL.Host, trustedAPI.Host) {
				req.Header.Del("Authorization")
			}
			return nil
		},
	}
}

func validateAssetURL(rawURL string) error {
	assetURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("parse release asset URL: %w", err)
	}
	apiURL, err := url.Parse(githubAPI)
	if err != nil {
		return fmt.Errorf("parse GitHub API URL: %w", err)
	}
	if assetURL.Scheme != apiURL.Scheme || !strings.EqualFold(assetURL.Host, apiURL.Host) {
		return fmt.Errorf("release asset URL is not hosted by the configured GitHub API")
	}
	return nil
}

func getReleases(token string) ([]Release, error) {
	req, err := newRequest(fmt.Sprintf("%s/repos/%s/releases", githubAPI, repo), token)
	if err != nil {
		return nil, err
	}
	resp, err := newHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("repository not found (if it is private, set GITHUB_TOKEN or run 'gh auth login')")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}
	return releases, nil
}

func latestRelease(releases []Release, currentVersion string) *Release {
	var latest *Release
	for i := range releases {
		if isNewerVersion(releases[i].TagName, currentVersion) && (latest == nil || isNewerVersion(releases[i].TagName, latest.TagName)) {
			latest = &releases[i]
		}
	}
	return latest
}

func platformAssetName() string {
	name := fmt.Sprintf("pith-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}

func findAsset(assets []ReleaseAsset, name string) (ReleaseAsset, error) {
	var found *ReleaseAsset
	for i := range assets {
		if assets[i].Name != name {
			continue
		}
		if found != nil {
			return ReleaseAsset{}, fmt.Errorf("release contains multiple assets named %q", name)
		}
		found = &assets[i]
	}
	if found == nil {
		return ReleaseAsset{}, fmt.Errorf("could not find %q in release", name)
	}
	return *found, nil
}

func downloadAsset(url, token string, maxBytes int64) ([]byte, error) {
	if err := validateAssetURL(url); err != nil {
		return nil, err
	}
	req, err := newRequest(url, token)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/octet-stream")
	resp, err := newHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub returned status %d while downloading asset", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("downloaded asset exceeds %d byte limit", maxBytes)
	}
	return data, nil
}

func checksumForAsset(checksums []byte, assetName string) (string, error) {
	var expected string
	for _, line := range strings.Split(string(checksums), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 || len(fields[0]) != sha256.Size*2 || fields[1] != assetName {
			continue
		}
		if _, err := hex.DecodeString(fields[0]); err != nil {
			continue
		}
		if expected != "" {
			return "", fmt.Errorf("checksums file contains multiple entries for %q", assetName)
		}
		expected = strings.ToLower(fields[0])
	}
	if expected == "" {
		return "", fmt.Errorf("checksums file does not contain a SHA-256 checksum for %q", assetName)
	}
	return expected, nil
}

func CheckAndApplyUpdate(currentVersion string) (bool, error) {
	// pith-public releases are public. Do not read GITHUB_TOKEN or invoke gh:
	// an update check must never silently acquire local credentials.
	token := ""
	releases, err := getReleases(token)
	if err != nil {
		return false, err
	}
	if len(releases) == 0 {
		return false, fmt.Errorf("no releases found")
	}
	release := latestRelease(releases, currentVersion)
	if release == nil {
		return false, nil
	}

	fmt.Printf("New version found: %s (current: %s)\n", release.TagName, currentVersion)
	if release.Body != "" {
		fmt.Printf("\n--- Changelog ---\n%s\n-----------------\n\n", release.Body)
	}

	binaryName := platformAssetName()
	binary, err := findAsset(release.Assets, binaryName)
	if err != nil {
		return false, err
	}
	checksumAsset, err := findAsset(release.Assets, "checksums.txt")
	if err != nil {
		return false, fmt.Errorf("release %s has no checksums.txt: refusing unverified update", release.TagName)
	}
	checksums, err := downloadAsset(checksumAsset.URL, token, maxChecksumFileSize)
	if err != nil {
		return false, fmt.Errorf("download release checksums: %w", err)
	}
	expectedChecksum, err := checksumForAsset(checksums, binaryName)
	if err != nil {
		return false, fmt.Errorf("verify release checksums: %w", err)
	}

	fmt.Println("Downloading and verifying update...")
	if err := downloadAndReplace(binary.URL, token, expectedChecksum); err != nil {
		return false, err
	}
	return true, nil
}

func CheckForUpdateSilent(currentVersion string) (string, error) {
	// Silent checks intentionally use anonymous public GitHub API requests.
	releases, err := getReleases("")
	if err != nil {
		return "", err
	}
	if release := latestRelease(releases, currentVersion); release != nil {
		return release.TagName, nil
	}
	return "", nil
}

func downloadAndReplace(url, token, expectedChecksum string) error {
	executablePath, err := osExecutable()
	if err != nil {
		return err
	}
	if err := validateAssetURL(url); err != nil {
		return err
	}

	req, err := newRequest(url, token)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/octet-stream")
	resp, err := newHTTPClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub returned status %d while downloading asset", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(executablePath), "."+filepath.Base(executablePath)+".update-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	if err := tmpFile.Chmod(0755); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	hasher := sha256.New()
	written, copyErr := io.Copy(io.MultiWriter(tmpFile, hasher), io.LimitReader(resp.Body, maxBinarySize+1))
	closeErr := tmpFile.Close()
	if copyErr != nil {
		_ = os.Remove(tmpPath)
		return copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tmpPath)
		return closeErr
	}
	if written > maxBinarySize {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("downloaded update exceeds %d byte limit", maxBinarySize)
	}
	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	if !strings.EqualFold(actualChecksum, expectedChecksum) {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("downloaded update checksum mismatch: got %s, expected %s", actualChecksum, expectedChecksum)
	}

	if runtime.GOOS != "windows" {
		// Same-directory rename atomically replaces the old executable on Unix.
		if err := os.Rename(tmpPath, executablePath); err != nil {
			_ = os.Remove(tmpPath)
			return err
		}
	} else {
		// Windows cannot atomically overwrite an executable in every deployment
		// configuration. Preserve a recoverable prior binary and report rollback
		// failure rather than silently leaving the installation broken.
		oldPath := executablePath + ".old"
		_ = os.Remove(oldPath)
		if err := os.Rename(executablePath, oldPath); err != nil {
			_ = os.Remove(tmpPath)
			return err
		}
		if err := os.Rename(tmpPath, executablePath); err != nil {
			if rollbackErr := os.Rename(oldPath, executablePath); rollbackErr != nil {
				return fmt.Errorf("install update: %w; rollback failed: %v", err, rollbackErr)
			}
			return err
		}
	}

	fmt.Println("Update successful! Please restart the application.")
	return nil
}
