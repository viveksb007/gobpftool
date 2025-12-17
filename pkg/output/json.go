package output

import (
	"encoding/json"
	"fmt"
)

// JSONFormatter formats output as JSON, compatible with bpftool JSON output.
type JSONFormatter struct {
	pretty bool
}

// programJSON represents a program in bpftool-compatible JSON format.
type programJSON struct {
	ID            uint32   `json:"id"`
	Type          string   `json:"type"`
	Name          string   `json:"name"`
	Tag           string   `json:"tag"`
	GPLCompatible bool     `json:"gpl_compatible"`
	LoadedAt      string   `json:"loaded_at"`
	UID           uint32   `json:"uid"`
	BytesXlated   uint32   `json:"bytes_xlated"`
	BytesJited    uint32   `json:"bytes_jited"`
	BytesMemlock  uint32   `json:"bytes_memlock"`
	MapIDs        []uint32 `json:"map_ids,omitempty"`
}

// programsJSON wraps programs for JSON output.
type programsJSON struct {
	Programs []programJSON `json:"programs"`
}

// mapJSON represents a map in bpftool-compatible JSON format.
type mapJSON struct {
	ID           uint32 `json:"id"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	KeySize      uint32 `json:"key_size"`
	ValueSize    uint32 `json:"value_size"`
	MaxEntries   uint32 `json:"max_entries"`
	Flags        uint32 `json:"flags"`
	BytesMemlock uint32 `json:"bytes_memlock"`
}

// mapsJSON wraps maps for JSON output.
type mapsJSON struct {
	Maps []mapJSON `json:"maps"`
}

// mapEntryJSON represents a map entry in JSON format.
type mapEntryJSON struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

// mapEntriesJSON wraps map entries for JSON output.
type mapEntriesJSON struct {
	Entries []mapEntryJSON `json:"entries"`
	Count   int            `json:"count"`
}

// nextKeyJSON represents a next key result in JSON format.
type nextKeyJSON struct {
	Key     []byte `json:"key,omitempty"`
	NextKey []byte `json:"next_key"`
}

// errorJSON represents an error in JSON format.
type errorJSON struct {
	Error string `json:"error"`
}

// FormatPrograms formats programs as JSON.
func (f *JSONFormatter) FormatPrograms(progs []ProgramInfo) string {
	programs := make([]programJSON, len(progs))
	for i, p := range progs {
		programs[i] = programJSON{
			ID:            p.ID,
			Type:          p.Type,
			Name:          p.Name,
			Tag:           p.Tag,
			GPLCompatible: p.GPL,
			LoadedAt:      p.LoadedAt.Format("2006-01-02T15:04:05-0700"),
			UID:           p.UID,
			BytesXlated:   p.BytesXlat,
			BytesJited:    p.BytesJIT,
			BytesMemlock:  p.MemLock,
			MapIDs:        p.MapIDs,
		}
	}

	return f.marshal(programsJSON{Programs: programs})
}

// FormatMaps formats maps as JSON.
func (f *JSONFormatter) FormatMaps(maps []MapInfo) string {
	jsonMaps := make([]mapJSON, len(maps))
	for i, m := range maps {
		jsonMaps[i] = mapJSON{
			ID:           m.ID,
			Type:         m.Type,
			Name:         m.Name,
			KeySize:      m.KeySize,
			ValueSize:    m.ValueSize,
			MaxEntries:   m.MaxEntries,
			Flags:        m.Flags,
			BytesMemlock: m.MemLock,
		}
	}

	return f.marshal(mapsJSON{Maps: jsonMaps})
}

// FormatMapEntries formats map entries as JSON.
func (f *JSONFormatter) FormatMapEntries(entries []MapEntry, keySize, valueSize uint32) string {
	jsonEntries := make([]mapEntryJSON, len(entries))
	for i, e := range entries {
		jsonEntries[i] = mapEntryJSON{
			Key:   e.Key,
			Value: e.Value,
		}
	}

	return f.marshal(mapEntriesJSON{
		Entries: jsonEntries,
		Count:   len(entries),
	})
}

// FormatMapEntry formats a single map entry as JSON.
func (f *JSONFormatter) FormatMapEntry(entry MapEntry, keySize, valueSize uint32) string {
	return f.marshal(mapEntryJSON{
		Key:   entry.Key,
		Value: entry.Value,
	})
}

// FormatNextKey formats the next key result as JSON.
func (f *JSONFormatter) FormatNextKey(currentKey, nextKey []byte) string {
	return f.marshal(nextKeyJSON{
		Key:     currentKey,
		NextKey: nextKey,
	})
}

// FormatError formats an error as JSON.
func (f *JSONFormatter) FormatError(err error) string {
	return f.marshal(errorJSON{Error: err.Error()})
}

// marshal converts data to JSON string, with optional pretty printing.
func (f *JSONFormatter) marshal(v interface{}) string {
	var data []byte
	var err error

	if f.pretty {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return fmt.Sprintf(`{"error":"failed to marshal JSON: %v"}`, err)
	}

	return string(data)
}
