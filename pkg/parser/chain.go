package parser

import (
	"strings"
)

type ChainParser struct {
}

func (c *ChainParser) Name() string {
	return "chain"
}

func (c *ChainParser) CanParse(cmd string, args []string) bool {
	fullCmd := strings.Join(append([]string{cmd}, args...), " ")
	return strings.ContainsAny(fullCmd, ";|&")
}

func (c *ChainParser) Parse(output string) string {
	return output
}

func (c *ChainParser) SplitSubCommands(fullCmd string) []string {
	// Simple split by shell operators
	delimiters := []string{";", "&&", "||", "|"}
	subcmds := []string{fullCmd}

	for _, delim := range delimiters {
		var newSubcmds []string
		for _, s := range subcmds {
			parts := strings.Split(s, delim)
			for _, p := range parts {
				trimmed := strings.TrimSpace(p)
				if trimmed != "" {
					newSubcmds = append(newSubcmds, trimmed)
				}
			}
		}
		subcmds = newSubcmds
	}
	return subcmds
}
