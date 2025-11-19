package cache

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// Manager manages the Go build cache
type BuildManager struct {
	cacheDir string
}

var _ CacheManager = (*BuildManager)(nil)

// NewManager creates a new build cache manager
func NewBuildManager(cacheDir string) (*BuildManager, error) {
	if cacheDir == "" {
		dir, err := GetGoEnv("GOCACHE")
		if err != nil {
			return nil, fmt.Errorf("failed to get GOCACHE: %w", err)
		}
		cacheDir = dir
	}

	// Verify cache directory exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("build cache directory does not exist: %s", cacheDir)
	}

	return &BuildManager{
		cacheDir: cacheDir,
	}, nil
}

// GetStats retrieves build cache statistics
func (m *BuildManager) GetStats() (Stats, error) {
	stats := &BuildCacheStats{
		Location:    m.cacheDir,
		OldestEntry: time.Now(),
	}

	// Walk the cache directory
	err := filepath.WalkDir(m.cacheDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Get file info
		info, err := d.Info()
		if err != nil {
			return nil // Skip files we can't read
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

		// Size distribution
		const (
			MB          = 1024 * 1024
			smallLimit  = MB
			mediumLimit = 10 * MB
		)

		if info.Size() < smallLimit {
			stats.Distribution.Small++
			stats.Distribution.SmallSize += info.Size()
		} else if info.Size() < mediumLimit {
			stats.Distribution.Medium++
			stats.Distribution.MediumSize += info.Size()
		} else {
			stats.Distribution.Large++
			stats.Distribution.LargeSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk build cache: %w", err)
	}

	// If no entries, reset oldest to zero time
	if stats.EntryCount == 0 {
		stats.OldestEntry = time.Time{}
	}

	return stats, nil
}

// Clear removes all build cache entries
func (m *BuildManager) Clear() (int, int64, error) {
	var deletedCount int
	var freedSpace int64

	err := filepath.WalkDir(m.cacheDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Skip the root directory
		if path == m.cacheDir {
			return nil
		}

		// Skip directories (will be removed if empty)
		if d.IsDir() {
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
		return deletedCount, freedSpace, fmt.Errorf("failed to clear build cache: %w", err)
	}

	return deletedCount, freedSpace, nil
}

// GetLocation returns the cache directory path
func (m *BuildManager) GetLocation() string {
	return m.cacheDir
}
