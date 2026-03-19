package parser

import (
	"fmt"
	"strings"
)

type GoParser struct{}

func (g *GoParser) Name() string { return "go" }

func (g *GoParser) CanParse(cmd string, args []string) bool {
	return cmd == "go" || strings.HasSuffix(cmd, "go.exe")
}

func (g *GoParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }
		
		// Summary lines (common in go test and go list)
		if strings.Contains(trimmed, "passed") || strings.Contains(trimmed, "failed") || 
		   strings.Contains(trimmed, "TOTAL") || strings.HasPrefix(trimmed, "ok  ") || 
		   strings.HasPrefix(trimmed, "FAIL") {
			result = append(result, trimmed)
			continue
		}
		
		// If it's go list -m all or go version
		if strings.Contains(trimmed, "go version") || (len(strings.Fields(trimmed)) == 2 && strings.Contains(trimmed, "/")) {
			result = append(result, trimmed)
			continue
		}
		
		// If it's help or build errors, show only a few
		if len(result) < 20 {
			result = append(result, trimmed)
		}
	}
	
	if len(lines) > 20 && len(result) > 20 {
		return strings.Join(result[:20], "\n") + fmt.Sprintf("\n... (+ %d more lines)", len(lines)-20)
	}
	return strings.Join(result, "\n")
}
