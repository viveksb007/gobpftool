// Package prog provides services for inspecting eBPF programs.
package prog

import "time"

// ProgramInfo contains information about a loaded eBPF program.
type ProgramInfo struct {
	// ID is the unique identifier of the program.
	ID uint32
	// Type is the program type (e.g., "sched_cls", "xdp", "kprobe").
	Type string
	// Name is the program name.
	Name string
	// Tag is the 8-byte program tag as a hex string.
	Tag string
	// GPL indicates if the program is GPL compatible.
	GPL bool
	// LoadedAt is the time when the program was loaded.
	LoadedAt time.Time
	// UID is the user ID that loaded the program.
	UID uint32
	// BytesXlated is the number of bytes in the translated eBPF bytecode.
	BytesXlated uint32
	// BytesJIT is the number of bytes in the JIT-compiled code.
	BytesJIT uint32
	// MemLock is the amount of memory locked for the program.
	MemLock uint32
	// MapIDs is the list of map IDs associated with this program.
	MapIDs []uint32
	// PinnedPaths contains the paths where this program is pinned in bpffs.
	PinnedPaths []string `json:"pinned_paths,omitempty"`
}

// Service defines the interface for inspecting eBPF programs.
type Service interface {
	// List returns all loaded eBPF programs.
	List() ([]ProgramInfo, error)

	// GetByID returns program info by ID.
	GetByID(id uint32) (*ProgramInfo, error)

	// GetByTag returns programs matching the tag.
	GetByTag(tag string) ([]ProgramInfo, error)

	// GetByName returns programs matching the name.
	GetByName(name string) ([]ProgramInfo, error)

	// GetByPinnedPath returns program at the pinned path.
	GetByPinnedPath(path string) (*ProgramInfo, error)
}
