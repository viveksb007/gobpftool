package output

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestJSONFormatter_FormatPrograms(t *testing.T) {
	loadedAt := time.Date(2025, 11, 24, 5, 50, 46, 0, time.UTC)

	tests := []struct {
		name   string
		pretty bool
		progs  []ProgramInfo
		check  func(t *testing.T, result string)
	}{
		{
			name:   "empty list compact",
			pretty: false,
			progs:  []ProgramInfo{},
			check: func(t *testing.T, result string) {
				expected := `{"programs":[]}`
				if result != expected {
					t.Errorf("got %q, want %q", result, expected)
				}
			},
		},
		{
			name:   "single program compact",
			pretty: false,
			progs: []ProgramInfo{
				{
					ID:        185,
					Type:      "sched_cls",
					Name:      "my_prog",
					Tag:       "f0055c08993fea1e",
					GPL:       true,
					LoadedAt:  loadedAt,
					UID:       0,
					BytesXlat: 5200,
					BytesJIT:  3263,
					MemLock:   8192,
					MapIDs:    []uint32{85, 39, 38},
				},
			},
			check: func(t *testing.T, result string) {
				var parsed programsJSON
				if err := json.Unmarshal([]byte(result), &parsed); err != nil {
					t.Fatalf("failed to parse JSON: %v", err)
				}
				if len(parsed.Programs) != 1 {
					t.Fatalf("expected 1 program, got %d", len(parsed.Programs))
				}
				p := parsed.Programs[0]
				if p.ID != 185 {
					t.Errorf("ID = %d, want 185", p.ID)
				}
				if p.Type != "sched_cls" {
					t.Errorf("Type = %q, want %q", p.Type, "sched_cls")
				}
				if p.Name != "my_prog" {
					t.Errorf("Name = %q, want %q", p.Name, "my_prog")
				}
				if !p.GPLCompatible {
					t.Error("GPLCompatible = false, want true")
				}
				if len(p.MapIDs) != 3 {
					t.Errorf("MapIDs length = %d, want 3", len(p.MapIDs))
				}
			},
		},
		{
			name:   "program without map_ids",
			pretty: false,
			progs: []ProgramInfo{
				{
					ID:        10,
					Type:      "kprobe",
					Name:      "test",
					Tag:       "abcd1234",
					GPL:       false,
					LoadedAt:  loadedAt,
					UID:       1000,
					BytesXlat: 100,
					BytesJIT:  80,
					MemLock:   4096,
					MapIDs:    nil,
				},
			},
			check: func(t *testing.T, result string) {
				var parsed programsJSON
				if err := json.Unmarshal([]byte(result), &parsed); err != nil {
					t.Fatalf("failed to parse JSON: %v", err)
				}
				p := parsed.Programs[0]
				if p.GPLCompatible {
					t.Error("GPLCompatible = true, want false")
				}
				if p.MapIDs != nil {
					t.Errorf("MapIDs = %v, want nil", p.MapIDs)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := &JSONFormatter{pretty: tt.pretty}
			result := formatter.FormatPrograms(tt.progs)
			tt.check(t, result)
		})
	}
}

func TestJSONFormatter_FormatPrograms_Pretty(t *testing.T) {
	loadedAt := time.Date(2025, 11, 24, 5, 50, 46, 0, time.UTC)
	formatter := &JSONFormatter{pretty: true}

	progs := []ProgramInfo{
		{
			ID:        1,
			Type:      "xdp",
			Name:      "test",
			Tag:       "12345678",
			GPL:       true,
			LoadedAt:  loadedAt,
			UID:       0,
			BytesXlat: 100,
			BytesJIT:  50,
			MemLock:   4096,
			MapIDs:    []uint32{1},
		},
	}

	result := formatter.FormatPrograms(progs)

	// Pretty printed JSON should contain newlines and indentation
	if len(result) == 0 {
		t.Fatal("result is empty")
	}

	// Verify it's valid JSON
	var parsed programsJSON
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("failed to parse pretty JSON: %v", err)
	}

	// Check for indentation (pretty print indicator)
	if result[0] != '{' {
		t.Error("expected JSON to start with '{'")
	}
}

