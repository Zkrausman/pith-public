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

	// High-signal keywords we want to preserve in comments
	highSignal := regexp.MustCompile(`(?i)(BUG|TODO|FIXME|NOTE|HACK|WARN|IMPORTANT|CRITICAL)`)

	// Track if we are inside a block comment
	inBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Handle block comments (/* ... */)
		if strings.Contains(trimmed, "/*") {
			// If the block comment is small or has high-signal info, keep it
			if strings.Contains(trimmed, "*/") {
				if highSignal.MatchString(trimmed) || len(trimmed) < 100 {
					result = append(result, line)
				}
				continue
			}
			inBlock = true
			if highSignal.MatchString(trimmed) {
				result = append(result, line)
			}
			continue
		}
		if inBlock {
			if strings.Contains(trimmed, "*/") {
				inBlock = false
				if highSignal.MatchString(trimmed) {
					result = append(result, line)
				}
			} else if highSignal.MatchString(trimmed) {
				result = append(result, line)
			}
			continue
		}

		// Process inline comments
		match := reInline.FindString(line)
		if match != "" {
			// If it's high signal, keep the whole line
			if highSignal.MatchString(match) {
				result = append(result, line)
				continue
			}
			// Otherwise, strip the comment but keep the code
			line = strings.TrimSpace(reInline.ReplaceAllString(line, ""))
		}

		// Collapse excessive spaces (Issue #582: balance savings vs context)
		// We preserve indentation (up to 8 spaces) but collapse huge gaps
		reHugeGaps := regexp.MustCompile(`\s{12,}`)
		cleaned := reHugeGaps.ReplaceAllString(line, "    ")

		if strings.TrimSpace(cleaned) != "" {
			result = append(result, cleaned)
		}
	}

	return strings.Join(result, "\n")
}
