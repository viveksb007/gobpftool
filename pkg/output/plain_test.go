package output

import (
	"fmt"
	"testing"
	"time"
)

func TestPlainFormatter_FormatPrograms(t *testing.T) {
	formatter := &PlainFormatter{}

	// Test with loaded_at time
	loadedAt := time.Date(2025, 11, 24, 5, 50, 46, 0, time.UTC)

	tests := []struct {
		name     string
		progs    []ProgramInfo
		expected string
	}{
		{
			name:     "empty list",
			progs:    []ProgramInfo{},
			expected: "",
		},
		{
			name: "single program with GPL and map_ids",
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
					MapIDs:    []uint32{85, 39, 38, 83, 84},
				},
			},
			expected: "185: sched_cls  name my_prog  tag f0055c08993fea1e  gpl\n" +
				"\tloaded_at 2025-11-24T05:50:46+0000  uid 0\n" +
				"\txlated 5200B  jited 3263B  memlock 8192B  map_ids 85,39,38,83,84",
		},
		{
			name: "program without GPL",
			progs: []ProgramInfo{
				{
					ID:        10,
					Type:      "kprobe",
					Name:      "test_prog",
					Tag:       "abcd1234abcd1234",
					GPL:       false,
					LoadedAt:  loadedAt,
					UID:       1000,
					BytesXlat: 100,
					BytesJIT:  80,
					MemLock:   4096,
					MapIDs:    nil,
				},
			},
			expected: "10: kprobe  name test_prog  tag abcd1234abcd1234\n" +
				"\tloaded_at 2025-11-24T05:50:46+0000  uid 1000\n" +
				"\txlated 100B  jited 80B  memlock 4096B",
		},
		{
			name: "multiple programs",
			progs: []ProgramInfo{
				{
					ID:        1,
					Type:      "xdp",
					Name:      "prog1",
					Tag:       "1111111111111111",
					GPL:       true,
					LoadedAt:  loadedAt,
					UID:       0,
					BytesXlat: 100,
					BytesJIT:  50,
					MemLock:   4096,
					MapIDs:    []uint32{1},
				},
				{
					ID:        2,
					Type:      "tc",
					Name:      "prog2",
					Tag:       "2222222222222222",
					GPL:       false,
					LoadedAt:  loadedAt,
					UID:       0,
					BytesXlat: 200,
					BytesJIT:  100,
					MemLock:   8192,
					MapIDs:    nil,
				},
			},
			expected: "1: xdp  name prog1  tag 1111111111111111  gpl\n" +
				"\tloaded_at 2025-11-24T05:50:46+0000  uid 0\n" +
				"\txlated 100B  jited 50B  memlock 4096B  map_ids 1\n" +
				"2: tc  name prog2  tag 2222222222222222\n" +
				"\tloaded_at 2025-11-24T05:50:46+0000  uid 0\n" +
				"\txlated 200B  jited 100B  memlock 8192B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatPrograms(tt.progs)
			if result != tt.expected {
				t.Errorf("FormatPrograms() =\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

func TestPlainFormatter_FormatMaps(t *testing.T) {
	formatter := &PlainFormatter{}

	tests := []struct {
		name     string
		maps     []MapInfo
		expected string
	}{
		{
			name:     "empty list",
			maps:     []MapInfo{},
			expected: "",
		},
		{
			name: "single map",
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
			expected: "10: hash  name some_map  flags 0x0\n" +
				"\tkey 4B  value 8B  max_entries 2048  memlock 167936B",
		},
		{
			name: "map with flags",
			maps: []MapInfo{
				{
					ID:         20,
					Type:       "array",
					Name:       "my_array",
					KeySize:    4,
					ValueSize:  16,
					MaxEntries: 100,
					Flags:      0x1,
					MemLock:    8192,
				},
			},
			expected: "20: array  name my_array  flags 0x1\n" +
				"\tkey 4B  value 16B  max_entries 100  memlock 8192B",
		},
		{
			name: "multiple maps",
			maps: []MapInfo{
				{
					ID:         1,
					Type:       "hash",
					Name:       "map1",
					KeySize:    4,
					ValueSize:  4,
					MaxEntries: 100,
					Flags:      0,
					MemLock:    4096,
				},
				{
					ID:         2,
					Type:       "array",
					Name:       "map2",
					KeySize:    4,
					ValueSize:  8,
					MaxEntries: 50,
					Flags:      0,
					MemLock:    2048,
				},
			},
			expected: "1: hash  name map1  flags 0x0\n" +
				"\tkey 4B  value 4B  max_entries 100  memlock 4096B\n" +
				"2: array  name map2  flags 0x0\n" +
				"\tkey 4B  value 8B  max_entries 50  memlock 2048B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatMaps(tt.maps)
			if result != tt.expected {
				t.Errorf("FormatMaps() =\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

func TestPlainFormatter_FormatMapEntries(t *testing.T) {
	formatter := &PlainFormatter{}

	tests := []struct {
		name      string
		entries   []MapEntry
		keySize   uint32
		valueSize uint32
		expected  string
	}{
		{
			name:      "empty entries",
			entries:   []MapEntry{},
			keySize:   4,
			valueSize: 8,
			expected:  "Found 0 elements",
		},
		{
			name: "single entry",
			entries: []MapEntry{
				{
					Key:   []byte{0x00, 0x01, 0x02, 0x03},
					Value: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
				},
			},
			keySize:   4,
			valueSize: 8,
			expected:  "key: 00 01 02 03  value: 00 01 02 03 04 05 06 07\nFound 1 element",
		},
		{
			name: "multiple entries",
			entries: []MapEntry{
				{
					Key:   []byte{0x00, 0x01, 0x02, 0x03},
					Value: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
				},
				{
					Key:   []byte{0x0d, 0x00, 0x07, 0x00},
					Value: []byte{0x02, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04},
				},
			},
			keySize:   4,
			valueSize: 8,
			expected: "key: 00 01 02 03  value: 00 01 02 03 04 05 06 07\n" +
				"key: 0d 00 07 00  value: 02 00 00 00 01 02 03 04\n" +
				"Found 2 elements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatMapEntries(tt.entries, tt.keySize, tt.valueSize)
			if result != tt.expected {
				t.Errorf("FormatMapEntries() =\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

func TestPlainFormatter_FormatMapEntry(t *testing.T) {
	formatter := &PlainFormatter{}

	tests := []struct {
		name      string
		entry     MapEntry
		keySize   uint32
		valueSize uint32
		expected  string
	}{
		{
			name: "standard entry",
			entry: MapEntry{
				Key:   []byte{0x00, 0x01, 0x02, 0x03},
				Value: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
			},
			keySize:   4,
			valueSize: 8,
			expected:  "key: 00 01 02 03 value: 00 01 02 03 04 05 06 07",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatMapEntry(tt.entry, tt.keySize, tt.valueSize)
			if result != tt.expected {
				t.Errorf("FormatMapEntry() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestPlainFormatter_FormatNextKey(t *testing.T) {
	formatter := &PlainFormatter{}

	tests := []struct {
		name       string
		currentKey []byte
		nextKey    []byte
		expected   string
	}{
		{
			name:       "first key (no current)",
			currentKey: nil,
			nextKey:    []byte{0x00, 0x01, 0x02, 0x03},
			expected:   "next key:\n00 01 02 03",
		},
		{
			name:       "with current key",
			currentKey: []byte{0x00, 0x01, 0x02, 0x03},
			nextKey:    []byte{0x0d, 0x00, 0x07, 0x00},
			expected:   "key:\n00 01 02 03\nnext key:\n0d 00 07 00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatNextKey(tt.currentKey, tt.nextKey)
			if result != tt.expected {
				t.Errorf("FormatNextKey() =\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

func TestPlainFormatter_FormatError(t *testing.T) {
	formatter := &PlainFormatter{}

	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "simple error",
			err:      fmt.Errorf("something went wrong"),
			expected: "Error: something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatError(tt.err)
			if result != tt.expected {
				t.Errorf("FormatError() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatHexBytes(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "empty",
			data:     []byte{},
			expected: "",
		},
		{
			name:     "single byte",
			data:     []byte{0x0a},
			expected: "0a",
		},
		{
			name:     "multiple bytes",
			data:     []byte{0x00, 0x01, 0x02, 0x03},
			expected: "00 01 02 03",
		},
		{
			name:     "high values",
			data:     []byte{0xff, 0xfe, 0xfd},
			expected: "ff fe fd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatHexBytes(tt.data)
			if result != tt.expected {
				t.Errorf("formatHexBytes() = %q, want %q", result, tt.expected)
			}
		})
	}
}
