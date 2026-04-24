package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SnagParser struct{}

func (s *SnagParser) Name() string { return "snag" }

func (s *SnagParser) CanParse(cmd string, args []string) bool {
	return MatchCommand(cmd, "snag")
}

func (s *SnagParser) Parse(output string) string {
	var snags []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &snags); err != nil {
		// If not JSON or single object, fallback to plain text filtering
		return s.parsePlain(output)
	}

	if len(snags) == 0 {
		return "Snag: No snags found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Snag found %d issues:\n", len(snags)))

	// Limit to top 15 snags
	limit := len(snags)
	if limit > 15 {
		limit = 15
	}

	for i := 0; i < limit; i++ {
		snag := snags[i]
		id, _ := snag["id"].(string)
		cmd, _ := snag["failed_command"].(string)
		advice, _ := snag["advice"].(string)
		errMsg, _ := snag["error_message"].(string)
		status, _ := snag["status"].(string)

		sb.WriteString(fmt.Sprintf("- [%s] %s: %s\n", status, id, cmd))
		if errMsg != "" {
			// Truncate error message if too long
			if len(errMsg) > 200 {
				errMsg = errMsg[:197] + "..."
			}
			sb.WriteString(fmt.Sprintf("  Error: %s\n", errMsg))
		}
		if advice != "" {
			sb.WriteString(fmt.Sprintf("  Advice: %s\n", advice))
		}
	}

	if len(snags) > limit {
		sb.WriteString(fmt.Sprintf("... and %d more snags hidden.\n", len(snags)-limit))
	}

	return sb.String()
}

func (s *SnagParser) parsePlain(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Basic noise reduction for non-JSON snag output if it exists
		result = append(result, trimmed)
	}
	if len(result) > 40 {
		return strings.Join(result[:20], "\n") + "\n...\n" + strings.Join(result[len(result)-20:], "\n")
	}
	return strings.Join(result, "\n")
}
