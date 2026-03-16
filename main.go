package main

import (
	"diet/pkg/config"
	"diet/pkg/parser"
	"diet/pkg/runner"
	"diet/pkg/telemetry"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "diet [command]",
	Short: "Diet is a token-optimized CLI proxy",
	Long:  `Diet intercepts terminal commands, compresses their output, and filters out noise to save tokens for LLMs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		// Check if first arg is a known subcommand
		subcmds := []string{"config", "gain", "discover"}
		isSub := false
		for _, s := range subcmds {
			if args[0] == s {
				isSub = true
				break
			}
		}

		if isSub {
			// This part shouldn't really be reached if Cobra matches subcommands,
			// but just in case or if we want to force routing.
			return nil
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		tel, err := telemetry.NewTelemetry()
		if err != nil {
			return err
		}
		defer tel.Close()

		run := runner.NewRunner(cfg, tel)
		return run.Run(args)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Interactive configuration tool",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		return cfg.InteractiveConfig()
	},
}

var gainCmd = &cobra.Command{
	Use:   "gain",
	Short: "Show token savings analytics",
	RunE: func(cmd *cobra.Command, args []string) error {
		tel, err := telemetry.NewTelemetry()
		if err != nil {
			return err
		}
		defer tel.Close()

		totalOrig, totalComp, err := tel.GetStats()
		if err != nil {
			return err
		}

		fmt.Printf("--- Diet Token Savings Report ---\n")
		fmt.Printf("Total Raw Tokens:        %d\n", totalOrig)
		fmt.Printf("Total Compressed Tokens: %d\n", totalComp)
		fmt.Printf("Total Tokens Saved:      %d (%.2f%%)\n", totalOrig-totalComp, float64(totalOrig-totalComp)/float64(totalOrig)*100)
		fmt.Println()

		byCmd, err := tel.GetStatsByCommand()
		if err == nil {
			fmt.Printf("%-30s | %-10s | %-10s | %-10s\n", "Command Pattern", "Raw", "Diet", "Savings")
			fmt.Println(strings.Repeat("-", 66))
			for _, r := range byCmd {
				savings := r.Original - r.Compressed
				fmt.Printf("%-30s | %-10d | %-10d | %-10d\n", r.Command, r.Original, r.Compressed, savings)
			}
		}
		return nil
	},
}

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Identify commands that could benefit from dedicated parsers",
	RunE: func(cmd *cobra.Command, args []string) error {
		tel, err := telemetry.NewTelemetry()
		if err != nil {
			return err
		}
		defer tel.Close()

		unparsed, err := tel.GetUnparsedCommands()
		if err != nil {
			return err
		}

		fmt.Printf("--- Opportunity Discovery (Unparsed Commands) ---\n")
		fmt.Printf("%-30s | %-10s | %-15s | %-15s\n", "Command Pattern", "Count", "Total Tokens", "Est. Savings (70%)")
		fmt.Println(strings.Repeat("-", 76))

		for _, r := range unparsed {
			estSavings := float64(r.TotalRawTokens) * 0.7 // Assume 70% average savings
			fmt.Printf("%-30s | %-10d | %-15d | %-15.0f\n", r.Pattern, r.InvocationCount, r.TotalRawTokens, estSavings)
		}

		return nil
	},
}

// Internal Hook Input/Output schemas
type HookInput struct {
	ToolResponse struct {
		LlmContent string `json:"llmContent"`
	} `json:"tool_response"`
	ToolCallRequest struct {
		Arguments string `json:"arguments"`
	} `json:"tool_call_request"`
}

type HookOutput struct {
	Decision      string `json:"decision"`
	Reason        string `json:"reason"`
	SystemMessage string `json:"systemMessage"`
}

type ToolArgs struct {
	Command string `json:"command"`
}

var hookCmd = &cobra.Command{
	Use:    "_hook",
	Short:  "Internal hook for Gemini CLI",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputData, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		var input HookInput
		if err := json.Unmarshal(inputData, &input); err != nil {
			return err
		}

		var toolArgs ToolArgs
		if err := json.Unmarshal([]byte(input.ToolCallRequest.Arguments), &toolArgs); err != nil {
			return err
		}

		cfg, _ := config.LoadConfig()
		tel, _ := telemetry.NewTelemetry()
		if tel != nil {
			defer tel.Close()
		}

		originalOutput := input.ToolResponse.LlmContent
		cmdParts := strings.Fields(toolArgs.Command)
		if len(cmdParts) == 0 {
			return respondAllow()
		}

		// Use the parser directly
		parsers := []parser.Parser{&parser.GitStatusParser{}}
		var p parser.Parser
		cmdName := cmdParts[0]
		pArgs := cmdParts[1:]

		for _, parser := range parsers {
			if enabled, ok := cfg.EnabledParsers[parser.Name()]; ok && enabled && parser.CanParse(cmdName, pArgs) {
				p = parser
				break
			}
		}

		if p == nil {
			if tel != nil {
				_ = tel.Record(telemetry.ExecutionRecord{
					Command:          toolArgs.Command,
					OriginalTokens:   len(originalOutput) / 4,
					CompressedTokens: len(originalOutput) / 4,
					ParserUsed:       "none",
					IsPassthrough:    true,
				})
			}
			return respondAllow()
		}

		compressed := p.Parse(originalOutput)
		if tel != nil {
			_ = tel.Record(telemetry.ExecutionRecord{
				Command:          toolArgs.Command,
				OriginalTokens:   len(originalOutput) / 4,
				CompressedTokens: len(compressed) / 4,
				ParserUsed:       p.Name(),
				IsPassthrough:    false,
			})
		}

		output := HookOutput{
			Decision:      "deny",
			Reason:        compressed,
			SystemMessage: fmt.Sprintf("Output compressed by Diet (%s parser)", p.Name()),
		}
		return json.NewEncoder(os.Stdout).Encode(output)
	},
}

func respondAllow() error {
	output := HookOutput{Decision: "allow"}
	return json.NewEncoder(os.Stdout).Encode(output)
}

func init() {
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(gainCmd)
	rootCmd.AddCommand(discoverCmd)
	rootCmd.AddCommand(hookCmd)
}

func main() {
	if len(os.Args) > 1 {
		subcmds := []string{"config", "gain", "discover", "_hook", "help", "--help", "-h"}
		isSub := false
		for _, s := range subcmds {
			if os.Args[1] == s {
				isSub = true
				break
			}
		}

		if !isSub {
			// Not a subcommand, treat as proxy target
			cfg, err := config.LoadConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}

			tel, err := telemetry.NewTelemetry()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading telemetry: %v\n", err)
				os.Exit(1)
			}
			defer tel.Close()

			run := runner.NewRunner(cfg, tel)
			if err := run.Run(os.Args[1:]); err != nil {
				// We don't exit with error here because the underlying command might have failed
				// and we already printed its output.
				os.Exit(0) 
			}
			os.Exit(0)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
