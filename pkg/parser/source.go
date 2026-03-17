package parser

import (
	"regexp"
	"strings"
)

type SourceParser struct{}

func (s *SourceParser) Name() string {
	return "source"
}

func (s *SourceParser) CanParse(cmd string, args []string) bool {
	return cmd == "cat" || cmd == "type"
}

func (s *SourceParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string

	// Regex for common comment styles (Go, JS, C, Python, etc.)
	reInline := regexp.MustCompile(`(//.*|#.*)`)
	
	// Track if we are inside a block comment
	inBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Handle block comments (/* ... */)
		if strings.Contains(trimmed, "/*") {
			inBlock = true
			if strings.Contains(trimmed, "*/") {
				inBlock = false
				continue // Simplified: skip the whole line if it contains a full block
			}
			continue
		}
		if inBlock {
			if strings.Contains(trimmed, "*/") {
				inBlock = false
			}
			continue
		}

		// Strip inline comments but preserve line structure
		cleaned := reInline.ReplaceAllString(line, "")
		
		// Collapse multiple spaces to one (Issue #582: balance savings vs context)
		// We keep some indentation for structure but remove excessive padding
		reSpaces := regexp.MustCompile(`\s{4,}`)
		cleaned = reSpaces.ReplaceAllString(cleaned, "  ")

		if strings.TrimSpace(cleaned) != "" {
			result = append(result, cleaned)
		}
	}

	return strings.Join(result, "\n")
}