func TestJSONFormatter_FormatMaps(t *testing.T) {
	tests := []struct {
		name   string
		pretty bool
		maps   []MapInfo
		check  func(t *testing.T, result string)
	}{
		{
			name:   "empty list",
			pretty: false,
			maps:   []MapInfo{},
			check: func(t *testing.T, result string) {
				expected := `{"maps":[]}`
				if result != expected {
					t.Errorf("got %q, want %q", result, expected)
				}
			},
		},
		{
			name:   "single map",
			pretty: false,
			maps: []MapInfo{
				{
					ID:         10,
					Type:       "hash",
					Name:       "some_map",
					KeySize:    4,
					ValueSize:  8,
					MaxEntries: 2048,
					Flags:      0,
					MemLock:    167936,
				},
			},
			check: func(t *testing.T, result string) {
				var parsed mapsJSON
				if err := json.Unmarshal([]byte(result), &parsed); err != nil {
					t.Fatalf("failed to parse JSON: %v", err)
				}
				if len(parsed.Maps) != 1 {
					t.Fatalf("expected 1 map, got %d", len(parsed.Maps))
				}
				m := parsed.Maps[0]
				if m.ID != 10 {
					t.Errorf("ID = %d, want 10", m.ID)
				}
				if m.Type != "hash" {
					t.Errorf("Type = %q, want %q", m.Type, "hash")
				}
				if m.KeySize != 4 {
					t.Errorf("KeySize = %d, want 4", m.KeySize)
				}
				if m.ValueSize != 8 {
					t.Errorf("ValueSize = %d, want 8", m.ValueSize)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := &JSONFormatter{pretty: tt.pretty}
			result := formatter.FormatMaps(tt.maps)
			tt.check(t, result)
		})
	}
}

func TestJSONFormatter_FormatMapEntries(t *testing.T) {
	formatter := &JSONFormatter{pretty: false}

	entries := []MapEntry{
		{
			Key:   []byte{0x00, 0x01, 0x02, 0x03},
			Value: []byte{0x10, 0x11, 0x12, 0x13},
		},
		{
			Key:   []byte{0x04, 0x05, 0x06, 0x07},
			Value: []byte{0x20, 0x21, 0x22, 0x23},
		},
	}

	result := formatter.FormatMapEntries(entries, 4, 4)

	var parsed mapEntriesJSON
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed.Count != 2 {
		t.Errorf("Count = %d, want 2", parsed.Count)
	}
	if len(parsed.Entries) != 2 {
		t.Errorf("Entries length = %d, want 2", len(parsed.Entries))
	}
}

func TestJSONFormatter_FormatMapEntry(t *testing.T) {
	formatter := &JSONFormatter{pretty: false}

	entry := MapEntry{
		Key:   []byte{0x00, 0x01, 0x02, 0x03},
		Value: []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17},
	}

	result := formatter.FormatMapEntry(entry, 4, 8)

	var parsed mapEntryJSON
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(parsed.Key) != 4 {
		t.Errorf("Key length = %d, want 4", len(parsed.Key))
	}
	if len(parsed.Value) != 8 {
		t.Errorf("Value length = %d, want 8", len(parsed.Value))
	}
}

func TestJSONFormatter_FormatNextKey(t *testing.T) {
	tests := []struct {
		name       string
		currentKey []byte
		nextKey    []byte
		check      func(t *testing.T, result string)
	}{
		{
			name:       "first key (no current)",
			currentKey: nil,
			nextKey:    []byte{0x00, 0x01, 0x02, 0x03},
			check: func(t *testing.T, result string) {
				var parsed nextKeyJSON
				if err := json.Unmarshal([]byte(result), &parsed); err != nil {
					t.Fatalf("failed to parse JSON: %v", err)
				}
				if parsed.Key != nil {
					t.Errorf("Key = %v, want nil", parsed.Key)
				}
				if len(parsed.NextKey) != 4 {
					t.Errorf("NextKey length = %d, want 4", len(parsed.NextKey))
				}
			},
		},
		{
			name:       "with current key",
			currentKey: []byte{0x00, 0x01, 0x02, 0x03},
			nextKey:    []byte{0x04, 0x05, 0x06, 0x07},
			check: func(t *testing.T, result string) {
				var parsed nextKeyJSON
				if err := json.Unmarshal([]byte(result), &parsed); err != nil {
					t.Fatalf("failed to parse JSON: %v", err)
				}
				if len(parsed.Key) != 4 {
					t.Errorf("Key length = %d, want 4", len(parsed.Key))
				}
				if len(parsed.NextKey) != 4 {
					t.Errorf("NextKey length = %d, want 4", len(parsed.NextKey))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := &JSONFormatter{pretty: false}
			result := formatter.FormatNextKey(tt.currentKey, tt.nextKey)
			tt.check(t, result)
		})
	}
}

func TestJSONFormatter_FormatError(t *testing.T) {
	formatter := &JSONFormatter{pretty: false}

	err := fmt.Errorf("something went wrong")
	result := formatter.FormatError(err)

	var parsed errorJSON
	if jsonErr := json.Unmarshal([]byte(result), &parsed); jsonErr != nil {
		t.Fatalf("failed to parse JSON: %v", jsonErr)
	}

	if parsed.Error != "something went wrong" {
		t.Errorf("Error = %q, want %q", parsed.Error, "something went wrong")
	}
}

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name     string
		format   Format
		wantType string
	}{
		{
			name:     "plain format",
			format:   FormatPlain,
			wantType: "*output.PlainFormatter",
		},
		{
			name:     "JSON format",
			format:   FormatJSON,
			wantType: "*output.JSONFormatter",
		},
		{
			name:     "JSON pretty format",
			format:   FormatJSONPretty,
			wantType: "*output.JSONFormatter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFormatter(tt.format)
			if f == nil {
				t.Fatal("NewFormatter returned nil")
			}
			// Just verify it implements the interface
			var _ Formatter = f
		})
	}
}
