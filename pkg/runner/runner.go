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
	return r.RunWithOptions(args, false)
}

func (r *Runner) RunWithOptions(args []string, skipParsing bool) error {
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
	if !skipParsing {
		for _, parser := range r.parsers {
			enabled, ok := r.cfg.EnabledParsers[parser.Name()]
			// If not in config, default to enabled. If in config, check boolean.
			if (!ok || enabled) && parser.CanParse(cmdName, cmdArgs) {
				p = parser
				break
			}
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

	// Apply Middle-Out Truncation
	finalOutput = r.ApplyMiddleOutTruncation(finalOutput)

	compressedTokens := estimateTokens(finalOutput)

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

func (r *Runner) ApplyMiddleOutTruncation(output string) string {
	lines := strings.Split(output, "\n")
	if len(lines) <= r.cfg.MaxLines {
		return output
	}

	head := r.cfg.HeadLines
	tail := r.cfg.TailLines

	// Safety check to ensure we don't try to keep more lines than exist
	if head+tail >= len(lines) {
		return output
	}

	resultLines := append([]string{}, lines[:head]...)
	resultLines = append(resultLines, fmt.Sprintf("\n... [%d lines removed by Diet middle-out truncation] ...\n", len(lines)-(head+tail)))
	resultLines = append(resultLines, lines[len(lines)-tail:]...)

	return strings.Join(resultLines, "\n")
}

func estimateTokens(s string) int {
	// Heuristic: 1 token ≈ 4 characters
	return len(s) / 4
}
