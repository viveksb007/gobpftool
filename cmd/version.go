package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display the version, git commit, and build date of gobpftool.`,
	Run: func(cmd *cobra.Command, args []string) {
		printVersionInfo()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// printVersionInfo prints detailed version information
func printVersionInfo() {
	fmt.Fprintf(os.Stdout, "gobpftool version %s\n", Version)
	if GitCommit != "unknown" {
		fmt.Fprintf(os.Stdout, "  git commit: %s\n", GitCommit)
	}
	if BuildDate != "unknown" {
		fmt.Fprintf(os.Stdout, "  build date: %s\n", BuildDate)
	}
}
