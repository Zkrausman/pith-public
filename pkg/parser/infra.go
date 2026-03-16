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
	return cmd == "env" || cmd == "set" || cmd == "export"
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
	return cmd == "docker" && len(args) > 0 && (args[0] == "ps" || args[0] == "images")
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
	return (cmd == "npm" && len(args) > 0 && args[0] == "list") ||
		   (cmd == "pip" && len(args) > 0 && args[0] == "list")
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
	return (cmd == "npm" && len(args) > 0 && args[0] == "test") ||
		   (cmd == "go" && len(args) > 0 && args[0] == "test") ||
		   cmd == "pytest"
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
