package parser

import (
	"fmt"
	"strings"
)

type PromptfooParser struct{}

func (p *PromptfooParser) Name() string { return "promptfoo" }

func (p *PromptfooParser) CanParse(cmd string, args []string) bool {
	isNpx := MatchCommand(cmd, "npx") && len(args) > 0 && (args[0] == "promptfoo" || MatchCommand(args[0], "promptfoo"))
	isPromptfoo := MatchCommand(cmd, "promptfoo")
	return isNpx || isPromptfoo
}

func (p *PromptfooParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }
		
		// Skip table borders
		if strings.HasPrefix(trimmed, "┌") || strings.HasPrefix(trimmed, "├") || 
		   strings.HasPrefix(trimmed, "└") || strings.Contains(trimmed, "───") {
			continue
		}
		
		// Progress bars
		if strings.Contains(trimmed, "████") || (strings.Contains(trimmed, "%") && strings.Contains(trimmed, "|")) {
			result = append(result, trimmed)
			continue
		}
		
		// Summary lines
		if strings.Contains(trimmed, "Starting evaluation") || strings.Contains(trimmed, "Evaluating") ||
		   strings.Contains(trimmed, "Tests") || strings.Contains(trimmed, "Duration") {
			result = append(result, trimmed)
			continue
		}
		
		// Handle table content (lines with pipe separators)
		if strings.Contains(trimmed, "│") {
			// Clean up pipe characters and extra spaces
			parts := strings.Split(trimmed, "│")
			var cleanedParts []string
			for _, part := range parts {
				p := strings.TrimSpace(part)
				if p != "" {
					cleanedParts = append(cleanedParts, p)
				}
			}
			if len(cleanedParts) > 0 {
				result = append(result, strings.Join(cleanedParts, " | "))
			}
			continue
		}

		// Limit the number of lines
		if len(result) < 40 {
			result = append(result, trimmed)
		}
	}
	
	if len(lines) > 40 && len(result) >= 40 {
		return strings.Join(result, "\n") + fmt.Sprintf("\n... (+ %d more lines)", len(lines)-len(result))
	}
	return strings.Join(result, "\n")
}
