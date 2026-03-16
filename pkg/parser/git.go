package parser

import (
	"fmt"
	"strings"
)

// GitStatusParser (Existing)
type GitStatusParser struct{}

func (g *GitStatusParser) Name() string { return "git_status" }
func (g *GitStatusParser) CanParse(cmd string, args []string) bool {
	return cmd == "git" && len(args) > 0 && args[0] == "status"
}
func (g *GitStatusParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "(use ") ||
			strings.HasPrefix(trimmed, "On branch") || strings.HasPrefix(trimmed, "Your branch") ||
			strings.Contains(trimmed, "nothing to commit") || strings.Contains(trimmed, "no changes added") ||
			strings.HasPrefix(trimmed, "Changes not staged") || strings.HasPrefix(trimmed, "Changes to be committed") ||
			strings.HasPrefix(trimmed, "Untracked files") {
			continue
		}
		result = append(result, trimmed)
	}
	return strings.Join(result, "\n")
}

// GitLogParser (Existing)
type GitLogParser struct{}

func (g *GitLogParser) Name() string { return "git_log" }
func (g *GitLogParser) CanParse(cmd string, args []string) bool {
	return cmd == "git" && len(args) > 0 && args[0] == "log"
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
			if len(currentCommit) > 7 { currentCommit = currentCommit[:7] }
			currentAuthor, currentDate, currentSubject = "", "", ""
		} else if strings.HasPrefix(line, "Author: ") {
			currentAuthor = strings.TrimPrefix(line, "Author: ")
			if idx := strings.Index(currentAuthor, " <"); idx != -1 { currentAuthor = currentAuthor[:idx] }
		} else if strings.HasPrefix(line, "Date: ") {
			currentDate = strings.TrimPrefix(line, "Date: ")
			fields := strings.Fields(currentDate)
			if len(fields) >= 3 { currentDate = fields[1] + " " + fields[2] + " " + fields[4] }
		} else if strings.HasPrefix(line, "    ") && currentSubject == "" {
			currentSubject = strings.TrimSpace(line)
		}
	}
	if currentCommit != "" {
		result = append(result, formatCommit(currentCommit, currentAuthor, currentDate, currentSubject))
	}
	return strings.Join(result, "\n")
}

func formatCommit(h, a, d, s string) string { return h + " | " + a + " | " + d + " | " + s }

// GitDiffParser (NEW)
type GitDiffParser struct{}

func (g *GitDiffParser) Name() string { return "git_diff" }
func (g *GitDiffParser) CanParse(cmd string, args []string) bool {
	return cmd == "git" && len(args) > 0 && args[0] == "diff"
}
func (g *GitDiffParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		if strings.HasPrefix(line, "index ") || strings.HasPrefix(line, "diff --git") {
			continue
		}
		// Condense hunk headers: @@ -1,4 +1,4 @@ -> @@
		if strings.HasPrefix(line, "@@") {
			result = append(result, "@@")
			continue
		}
		// Simplify file markers
		if strings.HasPrefix(line, "--- a/") {
			result = append(result, "--- "+strings.TrimPrefix(line, "--- a/"))
			continue
		}
		if strings.HasPrefix(line, "+++ b/") {
			result = append(result, "+++ "+strings.TrimPrefix(line, "+++ b/"))
			continue
		}
		// Only keep changes and minimal context if needed, but here we keep all changed lines
		if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// GitBranchParser (NEW)
type GitBranchParser struct{}

func (g *GitBranchParser) Name() string { return "git_branch" }
func (g *GitBranchParser) CanParse(cmd string, args []string) bool {
	return cmd == "git" && len(args) > 0 && args[0] == "branch"
}
func (g *GitBranchParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var current string
	var others []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }
		if strings.HasPrefix(trimmed, "*") {
			current = trimmed
		} else {
			others = append(others, trimmed)
		}
	}
	res := current
	if len(others) > 0 {
		res += fmt.Sprintf("\n(+ %d other branches)", len(others))
	}
	return res
}
