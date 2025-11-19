package cache

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manager manages the Go test cache (part of build cache)
type TestManager struct {
	cacheDir string
}

var _ CacheManager = (*TestManager)(nil)

// NewManager creates a new test cache manager
func NewTestManager(cacheDir string) (*TestManager, error) {
	if cacheDir == "" {
		dir, err := GetGoEnv("GOCACHE")
		if err != nil {
			return nil, fmt.Errorf("failed to get GOCACHE: %w", err)
		}
		cacheDir = dir
	}

	// Verify cache directory exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("test cache directory does not exist: %s", cacheDir)
	}

	return &TestManager{
		cacheDir: cacheDir,
	}, nil
}

// GetStats retrieves test cache statistics
func (m *TestManager) GetStats() (Stats, error) {
	stats := &TestCacheStats{
		Location:    m.cacheDir,
		OldestEntry: time.Now(),
	}

	// Walk the cache directory looking for test-related entries
	err := filepath.WalkDir(m.cacheDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check if this is a test-related entry
		// Test cache entries typically have specific patterns
		if !isTestEntry(path) {
			return nil
		}

		// Get file info
		info, err := d.Info()
		if err != nil {
			return nil
		}

		// Update statistics
		stats.EntryCount++
		stats.Size += info.Size()

		// Track oldest and newest
		modTime := info.ModTime()
		if modTime.Before(stats.OldestEntry) {
			stats.OldestEntry = modTime
		}
		if modTime.After(stats.NewestEntry) {
			stats.NewestEntry = modTime
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk test cache: %w", err)
	}

	// If no entries, reset oldest to zero time
	if stats.EntryCount == 0 {
		stats.OldestEntry = time.Time{}
	}

	return stats, nil
}

// Clear removes test cache entries
func (m *TestManager) Clear() (int, int64, error) {
	var deletedCount int
	var freedSpace int64

	err := filepath.WalkDir(m.cacheDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Only delete test entries
		if !isTestEntry(path) {
			return nil
		}

		// Get file size before deletion
		info, err := d.Info()
		if err == nil {
			freedSpace += info.Size()
		}

		// Delete file
		if err := os.Remove(path); err == nil {
			deletedCount++
		}

		return nil
	})

	if err != nil {
		return deletedCount, freedSpace, fmt.Errorf("failed to clear test cache: %w", err)
	}

	return deletedCount, freedSpace, nil
}

// GetLocation returns the cache directory path
func (m *TestManager) GetLocation() string {
	return m.cacheDir
}

// isTestEntry attempts to determine if a cache entry is test-related
// This is heuristic-based since Go's cache format is internal
func isTestEntry(path string) bool {
	base := filepath.Base(path)

	// Test cache entries often have specific patterns
	// They typically end with "-a" or "-d" and are in subdirectories
	if strings.HasSuffix(base, "-a") || strings.HasSuffix(base, "-d") {
		// Check if path contains "test" indicators
		if strings.Contains(path, string(os.PathSeparator)+"test"+string(os.PathSeparator)) {
			return true
		}

		// Check for longer hash patterns typical of test entries
		if len(base) > 40 {
			return true
		}
	}

	return false
}
