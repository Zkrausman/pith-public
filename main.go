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

const version = "v0.14.0"

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

func respondAllow(cmd *cobra.Command) error {
	output := HookOutput{Decision: "allow"}
	return json.NewEncoder(cmd.OutOrStdout()).Encode(output)
}

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "pith [command]",
		Short:   "Pith is a token-optimized CLI proxy",
		Version: version,
		Long:    `Pith intercepts terminal commands, compresses their output, and filters out noise to save tokens for LLMs.`,
		Args:    cobra.ArbitraryArgs,
		RunE: runRoot,
	}

	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Interactive configuration tool",
		RunE: runConfig,
	}

	var gainCmd = &cobra.Command{
		Use:   "gain",
		Short: "Show token savings analytics",
		RunE: runGain,
	}

	var discoverCmd = &cobra.Command{
		Use:   "discover",
		Short: "Identify commands that could benefit from dedicated parsers",
		RunE: runDiscover,
	}

	var hookCmd = &cobra.Command{
		Use:    "_hook",
		Short:  "Internal hook for Gemini CLI",
		Hidden: true,
		RunE: runHook,
	}

	var installCmd = &cobra.Command{
		Use:   "install",
		Short: "Install Pith to your system PATH and setup CLI hooks",
		RunE: runInstall,
	}

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update Pith to the latest version from GitHub",
		RunE: runUpdate,
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Pith",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Pith %s\n", version)
		},
	}

	var resetCmd = &cobra.Command{
		Use:   "reset",
		Short: "Reset telemetry data",
		RunE: runReset,
	}

	var rawCmd = &cobra.Command{
		Use:   "raw [command]",
		Short: "Run a command and bypass all parsers (escape hatch)",
		RunE: runRaw,
	}

	var dashboardCmd = &cobra.Command{
		Use:   "dashboard",
		Short: "Open the interactive analytics dashboard in your browser",
		RunE: runDashboard,
	}

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
	rootCmd.AddCommand(synapseSyncCmd)

	return rootCmd
}

func main() {
	// Trigger migration if needed
	cfg, err := config.LoadConfig()
	if err == nil {
		_ = config.MigrateStorage(cfg.StoragePath)
	}

	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
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
				cmd.PrintErrf("\n[Pith] A new version is available: %s. Run 'pith update' to upgrade!\n", newTag)
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
}

func runConfig(cmd *cobra.Command, args []string) error {
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
}

func runGain(cmd *cobra.Command, args []string) error {
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

	cmd.Printf("\n=== Pith: Overall Token Savings ===\n")
	if totalOrig == 0 {
		cmd.Println("No telemetry data recorded yet. Run some commands through Pith to start tracking savings!")
		cmd.Println()
		return nil
	}

	saved := totalOrig - totalComp
	percent := (float64(saved) / float64(totalOrig)) * 100
	usdSaved := (float64(saved) / 1000000.0) * cfg.USDPerMillionTokens

	cmd.Printf("Raw Tokens:        %d\n", totalOrig)
	cmd.Printf("Compressed:        %d\n", totalComp)
	cmd.Printf("Tokens Saved:      %d (%.2f%%)\n", saved, percent)
	cmd.Printf("Est. USD Saved:    $%.2f (at $%.2f/M tokens)\n", usdSaved, cfg.USDPerMillionTokens)

	byCmd, err := tel.GetStatsByCommand("")
	if err == nil && len(byCmd) > 0 {
		cmd.Printf("\n--- Breakdown by Command and Agent ---\n")
		cmd.Printf("%-25s | %-12s | %-10s | %-10s | %-10s\n", "Command Pattern", "Agent", "Raw", "Pith", "USD Saved")
		cmd.Println(strings.Repeat("-", 76))

		limit := len(byCmd)
		if limit > 20 {
			limit = 20
		}

		for i := 0; i < limit; i++ {
			r := byCmd[i]
			savings := r.Original - r.Compressed
			usdCmd := (float64(savings) / 1000000.0) * cfg.USDPerMillionTokens
			cmd.Printf("%-25s | %-12s | %-10d | %-10d | $%-10.3f\n", r.Command, r.Source, r.Original, r.Compressed, usdCmd)
		}

		if len(byCmd) > 20 {
			cmd.Printf("... and %d more. Run 'pith dashboard' for full interactive breakdown.\n", len(byCmd)-20)
		}
	}

	unparsed, err := tel.GetUnparsedCommands("")
	if err == nil && len(unparsed) > 0 {
		cmd.Printf("\n--- Top Unparsed Commands (Discovery) ---\n")
		cmd.Printf("%-25s | %-12s | %-8s | %-12s | %-12s\n", "Command Pattern", "Agent", "Count", "Raw Tokens", "Est. Savings")
		cmd.Println(strings.Repeat("-", 80))
		limit := len(unparsed)
		if limit > 10 {
			limit = 10
		}
		for i := 0; i < limit; i++ {
			r := unparsed[i]
			estSavings := float64(r.TotalRawTokens) * 0.7
			cmd.Printf("%-25s | %-12s | %-8d | %-12d | %-12.0f\n", r.Pattern, r.Source, r.InvocationCount, r.TotalRawTokens, estSavings)
		}
		if len(unparsed) > 10 {
			cmd.Printf("... and %d more. Run 'pith discover' for full list.\n", len(unparsed)-10)
		}
	}

	cmd.Println()
	return nil
}

