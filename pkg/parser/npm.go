package parser

import (
	"strings"
)

type NPMParser struct{}

func (n *NPMParser) Name() string { return "npm" }

func (n *NPMParser) CanParse(cmd string, args []string) bool {
	return MatchCommand(cmd, "npm")
}

func (n *NPMParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Skip progress bars and common boilerplate
		if trimmed == "" || 
		   strings.HasPrefix(trimmed, "npm WARN") ||
		   strings.Contains(trimmed, "added") && strings.Contains(trimmed, "packages from") ||
		   strings.Contains(trimmed, "audited") && strings.Contains(trimmed, "packages in") ||
		   strings.Contains(trimmed, "found") && strings.Contains(trimmed, "vulnerabilities") {
			continue
		}
		
		// Skip the long list of packages in npm install
		if strings.HasPrefix(trimmed, "+ ") || strings.HasPrefix(trimmed, "- ") {
			continue
		}

		result = append(result, trimmed)
	}

	if len(result) > 40 {
		return strings.Join(result[:20], "\n") + "\n...\n" + strings.Join(result[len(result)-20:], "\n")
	}
	return strings.Join(result, "\n")
}
