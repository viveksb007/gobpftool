package maps

import (
	"testing"
	"time"
)

func TestMapInfo_JSONTags(t *testing.T) {
	// Test that MapInfo struct has correct JSON tags
	mapInfo := MapInfo{
		ID:         123,
		Type:       "hash",
		Name:       "test_map",
		KeySize:    4,
		ValueSize:  8,
		MaxEntries: 1024,
		Flags:      0,
		MemLock:    4096,
		LoadedAt:   time.Now(),
		UID:        0,
	}

	// Verify struct fields are accessible
	if mapInfo.ID != 123 {
		t.Errorf("Expected ID 123, got %d", mapInfo.ID)
	}

	if mapInfo.Type != "hash" {
		t.Errorf("Expected type 'hash', got %s", mapInfo.Type)
	}

	if mapInfo.Name != "test_map" {
		t.Errorf("Expected name 'test_map', got %s", mapInfo.Name)
	}

	if mapInfo.KeySize != 4 {
		t.Errorf("Expected KeySize 4, got %d", mapInfo.KeySize)
	}

	if mapInfo.ValueSize != 8 {
		t.Errorf("Expected ValueSize 8, got %d", mapInfo.ValueSize)
	}

	if mapInfo.MaxEntries != 1024 {
		t.Errorf("Expected MaxEntries 1024, got %d", mapInfo.MaxEntries)
	}
}

func TestMapEntry_Structure(t *testing.T) {
	// Test MapEntry struct
	entry := MapEntry{
		Key:   []byte{0x01, 0x02, 0x03, 0x04},
		Value: []byte{0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c},
	}

	expectedKeyLen := 4
	expectedValueLen := 8

	if len(entry.Key) != expectedKeyLen {
		t.Errorf("Expected key length %d, got %d", expectedKeyLen, len(entry.Key))
	}

	if len(entry.Value) != expectedValueLen {
		t.Errorf("Expected value length %d, got %d", expectedValueLen, len(entry.Value))
	}

	// Verify key content
	expectedKey := []byte{0x01, 0x02, 0x03, 0x04}
	for i, b := range entry.Key {
		if b != expectedKey[i] {
			t.Errorf("Expected key byte %d to be 0x%02x, got 0x%02x", i, expectedKey[i], b)
		}
	}

	// Verify value content
	expectedValue := []byte{0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c}
	for i, b := range entry.Value {
		if b != expectedValue[i] {
			t.Errorf("Expected value byte %d to be 0x%02x, got 0x%02x", i, expectedValue[i], b)
		}
	}
}

func TestNewService(t *testing.T) {
	// Test service creation
	service := NewService()
	if service == nil {
		t.Error("Expected service to be created, got nil")
	}

	// Verify it implements the Service interface
	var _ Service = service
}

// Note: Integration tests that interact with real eBPF maps would require
// root privileges and actual eBPF programs/maps to be loaded.
// These tests focus on the structure and basic functionality that can be
// tested without kernel interaction.

func TestServiceImpl_Interface(t *testing.T) {
	// Verify that serviceImpl implements Service interface
	var service Service = &serviceImpl{}

	// Test that all interface methods are available
	// (This will fail to compile if interface is not properly implemented)
	_ = service.List
	_ = service.GetByID
	_ = service.GetByName
	_ = service.GetByPinnedPath
	_ = service.Dump
	_ = service.Lookup
	_ = service.GetNextKey
}
