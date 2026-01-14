package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/viveksb007/gobpftool/internal/utils"
	bpferrors "github.com/viveksb007/gobpftool/pkg/errors"
	"github.com/viveksb007/gobpftool/pkg/maps"
	"github.com/viveksb007/gobpftool/pkg/output"
)

var mapService maps.Service

// mapCmd represents the map command
var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Inspect and read eBPF maps",
	Long: `Inspect and read data from loaded eBPF maps.

Available commands:
  show      Show information about loaded maps
  dump      Dump all entries in a map
  lookup    Lookup a key in a map
  getnext   Get next key in a map
  help      Display help for map commands`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		cmd.Help()
	},
}

// mapShowCmd represents the map show command
var mapShowCmd = &cobra.Command{
	Use:     "show [MAP]",
	Aliases: []string{"list"},
	Short:   "Show information about loaded maps",
	Long: `Show information about loaded eBPF maps.

Without arguments, lists all loaded maps.
With arguments, shows specific map(s):

  gobpftool map show                    # List all maps
  gobpftool map show id 123             # Show map with ID 123
  gobpftool map show name my_map        # Show maps with name
  gobpftool map show pinned /sys/fs/bpf/my_map  # Show pinned map`,
	RunE: runMapShow,
}

// mapDumpCmd represents the map dump command
var mapDumpCmd = &cobra.Command{
	Use:   "dump MAP",
	Short: "Dump all entries in a map",
	Long: `Dump all key-value entries in an eBPF map.

  gobpftool map dump id 123             # Dump map with ID 123
  gobpftool map dump name my_map        # Dump maps with name
  gobpftool map dump pinned /sys/fs/bpf/my_map  # Dump pinned map`,
	RunE: runMapDump,
}

// mapLookupCmd represents the map lookup command
var mapLookupCmd = &cobra.Command{
	Use:   "lookup MAP key KEY_DATA",
	Short: "Lookup a key in a map",
	Long: `Lookup a specific key in an eBPF map.

Key data is specified as space-separated hex bytes.

  gobpftool map lookup id 123 key 0a 0b 0c 0d
  gobpftool map lookup pinned /sys/fs/bpf/my_map key 01 02 03 04`,
	RunE: runMapLookup,
}

// mapGetNextCmd represents the map getnext command
var mapGetNextCmd = &cobra.Command{
	Use:   "getnext MAP [key KEY_DATA]",
	Short: "Get next key in a map",
	Long: `Get the next key in an eBPF map.

Without a key, returns the first key in the map.
With a key, returns the next key after the specified key.

  gobpftool map getnext id 123                    # Get first key
  gobpftool map getnext id 123 key 0a 0b 0c 0d    # Get next key after specified key`,
	RunE: runMapGetNext,
}

// mapHelpCmd represents the map help command
var mapHelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Display help for map commands",
	Long: `Display help information for map commands.

Available map commands:
  show      Show information about loaded maps
  dump      Dump all entries in a map
  lookup    Lookup a key in a map
  getnext   Get next key in a map
  help      Display this help message

Examples:
  gobpftool map show                              # List all maps
  gobpftool map show id 123                       # Show map with ID 123
  gobpftool map show name my_map                  # Show maps with name
  gobpftool map show pinned /sys/fs/bpf/map       # Show pinned map
  gobpftool map dump id 123                       # Dump all entries
  gobpftool map lookup id 123 key 0a 0b 0c 0d     # Lookup key
  gobpftool map getnext id 123                    # Get first key
  gobpftool map getnext id 123 key 0a 0b 0c 0d    # Get next key

Global flags:
  -j, --json     Output in JSON format
  -p, --pretty   Output in pretty-printed JSON format`,
	Run: func(cmd *cobra.Command, args []string) {
		mapCmd.Help()
	},
}

