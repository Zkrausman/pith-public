package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type GitHubReleaseParser struct{}

func (g *GitHubReleaseParser) Name() string {
	return "gh_release"
}

func (g *GitHubReleaseParser) CanParse(cmd string, args []string) bool {
	if cmd != "gh" {
		return false
	}
	if len(args) < 2 {
		return false
	}
	// gh release view OR gh release list
	return args[0] == "release" && (args[1] == "view" || args[1] == "list")
}

func (g *GitHubReleaseParser) Parse(output string) string {
	lines := strings.Split(output, "\n")
	var result []string

	reSHA := regexp.MustCompile(`sha256:[a-f0-9]{64}`)
	reAsset := regexp.MustCompile(`^(?P<name>\S+)\s+(?P<digest>\S+)\s+(?P<size>[\d.]+\s\S+)`)

	inAssets := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Handle Headers
		if strings.HasPrefix(trimmed, "Tag:") || strings.HasPrefix(trimmed, "Title:") || strings.HasPrefix(trimmed, "Latest:") {
			result = append(result, trimmed)
			continue
		}

		// Transition to Assets
		if strings.HasPrefix(trimmed, "Assets") {
			inAssets = true
			result = append(result, "Assets:")
			continue
		}

		if inAssets {
			// Extract just the name and size from the asset line, skipping the SHA
			matches := reAsset.FindStringSubmatch(trimmed)
			if len(matches) > 0 {
				name := matches[1]
				size := matches[3]
				// Skip the header itself if it was captured
				if name != "NAME" {
					result = append(result, fmt.Sprintf("- %s (%s)", name, size))
				}
				continue
			}
			// If we are in assets but it's not a match, it's either the header or end of section
			if strings.HasPrefix(trimmed, "View on GitHub:") {
				result = append(result, trimmed)
			}
			continue
		}

		// Capture the first few lines of the release body if present
		if !inAssets && len(result) < 10 && !strings.HasPrefix(trimmed, "github-actions") {
			// Strip SHAs from the body just in case
			cleaned := reSHA.ReplaceAllString(trimmed, "[SHA]")
			result = append(result, cleaned)
		}
	}

	return strings.Join(result, "\n")
}
