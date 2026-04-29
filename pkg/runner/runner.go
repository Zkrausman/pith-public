package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pith/pkg/config"
	"pith/pkg/parser"
	"pith/pkg/telemetry"
	"strings"
	"time"
	"unicode/utf8"
)

type Runner struct {
	cfg       *config.Config
	telemetry *telemetry.Telemetry
	parsers   []parser.Parser
	Source    string
}

func NewRunner(cfg *config.Config, tel *telemetry.Telemetry) *Runner {
	return &Runner{
		cfg:       cfg,
		telemetry: tel,
		parsers:   parser.GetAllParsers(),
		Source:    DetectSource(),
	}
}

func DetectSource() string {
	if os.Getenv("GEMINI_CLI") != "" || os.Getenv("GOOGLE_API_KEY") != "" {
		return "gemini"
	}
	if os.Getenv("CLAUDE_CODE") != "" || os.Getenv("ANTHROPIC_API_KEY") != "" {
		return "claude"
	}
	// Check parent process names could be added here later
	return "unknown"
}

func (r *Runner) Run(args []string) error {
	return r.RunWithOptions(args, false)
}

func (r *Runner) LogForSnag(cmdStr string, output string, exitCode int) {
	logDir := r.cfg.StoragePath
	if logDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return // Silently fail if home dir can't be found
		}
		logDir = filepath.Join(home, ".pith")
	}

	_ = os.MkdirAll(logDir, 0755)

	logPath := filepath.Join(logDir, "pith.log")

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	// Truncate output for Snag to the last 50 lines to prevent log bloat
	lines := strings.Split(output, "\n")
	if len(lines) > 50 {
		output = "[... output truncated by Pith for Snag log ...]\n" + strings.Join(lines[len(lines)-50:], "\n")
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
			if len(parts) == 0 {
				return fmt.Errorf("no command provided")
			}
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

	originalTokens := r.EstimateTokens(fullOutput)

	var p parser.Parser
	if !skipParsing {
		cmdParts := strings.Fields(fullCmd)
		var cmdName string
		var cmdArgs []string
		if len(cmdParts) > 0 {
			cmdName = cmdParts[0]
			cmdArgs = cmdParts[1:]
		}

		// Special Case: ChainParser
		if strings.ContainsAny(fullCmd, ";|&><") {
			cp := &parser.ChainParser{}
			subcmds := cp.SplitSubCommands(fullCmd)
			if len(subcmds) > 1 {
				for _, sub := range subcmds {
					// Identify the best parser for this sub-command
					var bestP parser.Parser
					subParts := strings.Fields(sub)
					if len(subParts) == 0 {
						continue
					}
					subCmdName := subParts[0]
					subArgs := subParts[1:]

					for _, pCandidate := range r.parsers {
						enabled, ok := r.cfg.EnabledParsers[pCandidate.Name()]
						if (!ok || enabled) && pCandidate.CanParse(subCmdName, subArgs) {
							bestP = pCandidate
							break
						}
					}

					if bestP != nil {
						p = cp
						break
					}
				}
			}
		}

		if p == nil {
			for _, parser := range r.parsers {
				enabled, ok := r.cfg.EnabledParsers[parser.Name()]
				if (!ok || enabled) && parser.CanParse(cmdName, cmdArgs) {
					p = parser
					break
				}
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

	compressedTokens := r.EstimateTokens(finalOutput)

	// Output to stdout
	fmt.Print(finalOutput)

	// Record telemetry
	record := telemetry.ExecutionRecord{
		Command:           strings.Join(args, " "),
		OriginalTokens:    originalTokens,
		CompressedTokens:  compressedTokens,
		OriginalContent:   fullOutput,
		CompressedContent: finalOutput,
		DurationMs:        duration,
		ParserUsed:        parserUsed,
		IsPassthrough:     isPassthrough,
		Source:            r.Source,
	}

	_ = r.telemetry.Record(record)

	return err
}

func (r *Runner) ApplyMiddleOutTruncation(output string) string {
	lines := strings.Split(output, "\n")
	if len(lines) <= r.cfg.MaxLines { return output }
	head := r.cfg.HeadLines
	tail := r.cfg.TailLines
	if head+tail >= len(lines) { return output }
	middleStart := head
	middleEnd := len(lines) - tail
	hotZones := []int{}
	keywords := []string{"error", "fail", "fatal", "panic", "exception", "!!!", "error:", "[exit]"}
	for i := middleStart; i < middleEnd; i++ {
		lowerLine := strings.ToLower(lines[i])
		for _, k := range keywords {
			if strings.Contains(lowerLine, k) { hotZones = append(hotZones, i); break }
		}
	}
	resultLines := append([]string{}, lines[:head]...)
	lastIndex := head
	for _, zone := range hotZones {
		if zone > lastIndex+2 { resultLines = append(resultLines, fmt.Sprintf("\n... [%d lines of non-critical output removed by Pith] ...\n", zone-lastIndex-1)) }
		start := zone - 1
		if start < lastIndex { start = lastIndex }
		end := zone + 1
		if end >= middleEnd { end = middleEnd - 1 }
		for j := start; j <= end; j++ {
			if j >= lastIndex { resultLines = append(resultLines, lines[j]); lastIndex = j + 1 }
		}
	}
	if middleEnd > lastIndex { resultLines = append(resultLines, fmt.Sprintf("\n... [%d lines removed by Pith middle-out truncation] ...\n", middleEnd-lastIndex)) }
	resultLines = append(resultLines, lines[len(lines)-tail:]...)
	return strings.Join(resultLines, "\n")
}

func (r *Runner) EstimateTokens(s string) int {
	return EstimateTokensWithHeuristic(s, r.cfg.TokenHeuristic)
}

func EstimateTokensWithHeuristic(s string, heuristic float64) int {
	if heuristic <= 0 {
		heuristic = 4.0
	}
	return int(float64(utf8.RuneCountInString(s)) / heuristic)
}

// Deprecated: Use Runner.EstimateTokens or EstimateTokensWithHeuristic
func EstimateTokens(s string) int {
	return EstimateTokensWithHeuristic(s, 4.0)
}
