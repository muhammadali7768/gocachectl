package cachemgr

import (
	"fmt"

	"github.com/muhammadali7768/gocachectl/internal/cache"
)

// UnifiedManager acts as a high-level consumer that works with any cache.Manager.
type UnifiedManager struct {
	managers []cache.CacheManager
}

// NewUnifiedManager constructs all cache managers and registers them.
func NewUnifiedManager() (*UnifiedManager, error) {
	var managers []cache.CacheManager
	// Initialize build cache manager
	buildMgr, err := cache.NewBuildManager("")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize build cache: %w", err)
	}

	managers = append(managers, buildMgr)
	// Initialize module cache manager
	modMgr, err := cache.NewModManager("")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize module cache: %w", err)
	}

	managers = append(managers, modMgr)
	// Initialize test cache manager
	testMgr, err := cache.NewTestManager("")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test cache: %w", err)
	}

	managers = append(managers, testMgr)
	return &UnifiedManager{
		managers: managers,
	}, nil
}

// GetAllStats returns stats for all caches providers
func (m *UnifiedManager) GetAllStats() ([]cache.Stats, error) {
	var all []cache.Stats

	for _, mgr := range m.managers {
		stat, err := mgr.GetStats()
		if err != nil {
			return nil, fmt.Errorf("failed to get stats for %T: %w", mgr, err)
		}
		all = append(all, stat)
	}
	return all, nil
}

// GetStatsByType retrieves the stats for a specific type ("mod", "build", "test").
func (m *UnifiedManager) GetStatsByType(kind string) (cache.Stats, error) {
	for _, mgr := range m.managers {
		stats, err := mgr.GetStats()
		if err != nil {
			return nil, err
		}
		if stats.Type() == kind {
			return stats, nil
		}
	}
	return nil, fmt.Errorf("no stats found for type: %s", kind)
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

	for _, mgr := range m.managers {
		stats, err := mgr.GetStats()
		if err != nil {
			result.Errors++
			continue
		}

		kind := stats.Type()

		// Decide whether to clear this manager
		if opts.All ||
			(kind == "build" && opts.Build) ||
			(kind == "module" && opts.Modules) ||
			(kind == "test" && opts.Test) {

			deleted, freed, err := mgr.Clear()
			if err != nil {
				result.Errors++
				continue
			}

			// Update fields based on type
			switch kind {
			case "build":
				result.BuildDeleted += deleted
			case "module":
				result.ModulesDeleted += deleted
			case "test":
				result.TestDeleted += deleted
			}

			result.TotalFreed += freed
		}
	}

	return result, nil
}
