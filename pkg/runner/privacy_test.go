package runner

import (
	"os"
	"path/filepath"
	"pith/pkg/config"
	"runtime"
	"strings"
	"testing"
)

func TestLogForSnagDisabledByDefault(t *testing.T) {
	dir := t.TempDir()
	(&Runner{cfg: &config.Config{StoragePath: dir}}).LogForSnag("echo token=secret", "token=secret", 0)
	if _, err := os.Stat(filepath.Join(dir, "pith.log")); !os.IsNotExist(err) {
		t.Fatal("diagnostic log was created without opt-in")
	}
}

func TestLogForSnagRedactsAndUsesPrivateMode(t *testing.T) {
	dir := t.TempDir()
	(&Runner{cfg: &config.Config{StoragePath: dir, SnagLogging: true}}).LogForSnag("curl --token=secret", "authorization: bearer secret", 0)
	data, err := os.ReadFile(filepath.Join(dir, "pith.log"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "secret") {
		t.Fatalf("secret persisted: %s", data)
	}
	if runtime.GOOS != "windows" {
		if info, err := os.Stat(filepath.Join(dir, "pith.log")); err != nil || info.Mode().Perm() != 0600 {
			t.Fatalf("unexpected log mode: %v, %v", info.Mode(), err)
		}
	}
}
