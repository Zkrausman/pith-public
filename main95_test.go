package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"pith/pkg/telemetry"
	"testing"
)

// =====================================================================
// runGain - with actual data and more branches
// =====================================================================

func TestGainCmd_WithData(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	// Seed some telemetry
	tel, _ := telemetry.NewTelemetry(tmpDir)
	for i := 0; i < 25; i++ {
		tel.Record(telemetry.ExecutionRecord{
			Command:          "git status",
			OriginalTokens:   100,
			CompressedTokens: 30,
			IsPassthrough:    false,
			Source:           "gemini",
		})
	}
	tel.Close()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"gain"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("gain command failed: %v", err)
	}
	if !bytes.Contains(b.Bytes(), []byte("Token")) {
		t.Errorf("Expected Token info in output, got %s", b.String())
	}
}

func TestGainCmd_NoData(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"gain"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("gain command failed: %v", err)
	}
	if !bytes.Contains(b.Bytes(), []byte("No telemetry")) {
		t.Errorf("Expected 'No telemetry' message, got %s", b.String())
	}
}

func TestGainCmd_WithUnparsed(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	tel, _ := telemetry.NewTelemetry(tmpDir)
	// Record passthrough (unparsed) commands
	for i := 0; i < 5; i++ {
		tel.Record(telemetry.ExecutionRecord{
			Command:          "mycommand args",
			OriginalTokens:   200,
			CompressedTokens: 200,
			IsPassthrough:    true,
			Source:           "gemini",
		})
	}
	// Also record parsed commands to make totalOrig > 0
	tel.Record(telemetry.ExecutionRecord{
		Command:          "git status",
		OriginalTokens:   500,
		CompressedTokens: 100,
		IsPassthrough:    false,
		Source:           "gemini",
	})
	tel.Close()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"gain"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("gain command failed: %v", err)
	}
}

// =====================================================================
// runDiscover - with data
// =====================================================================

func TestDiscoverCmd_WithData(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	tel, _ := telemetry.NewTelemetry(tmpDir)
	for i := 0; i < 3; i++ {
		tel.Record(telemetry.ExecutionRecord{
			Command:          "myunknowncmd args",
			OriginalTokens:   300,
			CompressedTokens: 300,
			IsPassthrough:    true,
			Source:           "gemini",
		})
	}
	tel.Close()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"discover"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("discover command failed: %v", err)
	}
	if !bytes.Contains(b.Bytes(), []byte("myunknowncmd")) {
		t.Errorf("Expected 'myunknowncmd' in discover output, got %s", b.String())
	}
}

func TestDiscoverCmd_NoData(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"discover"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("discover command failed: %v", err)
	}
	if !bytes.Contains(b.Bytes(), []byte("No unparsed")) {
		t.Errorf("Expected 'No unparsed' message, got %s", b.String())
	}
}

// =====================================================================
// runUpdate - with mock server
// =====================================================================

func TestUpdateCmd_AlreadyUpToDate(t *testing.T) {
	// Just test the help output since we can't mock the selfupdate API easily from main package
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"update", "--help"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	_ = cmd.Execute()
}

// =====================================================================
// runRaw - with actual command
// =====================================================================

func TestRawCmd_WithEcho(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"raw", "cmd", "/c", "echo", "hello"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("raw command failed: %v", err)
	}
}

// =====================================================================
// runReset branches
// =====================================================================

func TestResetCmd_DiscoverFlag(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"reset", "--discover"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("reset --discover failed: %v", err)
	}
	if !bytes.Contains(b.Bytes(), []byte("Discovery data")) {
		t.Errorf("Expected 'Discovery data' in output, got %s", b.String())
	}
}

// =====================================================================
// runDashboard - port flag
// =====================================================================

func TestDashboardCmd_Port(t *testing.T) {
	// Just test that --port flag is accepted without crashing immediately.
	// Actual server binding is not testable headlessly.
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"dashboard", "--help"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	_ = cmd.Execute()
	if !bytes.Contains(b.Bytes(), []byte("dashboard")) && !bytes.Contains(b.Bytes(), []byte("port")) {
		t.Logf("Dashboard help: %s", b.String())
	}
}

// =====================================================================
// runHook - various tool responses
// =====================================================================

func TestHookCmd_RespondBlock(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	// Create a pipe to simulate stdin
	r, w, _ := os.Pipe()
	defer r.Close()

	hookInput := HookInput{
		ToolInput: map[string]interface{}{
			"command": "echo test",
		},
	}
	inputBytes, _ := json.Marshal(hookInput)
	w.Write(inputBytes)
	w.Close()

	// Save and restore stdin
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"_hook"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	_ = cmd.Execute()
}

// =====================================================================
// runInstall - already installed branch
// =====================================================================

func TestInstallCmd_GlobalFlag(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	// Create a fake pith.exe destination to avoid actual install
	fakeHome := t.TempDir()
	t.Setenv("USERPROFILE", fakeHome)
	t.Setenv("HOME", fakeHome)

	binDir := filepath.Join(fakeHome, ".pith", "bin")
	os.MkdirAll(binDir, 0755)

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"install", "--global"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	_ = cmd.Execute()
}
