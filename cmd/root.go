package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version information - can be set at build time using ldflags
var (
	Version   = "0.1.0"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// GlobalFlags holds the global CLI flags
type GlobalFlags struct {
	JSON   bool // -j, --json
	Pretty bool // -p, --pretty
}

var globalFlags GlobalFlags
var showVersion bool

var rootCmd = &cobra.Command{
	Use:   "gobpftool",
	Short: "Tool for inspection of eBPF programs and maps",
	Long: `gobpftool is a Go-based CLI tool that provides functionality similar to
the Linux bpftool utility for inspecting eBPF programs and maps.

It uses the cilium/ebpf library to interact with the kernel's eBPF subsystem.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			printVersion()
			return
		}
		// If no subcommand is provided, show help
		cmd.Help()
	},
	SilenceUsage: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.JSON, "json", "j", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.Pretty, "pretty", "p", false, "Output in pretty-printed JSON format")
	rootCmd.Flags().BoolVar(&showVersion, "version", false, "Display version information")

}

// GetGlobalFlags returns the global flags
func GetGlobalFlags() GlobalFlags {
	return globalFlags
}

// printVersion prints the version information
func printVersion() {
	fmt.Fprintf(os.Stdout, "gobpftool version %s\n", Version)
	if GitCommit != "unknown" {
		fmt.Fprintf(os.Stdout, "  git commit: %s\n", GitCommit)
	}
	if BuildDate != "unknown" {
		fmt.Fprintf(os.Stdout, "  build date: %s\n", BuildDate)
	}
}

// SetVersionInfo allows setting version info programmatically (useful for testing)
func SetVersionInfo(version, commit, date string) {
	Version = version
	GitCommit = commit
	BuildDate = date
}

// GetRootCmd returns the root command (useful for testing)
func GetRootCmd() *cobra.Command {
	return rootCmd
}

// ResetFlags resets all flags to their default values (useful for testing)
func ResetFlags() {
	globalFlags = GlobalFlags{}
	showVersion = false
}
