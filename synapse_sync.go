package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var synapseSyncCmd = &cobra.Command{
	Use:    "synapse-sync",
	Short:  "Internal hook for Synapse to update Pith",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Synapse syncing Pith via its built-in update command...")
		// Use Pith's own update logic which handles GitHub releases properly
		execCmd := exec.Command(os.Args[0], "update")
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		return execCmd.Run()
	},
}

func init() {
	rootCmd.AddCommand(synapseSyncCmd)
}
