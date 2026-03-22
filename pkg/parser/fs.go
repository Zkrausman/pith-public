package parser

import (
	"fmt"
	"strings"
)

// LsParser (Existing)
type LsParser struct{}

func (l *LsParser) Name() string { return "ls" }
func (l *LsParser) CanParse(cmd string, args []string) bool {
	return MatchCommand(cmd, "ls") || MatchCommand(cmd, "dir")
}
func (l *LsParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "total ") || strings.HasPrefix(trimmed, "Directory:") || strings.HasPrefix(trimmed, "Mode") || strings.HasPrefix(trimmed, "----") { 
			continue 
		}
		
		fields := strings.Fields(trimmed)
		if len(fields) >= 4 {
			// For standard ls -l: mode (0), links (1), user (2), group (3), size (4), month (5), day (6), time (7), name (8)
			// For Windows dir: mode (0), date (1), time (2), size (3), name (4)
			
			// Simple heuristic: if field 0 looks like a mode (-rwx... or d----)
			if strings.HasPrefix(fields[0], "d") || strings.HasPrefix(fields[0], "-") || strings.HasPrefix(fields[0], "l") {
				name := fields[len(fields)-1]
				size := ""
				
				if len(fields) >= 9 { // Standard Linux ls -l
					size = fields[4]
				} else if len(fields) >= 5 { // Windows Mode/Date/Time/Size/Name
					size = fields[3]
				}
				
				if size != "" {
					result = append(result, fmt.Sprintf("%s %s", size, name))
				} else {
					result = append(result, name)
				}
			} else {
				// Fallback to name only if we can't parse it reliably
				result = append(result, fields[len(fields)-1])
			}
		} else if len(fields) > 0 {
			result = append(result, fields[len(fields)-1])
		}
	}
	
	// If it's a long list, return as lines, otherwise join as space-separated
	if len(result) > 5 {
		return strings.Join(result, "\n")
	}
	return strings.Join(result, " ")
}

// FindParser (NEW)
type FindParser struct{}

func (f *FindParser) Name() string { return "find" }
func (f *FindParser) CanParse(cmd string, args []string) bool {
	return MatchCommand(cmd, "find") || MatchCommand(cmd, "where")
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
	return MatchCommand(cmd, "tree")
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
	return MatchCommand(cmd, "du")
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