func runDiscover(cmd *cobra.Command, args []string) error {
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
		cmd.Println("No unparsed commands discovered yet. Run some commands through Pith to gather data!")
		return nil
	}

	cmd.Printf("--- Opportunity Discovery (Unparsed Commands) ---\n")
	cmd.Printf("%-30s | %-12s | %-10s | %-15s | %-15s\n", "Command Pattern", "Agent", "Count", "Total Tokens", "Est. USD Saved")
	cmd.Println(strings.Repeat("-", 91))

	for _, r := range unparsed {
		estSavedTokens := float64(r.TotalRawTokens) * 0.7
		estSavingsUSD := (estSavedTokens / 1000000.0) * cfg.USDPerMillionTokens
		cmd.Printf("%-30s | %-12s | %-10d | %-15d | $%-15.3f\n", r.Pattern, r.Source, r.InvocationCount, r.TotalRawTokens, estSavingsUSD)
	}

	return nil
}

func runHook(cmd *cobra.Command, args []string) error {
	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var input HookInput
	if err := json.Unmarshal(inputData, &input); err != nil {
		return err
	}

	var command string
	if cmdVal, ok := input.ToolInput["command"].(string); ok {
		command = cmdVal
	} else if input.ToolCallRequest.Arguments != "" {
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
	if now-cfg.LastUpdateCheck > 86400 {
		cfg.LastUpdateCheck = now
		_ = cfg.Save()
		go func() {
			newTag, _ := selfupdate.CheckForUpdateSilent(version)
			if newTag != "" {
				cmd.PrintErrf("\n[Pith] A new version is available: %s. Run 'pith update' to upgrade!\n", newTag)
			}
		}()
	}

	originalOutput := input.ToolResponse.LlmContent
	if command == "" {
		return respondAllow(cmd)
	}

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
	run.LogForSnag(command, originalOutput, exitCode)

	source, _ := cmd.Flags().GetString("source")
	if source == "" {
		source = "unknown"
	}

	cmdParts := strings.Fields(command)
	if len(cmdParts) == 0 {
		return respondAllow(cmd)
	}

	parsers := parser.GetAllParsers()
	var p parser.Parser
	
	// Strategy 1: Direct match
	cmdName := cmdParts[0]
	pArgs := cmdParts[1:]
	for _, parser := range parsers {
		enabled, ok := cfg.EnabledParsers[parser.Name()]
		if (!ok || enabled) && parser.CanParse(cmdName, pArgs) {
			p = parser
			break
		}
	}

	// Strategy 2: PowerShell '&' wrapper
	if p == nil && cmdName == "&" && len(cmdParts) > 1 {
		cmdName = cmdParts[1]
		pArgs = cmdParts[2:]
		// Handle quoted path in & "path" format
		if strings.HasPrefix(cmdName, "\"") && strings.HasSuffix(cmdName, "\"") {
			cmdName = strings.Trim(cmdName, "\"")
		}
		for _, parser := range parsers {
			enabled, ok := cfg.EnabledParsers[parser.Name()]
			if (!ok || enabled) && parser.CanParse(cmdName, pArgs) {
				p = parser
				break
			}
		}
	}

	// Strategy 3: Leading quoted path (e.g. "C:\bin\go.exe" version)
	if p == nil && strings.HasPrefix(command, "\"") {
		endQuote := strings.Index(command[1:], "\"")
		if endQuote != -1 {
			cmdName = command[1 : endQuote+1]
			remaining := strings.TrimSpace(command[endQuote+2:])
			pArgs = strings.Fields(remaining)
			for _, parser := range parsers {
				enabled, ok := cfg.EnabledParsers[parser.Name()]
				if (!ok || enabled) && parser.CanParse(cmdName, pArgs) {
					p = parser
					break
				}
			}
		}
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

	truncated := run.ApplyMiddleOutTruncation(compressed)
	wasTruncated := truncated != compressed

	origTokens := runner.EstimateTokensWithHeuristic(originalOutput, cfg.TokenHeuristic)
	compTokens := runner.EstimateTokensWithHeuristic(truncated, cfg.TokenHeuristic)
	savedTokens := origTokens - compTokens

	if tel != nil {
		_ = tel.Record(telemetry.ExecutionRecord{
			Command:          command,
			OriginalTokens:   origTokens,
			CompressedTokens: compTokens,
			ParserUsed:       parserUsed,
			IsPassthrough:    isPassthrough,
			Source:           source,
		})
	}

	if p == nil && !wasTruncated {
		return respondAllow(cmd)
	}

	output := HookOutput{
		Decision:      "deny",
		Reason:        prefix + truncated,
		SystemMessage: fmt.Sprintf("Output optimized by Pith (parser: %s, tokens saved: %d, truncated: %v)", parserUsed, savedTokens, wasTruncated),
	}
	return json.NewEncoder(cmd.OutOrStdout()).Encode(output)
}

func runInstall(cmd *cobra.Command, args []string) error {
	if err := install.Install(); err != nil {
		return err
	}

	all, _ := cmd.Flags().GetBool("all")
	gemini, _ := cmd.Flags().GetBool("gemini")
	claude, _ := cmd.Flags().GetBool("claude")
	codex, _ := cmd.Flags().GetBool("codex")
	global, _ := cmd.Flags().GetBool("global")

	if !all && !gemini && !claude && !codex {
		all = true
		if !cmd.Flags().Changed("global") {
			global = true
		}
	}

	if all || gemini {
		if err := install.SetupGeminiHook(global); err != nil {
			cmd.PrintErrf("Gemini hook failed: %v\n", err)
		}
	}

	if all || claude {
		if err := install.SetupClaudeHook(global); err != nil {
			cmd.PrintErrf("Claude hook failed: %v\n", err)
		}
	}

	if all || codex {
		if err := install.SetupCodexHook(global); err != nil {
			cmd.PrintErrf("Codex hook failed: %v\n", err)
		}
	}

	return nil
}

func runUpdate(cmd *cobra.Command, args []string) error {
	updated, err := selfupdate.CheckAndApplyUpdate(version)
	if err != nil {
		return err
	}
	if !updated {
		cmd.Println("Pith is already up to date.")
	}
	return nil
}

func runReset(cmd *cobra.Command, args []string) error {
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
		cmd.Println("All telemetry data has been reset.")
		return nil
	}

	if passthrough {
		if err := tel.ResetPassthrough(); err != nil {
			return err
		}
		cmd.Println("Discovery data (passthrough commands) has been reset.")
		return nil
	}

	return fmt.Errorf("please specify what to reset using --all or --discover")
}

func runRaw(cmd *cobra.Command, args []string) error {
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
}

func runDashboard(cmd *cobra.Command, args []string) error {
	cfg, _ := config.LoadConfig()
	tel, err := telemetry.NewTelemetry(cfg.StoragePath)
	if err != nil {
		return err
	}
	defer tel.Close()

	port, _ := cmd.Flags().GetInt("port")
	return gui.StartDashboard(cfg, tel, port)
}
