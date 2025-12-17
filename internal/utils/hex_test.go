package utils

import (
	"reflect"
	"testing"
)

func TestParseHexBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
		wantErr  bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []byte{},
			wantErr:  false,
		},
		{
			name:     "single byte",
			input:    "0a",
			expected: []byte{0x0a},
			wantErr:  false,
		},
		{
			name:     "multiple bytes",
			input:    "0a 0b 0c 0d",
			expected: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			wantErr:  false,
		},
		{
			name:     "uppercase hex",
			input:    "0A 0B 0C 0D",
			expected: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			wantErr:  false,
		},
		{
			name:     "mixed case",
			input:    "0a 0B 0c 0D",
			expected: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			wantErr:  false,
		},
		{
			name:     "single digit hex",
			input:    "a b c d",
			expected: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			wantErr:  false,
		},
		{
			name:     "zero bytes",
			input:    "00 00 00",
			expected: []byte{0x00, 0x00, 0x00},
			wantErr:  false,
		},
		{
			name:     "max byte values",
			input:    "ff FF",
			expected: []byte{0xff, 0xff},
			wantErr:  false,
		},
		{
			name:     "extra whitespace",
			input:    "  0a   0b  0c  ",
			expected: []byte{0x0a, 0x0b, 0x0c},
			wantErr:  false,
		},
		{
			name:     "tabs and spaces",
			input:    "0a\t0b\n0c",
			expected: []byte{0x0a, 0x0b, 0x0c},
			wantErr:  false,
		},
		{
			name:     "invalid hex character",
			input:    "0g",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "too long hex byte",
			input:    "0abc",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "non-hex string",
			input:    "hello",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "empty after split",
			input:    "   ",
			expected: []byte{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseHexBytes(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseHexBytes() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseHexBytes() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseHexBytes() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestFormatHexBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty slice",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "single byte",
			input:    []byte{0x0a},
			expected: "0a",
		},
		{
			name:     "multiple bytes",
			input:    []byte{0x0a, 0x0b, 0x0c, 0x0d},
			expected: "0a 0b 0c 0d",
		},
		{
			name:     "zero bytes",
			input:    []byte{0x00, 0x00, 0x00},
			expected: "00 00 00",
		},
		{
			name:     "max byte values",
			input:    []byte{0xff, 0xff},
			expected: "ff ff",
		},
		{
			name:     "mixed values",
			input:    []byte{0x00, 0x0a, 0xff, 0x10},
			expected: "00 0a ff 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatHexBytes(tt.input)
			if result != tt.expected {
				t.Errorf("FormatHexBytes() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestFormatHexBytesWithPrefix(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		input    []byte
		expected string
	}{
		{
			name:     "empty slice with key prefix",
			prefix:   "key",
			input:    []byte{},
			expected: "key:",
		},
		{
			name:     "single byte with key prefix",
			prefix:   "key",
			input:    []byte{0x0a},
			expected: "key: 0a",
		},
		{
			name:     "multiple bytes with key prefix",
			prefix:   "key",
			input:    []byte{0x0a, 0x0b, 0x0c, 0x0d},
			expected: "key: 0a 0b 0c 0d",
		},
		{
			name:     "multiple bytes with value prefix",
			prefix:   "value",
			input:    []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			expected: "value: 01 02 03 04 05 06 07 08",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatHexBytesWithPrefix(tt.prefix, tt.input)
			if result != tt.expected {
				t.Errorf("FormatHexBytesWithPrefix() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestParseHexString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
		wantErr  bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []byte{},
			wantErr:  false,
		},
		{
			name:     "simple hex string",
			input:    "0a0b0c0d",
			expected: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			wantErr:  false,
		},
		{
			name:     "uppercase hex string",
			input:    "0A0B0C0D",
			expected: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			wantErr:  false,
		},
		{
			name:     "mixed case hex string",
			input:    "0a0B0c0D",
			expected: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			wantErr:  false,
		},
		{
			name:     "hex string with spaces (should be removed)",
			input:    "0a 0b 0c 0d",
			expected: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			wantErr:  false,
		},
		{
			name:     "program tag format",
			input:    "f0055c08993fea1e",
			expected: []byte{0xf0, 0x05, 0x5c, 0x08, 0x99, 0x3f, 0xea, 0x1e},
			wantErr:  false,
		},
		{
			name:     "odd length string",
			input:    "0a0b0",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "invalid hex character",
			input:    "0g0h",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseHexString(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseHexString() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseHexString() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseHexString() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestFormatHexString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty slice",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "simple bytes",
			input:    []byte{0x0a, 0x0b, 0x0c, 0x0d},
			expected: "0a0b0c0d",
		},
		{
			name:     "program tag format",
			input:    []byte{0xf0, 0x05, 0x5c, 0x08, 0x99, 0x3f, 0xea, 0x1e},
			expected: "f0055c08993fea1e",
		},
		{
			name:     "zero bytes",
			input:    []byte{0x00, 0x00, 0x00},
			expected: "000000",
		},
		{
			name:     "max byte values",
			input:    []byte{0xff, 0xff},
			expected: "ffff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatHexString(tt.input)
			if result != tt.expected {
				t.Errorf("FormatHexString() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// Test round-trip conversion for ParseHexBytes and FormatHexBytes
func TestHexBytesRoundTrip(t *testing.T) {
	testCases := [][]byte{
		{},
		{0x0a},
		{0x0a, 0x0b, 0x0c, 0x0d},
		{0x00, 0x00, 0x00},
		{0xff, 0xff},
		{0x00, 0x0a, 0xff, 0x10},
	}

	for _, original := range testCases {
		t.Run("round trip", func(t *testing.T) {
			// Convert to string and back
			hexStr := FormatHexBytes(original)
			parsed, err := ParseHexBytes(hexStr)

			if err != nil {
				t.Errorf("Round trip failed with error: %v", err)
				return
			}

			if !reflect.DeepEqual(parsed, original) {
				t.Errorf("Round trip failed: original=%v, parsed=%v", original, parsed)
			}
		})
	}
}

// Test round-trip conversion for ParseHexString and FormatHexString
func TestHexStringRoundTrip(t *testing.T) {
	testCases := [][]byte{
		{},
		{0x0a, 0x0b, 0x0c, 0x0d},
		{0xf0, 0x05, 0x5c, 0x08, 0x99, 0x3f, 0xea, 0x1e},
		{0x00, 0x00, 0x00},
		{0xff, 0xff},
	}

	for _, original := range testCases {
		t.Run("round trip", func(t *testing.T) {
			// Convert to string and back
			hexStr := FormatHexString(original)
			parsed, err := ParseHexString(hexStr)

			if err != nil {
				t.Errorf("Round trip failed with error: %v", err)
				return
			}

			if !reflect.DeepEqual(parsed, original) {
				t.Errorf("Round trip failed: original=%v, parsed=%v", original, parsed)
			}
		})
	}
}
