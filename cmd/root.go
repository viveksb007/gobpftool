package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	bpferrors "gobpftool/pkg/errors"
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
			printVersionInfo()
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

// handleError writes a formatted error message to stderr.
// It detects common error types (permission, BPF filesystem) and provides
// helpful guidance to the user.
func handleError(err error, context string) {
	if err == nil {
		return
	}

	// Check for permission errors first
	if bpferrors.IsPermissionError(err) {
		fmt.Fprintln(os.Stderr, bpferrors.FormatPermissionError())
		return
	}

	// Check for BPF filesystem issues
	if bpferrors.IsBpfFSNotMounted() {
		fmt.Fprintln(os.Stderr, bpferrors.FormatBpfFSError())
		return
	}

	// Check for specific error types
	if bpferrors.IsNoMoreKeysError(err) {
		fmt.Fprintln(os.Stderr, "Error: no more keys")
		return
	}

	// Default error formatting
	fmt.Fprintf(os.Stderr, "Error %s: %v\n", context, err)
}
