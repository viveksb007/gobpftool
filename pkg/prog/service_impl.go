package prog

import (
	"fmt"
	"os"
	"time"

	"github.com/cilium/ebpf"
)

// EBPFService implements the Service interface using cilium/ebpf.
type EBPFService struct{}

// NewService creates a new program service.
func NewService() Service {
	return &EBPFService{}
}

// List returns all loaded eBPF programs.
func (s *EBPFService) List() ([]ProgramInfo, error) {
	var programs []ProgramInfo

	var id ebpf.ProgramID
	for {
		nextID, err := ebpf.ProgramGetNextID(id)
		if err != nil {
			// No more programs
			break
		}
		id = nextID

		prog, err := ebpf.NewProgramFromID(id)
		if err != nil {
			// Skip programs we can't access
			continue
		}

		info, err := extractProgramInfo(prog)
		prog.Close()
		if err != nil {
			continue
		}

		programs = append(programs, *info)
	}

	return programs, nil
}

// GetByID returns program info by ID.
func (s *EBPFService) GetByID(id uint32) (*ProgramInfo, error) {
	prog, err := ebpf.NewProgramFromID(ebpf.ProgramID(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("program with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get program %d: %w", id, err)
	}
	defer prog.Close()

	return extractProgramInfo(prog)
}

// GetByTag returns programs matching the tag.
func (s *EBPFService) GetByTag(tag string) ([]ProgramInfo, error) {
	allProgs, err := s.List()
	if err != nil {
		return nil, err
	}

	var matched []ProgramInfo
	for _, p := range allProgs {
		if p.Tag == tag {
			matched = append(matched, p)
		}
	}

	return matched, nil
}

// GetByName returns programs matching the name.
func (s *EBPFService) GetByName(name string) ([]ProgramInfo, error) {
	allProgs, err := s.List()
	if err != nil {
		return nil, err
	}

	var matched []ProgramInfo
	for _, p := range allProgs {
		if p.Name == name {
			matched = append(matched, p)
		}
	}

	return matched, nil
}

// GetByPinnedPath returns program at the pinned path.
func (s *EBPFService) GetByPinnedPath(path string) (*ProgramInfo, error) {
	prog, err := ebpf.LoadPinnedProgram(path, nil)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no program pinned at %s", path)
		}
		return nil, fmt.Errorf("failed to load pinned program at %s: %w", path, err)
	}
	defer prog.Close()

	return extractProgramInfo(prog)
}

// extractProgramInfo extracts ProgramInfo from an ebpf.Program.
func extractProgramInfo(prog *ebpf.Program) (*ProgramInfo, error) {
	info, err := prog.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get program info: %w", err)
	}

	id, ok := info.ID()
	if !ok {
		return nil, fmt.Errorf("failed to get program ID")
	}

	tag := info.Tag

	// Get map IDs associated with this program
	mapIDs, _ := info.MapIDs()

	// Convert []ebpf.MapID to []uint32
	mapIDsUint32 := make([]uint32, len(mapIDs))
	for i, mid := range mapIDs {
		mapIDsUint32[i] = uint32(mid)
	}

	// Get loaded time - LoadTime returns a duration since boot
	var loadedAt time.Time
	if loadTime, ok := info.LoadTime(); ok {
		// Convert duration since boot to actual time
		loadedAt = time.Now().Add(-loadTime)
	}

	return &ProgramInfo{
		ID:          uint32(id),
		Type:        info.Type.String(),
		Name:        info.Name,
		Tag:         tag,
		GPL:         false, // GPL info not directly exposed in this version
		LoadedAt:    loadedAt,
		UID:         0, // UID is not directly exposed by cilium/ebpf
		BytesXlated: 0, // Not directly exposed in this API version
		BytesJIT:    0, // Not directly exposed in this API version
		MemLock:     0, // Not directly exposed in this API version
		MapIDs:      mapIDsUint32,
	}, nil
}
