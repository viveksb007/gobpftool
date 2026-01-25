// Package bpffs provides utilities for scanning the BPF filesystem.
package bpffs

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/cilium/ebpf"
)

const defaultBPFFS = "/sys/fs/bpf"

// Scanner discovers pinned BPF objects by scanning the BPF filesystem.
type Scanner struct {
	mu        sync.RWMutex
	progPaths map[uint32][]string // program ID -> pinned paths
	mapPaths  map[uint32][]string // map ID -> pinned paths
	bpffsRoot string
	scanned   bool
}

// Global scanner instance
var (
	globalScanner *Scanner
	scannerOnce   sync.Once
)

// GetScanner returns the global scanner instance, creating it if necessary.
func GetScanner() *Scanner {
	scannerOnce.Do(func() {
		globalScanner = &Scanner{
			progPaths: make(map[uint32][]string),
			mapPaths:  make(map[uint32][]string),
			bpffsRoot: defaultBPFFS,
		}
	})
	return globalScanner
}

// GetProgramPinnedPaths returns all pinned paths for a program ID.
func (s *Scanner) GetProgramPinnedPaths(id uint32) []string {
	s.ensureScanned()
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]string(nil), s.progPaths[id]...)
}

// GetMapPinnedPaths returns all pinned paths for a map ID.
func (s *Scanner) GetMapPinnedPaths(id uint32) []string {
	s.ensureScanned()
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]string(nil), s.mapPaths[id]...)
}

// Refresh forces a rescan of the BPF filesystem, updating the cache.
func (s *Scanner) Refresh() {
	s.mu.Lock()
	s.scanned = false
	s.mu.Unlock()
	s.ensureScanned()
}

// ensureScanned performs the scan if not already done.
func (s *Scanner) ensureScanned() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.scanned {
		return
	}

	// Clear existing data
	s.progPaths = make(map[uint32][]string)
	s.mapPaths = make(map[uint32][]string)
	s.scanned = true

	// Check if bpffs is mounted
	if _, err := os.Stat(s.bpffsRoot); os.IsNotExist(err) {
		return // bpffs not mounted, nothing to scan
	}

	// Walk the BPF filesystem
	_ = filepath.Walk(s.bpffsRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Try to open as a program first
		if prog, err := ebpf.LoadPinnedProgram(path, nil); err == nil {
			progInfo, err := prog.Info()
			prog.Close()
			if err == nil {
				if id, ok := progInfo.ID(); ok {
					s.progPaths[uint32(id)] = append(s.progPaths[uint32(id)], path)
				}
			}
			return nil
		}

		// Try to open as a map
		if m, err := ebpf.LoadPinnedMap(path, nil); err == nil {
			mapInfo, err := m.Info()
			m.Close()
			if err == nil {
				if id, ok := mapInfo.ID(); ok {
					s.mapPaths[uint32(id)] = append(s.mapPaths[uint32(id)], path)
				}
			}
		}

		return nil
	})
}
