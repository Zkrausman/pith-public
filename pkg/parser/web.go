package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

type WebParser struct{}

func (w *WebParser) Name() string { return "web-content" }

func (w *WebParser) CanParse(cmd string, args []string) bool {
	// curl, wget, Invoke-WebRequest, iwr
	return MatchCommand(cmd, "curl") || MatchCommand(cmd, "wget") || MatchCommand(cmd, "Invoke-WebRequest") || MatchCommand(cmd, "iwr")
}

func (w *WebParser) Parse(output string) string {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return output
	}

	// Try JSON first
	var js interface{}
	if err := json.Unmarshal([]byte(trimmed), &js); err == nil {
		minified, _ := json.Marshal(js)
		if len(minified) < len(trimmed) {
			// If it's a large JSON object, maybe just show the top-level keys
			if m, ok := js.(map[string]interface{}); ok && len(minified) > 500 {
				keys := make([]string, 0, len(m))
				for k := range m {
					keys = append(keys, k)
				}
				return fmt.Sprintf("JSON Object Keys: [%s] (Total minified: %d chars)",
					strings.Join(keys, ", "), len(minified))
			}
			return string(minified)
		}
	}

	// Try HTML basic extraction
	if strings.Contains(trimmed, "<html") || strings.Contains(trimmed, "<!DOCTYPE html") {
		// Just extract the title if possible
		titleStart := strings.Index(trimmed, "<title>")
		titleEnd := strings.Index(trimmed, "</title>")
		if titleStart != -1 && titleEnd != -1 && titleEnd > titleStart {
			title := trimmed[titleStart+7 : titleEnd]
			return fmt.Sprintf("HTML Content: [%s] (%d chars total)", strings.TrimSpace(title), len(trimmed))
		}
		return fmt.Sprintf("HTML Content (%d chars total)", len(trimmed))
	}

	// Default: if it's very long, summarize
	if len(trimmed) > 1000 {
		return fmt.Sprintf("%s\n... (Total: %d chars)", trimmed[:500], len(trimmed))
	}

	return output
}
