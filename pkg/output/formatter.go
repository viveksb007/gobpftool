// Package output provides formatters for displaying eBPF program and map information.
package output

import "time"

// Format represents the output format type.
type Format int

const (
	// FormatPlain outputs human-readable plain text.
	FormatPlain Format = iota
	// FormatJSON outputs compact JSON.
	FormatJSON
	// FormatJSONPretty outputs pretty-printed JSON with indentation.
	FormatJSONPretty
)

// ProgramInfo contains information about an eBPF program.
type ProgramInfo struct {
	ID        uint32
	Type      string
	Name      string
	Tag       string
	GPL       bool
	LoadedAt  time.Time
	UID       uint32
	BytesXlat uint32
	BytesJIT  uint32
	MemLock   uint32
	MapIDs    []uint32
}

// MapInfo contains information about an eBPF map.
type MapInfo struct {
	ID         uint32
	Type       string
	Name       string
	KeySize    uint32
	ValueSize  uint32
	MaxEntries uint32
	Flags      uint32
	MemLock    uint32
}

// MapEntry represents a key-value pair in an eBPF map.
type MapEntry struct {
	Key   []byte
	Value []byte
}

// Formatter defines the interface for formatting eBPF program and map output.
type Formatter interface {
	// FormatPrograms formats a list of programs for output.
	FormatPrograms(progs []ProgramInfo) string

	// FormatMaps formats a list of maps for output.
	FormatMaps(maps []MapInfo) string

	// FormatMapEntries formats map entries for output (used by dump).
	FormatMapEntries(entries []MapEntry, keySize, valueSize uint32) string

	// FormatMapEntry formats a single map entry (used by lookup).
	FormatMapEntry(entry MapEntry, keySize, valueSize uint32) string

	// FormatNextKey formats the next key result (used by getnext).
	FormatNextKey(currentKey, nextKey []byte) string

	// FormatError formats an error message.
	FormatError(err error) string
}

// NewFormatter creates a new Formatter based on the specified format.
func NewFormatter(format Format) Formatter {
	switch format {
	case FormatJSON:
		return &JSONFormatter{pretty: false}
	case FormatJSONPretty:
		return &JSONFormatter{pretty: true}
	default:
		return &PlainFormatter{}
	}
}
