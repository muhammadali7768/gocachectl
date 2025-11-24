package cache

import (
	"fmt"
	"io"
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
	// Test cache entries must end with "-d" (data/output)
	// Build artifacts also end with "-d" or "-a", so we need to check content
	if !strings.HasSuffix(path, "-d") {
		return false
	}

	// Open file to check header
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	// Read first 512 bytes (enough for header checks)
	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}
	header := buf[:n]

	// Check for known binary/metadata headers to EXCLUDE
	if len(header) >= 7 && string(header[:7]) == "!<arch>" {
		return false
	}
	if len(header) >= 4 && string(header[:4]) == "\x7fELF" {
		return false
	}
	if len(header) >= 2 && string(header[:2]) == "\r\xff" { // GOB/Binary data
		return false
	}
	if len(header) >= 9 && string(header[:9]) == "go object" {
		return false
	}
	if len(header) >= 8 && string(header[:8]) == "go index" {
		return false
	}
	if len(header) >= 3 && string(header[:3]) == "v1 " { // Action graph/metadata
		return false
	}
	if len(header) >= 2 && string(header[:2]) == "./" { // Source file list
		return false
	}

	// Check for known test output patterns to INCLUDE
	// Standard success: "ok \t" or "ok "
	if len(header) >= 3 && string(header[:3]) == "ok " {
		return true
	}
	if len(header) >= 3 && string(header[:3]) == "ok\t" {
		return true
	}
	// Verbose output: "=== RUN"
	if len(header) >= 7 && string(header[:7]) == "=== RUN" {
		return true
	}
	// Failure (rarely cached, but possible): "FAIL"
	if len(header) >= 4 && string(header[:4]) == "FAIL" {
		return true
	}

	return false
}
