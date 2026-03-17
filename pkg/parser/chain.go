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
	return strings.Contains(fullCmd, ";") || strings.Contains(fullCmd, "&&")
}

func (c *ChainParser) Parse(output string) string {
	// The ChainParser is unique because the runner needs to split the input command
	// and run the sub-commands through their respective parsers.
	// This Parse method will be used for the combined output if available.
	return output // Runner handles the splitting logic for better accuracy
}

// SplitSubCommands splits a raw command string into individual sub-commands
func (c *ChainParser) SplitSubCommands(fullCmd string) []string {
	var subcmds []string
	
	// Split by ; or && (simplified shell logic)
	parts := strings.Split(fullCmd, ";")
	for _, p := range parts {
		subparts := strings.Split(p, "&&")
		for _, sp := range subparts {
			trimmed := strings.TrimSpace(sp)
			if trimmed != "" {
				subcmds = append(subcmds, trimmed)
			}
		}
	}
	return subcmds
}
