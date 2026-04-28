package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var synapseSyncCmd = &cobra.Command{
	Use:    "synapse-sync",
	Short:  "Internal hook for Synapse to update Pith",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Synapse syncing Pith via its built-in update command...")

		exe, err := os.Executable()
		if err != nil {
			return err
		}

		execCmd := exec.Command(exe, "update")
		output, err := execCmd.CombinedOutput()
		fmt.Print(string(output))

		if err != nil {
			return err
		}

		// SCS-1 Status: Check if the tool reported it was already up to date
		if strings.Contains(string(output), "already up to date") {
			os.Exit(2)
		}

		return nil
	},
}

// Removed init() to avoid dependency on global rootCmd
// This command is now added in main.go's NewRootCmd
