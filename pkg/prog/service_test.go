package prog

import (
	"testing"
	"time"
)

// TestProgramInfoStruct tests that ProgramInfo struct has all required fields.
func TestProgramInfoStruct(t *testing.T) {
	info := ProgramInfo{
		ID:          123,
		Type:        "sched_cls",
		Name:        "test_prog",
		Tag:         "f0055c08993fea1e",
		GPL:         true,
		LoadedAt:    time.Now(),
		UID:         0,
		BytesXlated: 5200,
		BytesJIT:    3263,
		MemLock:     8192,
		MapIDs:      []uint32{85, 39, 38},
	}

	if info.ID != 123 {
		t.Errorf("expected ID 123, got %d", info.ID)
	}
	if info.Type != "sched_cls" {
		t.Errorf("expected Type sched_cls, got %s", info.Type)
	}
	if info.Name != "test_prog" {
		t.Errorf("expected Name test_prog, got %s", info.Name)
	}
	if info.Tag != "f0055c08993fea1e" {
		t.Errorf("expected Tag f0055c08993fea1e, got %s", info.Tag)
	}
	if !info.GPL {
		t.Error("expected GPL to be true")
	}
	if len(info.MapIDs) != 3 {
		t.Errorf("expected 3 MapIDs, got %d", len(info.MapIDs))
	}
}

// TestServiceInterface tests that EBPFService implements Service interface.
func TestServiceInterface(t *testing.T) {
	var _ Service = (*EBPFService)(nil)
	var _ Service = NewService()
}

// TestNewService tests that NewService returns a valid service.
func TestNewService(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Error("NewService returned nil")
	}
}

// MockService is a mock implementation of Service for testing.
type MockService struct {
	programs       []ProgramInfo
	listErr        error
	getByIDErr     error
	getByTagErr    error
	getByNameErr   error
	getByPinnedErr error
}

func (m *MockService) List() ([]ProgramInfo, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.programs, nil
}

func (m *MockService) GetByID(id uint32) (*ProgramInfo, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	for _, p := range m.programs {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, nil
}

func (m *MockService) GetByTag(tag string) ([]ProgramInfo, error) {
	if m.getByTagErr != nil {
		return nil, m.getByTagErr
	}
	var result []ProgramInfo
	for _, p := range m.programs {
		if p.Tag == tag {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockService) GetByName(name string) ([]ProgramInfo, error) {
	if m.getByNameErr != nil {
		return nil, m.getByNameErr
	}
	var result []ProgramInfo
	for _, p := range m.programs {
		if p.Name == name {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockService) GetByPinnedPath(path string) (*ProgramInfo, error) {
	if m.getByPinnedErr != nil {
		return nil, m.getByPinnedErr
	}
	// Mock doesn't support pinned paths
	return nil, nil
}

// TestMockServiceList tests the mock service List method.
func TestMockServiceList(t *testing.T) {
	mock := &MockService{
		programs: []ProgramInfo{
			{ID: 1, Name: "prog1", Type: "xdp"},
			{ID: 2, Name: "prog2", Type: "kprobe"},
		},
	}

	progs, err := mock.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(progs) != 2 {
		t.Errorf("expected 2 programs, got %d", len(progs))
	}
}

// TestMockServiceGetByID tests the mock service GetByID method.
func TestMockServiceGetByID(t *testing.T) {
	mock := &MockService{
		programs: []ProgramInfo{
			{ID: 1, Name: "prog1", Type: "xdp"},
			{ID: 2, Name: "prog2", Type: "kprobe"},
		},
	}

	prog, err := mock.GetByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prog == nil {
		t.Fatal("expected program, got nil")
	}
	if prog.Name != "prog1" {
		t.Errorf("expected prog1, got %s", prog.Name)
	}
}

// TestMockServiceGetByTag tests the mock service GetByTag method.
func TestMockServiceGetByTag(t *testing.T) {
	mock := &MockService{
		programs: []ProgramInfo{
			{ID: 1, Name: "prog1", Tag: "abc123"},
			{ID: 2, Name: "prog2", Tag: "def456"},
			{ID: 3, Name: "prog3", Tag: "abc123"},
		},
	}

	progs, err := mock.GetByTag("abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(progs) != 2 {
		t.Errorf("expected 2 programs with tag abc123, got %d", len(progs))
	}
}

// TestMockServiceGetByName tests the mock service GetByName method.
func TestMockServiceGetByName(t *testing.T) {
	mock := &MockService{
		programs: []ProgramInfo{
			{ID: 1, Name: "my_prog"},
			{ID: 2, Name: "other_prog"},
			{ID: 3, Name: "my_prog"},
		},
	}

	progs, err := mock.GetByName("my_prog")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(progs) != 2 {
		t.Errorf("expected 2 programs named my_prog, got %d", len(progs))
	}
}
