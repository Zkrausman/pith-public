package parser

import (
	"regexp"
	"strings"
)

type Parser interface {
	Name() string
	CanParse(cmd string, args []string) bool
	Parse(output string) string
}

func GetAllParsers() []Parser {
	return []Parser{
		&GitStatusParser{},
		&LsParser{},
		&GitLogParser{},
		&EnvParser{},
	}
}

// GitStatusParser
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

// LsParser
type LsParser struct{}

func (l *LsParser) Name() string {
	return "ls"
}

func (l *LsParser) CanParse(cmd string, args []string) bool {
	return cmd == "ls" || cmd == "dir"
}

func (l *LsParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "total ") {
			continue
		}
		// Handle 'ls -l' style: extract last field
		fields := strings.Fields(trimmed)
		if len(fields) >= 9 { // Typically 9 fields in 'ls -l'
			result = append(result, fields[len(fields)-1])
		} else if len(fields) > 0 {
			// Single file per line or other format
			result = append(result, fields[len(fields)-1])
		}
	}
	return strings.Join(result, " ")
}

// GitLogParser
type GitLogParser struct{}

func (g *GitLogParser) Name() string {
	return "git_log"
}

func (g *GitLogParser) CanParse(cmd string, args []string) bool {
	if cmd == "git" && len(args) > 0 && args[0] == "log" {
		return true
	}
	return false
}

func (g *GitLogParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	var currentCommit, currentAuthor, currentDate, currentSubject string

	for _, line := range lines {
		if strings.HasPrefix(line, "commit ") {
			if currentCommit != "" {
				result = append(result, formatCommit(currentCommit, currentAuthor, currentDate, currentSubject))
			}
			currentCommit = strings.TrimPrefix(line, "commit ")
			if len(currentCommit) > 7 {
				currentCommit = currentCommit[:7]
			}
			currentAuthor, currentDate, currentSubject = "", "", ""
		} else if strings.HasPrefix(line, "Author: ") {
			currentAuthor = strings.TrimPrefix(line, "Author: ")
			if idx := strings.Index(currentAuthor, " <"); idx != -1 {
				currentAuthor = currentAuthor[:idx]
			}
		} else if strings.HasPrefix(line, "Date: ") {
			currentDate = strings.TrimPrefix(line, "Date: ")
			fields := strings.Fields(currentDate)
			if len(fields) >= 3 {
				currentDate = fields[1] + " " + fields[2] + " " + fields[4] // Mon Mar 2026
			}
		} else if strings.HasPrefix(line, "    ") && currentSubject == "" {
			currentSubject = strings.TrimSpace(line)
		}
	}
	if currentCommit != "" {
		result = append(result, formatCommit(currentCommit, currentAuthor, currentDate, currentSubject))
	}
	return strings.Join(result, "\n")
}

func formatCommit(h, a, d, s string) string {
	return h + " | " + a + " | " + d + " | " + s
}

// EnvParser
type EnvParser struct{}

func (e *EnvParser) Name() string {
	return "env"
}

func (e *EnvParser) CanParse(cmd string, args []string) bool {
	return cmd == "env" || cmd == "set" || cmd == "export"
}

var skipEnvRegex = regexp.MustCompile(`(?i)(token|key|secret|password|auth|ssh_auth_sock|ls_colors|color)`)

func (e *EnvParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		if line == "" || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := parts[0]
		if skipEnvRegex.MatchString(key) {
			continue
		}
		val := parts[1]
		if len(val) > 100 {
			val = val[:100] + "..."
		}
		result = append(result, key+"="+val)
	}
	return strings.Join(result, "\n")
}
