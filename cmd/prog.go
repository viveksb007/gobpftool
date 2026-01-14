package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	bpferrors "github.com/viveksb007/gobpftool/pkg/errors"
	"github.com/viveksb007/gobpftool/pkg/output"
	"github.com/viveksb007/gobpftool/pkg/prog"
)

var progService prog.Service

// progCmd represents the prog command
var progCmd = &cobra.Command{
	Use:   "prog",
	Short: "Inspect eBPF programs",
	Long: `Inspect eBPF programs loaded in the kernel.

Available commands:
  show    Show information about loaded programs
  help    Display help for prog commands`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		cmd.Help()
	},
}

// progShowCmd represents the prog show command
var progShowCmd = &cobra.Command{
	Use:     "show [PROG]",
	Aliases: []string{"list"},
	Short:   "Show information about loaded programs",
	Long: `Show information about loaded eBPF programs.

Without arguments, lists all loaded programs.
With arguments, shows specific program(s):

  gobpftool prog show                    # List all programs
  gobpftool prog show id 123             # Show program with ID 123
  gobpftool prog show tag f0055c08993fea1e  # Show programs with tag
  gobpftool prog show name my_prog       # Show programs with name
  gobpftool prog show pinned /sys/fs/bpf/my_prog  # Show pinned program`,
	RunE: runProgShow,
}

func runProgShow(cmd *cobra.Command, args []string) error {
	// Determine output format
	format := getOutputFormat()
	formatter := output.NewFormatter(format)

	var programs []prog.ProgramInfo
	var err error

	if len(args) == 0 {
		// List all programs
		programs, err = progService.List()
		if err != nil {
			handleError(err, "listing programs")
			return err
		}
	} else if len(args) >= 2 {
		// Parse program identifier
		identifier := args[0]
		value := args[1]

		switch identifier {
		case "id":
			id, parseErr := strconv.ParseUint(value, 10, 32)
			if parseErr != nil {
				fmt.Fprintf(os.Stderr, "Error: invalid program ID: %s\n", value)
				return bpferrors.ErrInvalidID
			}

			program, getErr := progService.GetByID(uint32(id))
			if getErr != nil {
				handleError(getErr, fmt.Sprintf("getting program with ID %d", id))
				return getErr
			}
			programs = []prog.ProgramInfo{*program}

		case "tag":
			programs, err = progService.GetByTag(value)
			if err != nil {
				handleError(err, fmt.Sprintf("getting programs with tag %s", value))
				return err
			}

		case "name":
			programs, err = progService.GetByName(value)
			if err != nil {
				handleError(err, fmt.Sprintf("getting programs with name %s", value))
				return err
			}

		case "pinned":
			program, getErr := progService.GetByPinnedPath(value)
			if getErr != nil {
				handleError(getErr, fmt.Sprintf("getting pinned program at %s", value))
				return getErr
			}
			programs = []prog.ProgramInfo{*program}

		default:
			fmt.Fprintf(os.Stderr, "Error: invalid program identifier: %s. Use 'id', 'tag', 'name', or 'pinned'\n", identifier)
			return fmt.Errorf("invalid identifier: %s", identifier)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: invalid arguments. Use 'gobpftool prog show' or 'gobpftool prog show <identifier> <value>'\n")
		return fmt.Errorf("invalid arguments")
	}

	// Convert prog.ProgramInfo to output.ProgramInfo
	outputPrograms := make([]output.ProgramInfo, len(programs))
	for i, p := range programs {
		outputPrograms[i] = output.ProgramInfo{
			ID:        p.ID,
			Type:      p.Type,
			Name:      p.Name,
			Tag:       p.Tag,
			GPL:       p.GPL,
			LoadedAt:  p.LoadedAt,
			UID:       p.UID,
			BytesXlat: p.BytesXlated,
			BytesJIT:  p.BytesJIT,
			MemLock:   p.MemLock,
			MapIDs:    p.MapIDs,
		}
	}

	// Format and output the results
	result := formatter.FormatPrograms(outputPrograms)
	fmt.Print(result)

	return nil
}

// progHelpCmd represents the prog help command
var progHelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Display help for prog commands",
	Long: `Display help information for prog commands.

Available prog commands:
  show    Show information about loaded programs
  help    Display this help message

Examples:
  gobpftool prog show                           # List all programs
  gobpftool prog show id 123                    # Show program with ID 123
  gobpftool prog show tag f0055c08993fea1e      # Show programs with tag
  gobpftool prog show name my_prog              # Show programs with name
  gobpftool prog show pinned /sys/fs/bpf/prog   # Show pinned program

Global flags:
  -j, --json     Output in JSON format
  -p, --pretty   Output in pretty-printed JSON format`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show the help for the prog command
		progCmd.Help()
	},
}

// getOutputFormat determines the output format based on global flags
func getOutputFormat() output.Format {
	flags := GetGlobalFlags()
	if flags.Pretty {
		return output.FormatJSONPretty
	} else if flags.JSON {
		return output.FormatJSON
	}
	return output.FormatPlain
}

func init() {
	// Initialize the program service
	progService = prog.NewService()

	// Add subcommands to prog command
	progCmd.AddCommand(progShowCmd)
	progCmd.AddCommand(progHelpCmd)

	// Add prog command to root command
	rootCmd.AddCommand(progCmd)
}
