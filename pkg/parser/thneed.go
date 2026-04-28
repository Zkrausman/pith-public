package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ThneedParser struct{}

func (t *ThneedParser) Name() string { return "thneed" }

func (t *ThneedParser) CanParse(cmd string, args []string) bool {
	return MatchCommand(cmd, "thneed")
}

func (t *ThneedParser) Parse(output string) string {
	// Attempt to parse as JSON
	var results []interface{}
	if err := json.Unmarshal([]byte(output), &results); err == nil {
		return t.parseJson(results)
	}

	// Try single object JSON (error response or single result)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(output), &obj); err == nil {
		if errMsg, ok := obj["error"].(string); ok {
			return fmt.Sprintf("Thneed Error: %s", errMsg)
		}
		return t.parseJsonObject(obj)
	}

	// Fallback to standard line compression
	return t.parsePlain(output)
}

func (t *ThneedParser) parseJson(results []interface{}) string {
	if len(results) == 0 {
		return "Thneed: No results found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Thneed found %d nodes:\n", len(results)))

	// Limit to top 10 results to save tokens
	limit := len(results)
	if limit > 10 {
		limit = 10
	}

	for i := 0; i < limit; i++ {
		res, ok := results[i].(map[string]interface{})
		if !ok {
			continue
		}

		path, _ := res["path"].(string)
		id, _ := res["id"].(string)
		content, _ := res["content"].(string)

		// Shorten path for readability
		shortPath := path
		if idx := strings.LastIndex(path, "\\"); idx != -1 {
			shortPath = path[idx+1:]
		}

		// Truncate content to first 150 characters
		snippet := strings.TrimSpace(content)
		if len(snippet) > 150 {
			snippet = snippet[:147] + "..."
		}
		snippet = strings.ReplaceAll(snippet, "\n", " ")

		sb.WriteString(fmt.Sprintf("- [%s] (%s): %s\n", shortPath, id, snippet))
	}

	if len(results) > limit {
		sb.WriteString(fmt.Sprintf("... and %d more results hidden to save tokens.\n", len(results)-limit))
	}

	return sb.String()
}

func (t *ThneedParser) parseJsonObject(obj map[string]interface{}) string {
	// If it's a single result object (or metadata)
	path, _ := obj["path"].(string)
	id, _ := obj["id"].(string)
	if path != "" || id != "" {
		return fmt.Sprintf("Thneed Node: %s (%s)", path, id)
	}

	// Just return a simplified JSON summary if it's some other object
	return "Thneed: [Structured Metadata Object]"
}

func (t *ThneedParser) parsePlain(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "Indexing") || strings.HasPrefix(trimmed, "Scanning") {
			continue
		}
		result = append(result, trimmed)
	}
	// Middle-out truncation will handle large text anyway, so we just filter noise here
	if len(result) > 50 {
		return strings.Join(result[:25], "\n") + "\n...\n" + strings.Join(result[len(result)-25:], "\n")
	}
	return strings.Join(result, "\n")
}
