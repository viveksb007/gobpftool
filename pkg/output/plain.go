package output

import (
	"fmt"
	"strings"
)

// PlainFormatter formats output as human-readable plain text matching bpftool format.
type PlainFormatter struct{}

// FormatPrograms formats programs in bpftool-compatible plain text format.
// Format:
//
//	<ID>: <type>  name <name>  tag <tag>  gpl
//	        loaded_at <timestamp>  uid <uid>
//	        xlated <bytes>B  jited <bytes>B  memlock <bytes>B  map_ids <id1>,<id2>,...
func (f *PlainFormatter) FormatPrograms(progs []ProgramInfo) string {
	if len(progs) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, p := range progs {
		if i > 0 {
			sb.WriteString("\n")
		}
		f.formatProgram(&sb, p)
	}
	return sb.String()
}

func (f *PlainFormatter) formatProgram(sb *strings.Builder, p ProgramInfo) {
	// First line: ID, type, name, tag, gpl
	gplStr := ""
	if p.GPL {
		gplStr = "  gpl"
	}
	fmt.Fprintf(sb, "%d: %s  name %s  tag %s%s\n",
		p.ID, p.Type, p.Name, p.Tag, gplStr)

	// Second line: loaded_at, uid
	loadedAt := p.LoadedAt.Format("2006-01-02T15:04:05-0700")
	fmt.Fprintf(sb, "\tloaded_at %s  uid %d\n", loadedAt, p.UID)

	// Third line: xlated, jited, memlock, map_ids
	fmt.Fprintf(sb, "\txlated %dB  jited %dB  memlock %dB",
		p.BytesXlat, p.BytesJIT, p.MemLock)

	if len(p.MapIDs) > 0 {
		mapIDStrs := make([]string, len(p.MapIDs))
		for i, id := range p.MapIDs {
			mapIDStrs[i] = fmt.Sprintf("%d", id)
		}
		fmt.Fprintf(sb, "  map_ids %s", strings.Join(mapIDStrs, ","))
	}
}

// FormatMaps formats maps in bpftool-compatible plain text format.
// Format:
//
//	<ID>: <type>  name <name>  flags 0x<flags>
//	        key <size>B  value <size>B  max_entries <count>  memlock <bytes>B
func (f *PlainFormatter) FormatMaps(maps []MapInfo) string {
	if len(maps) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, m := range maps {
		if i > 0 {
			sb.WriteString("\n")
		}
		f.formatMap(&sb, m)
	}
	return sb.String()
}

func (f *PlainFormatter) formatMap(sb *strings.Builder, m MapInfo) {
	// First line: ID, type, name, flags
	fmt.Fprintf(sb, "%d: %s  name %s  flags 0x%x\n",
		m.ID, m.Type, m.Name, m.Flags)

	// Second line: key, value, max_entries, memlock
	fmt.Fprintf(sb, "\tkey %dB  value %dB  max_entries %d  memlock %dB",
		m.KeySize, m.ValueSize, m.MaxEntries, m.MemLock)
}

// FormatMapEntries formats all map entries for dump output.
// Format:
//
//	key: <hex bytes>  value: <hex bytes>
//	...
//	Found <n> elements
func (f *PlainFormatter) FormatMapEntries(entries []MapEntry, keySize, valueSize uint32) string {
	var sb strings.Builder

	for _, entry := range entries {
		keyHex := formatHexBytes(entry.Key)
		valueHex := formatHexBytes(entry.Value)
		fmt.Fprintf(&sb, "key: %s  value: %s\n", keyHex, valueHex)
	}

	fmt.Fprintf(&sb, "Found %d element", len(entries))
	if len(entries) != 1 {
		sb.WriteString("s")
	}

	return sb.String()
}

// FormatMapEntry formats a single map entry for lookup output.
// Format: key: <hex bytes> value: <hex bytes>
func (f *PlainFormatter) FormatMapEntry(entry MapEntry, keySize, valueSize uint32) string {
	keyHex := formatHexBytes(entry.Key)
	valueHex := formatHexBytes(entry.Value)
	return fmt.Sprintf("key: %s value: %s", keyHex, valueHex)
}

// FormatNextKey formats the next key result for getnext output.
// Format:
//
//	key:
//	<hex bytes>
//	next key:
//	<hex bytes>
func (f *PlainFormatter) FormatNextKey(currentKey, nextKey []byte) string {
	var sb strings.Builder

	if currentKey != nil {
		sb.WriteString("key:\n")
		sb.WriteString(formatHexBytes(currentKey))
		sb.WriteString("\n")
	}

	sb.WriteString("next key:\n")
	sb.WriteString(formatHexBytes(nextKey))

	return sb.String()
}

// FormatError formats an error message for stderr output.
func (f *PlainFormatter) FormatError(err error) string {
	return fmt.Sprintf("Error: %v", err)
}

// formatHexBytes converts a byte slice to space-separated hex string.
func formatHexBytes(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	hexParts := make([]string, len(data))
	for i, b := range data {
		hexParts[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(hexParts, " ")
}
