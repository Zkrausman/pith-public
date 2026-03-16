package parser

import (
	"fmt"
	"strings"
)

// LsParser (Existing)
type LsParser struct{}

func (l *LsParser) Name() string { return "ls" }
func (l *LsParser) CanParse(cmd string, args []string) bool {
	return cmd == "ls" || cmd == "dir"
}
func (l *LsParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "total ") { continue }
		fields := strings.Fields(trimmed)
		if len(fields) > 0 {
			result = append(result, fields[len(fields)-1])
		}
	}
	return strings.Join(result, " ")
}

// FindParser (NEW)
type FindParser struct{}

func (f *FindParser) Name() string { return "find" }
func (f *FindParser) CanParse(cmd string, args []string) bool {
	return cmd == "find" || cmd == "where" || cmd == "where.exe"
}
func (f *FindParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	count := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }
		count++
		if count > 50 { continue } // Limit to 50 results
		result = append(result, trimmed)
	}
	res := strings.Join(result, "\n")
	if count > 50 {
		res += fmt.Sprintf("\n... (+ %d more results)", count-50)
	}
	return res
}

// TreeParser (NEW)
type TreeParser struct{}

func (t *TreeParser) Name() string { return "tree" }
func (t *TreeParser) CanParse(cmd string, args []string) bool {
	return cmd == "tree"
}
func (t *TreeParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" { continue }
		// Replace ASCII tree characters with simple spaces
		r := strings.NewReplacer("│", " ", "├", " ", "└", " ", "─", " ", "──", " ", "   ", "  ")
		cleaned := r.Replace(line)
		result = append(result, cleaned)
	}
	return strings.Join(result, "\n")
}

// DuParser (NEW)
type DuParser struct{}

func (d *DuParser) Name() string { return "du" }
func (d *DuParser) CanParse(cmd string, args []string) bool {
	return cmd == "du"
}
func (d *DuParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }
		fields := strings.Fields(trimmed)
		if len(fields) >= 2 {
			size := fields[0]
			path := fields[1]
			// Only show the size and the path
			result = append(result, fmt.Sprintf("%s\t%s", size, path))
		}
	}
	// If too many lines, show only top 20
	if len(result) > 20 {
		summary := result[:20]
		return strings.Join(summary, "\n") + fmt.Sprintf("\n... (+ %d more)", len(result)-20)
	}
	return strings.Join(result, "\n")
}