// runMapShow handles the map show command
func runMapShow(cmd *cobra.Command, args []string) error {
	format := getOutputFormat()
	formatter := output.NewFormatter(format)

	var mapInfos []maps.MapInfo
	var err error

	if len(args) == 0 {
		// List all maps
		mapInfos, err = mapService.List()
		if err != nil {
			handleError(err, "listing maps")
			return err
		}
	} else if len(args) >= 2 {
		// Parse map identifier
		identifier := args[0]
		value := args[1]

		switch identifier {
		case "id":
			id, parseErr := strconv.ParseUint(value, 10, 32)
			if parseErr != nil {
				fmt.Fprintf(os.Stderr, "Error: invalid map ID: %s\n", value)
				return bpferrors.ErrInvalidID
			}

			mapInfo, getErr := mapService.GetByID(uint32(id))
			if getErr != nil {
				handleError(getErr, fmt.Sprintf("getting map with ID %d", id))
				return getErr
			}
			mapInfos = []maps.MapInfo{*mapInfo}

		case "name":
			mapInfos, err = mapService.GetByName(value)
			if err != nil {
				handleError(err, fmt.Sprintf("getting maps with name %s", value))
				return err
			}

		case "pinned":
			mapInfo, getErr := mapService.GetByPinnedPath(value)
			if getErr != nil {
				handleError(getErr, fmt.Sprintf("getting pinned map at %s", value))
				return getErr
			}
			mapInfos = []maps.MapInfo{*mapInfo}

		default:
			fmt.Fprintf(os.Stderr, "Error: invalid map identifier: %s. Use 'id', 'name', or 'pinned'\n", identifier)
			return fmt.Errorf("invalid identifier: %s", identifier)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: invalid arguments. Use 'gobpftool map show' or 'gobpftool map show <identifier> <value>'\n")
		return fmt.Errorf("invalid arguments")
	}

	// Convert maps.MapInfo to output.MapInfo
	outputMaps := make([]output.MapInfo, len(mapInfos))
	for i, m := range mapInfos {
		outputMaps[i] = output.MapInfo{
			ID:         m.ID,
			Type:       m.Type,
			Name:       m.Name,
			KeySize:    m.KeySize,
			ValueSize:  m.ValueSize,
			MaxEntries: m.MaxEntries,
			Flags:      m.Flags,
			MemLock:    m.MemLock,
		}
	}

	result := formatter.FormatMaps(outputMaps)
	fmt.Print(result)

	return nil
}

// runMapDump handles the map dump command
func runMapDump(cmd *cobra.Command, args []string) error {
	format := getOutputFormat()
	formatter := output.NewFormatter(format)

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: map identifier required. Use 'gobpftool map dump <identifier> <value>'\n")
		return fmt.Errorf("map identifier required")
	}

	identifier := args[0]
	value := args[1]

	// Get map info first to get key/value sizes
	var mapInfo *maps.MapInfo
	var mapID uint32
	var err error

	switch identifier {
	case "id":
		id, parseErr := strconv.ParseUint(value, 10, 32)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid map ID: %s\n", value)
			return bpferrors.ErrInvalidID
		}
		mapID = uint32(id)
		mapInfo, err = mapService.GetByID(mapID)
		if err != nil {
			handleError(err, fmt.Sprintf("getting map with ID %d", mapID))
			return err
		}

	case "name":
		mapInfos, getErr := mapService.GetByName(value)
		if getErr != nil {
			handleError(getErr, fmt.Sprintf("getting maps with name %s", value))
			return getErr
		}
		if len(mapInfos) == 0 {
			fmt.Fprintf(os.Stderr, "Error: no maps found with name: %s\n", value)
			return bpferrors.ErrNotFound
		}
		mapInfo = &mapInfos[0]
		mapID = mapInfo.ID

	case "pinned":
		mapInfo, err = mapService.GetByPinnedPath(value)
		if err != nil {
			handleError(err, fmt.Sprintf("getting pinned map at %s", value))
			return err
		}
		mapID = mapInfo.ID

	default:
		fmt.Fprintf(os.Stderr, "Error: invalid map identifier: %s. Use 'id', 'name', or 'pinned'\n", identifier)
		return fmt.Errorf("invalid identifier: %s", identifier)
	}

	// Dump all entries
	entries, err := mapService.Dump(mapID)
	if err != nil {
		handleError(err, fmt.Sprintf("dumping map %d", mapID))
		return err
	}

	// Convert to output.MapEntry
	outputEntries := make([]output.MapEntry, len(entries))
	for i, e := range entries {
		outputEntries[i] = output.MapEntry{
			Key:   e.Key,
			Value: e.Value,
		}
	}

	result := formatter.FormatMapEntries(outputEntries, mapInfo.KeySize, mapInfo.ValueSize)
	fmt.Print(result)

	return nil
}

