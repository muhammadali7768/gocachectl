package cachemgr

import (
	"testing"

	"github.com/muhammadali7768/gocachectl/internal/cache"
)

// MockCacheManager implements cache.CacheManager for testing
type MockCacheManager struct {
	stats    cache.Stats
	deleted  int
	freed    int64
	location string
	err      error
}

func (m *MockCacheManager) GetStats() (cache.Stats, error) {
	return m.stats, m.err
}

func (m *MockCacheManager) Clear() (int, int64, error) {
	return m.deleted, m.freed, m.err
}

func (m *MockCacheManager) GetLocation() string {
	return m.location
}

// MockStats implements cache.Stats
type MockStats struct {
	typeStr string
}

func (s MockStats) Type() string {
	return s.typeStr
}

func TestUnifiedManager_GetAllStats(t *testing.T) {
	mockBuild := &MockCacheManager{
		stats: MockStats{typeStr: "build"},
	}
	mockTest := &MockCacheManager{
		stats: MockStats{typeStr: "test"},
	}

	mgr := &UnifiedManager{
		managers: []cache.CacheManager{mockBuild, mockTest},
	}

	stats, err := mgr.GetAllStats()
	if err != nil {
		t.Fatalf("GetAllStats failed: %v", err)
	}

	if len(stats) != 2 {
		t.Errorf("Expected 2 stats, got %d", len(stats))
	}
}

func TestUnifiedManager_Clear(t *testing.T) {
	mockBuild := &MockCacheManager{
		stats:   MockStats{typeStr: "build"},
		deleted: 10,
		freed:   1000,
	}
	mockTest := &MockCacheManager{
		stats:   MockStats{typeStr: "test"},
		deleted: 5,
		freed:   500,
	}

	mgr := &UnifiedManager{
		managers: []cache.CacheManager{mockBuild, mockTest},
	}

	// Test clearing only build
	opts := cache.ClearOptions{
		Build: true,
	}

	result, err := mgr.Clear(opts)
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	if result.BuildDeleted != 10 {
		t.Errorf("Expected 10 build deleted, got %d", result.BuildDeleted)
	}
	if result.TestDeleted != 0 {
		t.Errorf("Expected 0 test deleted, got %d", result.TestDeleted)
	}
	if result.TotalFreed != 1000 {
		t.Errorf("Expected 1000 freed, got %d", result.TotalFreed)
	}

	// Test clearing all
	opts = cache.ClearOptions{
		All: true,
	}

	result, err = mgr.Clear(opts)
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	if result.BuildDeleted != 10 {
		t.Errorf("Expected 10 build deleted, got %d", result.BuildDeleted)
	}
	if result.TestDeleted != 5 {
		t.Errorf("Expected 5 test deleted, got %d", result.TestDeleted)
	}
	if result.TotalFreed != 1500 {
		t.Errorf("Expected 1500 freed, got %d", result.TotalFreed)
	}
}
