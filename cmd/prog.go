package cmd

import (
	"github.com/spf13/cobra"
)

var progCmd = &cobra.Command{
	Use:   "prog",
	Short: "Inspect eBPF programs",
	Long:  `Commands to inspect and display information about loaded eBPF programs.`,
}

func init() {
	rootCmd.AddCommand(progCmd)
}
