package utils

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// ParseHexBytes parses space-separated hex bytes into a byte slice.
// Input format: "0a 0b 0c 0d" or "0A 0B 0C 0D"
// Returns the parsed bytes or an error if parsing fails.
func ParseHexBytes(hexStr string) ([]byte, error) {
	if hexStr == "" {
		return []byte{}, nil
	}

	// Split by whitespace and filter out empty strings
	parts := strings.Fields(hexStr)
	if len(parts) == 0 {
		return []byte{}, nil
	}

	result := make([]byte, len(parts))
	for i, part := range parts {
		// Parse each hex byte (should be 1-2 characters)
		if len(part) > 2 {
			return nil, fmt.Errorf("invalid hex byte '%s': too long", part)
		}

		val, err := strconv.ParseUint(part, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid hex byte '%s': %w", part, err)
		}

		result[i] = byte(val)
	}

	return result, nil
}

// FormatHexBytes formats a byte slice as space-separated hex bytes.
// Output format: "0a 0b 0c 0d"
func FormatHexBytes(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	parts := make([]string, len(data))
	for i, b := range data {
		parts[i] = fmt.Sprintf("%02x", b)
	}

	return strings.Join(parts, " ")
}

// FormatHexBytesWithPrefix formats a byte slice as space-separated hex bytes with a prefix.
// Output format: "key: 0a 0b 0c 0d" or "value: 01 02 03 04"
func FormatHexBytesWithPrefix(prefix string, data []byte) string {
	if len(data) == 0 {
		return prefix + ":"
	}
	return prefix + ": " + FormatHexBytes(data)
}

// ParseHexString parses a continuous hex string (without spaces) into a byte slice.
// Input format: "0a0b0c0d"
// This is useful for parsing program tags and other continuous hex values.
func ParseHexString(hexStr string) ([]byte, error) {
	if hexStr == "" {
		return []byte{}, nil
	}

	// Remove any whitespace
	cleaned := strings.ReplaceAll(hexStr, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "\t", "")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")

	// Check if length is even
	if len(cleaned)%2 != 0 {
		return nil, fmt.Errorf("hex string must have even length, got %d", len(cleaned))
	}

	return hex.DecodeString(cleaned)
}

// FormatHexString formats a byte slice as a continuous hex string.
// Output format: "0a0b0c0d"
func FormatHexString(data []byte) string {
	return hex.EncodeToString(data)
}
