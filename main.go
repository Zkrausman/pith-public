package main

import (
	"diet/pkg/config"
	"diet/pkg/install"
	"diet/pkg/parser"
	"diet/pkg/runner"
	"diet/pkg/selfupdate"
	"diet/pkg/telemetry"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const version = "v0.5.2"

var rootCmd = &cobra.Command{
	Use:     "diet [command]",
	Short:   "Diet is a token-optimized CLI proxy",
	Version: version,
	Long:    `Diet intercepts terminal commands, compresses their output, and filters out noise to save tokens for LLMs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		// Check if first arg is a known subcommand
		subcmds := []string{"config", "gain", "discover", "reset"}
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

		// Background update check (once a day)
		now := time.Now().Unix()
		if now-cfg.LastUpdateCheck > 86400 { // 24 hours
			cfg.LastUpdateCheck = now
			_ = cfg.Save()
			go func() {
				newTag, _ := selfupdate.CheckForUpdateSilent(version)
				if newTag != "" {
					fmt.Fprintf(os.Stderr, "\n[Diet] A new version is available: %s. Run 'diet update' to upgrade!\n", newTag)
				}
			}()
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

		parsers := parser.GetAllParsers()
		var names []string
		for _, p := range parsers {
			names = append(names, p.Name())
		}

		return cfg.InteractiveConfig(names)
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

		fmt.Printf("\n=== Diet: Overall Token Savings ===\n")
		if totalOrig == 0 {
			fmt.Println("No telemetry data recorded yet. Run some commands through Diet to start tracking savings!")
			fmt.Println()
			return nil
		}

		saved := totalOrig - totalComp
		percent := (float64(saved) / float64(totalOrig)) * 100
		fmt.Printf("Raw Tokens:        %d\n", totalOrig)
		fmt.Printf("Compressed:        %d\n", totalComp)
		fmt.Printf("Tokens Saved:      %d (%.2f%%)\n", saved, percent)

		byCmd, err := tel.GetStatsByCommand()
		if err == nil && len(byCmd) > 0 {
			fmt.Printf("\n--- Breakdown by Command ---\n")
			fmt.Printf("%-25s | %-10s | %-10s | %-10s\n", "Command Pattern", "Raw", "Diet", "Savings")
			fmt.Println(strings.Repeat("-", 61))
			for _, r := range byCmd {
				savings := r.Original - r.Compressed
				fmt.Printf("%-25s | %-10d | %-10d | %-10d\n", r.Command, r.Original, r.Compressed, savings)
			}
		}

		unparsed, err := tel.GetUnparsedCommands()
		if err == nil && len(unparsed) > 0 {
			fmt.Printf("\n--- Top Unparsed Commands (Discovery) ---\n")
			fmt.Printf("%-25s | %-8s | %-12s | %-12s\n", "Command Pattern", "Count", "Raw Tokens", "Est. Savings")
			fmt.Println(strings.Repeat("-", 65))
			// Only show top 5 in gain summary
			limit := len(unparsed)
			if limit > 5 { limit = 5 }
			for i := 0; i < limit; i++ {
				r := unparsed[i]
				estSavings := float64(r.TotalRawTokens) * 0.7
				fmt.Printf("%-25s | %-8d | %-12d | %-12.0f\n", r.Pattern, r.InvocationCount, r.TotalRawTokens, estSavings)
			}
			if len(unparsed) > 5 {
				fmt.Printf("... and %d more. Run 'diet discover' for full list.\n", len(unparsed)-5)
			}
		}

		fmt.Println()
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

		if len(unparsed) == 0 {
			fmt.Println("No unparsed commands discovered yet. Run some commands through Diet to gather data!")
			return nil
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
	// Older schema
	ToolCallRequest struct {
		Arguments string `json:"arguments"`
	} `json:"tool_call_request"`
	// Newer schema
	ToolInput map[string]interface{} `json:"tool_input"`
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

		var command string
		// Try newer schema first
		if cmdVal, ok := input.ToolInput["command"].(string); ok {
			command = cmdVal
		} else if input.ToolCallRequest.Arguments != "" {
			// Fallback to older schema
			var toolArgs ToolArgs
			if err := json.Unmarshal([]byte(input.ToolCallRequest.Arguments), &toolArgs); err == nil {
				command = toolArgs.Command
			}
		}

		cfg, _ := config.LoadConfig()
		tel, _ := telemetry.NewTelemetry()
		if tel != nil {
			defer tel.Close()
		}

		// Background update check (once a day)
		now := time.Now().Unix()
		if now-cfg.LastUpdateCheck > 86400 { // 24 hours
			cfg.LastUpdateCheck = now
			_ = cfg.Save()
			go func() {
				newTag, _ := selfupdate.CheckForUpdateSilent(version)
				if newTag != "" {
					fmt.Fprintf(os.Stderr, "\n[Diet] A new version is available: %s. Run 'diet update' to upgrade!\n", newTag)
				}
			}()
		}

		originalOutput := input.ToolResponse.LlmContent
		if command == "" {
			return respondAllow()
		}

		// Handle common tool output prefixes from Gemini CLI (Output: / Error: )
		var prefix string
		exitCode := 0
		if strings.HasPrefix(originalOutput, "Output: ") {
			prefix = "Output: "
			originalOutput = strings.TrimPrefix(originalOutput, "Output: ")
		} else if strings.HasPrefix(originalOutput, "Error: ") {
			prefix = "Error: "
			originalOutput = strings.TrimPrefix(originalOutput, "Error: ")
			exitCode = 1
		}

		if strings.Contains(originalOutput, "Exit Code: ") && !strings.Contains(originalOutput, "Exit Code: 0") {
			exitCode = 1
		}

		run := runner.NewRunner(cfg, tel)
		// Log to Snag using the raw intercepted hook data
		run.LogForSnag(command, originalOutput, exitCode)

		cmdParts := strings.Fields(command)
		if len(cmdParts) == 0 {
			return respondAllow()
		}

		// Use the parser directly
		parsers := parser.GetAllParsers()
		var p parser.Parser
		cmdName := cmdParts[0]
		pArgs := cmdParts[1:]

		for _, parser := range parsers {
			enabled, ok := cfg.EnabledParsers[parser.Name()]
			if (!ok || enabled) && parser.CanParse(cmdName, pArgs) {
				p = parser
				break
			}
		}

		if p == nil {
			if tel != nil {
				_ = tel.Record(telemetry.ExecutionRecord{
					Command:          command,
					OriginalTokens:   len(originalOutput) / 4,
					CompressedTokens: len(originalOutput) / 4,
					ParserUsed:       "none",
					IsPassthrough:    true,
				})
			}
			return respondAllow()
		}

		compressed := p.Parse(originalOutput)
		
		// Apply truncation to hook output as well
		compressed = run.ApplyMiddleOutTruncation(compressed)

		if tel != nil {
			_ = tel.Record(telemetry.ExecutionRecord{
				Command:          command,
				OriginalTokens:   len(originalOutput) / 4,
				CompressedTokens: len(compressed) / 4,
				ParserUsed:       p.Name(),
				IsPassthrough:    false,
			})
		}

		output := HookOutput{
			Decision:      "deny",
			Reason:        prefix + compressed,
			SystemMessage: fmt.Sprintf("Output compressed by Diet (%s parser)", p.Name()),
		}
		return json.NewEncoder(os.Stdout).Encode(output)
	},
}

func respondAllow() error {
	output := HookOutput{Decision: "allow"}
	return json.NewEncoder(os.Stdout).Encode(output)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Diet to your system PATH and setup CLI hooks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := install.Install(); err != nil {
			return err
		}

		all, _ := cmd.Flags().GetBool("all")
		gemini, _ := cmd.Flags().GetBool("gemini")
		claude, _ := cmd.Flags().GetBool("claude")
		global, _ := cmd.Flags().GetBool("global")

		// Default to Gemini if nothing specified (for backward compatibility)
		if !all && !gemini && !claude {
			gemini = true
		}

		if all || gemini {
			if err := install.SetupGeminiHook(global); err != nil {
				fmt.Fprintf(os.Stderr, "Gemini hook failed: %v\n", err)
			}
		}

		if all || claude {
			if err := install.SetupClaudeHook(global); err != nil {
				fmt.Fprintf(os.Stderr, "Claude hook failed: %v\n", err)
			}
		}

		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Diet to the latest version from GitHub",
	RunE: func(cmd *cobra.Command, args []string) error {
		updated, err := selfupdate.CheckAndApplyUpdate(version)
		if err != nil {
			return err
		}
		if !updated {
			fmt.Println("Diet is already up to date.")
		}
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Diet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Diet %s\n", version)
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset telemetry data",
	RunE: func(cmd *cobra.Command, args []string) error {
		tel, err := telemetry.NewTelemetry()
		if err != nil {
			return err
		}
		defer tel.Close()

		passthrough, _ := cmd.Flags().GetBool("discover")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			if err := tel.ResetAll(); err != nil {
				return err
			}
			fmt.Println("All telemetry data has been reset.")
			return nil
		}

		if passthrough {
			if err := tel.ResetPassthrough(); err != nil {
				return err
			}
			fmt.Println("Discovery data (passthrough commands) has been reset.")
			return nil
		}

		return fmt.Errorf("please specify what to reset using --all or --discover")
	},
}

var rawCmd = &cobra.Command{
	Use:   "raw [command]",
	Short: "Run a command and bypass all parsers (escape hatch)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
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
		return run.RunWithOptions(args, true)
	},
}

func init() {
	resetCmd.Flags().Bool("all", false, "Reset ALL telemetry data (gain and discover)")
	resetCmd.Flags().Bool("discover", false, "Reset only discovery data (passthrough commands)")

	installCmd.Flags().Bool("all", false, "Setup hooks for all supported CLIs")
	installCmd.Flags().Bool("gemini", false, "Setup hook for Gemini CLI")
	installCmd.Flags().Bool("claude", false, "Setup hook for Claude Code")
	installCmd.Flags().BoolP("global", "g", false, "Install hooks globally in the home directory")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(gainCmd)
	rootCmd.AddCommand(discoverCmd)
	rootCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(rawCmd)
}

func main() {
	if len(os.Args) > 1 {
		// List of internal subcommands and flags that SHOULD NOT be proxied
		subcmds := []string{"config", "gain", "discover", "reset", "raw", "_hook", "install", "update", "version", "--version", "-v", "help", "--help", "-h"}
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
