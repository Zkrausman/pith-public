package parser

import (
	"regexp"
	"strings"
)

// GrepParser (NEW)
type GrepParser struct{}

func (g *GrepParser) Name() string { return "grep" }
func (g *GrepParser) CanParse(cmd string, args []string) bool {
	return cmd == "grep" || cmd == "rg" || cmd == "ripgrep"
}
func (g *GrepParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	currentFile := ""
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }
		// Try to detect file:line:content
		parts := strings.SplitN(trimmed, ":", 3)
		if len(parts) >= 2 {
			file := parts[0]
			if file != currentFile {
				result = append(result, "\n"+file+":")
				currentFile = file
			}
			if len(parts) == 3 {
				result = append(result, "  "+parts[1]+": "+parts[2])
			} else {
				result = append(result, "  "+parts[1])
			}
		} else {
			result = append(result, trimmed)
		}
	}
	return strings.TrimSpace(strings.Join(result, "\n"))
}

// MinifyParser (NEW)
type MinifyParser struct{}

func (m *MinifyParser) Name() string { return "minify" }
func (m *MinifyParser) CanParse(cmd string, args []string) bool {
	if cmd != "cat" && cmd != "type" { return false }
	if len(args) == 0 { return false }
	file := args[0]
	return strings.HasSuffix(file, ".json") || strings.HasSuffix(file, ".xml") || 
		   strings.HasSuffix(file, ".html") || strings.HasSuffix(file, ".css")
}

var whitespaceRegex = regexp.MustCompile(`\s+`)

func (m *MinifyParser) Parse(output string) string {
	// Simple minification: strip newlines and collapse spaces
	res := whitespaceRegex.ReplaceAllString(output, " ")
	return strings.TrimSpace(res)
}
