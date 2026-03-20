package parser

import (
	"fmt"
	"strings"
)

type BDParser struct{}

func (b *BDParser) Name() string { return "bd" }

func (b *BDParser) CanParse(cmd string, args []string) bool {
	return cmd == "bd" || strings.HasSuffix(cmd, "bd.exe")
}

func (b *BDParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	
	// Keep essential information, strip help or long lists
	isHelp := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }
		
		// If it's help output (very common in discovery list)
		if strings.HasPrefix(trimmed, "Usage:") || strings.Contains(trimmed, "Available Commands:") || strings.Contains(trimmed, "Flags:") {
			isHelp = true
		}
		
		// Summary lines (e.g., Total: 1 issues (1 open, 0 in progress))
		if strings.HasPrefix(trimmed, "Total:") || strings.HasPrefix(trimmed, "Ready:") || strings.Contains(trimmed, "Summary:") {
			result = append(result, trimmed)
			continue
		}
		
		// Keep issue lines (e.g., ○ pith-rpm ● P1 Test issue for parser)
		if strings.HasPrefix(trimmed, "○ ") || strings.HasPrefix(trimmed, "◐ ") || strings.HasPrefix(trimmed, "● ") || strings.HasPrefix(trimmed, "✓ ") || strings.HasPrefix(trimmed, "❄ ") {
			// But only if we don't have too many already
			if len(result) < 30 {
				result = append(result, trimmed)
			}
			continue
		}
		
		// If it's a short help command or info
		if len(result) < 15 {
			// Truncate long lines
			if len(trimmed) > 100 {
				result = append(result, trimmed[:100] + "...")
			} else {
				result = append(result, trimmed)
			}
		}
	}
	
	if isHelp && len(lines) > 20 {
		// Just show the top and summary
		return strings.Join(result, "\n") + "\n... (truncated bd help)"
	}
	
	if len(lines) > 30 {
		return strings.Join(result, "\n") + fmt.Sprintf("\n... (+ %d more lines)", len(lines)-len(result))
	}
	
	return strings.Join(result, "\n")
}
