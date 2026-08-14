package anomaly

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Diagnose performs an autonomous root-cause analysis of a given anomaly
// using Thneed GraphRAG and a diagnostic LLM prompt.
func Diagnose(a Anomaly) (string, error) {
	fmt.Printf("\n    [Pith Diagnostics] Querying Thneed for context related to '%s'...\n", a.Project)

	// 1. Gather only a bounded, redacted prompt summary for local context.
	query := fmt.Sprintf("%s anomaly: %s", diagnosticSnippet(a.Project), diagnosticSnippet(a.Prompt))
	thneedCmd := exec.Command("thneed", "query", query, "--depth", "1")
	var thneedOut bytes.Buffer
	thneedCmd.Stdout = &thneedOut
	_ = thneedCmd.Run() // Ignore errors, partial context is fine

	contextPack := thneedOut.String()
	if len(contextPack) > 2000 {
		contextPack = contextPack[:2000] // Truncate to avoid massive prompts
	}

	// 2. Build the diagnostic prompt
	prompt := fmt.Sprintf(`You are an expert AI Forensics Engineer.
An anomaly was detected in our LLM telemetry stream. Your task is to diagnose the root cause and provide actionable advice.

## Anomaly Details
- Project: %s
- Severity: %s
- Flag Reason: %s
- Model Used: %s

## Offending Interaction
**Prompt:** %s
**Response:** %s

## Related Codebase Context (via Thneed)
%s

Please provide a concise, structured Root Cause Analysis and a proposed fix.`,
		a.Project, a.Severity, a.Reason, a.Model,
		diagnosticSnippet(a.Prompt), diagnosticSnippet(a.Response),
		diagnosticSnippet(contextPack))

	fmt.Println("    [Pith Diagnostics] Consulting the Oracle for Root Cause Analysis...")

	// 3. Ask the LLM (Using gemini CLI, which Whet will intercept and route via Overseer)
	geminiCmd := exec.Command("gemini", "chat", "--prompt", prompt)

	// We want to force the oracle so Overseer handles it
	geminiCmd.Args = append(geminiCmd.Args, "--oracle")

	var out bytes.Buffer
	var stderr bytes.Buffer
	geminiCmd.Stdout = &out
	geminiCmd.Stderr = &stderr

	err := geminiCmd.Run()
	if err != nil {
		return "", fmt.Errorf("diagnostic failed: %v\nstderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}

var diagnosticSecretPattern = regexp.MustCompile(`(?i)(api[_-]?key|token|password|secret)\s*[:=]\s*[^\s]+`)

// diagnosticSnippet keeps external diagnostic payloads small and removes
// common credential assignments. It is deliberately conservative: callers
// must still explicitly opt in before anything leaves the machine.
func diagnosticSnippet(value string) string {
	value = diagnosticSecretPattern.ReplaceAllString(value, "$1=[REDACTED]")
	const maxBytes = 512
	if len(value) > maxBytes {
		return value[:maxBytes] + "…[truncated]"
	}
	return value
}
