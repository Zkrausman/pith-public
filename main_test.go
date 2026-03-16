package main

import (
	"encoding/json"
	"testing"
)

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
