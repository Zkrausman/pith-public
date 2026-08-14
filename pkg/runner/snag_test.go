package runner

import (
	"os"
	"path/filepath"
	"pith/pkg/config"
	"strings"
	"testing"
)

func TestLogForSnag(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pith-snag-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		StoragePath: tmpDir,
		SnagLogging: true,
	}
	r := &Runner{cfg: cfg}

	r.LogForSnag("echo test", "test output", 0)

	logPath := filepath.Join(tmpDir, "pith.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Fatal("pith.log was not created")
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "[CMD] echo test") {
		t.Errorf("Expected command in log, got: %s", string(content))
	}
	if !strings.Contains(string(content), "test output") {
		t.Errorf("Expected output in log, got: %s", string(content))
	}
	if !strings.Contains(string(content), "[EXIT] 0") {
		t.Errorf("Expected exit code in log, got: %s", string(content))
	}
}

func TestLogForSnagTruncation(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "pith-snag-trunc-test-*")
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		StoragePath: tmpDir,
		SnagLogging: true,
	}
	r := &Runner{cfg: cfg}

	// Create long output (more than 50 lines)
	var longOutput strings.Builder
	for i := 0; i < 100; i++ {
		longOutput.WriteString("line\n")
	}

	r.LogForSnag("long command", longOutput.String(), 1)

	logPath := filepath.Join(tmpDir, "pith.log")
	content, _ := os.ReadFile(logPath)

	if !strings.Contains(string(content), "truncated by Pith") {
		t.Error("Expected log to be truncated")
	}
}
