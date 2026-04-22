package main

import (
	"bytes"
	"encoding/json"
	"os"
	"pith/pkg/telemetry"
	"testing"
)

func TestResetCmd_AllFlag(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"reset", "--all"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("reset --all failed: %v", err)
	}
	if !bytes.Contains(b.Bytes(), []byte("All telemetry data")) {
		t.Errorf("Expected 'All telemetry data' in output, got %s", b.String())
	}
}

func TestInstallCmd_Local(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"install"}) // No --global
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	_ = cmd.Execute()
}

func TestRawCmd_NoArgs(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"raw"}) // No args after "raw"
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	_ = cmd.Execute()
	if !bytes.Contains(b.Bytes(), []byte("Usage:")) {
		t.Errorf("Expected help output, got %s", b.String())
	}
}

func TestGainCmd_ManyRecords(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	tel, _ := telemetry.NewTelemetry(tmpDir)
	// 25 different commands to hit the >20 limit
	for i := 0; i < 25; i++ {
		tel.Record(telemetry.ExecutionRecord{
			Command:        "cmd-" + string(rune('a'+i)),
			OriginalTokens: 100,
			CompressedTokens: 50,
			Source:         "gemini",
		})
	}
	// 15 unparsed to hit the >10 limit
	for i := 0; i < 15; i++ {
		tel.Record(telemetry.ExecutionRecord{
			Command:        "unparsed-" + string(rune('a'+i)),
			OriginalTokens: 100,
			CompressedTokens: 100,
			IsPassthrough:  true,
			Source:         "gemini",
		})
	}
	tel.Close()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"gain"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("gain failed: %v", err)
	}
	if !bytes.Contains(b.Bytes(), []byte("and 5 more")) {
		t.Errorf("Expected 'and 5 more' for commands, got %s", b.String())
	}
	if !bytes.Contains(b.Bytes(), []byte("and 5 more")) {
		t.Errorf("Expected 'and 5 more' for unparsed, got %s", b.String())
	}
}

func TestHookCmd_OldSchemas(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	// Test ToolResponse schema
	r1, w1, _ := os.Pipe()
	input1 := HookInput{}
	input1.ToolResponse.LlmContent = "some content"
	json.NewEncoder(w1).Encode(input1)
	w1.Close()

	oldStdin := os.Stdin
	os.Stdin = r1
	
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"_hook"})
	_ = cmd.Execute()
	r1.Close()

	// Test ToolCallRequest schema
	r2, w2, _ := os.Pipe()
	input2 := HookInput{}
	input2.ToolCallRequest.Arguments = `{"command":"ls"}`
	json.NewEncoder(w2).Encode(input2)
	w2.Close()

	os.Stdin = r2
	cmd.Execute()
	r2.Close()

	os.Stdin = oldStdin
}

func TestUpdateCmd_Run(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"update", "--help"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("update failed: %v", err)
	}
}