// runMapLookup handles the map lookup command
func runMapLookup(cmd *cobra.Command, args []string) error {
	format := getOutputFormat()
	formatter := output.NewFormatter(format)

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: map identifier required. Use 'gobpftool map lookup <identifier> <value> key <key_data>'\n")
		return fmt.Errorf("map identifier required")
	}

	identifier := args[0]
	value := args[1]

	// Find the "key" keyword and parse key data
	keyIndex := -1
	for i, arg := range args {
		if arg == "key" {
			keyIndex = i
			break
		}
	}

	if keyIndex == -1 || keyIndex >= len(args)-1 {
		fmt.Fprintf(os.Stderr, "Error: key data required. Use 'gobpftool map lookup <identifier> <value> key <hex_bytes>'\n")
		return bpferrors.ErrInvalidKey
	}

	// Parse key data (space-separated hex bytes after "key")
	keyDataStr := strings.Join(args[keyIndex+1:], " ")
	keyData, err := utils.ParseHexBytes(keyDataStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid key format: %v\n", err)
		return bpferrors.ErrInvalidKey
	}

	// Get map info and lookup
	var mapInfo *maps.MapInfo
	var mapID uint32

	switch identifier {
	case "id":
		id, parseErr := strconv.ParseUint(value, 10, 32)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid map ID: %s\n", value)
			return bpferrors.ErrInvalidID
		}
		mapID = uint32(id)
		mapInfo, err = mapService.GetByID(mapID)
		if err != nil {
			handleError(err, fmt.Sprintf("getting map with ID %d", mapID))
			return err
		}

	case "name":
		mapInfos, getErr := mapService.GetByName(value)
		if getErr != nil {
			handleError(getErr, fmt.Sprintf("getting maps with name %s", value))
			return getErr
		}
		if len(mapInfos) == 0 {
			fmt.Fprintf(os.Stderr, "Error: no maps found with name: %s\n", value)
			return bpferrors.ErrNotFound
		}
		mapInfo = &mapInfos[0]
		mapID = mapInfo.ID

	case "pinned":
		mapInfo, err = mapService.GetByPinnedPath(value)
		if err != nil {
			handleError(err, fmt.Sprintf("getting pinned map at %s", value))
			return err
		}
		mapID = mapInfo.ID

	default:
		fmt.Fprintf(os.Stderr, "Error: invalid map identifier: %s. Use 'id', 'name', or 'pinned'\n", identifier)
		return fmt.Errorf("invalid identifier: %s", identifier)
	}

	// Lookup the key
	valueData, err := mapService.Lookup(mapID, keyData)
	if err != nil {
		if bpferrors.IsNotFoundError(err) {
			fmt.Fprintf(os.Stderr, "Error: key not found in map\n")
			return bpferrors.ErrKeyNotFound
		}
		handleError(err, "looking up key")
		return err
	}

	entry := output.MapEntry{
		Key:   keyData,
		Value: valueData,
	}

	result := formatter.FormatMapEntry(entry, mapInfo.KeySize, mapInfo.ValueSize)
	fmt.Print(result)

	return nil
}

// runMapGetNext handles the map getnext command
func runMapGetNext(cmd *cobra.Command, args []string) error {
	format := getOutputFormat()
	formatter := output.NewFormatter(format)

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: map identifier required. Use 'gobpftool map getnext <identifier> <value> [key <key_data>]'\n")
		return fmt.Errorf("map identifier required")
	}

	identifier := args[0]
	value := args[1]

	// Find the "key" keyword and parse key data (optional)
	var keyData []byte
	keyIndex := -1
	for i, arg := range args {
		if arg == "key" {
			keyIndex = i
			break
		}
	}

	if keyIndex != -1 && keyIndex < len(args)-1 {
		// Parse key data (space-separated hex bytes after "key")
		keyDataStr := strings.Join(args[keyIndex+1:], " ")
		var err error
		keyData, err = utils.ParseHexBytes(keyDataStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid key format: %v\n", err)
			return bpferrors.ErrInvalidKey
		}
	}

	// Get map info
	var mapID uint32
	var err error

	switch identifier {
	case "id":
		id, parseErr := strconv.ParseUint(value, 10, 32)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid map ID: %s\n", value)
			return bpferrors.ErrInvalidID
		}
		mapID = uint32(id)
		_, err = mapService.GetByID(mapID)
		if err != nil {
			handleError(err, fmt.Sprintf("getting map with ID %d", mapID))
			return err
		}

	case "name":
		mapInfos, getErr := mapService.GetByName(value)
		if getErr != nil {
			handleError(getErr, fmt.Sprintf("getting maps with name %s", value))
			return getErr
		}
		if len(mapInfos) == 0 {
			fmt.Fprintf(os.Stderr, "Error: no maps found with name: %s\n", value)
			return bpferrors.ErrNotFound
		}
		mapID = mapInfos[0].ID

	case "pinned":
		mapInfo, getErr := mapService.GetByPinnedPath(value)
		if getErr != nil {
			handleError(getErr, fmt.Sprintf("getting pinned map at %s", value))
			return getErr
		}
		mapID = mapInfo.ID

	default:
		fmt.Fprintf(os.Stderr, "Error: invalid map identifier: %s. Use 'id', 'name', or 'pinned'\n", identifier)
		return fmt.Errorf("invalid identifier: %s", identifier)
	}

	// Get next key
	nextKey, err := mapService.GetNextKey(mapID, keyData)
	if err != nil {
		// Check if it's a "no more keys" error
		if bpferrors.IsNoMoreKeysError(err) {
			if keyData == nil {
				fmt.Fprintf(os.Stderr, "Error: map is empty\n")
				return bpferrors.ErrMapEmpty
			}
			fmt.Fprintf(os.Stderr, "Error: no more keys\n")
			return bpferrors.ErrNoMoreKeys
		}
		handleError(err, "getting next key")
		return err
	}

	result := formatter.FormatNextKey(keyData, nextKey)
	fmt.Print(result)

	return nil
}

func init() {
	// Initialize the map service
	mapService = maps.NewService()

	// Add subcommands to map command
	mapCmd.AddCommand(mapShowCmd)
	mapCmd.AddCommand(mapDumpCmd)
	mapCmd.AddCommand(mapLookupCmd)
	mapCmd.AddCommand(mapGetNextCmd)
	mapCmd.AddCommand(mapHelpCmd)

	// Add map command to root command
	rootCmd.AddCommand(mapCmd)
}
