package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// EnvParser (Existing)
type EnvParser struct{}

func (e *EnvParser) Name() string { return "env" }
func (e *EnvParser) CanParse(cmd string, args []string) bool {
	return MatchCommand(cmd, "env") || MatchCommand(cmd, "set") || MatchCommand(cmd, "export")
}

var skipEnvRegex = regexp.MustCompile(`(?i)(token|key|secret|password|auth|ssh_auth_sock|ls_colors|color)`)

func (e *EnvParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		if line == "" || !strings.Contains(line, "=") { continue }
		parts := strings.SplitN(line, "=", 2)
		key := parts[0]
		if skipEnvRegex.MatchString(key) { continue }
		val := parts[1]
		if len(val) > 100 { val = val[:100] + "..." }
		result = append(result, key+"="+val)
	}
	return strings.Join(result, "\n")
}

// DockerPsParser (NEW)
type DockerPsParser struct{}

func (d *DockerPsParser) Name() string { return "docker_ps" }
func (d *DockerPsParser) CanParse(cmd string, args []string) bool {
	return MatchCommand(cmd, "docker") && len(args) > 0 && (args[0] == "ps" || args[0] == "images")
}
func (d *DockerPsParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }
		// Docker output is column-based. We extract the first and last few columns.
		fields := strings.Fields(trimmed)
		if len(fields) >= 2 {
			// ID + Name/Repo
			result = append(result, fmt.Sprintf("%s\t%s", fields[0], fields[len(fields)-1]))
		}
	}
	return strings.Join(result, "\n")
}

// DependencyParser (NEW)
type DependencyParser struct{}

func (d *DependencyParser) Name() string { return "dependencies" }
func (d *DependencyParser) CanParse(cmd string, args []string) bool {
	return (MatchCommand(cmd, "npm") && len(args) > 0 && args[0] == "list") ||
		   (MatchCommand(cmd, "pip") && len(args) > 0 && args[0] == "list")
}
func (d *DependencyParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	count := 0
	
	// Characters to strip from the beginning of dependency lines
	replacer := strings.NewReplacer("├── ", "", "└── ", "", "│   ", "", "├─ ", "", "└─ ", "")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.Contains(trimmed, "empty") { continue }
		
		cleaned := replacer.Replace(trimmed)
		if cleaned == trimmed && (strings.Contains(trimmed, "├──") || strings.Contains(trimmed, "└──")) {
			// Fallback if formatting is slightly different
			cleaned = strings.Map(func(r rune) rune {
				if r == '├' || r == '─' || r == '└' || r == '│' {
					return -1
				}
				return r
			}, trimmed)
		}
		
		result = append(result, strings.TrimSpace(cleaned))
		count++
		if count > 30 { break }
	}
	res := strings.Join(result, "\n")
	if len(lines) > 30 {
		res += "\n... (truncated dependency list)"
	}
	return res
}

// TestParser (NEW)
type TestParser struct{}

func (t *TestParser) Name() string { return "tests" }
func (t *TestParser) CanParse(cmd string, args []string) bool {
	return (MatchCommand(cmd, "npm") && len(args) > 0 && args[0] == "test") ||
		   (MatchCommand(cmd, "go") && len(args) > 0 && args[0] == "test") ||
		   MatchCommand(cmd, "pytest")
}
func (t *TestParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	isFailureBlock := false
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Capture summary lines
		if strings.Contains(trimmed, "passed") || strings.Contains(trimmed, "failed") || 
		   strings.Contains(trimmed, "TOTAL") || strings.HasPrefix(trimmed, "DONE") {
			result = append(result, trimmed)
		}
		// Capture failure details (heuristic)
		if strings.Contains(line, "FAIL") || strings.Contains(line, "Error:") {
			isFailureBlock = true
		}
		if isFailureBlock {
			result = append(result, line)
			if trimmed == "" { isFailureBlock = false } // End of failure block
		}
	}
	if len(result) == 0 { return "Tests finished. (No summary captured)" }
return strings.Join(result, "\n")
}

// GitHubParser (NEW)
type GitHubParser struct{}

func (g *GitHubParser) Name() string { return "github" }
func (g *GitHubParser) CanParse(cmd string, args []string) bool {
	return MatchCommand(cmd, "gh") && len(args) > 0 && 
		   (args[0] == "issue" || args[0] == "pr" || args[0] == "release" || args[0] == "repo" || args[0] == "run") &&
		   (strings.Contains(strings.Join(args, " "), " list") || strings.Contains(strings.Join(args, " "), " view"))
}
func (g *GitHubParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }
		
		// Skip header lines if present
		if strings.HasPrefix(trimmed, "Showing ") || strings.HasPrefix(trimmed, "NAME") || strings.HasPrefix(trimmed, "TITLE") {
			continue
		}

		fields := strings.Fields(trimmed)
		if len(fields) >= 2 {
			// For issues/PRs/repos: ID/NAME  TITLE/DESCRIPTION  STATUS/DATE
			id := fields[0]
			title := fields[1]
			
			status := ""
			if len(fields) >= 3 {
				status = fields[2]
				if len(status) > 15 { status = status[:15] }
			}
			
			if status != "" {
				result = append(result, fmt.Sprintf("%s | %s | %s", id, title, status))
			} else {
				result = append(result, fmt.Sprintf("%s | %s", id, title))
			}
		} else {
			result = append(result, trimmed)
		}
		
		if len(result) > 25 {
			result = append(result, "... (truncated gh list)")
			break
		}
	}
	
	return strings.Join(result, "\n")
}
