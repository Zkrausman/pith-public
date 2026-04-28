package parser

import (
	"fmt"
	"strings"
)

type PithParser struct{}

func (p *PithParser) Name() string { return "pith-internal" }

func (p *PithParser) CanParse(cmd string, args []string) bool {
	// Either 'pith' or 'go run main.go' (common dev command)
	if MatchCommand(cmd, "pith") {
		return true
	}
	if MatchCommand(cmd, "go") && len(args) >= 2 && args[0] == "run" && args[1] == "main.go" {
		return true
	}
	return false
}

func (p *PithParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// If it's help output (common for unparsed commands)
		if strings.HasPrefix(trimmed, "Usage:") || strings.Contains(trimmed, "Available Commands:") || strings.Contains(trimmed, "Flags:") {
			// Extract just the command list or flags
			result = append(result, trimmed)
			continue
		}

		// If it's log output, summarize
		if strings.Contains(trimmed, "[CMD]") || strings.Contains(trimmed, "[EXIT]") {
			// Compact the snag log format
			if strings.HasPrefix(trimmed, "[CMD]") {
				result = append(result, "Command: "+strings.TrimSpace(trimmed[5:]))
			} else if strings.HasPrefix(trimmed, "[EXIT]") {
				result = append(result, "Exit Code: "+strings.TrimSpace(trimmed[6:]))
			} else {
				// Potential middle output, truncate if long
				if len(trimmed) > 100 {
					result = append(result, trimmed[:100]+"...")
				} else {
					result = append(result, trimmed)
				}
			}
			continue
		}

		// Standard output lines
		if len(result) < 20 {
			result = append(result, trimmed)
		}
	}

	if len(lines) > 20 {
		return strings.Join(result, "\n") + fmt.Sprintf("\n... (+ %d more lines)", len(lines)-20)
	}
	return strings.Join(result, "\n")
}
