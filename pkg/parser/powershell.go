package parser

import (
	"fmt"
	"strings"
)

type PowerShellParser struct{}

func (ps *PowerShellParser) Name() string { return "windows-shell" }

func (ps *PowerShellParser) CanParse(cmd string, args []string) bool {
	return cmd == "powershell" || cmd == "pwsh" || cmd == "cmd"
}

func (ps *PowerShellParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	
	isLsh := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }
		
		// Headers to skip
		if strings.HasPrefix(trimmed, "Directory:") || strings.HasPrefix(trimmed, "Mode ") || strings.HasPrefix(trimmed, "----") {
			isLsh = true
			continue
		}

		if strings.Contains(trimmed, "Exit Code:") || strings.Contains(trimmed, "Process Group PGID:") {
			continue
		}
		
		// If it's a generic error from PowerShell (like CommandNotFoundException), minify it
		if strings.Contains(trimmed, "+ CategoryInfo") || strings.Contains(trimmed, "+ FullyQualifiedErrorId") {
			continue
		}
		
		// Handle Get-ChildItem / ls / dir output format
		// Mode                LastWriteTime         Length Name
		// ----                -------------         ------ ----
		// d-----        3/16/2026   1:55 AM                pkg
		if isLsh {
			fields := strings.Fields(trimmed)
			if len(fields) >= 4 {
				mode := fields[0]
				name := fields[len(fields)-1]
				
				// PowerShell output: Mode, LastWriteTime (Date, Time, AM/PM), Length (optional), Name
				// d-----        3/16/2026   1:55 AM                pkg (4 fields: mode, date, time, am/pm, name)
				// -a----        3/19/2026  12:23 AM           1990 web.go (5 fields: mode, date, time, am/pm, size, name)
				
				// If it's a directory, it usually has 4 fields (excluding Name) if we count AM/PM as a separate field.
				// Wait, let's look at the fields more carefully.
				// d----- | 3/16/2026 | 1:55 | AM | pkg  -> 5 fields
				// -a---- | 3/19/2026 | 12:23 | AM | 1990 | web.go -> 6 fields
				
				if len(fields) >= 6 {
					size := fields[len(fields)-2]
					result = append(result, fmt.Sprintf("%s %s %s", mode, size, name))
				} else {
					result = append(result, fmt.Sprintf("%s %s", mode, name))
				}
				continue
			}
		}

		if len(result) < 40 {
			result = append(result, trimmed)
		}
	}
	
	if len(result) == 0 {
		return ""
	}
	
	if len(lines) > 40 {
		return strings.Join(result, "\n") + fmt.Sprintf("\n... (+ %d more lines)", len(lines)-len(result))
	}
	return strings.Join(result, "\n")
}
