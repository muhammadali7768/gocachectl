package cache

import "time"

// UnifiedStats contains statistics for all Go caches
type UnifiedStats struct {
	BuildCache BuildCacheStats `json:"build_cache"`
	ModCache   ModCacheStats   `json:"module_cache"`
	TestCache  TestCacheStats  `json:"test_cache"`
	TotalSize  int64           `json:"total_size"`
	TotalCount int             `json:"total_count"`
}

// BuildCacheStats contains build cache statistics
type BuildCacheStats struct {
	Location     string           `json:"location"`
	Size         int64            `json:"size"`
	EntryCount   int              `json:"entry_count"`
	OldestEntry  time.Time        `json:"oldest_entry"`
	NewestEntry  time.Time        `json:"newest_entry"`
	Distribution SizeDistribution `json:"distribution"`
}

// ModCacheStats contains module cache statistics
type ModCacheStats struct {
	Location      string       `json:"location"`
	Size          int64        `json:"size"`
	ModuleCount   int          `json:"module_count"`
	DirectCount   int          `json:"direct_count"`
	IndirectCount int          `json:"indirect_count"`
	TopModules    []ModuleInfo `json:"top_modules,omitempty"`
}

// TestCacheStats contains test cache statistics
type TestCacheStats struct {
	Location    string    `json:"location"`
	Size        int64     `json:"size"`
	EntryCount  int       `json:"entry_count"`
	OldestEntry time.Time `json:"oldest_entry"`
	NewestEntry time.Time `json:"newest_entry"`
}

// SizeDistribution tracks distribution of cache entries by size
type SizeDistribution struct {
	Small      int   `json:"small_count"` // < 1MB
	SmallSize  int64 `json:"small_size"`
	Medium     int   `json:"medium_count"` // 1-10MB
	MediumSize int64 `json:"medium_size"`
	Large      int   `json:"large_count"` // > 10MB
	LargeSize  int64 `json:"large_size"`
}

// ModuleInfo contains information about a module
type ModuleInfo struct {
	Path    string `json:"path"`
	Version string `json:"version"`
	Size    int64  `json:"size"`
	Direct  bool   `json:"direct"`
}

// ClearOptions contains options for clearing cache
type ClearOptions struct {
	Build   bool
	Modules bool
	Test    bool
	All     bool
	Force   bool
	DryRun  bool
}

// ClearResult contains the result of a clear operation
type ClearResult struct {
	BuildDeleted   int   `json:"build_deleted"`
	ModulesDeleted int   `json:"modules_deleted"`
	TestDeleted    int   `json:"test_deleted"`
	TotalFreed     int64 `json:"total_freed"`
	Errors         int   `json:"errors"`
}

// CacheInfo contains information about cache locations
type CacheInfo struct {
	GOCACHE      string `json:"gocache"`
	GOMODCACHE   string `json:"gomodcache"`
	GoVersion    string `json:"go_version"`
	BuildCacheOK bool   `json:"build_cache_ok"`
	ModCacheOK   bool   `json:"mod_cache_ok"`
}
