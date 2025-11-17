package modcache

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/muhammadali7768/gocachectl/internal/cache"
)

// Manager manages the Go module cache
type Manager struct {
	cacheDir string
}

// NewManager creates a new module cache manager
func NewManager(cacheDir string) (*Manager, error) {
	if cacheDir == "" {
		dir, err := cache.GetGoEnv("GOMODCACHE")
		if err != nil {
			return nil, fmt.Errorf("failed to get GOMODCACHE: %w", err)
		}
		cacheDir = dir
	}

	// Verify cache directory exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("module cache directory does not exist: %s", cacheDir)
	}

	return &Manager{
		cacheDir: cacheDir,
	}, nil
}

// GetStats retrieves module cache statistics
func (m *Manager) GetStats() (*cache.ModCacheStats, error) {
	stats := &cache.ModCacheStats{
		Location: m.cacheDir,
	}

	moduleMap := make(map[string]*cache.ModuleInfo)

	// Walk the cache directory
	err := filepath.WalkDir(m.cacheDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Get file info
		info, err := d.Info()
		if err != nil {
			return nil
		}

		// Update total size
		stats.Size += info.Size()

		// Try to determine module from path
		// Module cache structure: $GOMODCACHE/module/path@version/...
		relPath, err := filepath.Rel(m.cacheDir, path)
		if err != nil {
			return nil
		}

		// Extract module path (simplified)
		parts := strings.Split(relPath, string(os.PathSeparator))
		if len(parts) >= 2 {
			// Combine path parts until we find @version
			var modulePath string
			for i, part := range parts {
				if strings.Contains(part, "@") {
					modulePath = filepath.Join(parts[:i+1]...)
					break
				}
			}

			if modulePath != "" {
				if mod, exists := moduleMap[modulePath]; exists {
					mod.Size += info.Size()
				} else {
					moduleMap[modulePath] = &cache.ModuleInfo{
						Path: modulePath,
						Size: info.Size(),
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk module cache: %w", err)
	}

	// Count modules
	stats.ModuleCount = len(moduleMap)

	// Get top modules by size
	stats.TopModules = getTopModules(moduleMap, 10)

	return stats, nil
}

// Clear removes all module cache entries
func (m *Manager) Clear() (int, int64, error) {
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

		// Skip directories
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
		return deletedCount, freedSpace, fmt.Errorf("failed to clear module cache: %w", err)
	}

	return deletedCount, freedSpace, nil
}

// GetLocation returns the cache directory path
func (m *Manager) GetLocation() string {
	return m.cacheDir
}

// getTopModules returns the N largest modules by size
func getTopModules(modules map[string]*cache.ModuleInfo, n int) []cache.ModuleInfo {
	// Convert map to slice
	var moduleList []cache.ModuleInfo
	for _, mod := range modules {
		moduleList = append(moduleList, *mod)
	}

	// Simple bubble sort for top N (sufficient for small lists)
	for i := 0; i < len(moduleList) && i < n; i++ {
		for j := i + 1; j < len(moduleList); j++ {
			if moduleList[j].Size > moduleList[i].Size {
				moduleList[i], moduleList[j] = moduleList[j], moduleList[i]
			}
		}
	}

	// Return top N
	if len(moduleList) > n {
		return moduleList[:n]
	}
	return moduleList
}
