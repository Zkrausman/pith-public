package runner

import (
	"bytes"
	"diet/pkg/config"
	"diet/pkg/parser"
	"diet/pkg/telemetry"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func (r *Runner) LogForSnag(cmdStr string, output string, exitCode int) {
	home, err := os.UserHomeDir()
	if err != nil {
		return // Silently fail if home dir can't be found
	}
	
	logDir := filepath.Join(home, ".diet")
	_ = os.MkdirAll(logDir, 0755)
	
	logPath := filepath.Join(logDir, "diet.log")
	
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	// Truncate output for Snag to the last 50 lines to prevent log bloat
	lines := strings.Split(output, "\n")
	if len(lines) > 50 {
		output = "[... output truncated by Diet for Snag log ...]\n" + strings.Join(lines[len(lines)-50:], "\n")
	}

	// Format: [CMD] command\n<output>\n[EXIT] code\n
	entry := fmt.Sprintf("[CMD] %s\n", cmdStr)
	if output != "" {
		// Ensure output ends with a newline so [EXIT] is on its own line
		if !strings.HasSuffix(output, "\n") {
			output += "\n"
		}
		entry += output
	}
	entry += fmt.Sprintf("[EXIT] %d\n", exitCode)

	_, _ = f.WriteString(entry)
}

func (r *Runner) RunWithOptions(args []string, skipParsing bool) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	fullCmd := strings.Join(args, " ")
	start := time.Now()
	
	var cmd *exec.Cmd
	// If it's a composite command or has shell redirects, run through shell
	if strings.ContainsAny(fullCmd, ";|&><") {
		// Windows
		cmd = exec.Command("cmd", "/c", fullCmd)
		// On Linux/macOS you'd use "sh", "-c", fullCmd
	} else {
		// IMPORTANT: If 'args' has more than one element, they are the arguments.
		// If 'args' has only one element but it contains spaces, it's a combined string 
		// that needs splitting (often happens when called via proxy).
		if len(args) == 1 && strings.Contains(args[0], " ") {
			parts := strings.Fields(args[0])
			cmd = exec.Command(parts[0], parts[1:]...)
		} else {
			cmd = exec.Command(args[0], args[1:]...)
		}
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	duration := time.Since(start).Milliseconds()
	
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1 // Generic failure
		}
	}

	stdoutStr := out.String()
	stderrStr := stderr.String()
	fullOutput := stdoutStr + stderrStr
	
	// Log for Snag before doing any compression/truncation
	r.LogForSnag(fullCmd, fullOutput, exitCode)

	originalTokens := estimateTokens(fullOutput)

	var p parser.Parser
	if !skipParsing {
		for _, parser := range r.parsers {
			enabled, ok := r.cfg.EnabledParsers[parser.Name()]
			if (!ok || enabled) && parser.CanParse(fullCmd, []string{}) {
				p = parser
				break
			}
		}
	}

	var finalOutput string
	parserUsed := "none"
	isPassthrough := true

	if p != nil {
		// Pass fullOutput to parser instead of just stdout
		finalOutput = p.Parse(fullOutput)
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
