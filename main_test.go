package main

import (
	"bytes"
	"encoding/json"
	"os"
	"pith/pkg/telemetry"
	"strings"
	"testing"
)


func TestRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--version"})
	
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() failed: %v", err)
	}

	out := b.String()
	if !strings.Contains(out, version) {
		t.Errorf("Expected version %s, got %s", version, out)
	}
}

func TestVersionCmd(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"version"})

	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	out := b.String()
	if !strings.Contains(out, version) {
		t.Errorf("Expected version %s, got %s", version, out)
	}
}

func TestGainCmd_Help(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"gain", "--help"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	if err := cmd.Execute(); err != nil {
		t.Errorf("gain --help failed: %v", err)
	}
}

func TestDiscoverCmd_Help(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"discover", "--help"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	if err := cmd.Execute(); err != nil {
		t.Errorf("discover --help failed: %v", err)
	}
}

func TestConfigCmd_Help(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"config", "--help"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	if err := cmd.Execute(); err != nil {
		t.Errorf("config --help failed: %v", err)
	}
}

func TestHookInputSchema(t *testing.T) {
	// Test newer schema
	data := []byte(`{"tool_input": {"command": "git status"}, "tool_response": {"llmContent": "hi"}}`)
	var input HookInput
	if err := json.Unmarshal(data, &input); err != nil {
		t.Fatal(err)
	}
	if input.ToolInput["command"] != "git status" {
		t.Errorf("Expected git status, got %v", input.ToolInput["command"])
	}

	// Test older schema
	dataOld := []byte(`{"tool_call_request": {"arguments": "{\"command\": \"ls\"}"}, "tool_response": {"llmContent": "hi"}}`)
	var inputOld HookInput
	if err := json.Unmarshal(dataOld, &inputOld); err != nil {
		t.Fatal(err)
	}
	if inputOld.ToolCallRequest.Arguments == "" {
		t.Error("Expected arguments to be populated")
	}
}

func TestHookOutputSchema(t *testing.T) {
	output := HookOutput{
		Decision: "deny",
		Reason:   "compressed",
	}
	data, err := json.Marshal(output)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"decision":"deny","reason":"compressed","systemMessage":""}` {
		t.Errorf("Unexpected JSON output: %s", string(data))
	}
}

func TestGainCmd(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	// Add some mock telemetry data
	tel, err := telemetry.NewTelemetry(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	tel.Record(telemetry.ExecutionRecord{
		Command:          "git status",
		OriginalTokens:   100,
		CompressedTokens: 50,
		ParserUsed:       "git",
	})
	tel.Close()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"gain"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	if err := cmd.Execute(); err != nil {
		t.Errorf("gain failed: %v", err)
	}
	out := b.String()
	if !strings.Contains(out, "Pith: Overall Token Savings") {
		t.Errorf("Unexpected output: %s", out)
	}
}

func TestDiscoverCmd(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	// Add some mock telemetry data
	tel, err := telemetry.NewTelemetry(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	tel.Record(telemetry.ExecutionRecord{
		Command:          "unknown-cmd",
		OriginalTokens:   100,
		CompressedTokens: 100,
		IsPassthrough:    true,
	})
	tel.Close()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"discover"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	if err := cmd.Execute(); err != nil {
		t.Errorf("discover failed: %v", err)
	}
	out := b.String()
	if !strings.Contains(out, "Opportunity Discovery") {
		t.Errorf("Unexpected output: %s", out)
	}
}

func TestResetCmd(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"reset", "--all"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	if err := cmd.Execute(); err != nil {
		t.Errorf("reset failed: %v", err)
	}
	out := b.String()
	if !strings.Contains(out, "All telemetry data has been reset") {
		t.Errorf("Unexpected output: %s", out)
	}
}

func TestHookCmd(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PITH_STORAGE", tmpDir)

	input := HookInput{}
	input.ToolInput = make(map[string]interface{})
	input.ToolInput["command"] = "git status"
	input.ToolResponse.LlmContent = "Output: On branch main\nnothing to commit"

	data, _ := json.Marshal(input)
	
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	
	r, w, _ := os.Pipe()
	os.Stdin = r
	go func() {
		w.Write(data)
		w.Close()
	}()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"_hook"})
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	
	if err := cmd.Execute(); err != nil {
		t.Errorf("_hook failed: %v", err)
	}
	
	var output HookOutput
	if err := json.Unmarshal(b.Bytes(), &output); err != nil {
		t.Fatalf("Failed to parse hook output: %v", err)
	}
	
	// Since git status should be parsed, decision should be deny
	if output.Decision != "deny" {
		t.Errorf("Expected decision deny, got %s", output.Decision)
	}
}


