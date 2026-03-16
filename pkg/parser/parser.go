package parser

import (
	"strings"
)

type Parser interface {
	Name() string
	CanParse(cmd string, args []string) bool
	Parse(output string) string
}

type GitStatusParser struct{}

func (g *GitStatusParser) Name() string {
	return "git_status"
}

func (g *GitStatusParser) CanParse(cmd string, args []string) bool {
	if cmd == "git" && len(args) > 0 && args[0] == "status" {
		return true
	}
	return false
}

func (g *GitStatusParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Skip common git status boilerplate
		if strings.HasPrefix(trimmed, "(use ") || 
		   strings.HasPrefix(trimmed, "On branch") || 
		   strings.HasPrefix(trimmed, "Your branch") ||
		   strings.Contains(trimmed, "nothing to commit") ||
           strings.Contains(trimmed, "no changes added to commit") ||
           strings.HasPrefix(trimmed, "Changes not staged for commit") ||
           strings.HasPrefix(trimmed, "Changes to be committed") ||
           strings.HasPrefix(trimmed, "Untracked files") {
			continue
		}
		result = append(result, trimmed)
	}
	return strings.Join(result, "\n")
}
