package cachemgr

import (
	"fmt"

	"github.com/muhammadali7768/gocachectl/internal/buildcache"
	"github.com/muhammadali7768/gocachectl/internal/cache"
	"github.com/muhammadali7768/gocachectl/internal/modcache"
	"github.com/muhammadali7768/gocachectl/internal/testcache"
)

// UnifiedManager coordinates all cache operations
type UnifiedManager struct {
	buildCache *buildcache.Manager
	modCache   *modcache.Manager
	testCache  *testcache.Manager
}

// NewUnifiedManager creates a new unified cache manager
func NewUnifiedManager() (*UnifiedManager, error) {
	// Initialize build cache manager
	buildMgr, err := buildcache.NewManager("")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize build cache: %w", err)
	}

	// Initialize module cache manager
	modMgr, err := modcache.NewManager("")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize module cache: %w", err)
	}

	// Initialize test cache manager
	testMgr, err := testcache.NewManager("")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test cache: %w", err)
	}

	return &UnifiedManager{
		buildCache: buildMgr,
		modCache:   modMgr,
		testCache:  testMgr,
	}, nil
}

// GetAllStats retrieves statistics for all caches
func (m *UnifiedManager) GetAllStats() (*cache.UnifiedStats, error) {
	stats := &cache.UnifiedStats{}

	// Get build cache stats
	buildStats, err := m.buildCache.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get build cache stats: %w", err)
	}
	stats.BuildCache = *buildStats

	// Get module cache stats
	modStats, err := m.modCache.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get module cache stats: %w", err)
	}
	stats.ModCache = *modStats

	// Get test cache stats
	testStats, err := m.testCache.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get test cache stats: %w", err)
	}
	stats.TestCache = *testStats

	// Calculate totals
	stats.TotalSize = buildStats.Size + modStats.Size + testStats.Size
	stats.TotalCount = buildStats.EntryCount + modStats.ModuleCount + testStats.EntryCount

	return stats, nil
}

// GetBuildStats retrieves only build cache statistics
func (m *UnifiedManager) GetBuildStats() (*cache.BuildCacheStats, error) {
	return m.buildCache.GetStats()
}

// GetModuleStats retrieves only module cache statistics
func (m *UnifiedManager) GetModuleStats() (*cache.ModCacheStats, error) {
	return m.modCache.GetStats()
}

// GetTestStats retrieves only test cache statistics
func (m *UnifiedManager) GetTestStats() (*cache.TestCacheStats, error) {
	return m.testCache.GetStats()
}

// GetCacheInfo retrieves cache location information
func (m *UnifiedManager) GetCacheInfo() (*cache.CacheInfo, error) {
	info := &cache.CacheInfo{}

	// Get GOCACHE
	gocache, err := cache.GetGoEnv("GOCACHE")
	if err != nil {
		info.BuildCacheOK = false
	} else {
		info.GOCACHE = gocache
		info.BuildCacheOK = true
	}

	// Get GOMODCACHE
	gomodcache, err := cache.GetGoEnv("GOMODCACHE")
	if err != nil {
		info.ModCacheOK = false
	} else {
		info.GOMODCACHE = gomodcache
		info.ModCacheOK = true
	}

	// Get Go version
	version, err := cache.GetGoVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get Go version: %w", err)
	}
	info.GoVersion = version

	return info, nil
}

// Clear removes cache entries based on options
func (m *UnifiedManager) Clear(opts cache.ClearOptions) (*cache.ClearResult, error) {
	result := &cache.ClearResult{}

	// Clear build cache
	if opts.Build || opts.All {
		deleted, freed, err := m.buildCache.Clear()
		if err != nil {
			result.Errors++
		} else {
			result.BuildDeleted = deleted
			result.TotalFreed += freed
		}
	}

	// Clear module cache
	if opts.Modules || opts.All {
		deleted, freed, err := m.modCache.Clear()
		if err != nil {
			result.Errors++
		} else {
			result.ModulesDeleted = deleted
			result.TotalFreed += freed
		}
	}

	// Clear test cache
	if opts.Test || opts.All {
		deleted, freed, err := m.testCache.Clear()
		if err != nil {
			result.Errors++
		} else {
			result.TestDeleted = deleted
			result.TotalFreed += freed
		}
	}

	return result, nil
}
