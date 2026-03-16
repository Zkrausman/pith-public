package runner

import (
	"bytes"
	"diet/pkg/config"
	"diet/pkg/parser"
	"diet/pkg/telemetry"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Runner struct {
	cfg       *config.Config
	telemetry *telemetry.Telemetry
	parsers   []parser.Parser
}

func NewRunner(cfg *config.Config, tel *telemetry.Telemetry) *Runner {
	return &Runner{
		cfg:       cfg,
		telemetry: tel,
		parsers:   parser.GetAllParsers(),
	}
}

func (r *Runner) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	start := time.Now()
	
	cmd := exec.Command(cmdName, cmdArgs...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	duration := time.Since(start).Milliseconds()

	stdoutStr := out.String()
	stderrStr := stderr.String()
	fullOutput := stdoutStr + stderrStr
	
	originalTokens := estimateTokens(fullOutput)

	var p parser.Parser
	for _, parser := range r.parsers {
		if enabled, ok := r.cfg.EnabledParsers[parser.Name()]; ok && enabled && parser.CanParse(cmdName, cmdArgs) {
			p = parser
			break
		}
	}

	var finalOutput string
	parserUsed := "none"
	isPassthrough := true

	if p != nil {
		finalOutput = p.Parse(stdoutStr)
		if stderrStr != "" {
			finalOutput += "\nERRORS:\n" + stderrStr
		}
		parserUsed = p.Name()
		isPassthrough = false
	} else {
		finalOutput = fullOutput
	}

	compressedTokens := estimateTokens(finalOutput)

	// Enforce max tokens if needed (simplified truncation for now)
	if compressedTokens > r.cfg.MaxTokens {
		limit := r.cfg.MaxTokens * 4
		if len(finalOutput) > limit {
			finalOutput = finalOutput[:limit] + "\n...[Output Truncated by Diet]..."
		}
	}

	// Output to stdout
	fmt.Print(finalOutput)

	// Record telemetry
	record := telemetry.ExecutionRecord{
		Command:          strings.Join(args, " "),
		OriginalTokens:   originalTokens,
		CompressedTokens: compressedTokens,
		DurationMs:       duration,
		ParserUsed:       parserUsed,
		IsPassthrough:    isPassthrough,
	}
	
	_ = r.telemetry.Record(record)

	return err
}

func estimateTokens(s string) int {
	// Heuristic: 1 token ≈ 4 characters
	return len(s) / 4
}
