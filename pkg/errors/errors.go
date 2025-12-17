// Package errors provides error types and utilities for gobpftool.
package errors

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
)

// Sentinel errors for common error conditions.
var (
	// ErrNotFound indicates a resource was not found.
	ErrNotFound = errors.New("not found")

	// ErrPermission indicates insufficient permissions.
	ErrPermission = errors.New("permission denied: requires CAP_SYS_ADMIN or CAP_BPF")

	// ErrBpfFSNotMounted indicates the BPF filesystem is not mounted.
	ErrBpfFSNotMounted = errors.New("BPF filesystem not mounted at /sys/fs/bpf")

	// ErrInvalidID indicates an invalid ID was provided.
	ErrInvalidID = errors.New("invalid ID")

	// ErrInvalidKey indicates an invalid key format.
	ErrInvalidKey = errors.New("invalid key format")

	// ErrKeyNotFound indicates a key was not found in a map.
	ErrKeyNotFound = errors.New("key not found in map")

	// ErrNoMoreKeys indicates there are no more keys in a map.
	ErrNoMoreKeys = errors.New("no more keys")

	// ErrMapEmpty indicates the map is empty.
	ErrMapEmpty = errors.New("map is empty")
)

// IsPermissionError checks if the error is a permission-related error.
func IsPermissionError(err error) bool {
	if err == nil {
		return false
	}

	// Check for our sentinel error
	if errors.Is(err, ErrPermission) {
		return true
	}

	// Check for syscall permission errors
	if errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EACCES) {
		return true
	}

	// Check for os permission errors
	if os.IsPermission(err) {
		return true
	}

	// Check error message for permission-related strings
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "permission denied") ||
		strings.Contains(errStr, "operation not permitted")
}

// IsNotFoundError checks if the error indicates a resource was not found.
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	// Check for our sentinel errors
	if errors.Is(err, ErrNotFound) || errors.Is(err, ErrKeyNotFound) {
		return true
	}

	// Check for syscall not found errors
	if errors.Is(err, syscall.ENOENT) {
		return true
	}

	// Check for os not exist errors
	if os.IsNotExist(err) {
		return true
	}

	// Check error message
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "no such file or directory")
}

// IsNoMoreKeysError checks if the error indicates no more keys in iteration.
func IsNoMoreKeysError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrNoMoreKeys) {
		return true
	}

	// ENOENT is returned when there are no more keys
	if errors.Is(err, syscall.ENOENT) {
		return true
	}

	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "no more keys") ||
		strings.Contains(errStr, "no such file or directory")
}

// IsBpfFSNotMounted checks if the BPF filesystem is mounted.
func IsBpfFSNotMounted() bool {
	_, err := os.Stat("/sys/fs/bpf")
	return os.IsNotExist(err)
}

// WrapError wraps an error with additional context and converts
// common system errors to our sentinel errors.
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}

	// Convert permission errors
	if IsPermissionError(err) {
		return fmt.Errorf("%s: %w", context, ErrPermission)
	}

	// Check for BPF filesystem issues
	if IsBpfFSNotMounted() && strings.Contains(context, "pinned") {
		return fmt.Errorf("%s: %w", context, ErrBpfFSNotMounted)
	}

	// Convert not found errors
	if IsNotFoundError(err) {
		return fmt.Errorf("%s: %w", context, ErrNotFound)
	}

	// Default wrapping
	return fmt.Errorf("%s: %w", context, err)
}

// FormatPermissionError returns a user-friendly permission error message.
func FormatPermissionError() string {
	return `Error: Permission denied.

This operation requires elevated privileges. You need one of the following:
  - Run as root (sudo gobpftool ...)
  - Have CAP_SYS_ADMIN capability
  - Have CAP_BPF capability (Linux 5.8+)

To grant CAP_BPF capability to the binary:
  sudo setcap cap_bpf=ep /path/to/gobpftool`
}

// FormatBpfFSError returns a user-friendly BPF filesystem error message.
func FormatBpfFSError() string {
	return `Error: BPF filesystem not mounted at /sys/fs/bpf.

To mount the BPF filesystem, run:
  sudo mount -t bpf bpf /sys/fs/bpf

To mount it permanently, add to /etc/fstab:
  bpf /sys/fs/bpf bpf defaults 0 0`
}

// FormatError returns a user-friendly error message for the given error.
func FormatError(err error) string {
	if err == nil {
		return ""
	}

	if errors.Is(err, ErrPermission) || IsPermissionError(err) {
		return FormatPermissionError()
	}

	if errors.Is(err, ErrBpfFSNotMounted) || IsBpfFSNotMounted() {
		return FormatBpfFSError()
	}

	if errors.Is(err, ErrKeyNotFound) {
		return "Error: key not found in map"
	}

	if errors.Is(err, ErrNoMoreKeys) {
		return "Error: no more keys"
	}

	if errors.Is(err, ErrMapEmpty) {
		return "Error: map is empty"
	}

	if errors.Is(err, ErrNotFound) {
		return fmt.Sprintf("Error: %v", err)
	}

	return fmt.Sprintf("Error: %v", err)
}

// ExitCode returns the appropriate exit code for the given error.
// Returns 0 for nil (success), 1 for any error (failure).
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	return 1
}
