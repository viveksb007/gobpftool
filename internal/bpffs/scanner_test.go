package bpffs

import "testing"

func TestGetScanner(t *testing.T) {
	s := GetScanner()
	if s == nil {
		t.Fatal("GetScanner returned nil")
	}

	// Should return same instance
	s2 := GetScanner()
	if s != s2 {
		t.Error("GetScanner should return singleton")
	}
}

func TestGetPinnedPaths_ReturnsSliceCopy(t *testing.T) {
	s := &Scanner{
		progPaths: map[uint32][]string{1: {"/sys/fs/bpf/test"}},
		mapPaths:  map[uint32][]string{2: {"/sys/fs/bpf/map"}},
		scanned:   true,
	}

	// Get paths and modify returned slice
	progPaths := s.GetProgramPinnedPaths(1)
	if len(progPaths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(progPaths))
	}
	progPaths[0] = "modified"

	// Original should be unchanged
	if s.progPaths[1][0] != "/sys/fs/bpf/test" {
		t.Error("internal slice was modified")
	}

	// Same for maps
	mapPaths := s.GetMapPinnedPaths(2)
	mapPaths[0] = "modified"
	if s.mapPaths[2][0] != "/sys/fs/bpf/map" {
		t.Error("internal map slice was modified")
	}
}

func TestGetPinnedPaths_NonExistentID(t *testing.T) {
	s := &Scanner{
		progPaths: make(map[uint32][]string),
		mapPaths:  make(map[uint32][]string),
		scanned:   true,
	}

	paths := s.GetProgramPinnedPaths(999)
	if len(paths) != 0 {
		t.Errorf("expected 0 paths, got %d", len(paths))
	}
}

func TestRefresh(t *testing.T) {
	s := &Scanner{
		progPaths: map[uint32][]string{1: {"/old/path"}},
		mapPaths:  make(map[uint32][]string),
		bpffsRoot: "/nonexistent/path",
		scanned:   true,
	}

	s.Refresh()

	// After refresh with nonexistent path, maps should be cleared
	if len(s.progPaths) != 0 {
		t.Error("expected progPaths to be cleared after refresh")
	}
}
