package main

import (
	"pith/pkg/config"
	"pith/pkg/gui"
	"pith/pkg/install"
	"pith/pkg/parser"
	"pith/pkg/runner"
	"pith/pkg/selfupdate"
	"pith/pkg/telemetry"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const version = "v0.11.2"

var rootCmd = &cobra.Command{
	Use:     "pith [command]",
	Short:   "Pith is a token-optimized CLI proxy",
	Version: version,
	Long:    `Pith intercepts terminal commands, compresses their output, and filters out noise to save tokens for LLMs.`,
	Args:    cobra.ArbitraryArgs,
	// By default, if no subcommand matches, Cobra runs Run or RunE
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		// Since RunE is only called if NO other subcommand matches,
		// and we handled the empty args case above, any args here
		// MUST be intended as a proxy target.
		
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
					fmt.Fprintf(os.Stderr, "\n[Pith] A new version is available: %s. Run 'pith update' to upgrade!\n", newTag)
				}
			}()
		}

		tel, err := telemetry.NewTelemetry(cfg.StoragePath)
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
		cfg, _ := config.LoadConfig()
		tel, err := telemetry.NewTelemetry(cfg.StoragePath)
		if err != nil {
			return err
		}
		defer tel.Close()

		totalOrig, totalComp, err := tel.GetStats("")
		if err != nil {
			return err
		}

		fmt.Printf("\n=== Pith: Overall Token Savings ===\n")
		if totalOrig == 0 {
			fmt.Println("No telemetry data recorded yet. Run some commands through Pith to start tracking savings!")
			fmt.Println()
			return nil
		}

		saved := totalOrig - totalComp
		percent := (float64(saved) / float64(totalOrig)) * 100
		fmt.Printf("Raw Tokens:        %d\n", totalOrig)
		fmt.Printf("Compressed:        %d\n", totalComp)
		fmt.Printf("Tokens Saved:      %d (%.2f%%)\n", saved, percent)

		byCmd, err := tel.GetStatsByCommand("")
		if err == nil && len(byCmd) > 0 {
			fmt.Printf("\n--- Breakdown by Command and Agent ---\n")
			fmt.Printf("%-25s | %-12s | %-10s | %-10s | %-10s\n", "Command Pattern", "Agent", "Raw", "Pith", "Savings")
			fmt.Println(strings.Repeat("-", 76))
			
			// Limit to top 20 results in summary
			limit := len(byCmd)
			if limit > 20 { limit = 20 }
			
			for i := 0; i < limit; i++ {
				r := byCmd[i]
				savings := r.Original - r.Compressed
				fmt.Printf("%-25s | %-12s | %-10d | %-10d | %-10d\n", r.Command, r.Source, r.Original, r.Compressed, savings)
			}
			
			if len(byCmd) > 20 {
				fmt.Printf("... and %d more. Run 'pith dashboard' for full interactive breakdown.\n", len(byCmd)-20)
			}
		}

		unparsed, err := tel.GetUnparsedCommands("")
		if err == nil && len(unparsed) > 0 {
			fmt.Printf("\n--- Top Unparsed Commands (Discovery) ---\n")
			fmt.Printf("%-25s | %-12s | %-8s | %-12s | %-12s\n", "Command Pattern", "Agent", "Count", "Raw Tokens", "Est. Savings")
			fmt.Println(strings.Repeat("-", 80))
			// Only show top 10 in gain summary
			limit := len(unparsed)
			if limit > 10 { limit = 10 }
			for i := 0; i < limit; i++ {
				r := unparsed[i]
				estSavings := float64(r.TotalRawTokens) * 0.7
				fmt.Printf("%-25s | %-12s | %-8d | %-12d | %-12.0f\n", r.Pattern, r.Source, r.InvocationCount, r.TotalRawTokens, estSavings)
			}
			if len(unparsed) > 10 {
				fmt.Printf("... and %d more. Run 'pith discover' for full list.\n", len(unparsed)-10)
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
		cfg, _ := config.LoadConfig()
		tel, err := telemetry.NewTelemetry(cfg.StoragePath)
		if err != nil {
			return err
		}
		defer tel.Close()

		unparsed, err := tel.GetUnparsedCommands("")
		if err != nil {
			return err
		}

		if len(unparsed) == 0 {
			fmt.Println("No unparsed commands discovered yet. Run some commands through Pith to gather data!")
			return nil
		}

		fmt.Printf("--- Opportunity Discovery (Unparsed Commands) ---\n")
		fmt.Printf("%-30s | %-12s | %-10s | %-15s | %-15s\n", "Command Pattern", "Agent", "Count", "Total Tokens", "Est. Savings (70%)")
		fmt.Println(strings.Repeat("-", 91))

		for _, r := range unparsed {
			estSavings := float64(r.TotalRawTokens) * 0.7 // Assume 70% average savings
			fmt.Printf("%-30s | %-12s | %-10d | %-15d | %-15.0f\n", r.Pattern, r.Source, r.InvocationCount, r.TotalRawTokens, estSavings)
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
		tel, _ := telemetry.NewTelemetry(cfg.StoragePath)
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
					fmt.Fprintf(os.Stderr, "\n[Pith] A new version is available: %s. Run 'pith update' to upgrade!\n", newTag)
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

		source, _ := cmd.Flags().GetString("source")
		if source == "" {
			source = "unknown"
		}

		var compressed string
		parserUsed := "none"
		isPassthrough := true

		if p != nil {
			compressed = p.Parse(originalOutput)
			parserUsed = p.Name()
			isPassthrough = false
		} else {
			compressed = originalOutput
		}
		
		// Apply truncation to all hook output
		truncated := run.ApplyMiddleOutTruncation(compressed)
		wasTruncated := truncated != compressed

		if tel != nil {
			_ = tel.Record(telemetry.ExecutionRecord{
				Command:          command,
				OriginalTokens:   runner.EstimateTokens(originalOutput),
				CompressedTokens: runner.EstimateTokens(truncated),
				ParserUsed:       parserUsed,
				IsPassthrough:    isPassthrough,
				Source:           source,
			})
		}

		// If nothing was parsed and nothing was truncated, allow original
		if p == nil && !wasTruncated {
			return respondAllow()
		}

		// Otherwise, deny and provide the optimized output
		output := HookOutput{
			Decision:      "deny",
			Reason:        prefix + truncated,
			SystemMessage: fmt.Sprintf("Output optimized by Pith (parser: %s, truncated: %v)", parserUsed, wasTruncated),
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
	Short: "Install Pith to your system PATH and setup CLI hooks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := install.Install(); err != nil {
			return err
		}

		all, _ := cmd.Flags().GetBool("all")
		gemini, _ := cmd.Flags().GetBool("gemini")
		claude, _ := cmd.Flags().GetBool("claude")
		codex, _ := cmd.Flags().GetBool("codex")
		global, _ := cmd.Flags().GetBool("global")

		// Default to global install for all CLIs if no specific agent is specified
		if !all && !gemini && !claude && !codex {
			all = true
			if !cmd.Flags().Changed("global") {
				global = true
			}
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

		if all || codex {
			if err := install.SetupCodexHook(global); err != nil {
				fmt.Fprintf(os.Stderr, "Codex hook failed: %v\n", err)
			}
		}

		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Pith to the latest version from GitHub",
	RunE: func(cmd *cobra.Command, args []string) error {
		updated, err := selfupdate.CheckAndApplyUpdate(version)
		if err != nil {
			return err
		}
		if !updated {
			fmt.Println("Pith is already up to date.")
		}
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Pith",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Pith %s\n", version)
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset telemetry data",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()
		tel, err := telemetry.NewTelemetry(cfg.StoragePath)
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

		tel, err := telemetry.NewTelemetry(cfg.StoragePath)
		if err != nil {
			return err
		}
		defer tel.Close()

		run := runner.NewRunner(cfg, tel)
		return run.RunWithOptions(args, true)
	},
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Open the interactive analytics dashboard in your browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()
		tel, err := telemetry.NewTelemetry(cfg.StoragePath)
		if err != nil {
			return err
		}
		defer tel.Close()

		port, _ := cmd.Flags().GetInt("port")
		return gui.StartDashboard(tel, port)
	},
}

func init() {
	resetCmd.Flags().Bool("all", false, "Reset ALL telemetry data (gain and discover)")
	resetCmd.Flags().Bool("discover", false, "Reset only discovery data (passthrough commands)")

	installCmd.Flags().Bool("all", false, "Setup hooks for all supported CLIs")
	installCmd.Flags().Bool("gemini", false, "Setup hook for Gemini CLI")
	installCmd.Flags().Bool("claude", false, "Setup hook for Claude Code")
	installCmd.Flags().Bool("codex", false, "Setup hook for Codex")
	installCmd.Flags().BoolP("global", "g", false, "Install hooks globally in the home directory")

	dashboardCmd.Flags().IntP("port", "p", 8080, "Port to run the dashboard server on")

	hookCmd.Flags().String("source", "unknown", "The agent or CLI triggering the hook")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(gainCmd)
	rootCmd.AddCommand(discoverCmd)
	rootCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(rawCmd)
	rootCmd.AddCommand(dashboardCmd)
}

func main() {
	// Trigger migration if needed
	cfg, err := config.LoadConfig()
	if err == nil {
		_ = config.MigrateStorage(cfg.StoragePath)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
