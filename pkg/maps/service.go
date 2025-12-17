package maps

import (
	"time"
)

// MapInfo represents information about an eBPF map
type MapInfo struct {
	ID         uint32    `json:"id"`
	Type       string    `json:"type"`
	Name       string    `json:"name"`
	KeySize    uint32    `json:"key_size"`
	ValueSize  uint32    `json:"value_size"`
	MaxEntries uint32    `json:"max_entries"`
	Flags      uint32    `json:"flags"`
	MemLock    uint32    `json:"bytes_memlock"`
	LoadedAt   time.Time `json:"loaded_at,omitempty"`
	UID        uint32    `json:"uid,omitempty"`
}

// MapEntry represents a key-value pair in an eBPF map
type MapEntry struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

// Service provides operations for inspecting eBPF maps
type Service interface {
	// List returns all loaded eBPF maps
	List() ([]MapInfo, error)

	// GetByID returns map info by ID
	GetByID(id uint32) (*MapInfo, error)

	// GetByName returns maps matching the name
	GetByName(name string) ([]MapInfo, error)

	// GetByPinnedPath returns map at the pinned path
	GetByPinnedPath(path string) (*MapInfo, error)

	// Dump returns all entries in the map
	Dump(id uint32) ([]MapEntry, error)

	// Lookup returns the value for a key in the map
	Lookup(id uint32, key []byte) ([]byte, error)

	// GetNextKey returns the next key after the given key
	// If key is nil, returns the first key
	GetNextKey(id uint32, key []byte) ([]byte, error)
}
