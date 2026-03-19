package parser

import (
	"fmt"
	"strings"
)

type PowerShellParser struct{}

func (ps *PowerShellParser) Name() string { return "windows-shell" }

func (ps *PowerShellParser) CanParse(cmd string, args []string) bool {
	return cmd == "powershell" || cmd == "pwsh" || cmd == "cmd"
}

func (ps *PowerShellParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	
	// Headers to skip (exact match or prefix where appropriate)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }
		
		if strings.HasPrefix(trimmed, "Mode ") || strings.HasPrefix(trimmed, "----") || 
		   strings.HasPrefix(trimmed, "Directory:") || strings.Contains(trimmed, "Exit Code:") || 
		   strings.Contains(trimmed, "Process Group PGID:") {
			continue
		}
		
		// If it's a generic error from PowerShell (like CommandNotFoundException), minify it
		if strings.Contains(trimmed, "+ CategoryInfo") || strings.Contains(trimmed, "+ FullyQualifiedErrorId") {
			continue
		}
		
		if len(result) < 30 {
			result = append(result, trimmed)
		}
	}
	
	if len(result) == 0 {
		return ""
	}
	
	if len(lines) > 30 {
		return strings.Join(result, "\n") + fmt.Sprintf("\n... (+ %d lines removed)", len(lines)-len(result))
	}
	return strings.Join(result, "\n")
}
