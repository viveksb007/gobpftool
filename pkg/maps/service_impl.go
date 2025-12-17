package maps

import (
	"fmt"
	"strings"

	"github.com/cilium/ebpf"
)

// serviceImpl implements the Service interface using cilium/ebpf
type serviceImpl struct{}

// NewService creates a new map service instance
func NewService() Service {
	return &serviceImpl{}
}

// List returns all loaded eBPF maps
func (s *serviceImpl) List() ([]MapInfo, error) {
	var maps []MapInfo

	var id ebpf.MapID
	firstIteration := true

	for {
		nextID, err := ebpf.MapGetNextID(id)
		if err != nil {
			// If this is the first iteration and we get an error, it's likely a permission issue
			if firstIteration {
				return nil, fmt.Errorf("failed to list maps: %w", err)
			}
			// Otherwise, no more maps
			break
		}
		firstIteration = false
		id = nextID

		m, err := ebpf.NewMapFromID(id)
		if err != nil {
			// Skip maps we can't access
			continue
		}

		mapInfo, err := s.mapToMapInfo(m)
		m.Close()
		if err != nil {
			continue
		}

		maps = append(maps, *mapInfo)
	}

	return maps, nil
}

// GetByID returns map info by ID
func (s *serviceImpl) GetByID(id uint32) (*MapInfo, error) {
	m, err := ebpf.NewMapFromID(ebpf.MapID(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get map by ID %d: %w", id, err)
	}
	defer m.Close()

	return s.mapToMapInfo(m)
}

// GetByName returns maps matching the name
func (s *serviceImpl) GetByName(name string) ([]MapInfo, error) {
	allMaps, err := s.List()
	if err != nil {
		return nil, err
	}

	var matchingMaps []MapInfo
	for _, mapInfo := range allMaps {
		if mapInfo.Name == name {
			matchingMaps = append(matchingMaps, mapInfo)
		}
	}

	return matchingMaps, nil
}

// GetByPinnedPath returns map at the pinned path
func (s *serviceImpl) GetByPinnedPath(path string) (*MapInfo, error) {
	m, err := ebpf.LoadPinnedMap(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load pinned map at %s: %w", path, err)
	}
	defer m.Close()

	return s.mapToMapInfo(m)
}

// Dump returns all entries in the map
func (s *serviceImpl) Dump(id uint32) ([]MapEntry, error) {
	m, err := ebpf.NewMapFromID(ebpf.MapID(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get map by ID %d: %w", id, err)
	}
	defer m.Close()

	var entries []MapEntry

	// Get map info to determine key and value sizes
	info, err := m.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get map info: %w", err)
	}

	keySize := info.KeySize
	valueSize := info.ValueSize

	// Create buffers for keys and values
	key := make([]byte, keySize)
	value := make([]byte, valueSize)

	// Iterate through all entries
	iter := m.Iterate()
	for iter.Next(&key, &value) {
		// Make copies of the key and value since they're reused
		keyCopy := make([]byte, len(key))
		valueCopy := make([]byte, len(value))
		copy(keyCopy, key)
		copy(valueCopy, value)

		entries = append(entries, MapEntry{
			Key:   keyCopy,
			Value: valueCopy,
		})
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate map entries: %w", err)
	}

	return entries, nil
}

// Lookup returns the value for a key in the map
func (s *serviceImpl) Lookup(id uint32, key []byte) ([]byte, error) {
	m, err := ebpf.NewMapFromID(ebpf.MapID(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get map by ID %d: %w", id, err)
	}
	defer m.Close()

	// Get map info to determine value size
	info, err := m.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get map info: %w", err)
	}

	// Create buffer for value
	value := make([]byte, info.ValueSize)

	// Lookup the key
	err = m.Lookup(key, &value)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup key: %w", err)
	}

	return value, nil
}

// GetNextKey returns the next key after the given key
// If key is nil, returns the first key
func (s *serviceImpl) GetNextKey(id uint32, key []byte) ([]byte, error) {
	m, err := ebpf.NewMapFromID(ebpf.MapID(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get map by ID %d: %w", id, err)
	}
	defer m.Close()

	// Get map info to determine key size
	info, err := m.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get map info: %w", err)
	}

	// Create buffer for next key
	nextKey := make([]byte, info.KeySize)

	// Get next key
	err = m.NextKey(key, &nextKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get next key: %w", err)
	}

	return nextKey, nil
}

// mapToMapInfo converts an ebpf.Map to MapInfo
func (s *serviceImpl) mapToMapInfo(m *ebpf.Map) (*MapInfo, error) {
	info, err := m.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get map info: %w", err)
	}

	// Convert map type to string
	mapType := strings.ToLower(info.Type.String())

	// Get the map ID - info.ID() returns (MapID, bool)
	mapID, _ := info.ID()

	mapInfo := &MapInfo{
		ID:         uint32(mapID),
		Type:       mapType,
		Name:       info.Name,
		KeySize:    info.KeySize,
		ValueSize:  info.ValueSize,
		MaxEntries: info.MaxEntries,
		Flags:      uint32(info.Flags),
	}

	return mapInfo, nil
}
