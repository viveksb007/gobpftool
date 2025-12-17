package cmd

import (
	"github.com/spf13/cobra"
)

var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Inspect and read eBPF maps",
	Long:  `Commands to inspect and read data from loaded eBPF maps.`,
}

func init() {
	rootCmd.AddCommand(mapCmd)
}
