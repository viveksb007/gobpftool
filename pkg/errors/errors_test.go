package errors

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
)

func TestIsPermissionError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "ErrPermission sentinel",
			err:      ErrPermission,
			expected: true,
		},
		{
			name:     "wrapped ErrPermission",
			err:      fmt.Errorf("failed: %w", ErrPermission),
			expected: true,
		},
		{
			name:     "syscall EPERM",
			err:      syscall.EPERM,
			expected: true,
		},
		{
			name:     "syscall EACCES",
			err:      syscall.EACCES,
			expected: true,
		},
		{
			name:     "wrapped syscall EPERM",
			err:      fmt.Errorf("operation failed: %w", syscall.EPERM),
			expected: true,
		},
		{
			name:     "permission denied in message",
			err:      errors.New("permission denied"),
			expected: true,
		},
		{
			name:     "operation not permitted in message",
			err:      errors.New("operation not permitted"),
			expected: true,
		},
		{
			name:     "unrelated error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "ErrNotFound",
			err:      ErrNotFound,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPermissionError(tt.err)
			if result != tt.expected {
				t.Errorf("IsPermissionError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "ErrNotFound sentinel",
			err:      ErrNotFound,
			expected: true,
		},
		{
			name:     "ErrKeyNotFound sentinel",
			err:      ErrKeyNotFound,
			expected: true,
		},
		{
			name:     "wrapped ErrNotFound",
			err:      fmt.Errorf("failed: %w", ErrNotFound),
			expected: true,
		},
		{
			name:     "syscall ENOENT",
			err:      syscall.ENOENT,
			expected: true,
		},
		{
			name:     "not found in message",
			err:      errors.New("resource not found"),
			expected: true,
		},
		{
			name:     "no such file or directory",
			err:      errors.New("no such file or directory"),
			expected: true,
		},
		{
			name:     "unrelated error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "os.ErrNotExist",
			err:      os.ErrNotExist,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFoundError(tt.err)
			if result != tt.expected {
				t.Errorf("IsNotFoundError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsNoMoreKeysError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "ErrNoMoreKeys sentinel",
			err:      ErrNoMoreKeys,
			expected: true,
		},
		{
			name:     "wrapped ErrNoMoreKeys",
			err:      fmt.Errorf("iteration: %w", ErrNoMoreKeys),
			expected: true,
		},
		{
			name:     "syscall ENOENT",
			err:      syscall.ENOENT,
			expected: true,
		},
		{
			name:     "no more keys in message",
			err:      errors.New("no more keys"),
			expected: true,
		},
		{
			name:     "unrelated error",
			err:      errors.New("some other error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNoMoreKeysError(tt.err)
			if result != tt.expected {
				t.Errorf("IsNoMoreKeysError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		context        string
		expectNil      bool
		expectContains string
		expectSentinel error
	}{
		{
			name:      "nil error",
			err:       nil,
			context:   "test",
			expectNil: true,
		},
		{
			name:           "permission error",
			err:            syscall.EPERM,
			context:        "listing programs",
			expectContains: "listing programs",
			expectSentinel: ErrPermission,
		},
		{
			name:           "not found error",
			err:            syscall.ENOENT,
			context:        "getting program",
			expectContains: "getting program",
			expectSentinel: ErrNotFound,
		},
		{
			name:           "generic error",
			err:            errors.New("something went wrong"),
			context:        "operation",
			expectContains: "operation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.context)

			if tt.expectNil {
				if result != nil {
					t.Errorf("WrapError() = %v, want nil", result)
				}
				return
			}

			if result == nil {
				t.Error("WrapError() = nil, want non-nil")
				return
			}

			if !strings.Contains(result.Error(), tt.expectContains) {
				t.Errorf("WrapError() = %v, want to contain %q", result, tt.expectContains)
			}

			if tt.expectSentinel != nil && !errors.Is(result, tt.expectSentinel) {
				t.Errorf("WrapError() should wrap %v", tt.expectSentinel)
			}
		})
	}
}

func TestFormatPermissionError(t *testing.T) {
	result := FormatPermissionError()

	expectedPhrases := []string{
		"Permission denied",
		"CAP_SYS_ADMIN",
		"CAP_BPF",
		"sudo",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(result, phrase) {
			t.Errorf("FormatPermissionError() should contain %q", phrase)
		}
	}
}

func TestFormatBpfFSError(t *testing.T) {
	result := FormatBpfFSError()

	expectedPhrases := []string{
		"BPF filesystem",
		"/sys/fs/bpf",
		"mount",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(result, phrase) {
			t.Errorf("FormatBpfFSError() should contain %q", phrase)
		}
	}
}

func TestFormatError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectContains string
	}{
		{
			name:           "nil error",
			err:            nil,
			expectContains: "",
		},
		{
			name:           "permission error",
			err:            ErrPermission,
			expectContains: "Permission denied",
		},
		{
			name:           "key not found",
			err:            ErrKeyNotFound,
			expectContains: "key not found",
		},
		{
			name:           "no more keys",
			err:            ErrNoMoreKeys,
			expectContains: "no more keys",
		},
		{
			name:           "map empty",
			err:            ErrMapEmpty,
			expectContains: "map is empty",
		},
		{
			name:           "generic error",
			err:            errors.New("something failed"),
			expectContains: "something failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatError(tt.err)

			if tt.expectContains == "" {
				if result != "" {
					t.Errorf("FormatError(%v) = %q, want empty string", tt.err, result)
				}
				return
			}

			if !strings.Contains(result, tt.expectContains) {
				t.Errorf("FormatError(%v) = %q, want to contain %q", tt.err, result, tt.expectContains)
			}
		})
	}
}

func TestExitCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: 0,
		},
		{
			name:     "permission error",
			err:      ErrPermission,
			expected: 1,
		},
		{
			name:     "not found error",
			err:      ErrNotFound,
			expected: 1,
		},
		{
			name:     "generic error",
			err:      errors.New("something failed"),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExitCode(tt.err)
			if result != tt.expected {
				t.Errorf("ExitCode(%v) = %d, want %d", tt.err, result, tt.expected)
			}
		})
	}
}

func TestSentinelErrors(t *testing.T) {
	// Test that sentinel errors have expected messages
	tests := []struct {
		err      error
		contains string
	}{
		{ErrNotFound, "not found"},
		{ErrPermission, "permission denied"},
		{ErrBpfFSNotMounted, "BPF filesystem"},
		{ErrInvalidID, "invalid ID"},
		{ErrInvalidKey, "invalid key"},
		{ErrKeyNotFound, "key not found"},
		{ErrNoMoreKeys, "no more keys"},
		{ErrMapEmpty, "empty"},
	}

	for _, tt := range tests {
		t.Run(tt.err.Error(), func(t *testing.T) {
			if !strings.Contains(strings.ToLower(tt.err.Error()), strings.ToLower(tt.contains)) {
				t.Errorf("Error %q should contain %q", tt.err.Error(), tt.contains)
			}
		})
	}
}
