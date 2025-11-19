package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/muhammadali7768/gocachectl/internal/cache"
	"github.com/muhammadali7768/gocachectl/internal/cachemgr"
	"github.com/spf13/cobra"
)

var (
	showBuild   bool
	showModules bool
	showTest    bool
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show cache statistics",
	Long: `Display statistics about Go caches including:
- Build cache size and entries
- Module cache size and module count
- Test cache size and entries
- Total size across all caches

Use flags to show specific cache statistics.`,
	Example: `  gocachectl stats              # Show all cache stats
  gocachectl stats --build      # Show only build cache
  gocachectl stats --modules    # Show only module cache
  gocachectl stats --json       # Output as JSON`,
	RunE: runStats,
}

func init() {
	rootCmd.AddCommand(statsCmd)

	statsCmd.Flags().BoolVar(&showBuild, "build", false, "show only build cache statistics")
	statsCmd.Flags().BoolVar(&showModules, "modules", false, "show only module cache statistics")
	statsCmd.Flags().BoolVar(&showTest, "test", false, "show only test cache statistics")
}

func runStats(cmd *cobra.Command, args []string) error {
	// Create unified manager
	manager, err := cachemgr.NewUnifiedManager()
	if err != nil {
		return fmt.Errorf("failed to initialize cache manager: %w", err)
	}

	// Determine what to show
	showAll := !showBuild && !showModules && !showTest

	if jsonOutput {
		return outputStatsJSON(cmd, manager, showAll)
	}

	return outputStatsHuman(manager, showAll)
}

func outputStatsJSON(cmd *cobra.Command, manager *cachemgr.UnifiedManager, showAll bool) error {
	var data interface{}

	if showAll {
		stats, err := manager.GetAllStats()
		if err != nil {
			return err
		}
		data = stats
	} else if showBuild {
		stats, err := manager.GetStatsByType("build")
		if err != nil {
			return err
		}
		data = stats
	} else if showModules {
		stats, err := manager.GetStatsByType("module")
		if err != nil {
			return err
		}
		data = stats
	} else if showTest {
		stats, err := manager.GetStatsByType("tests")
		if err != nil {
			return err
		}
		data = stats
	}

	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func outputStatsHuman(manager *cachemgr.UnifiedManager, showAll bool) error {
	if showAll {
		return outputAllStats(manager)
	}

	if showBuild {
		return outputBuildStats(manager)
	}

	if showModules {
		return outputModuleStats(manager)
	}

	if showTest {
		return outputTestStats(manager)
	}

	return nil
}

func outputAllStats(manager *cachemgr.UnifiedManager) error {
	all, err := manager.GetAllStats()
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Println("Go Cache Statistics")
		fmt.Println("===================")
		fmt.Println()
	}

	var totalSize int64
	var totalCount int
	for _, stats := range all {

		switch stat := stats.(type) {
		case *cache.BuildCacheStats:
			totalCount += stat.EntryCount
			totalSize += stat.Size
			// Build Cache
			fmt.Println("Build Cache")
			fmt.Printf("   Location:     %s\n", stat.Location)
			fmt.Printf("   Size:         %s\n", cache.FormatBytes(stat.Size))
			fmt.Printf("   Entries:      %s\n", cache.FormatCount(stat.EntryCount))
			if !stat.OldestEntry.IsZero() {
				fmt.Printf("   Oldest:       %s\n", stat.OldestEntry.Format("2006-01-02 15:04:05"))
				fmt.Printf("   Newest:       %s\n", stat.NewestEntry.Format("2006-01-02 15:04:05"))
			}
			fmt.Println()
		case *cache.ModCacheStats:
			totalCount += stat.ModuleCount
			totalSize += stat.Size
			// Module Cache
			fmt.Println("Module Cache")
			fmt.Printf("   Location:     %s\n", stat.Location)
			fmt.Printf("   Size:         %s\n", cache.FormatBytes(stat.Size))
			fmt.Printf("   Modules:      %s\n", cache.FormatCount(stat.ModuleCount))
			fmt.Println()
		case *cache.TestCacheStats:
			totalCount += stat.EntryCount
			totalSize += stat.Size
			// Test Cache
			fmt.Println("Test Cache")
			fmt.Printf("   Location:     %s\n", stat.Location)
			fmt.Printf("   Size:         %s\n", cache.FormatBytes(stat.Size))
			fmt.Printf("   Entries:      %s\n", cache.FormatCount(stat.EntryCount))
			if !stat.OldestEntry.IsZero() {
				fmt.Printf("   Oldest:       %s\n", stat.OldestEntry.Format("2006-01-02 15:04:05"))
				fmt.Printf("   Newest:       %s\n", stat.NewestEntry.Format("2006-01-02 15:04:05"))
			}
			fmt.Println()
		}
	}
	// Total
	fmt.Println("Total")
	fmt.Printf("   Total Size:   %s\n", cache.FormatBytes(totalSize))
	fmt.Printf("   Total Items:  %s\n", cache.FormatCount(totalCount))

	// Verbose output
	if verbose {
		fmt.Println()
		fmt.Println("Size Distribution (Build Cache):")
		for _, stats := range all {
			switch stat := stats.(type) {
			case *cache.BuildCacheStats:
				fmt.Printf("   Small (<1MB):    %d entries (%s)\n",
					stat.Distribution.Small, cache.FormatBytes(stat.Distribution.SmallSize))
				fmt.Printf("   Medium (1-10MB): %d entries (%s)\n",
					stat.Distribution.Medium, cache.FormatBytes(stat.Distribution.MediumSize))
				fmt.Printf("   Large (>10MB):   %d entries (%s)\n",
					stat.Distribution.Large, cache.FormatBytes(stat.Distribution.LargeSize))
			case *cache.ModCacheStats:
				if len(stat.TopModules) > 0 {
					fmt.Println()
					fmt.Println("Top Modules by Size:")
					for i, mod := range stat.TopModules {
						if i >= 5 {
							break
						}
						fmt.Printf("   %d. %s (%s)\n", i+1, mod.Path, cache.FormatBytes(mod.Size))
					}
				}
			}
		}
	}

	return nil
}

func outputBuildStats(manager *cachemgr.UnifiedManager) error {
	stats, err := manager.GetStatsByType("build")
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Println("Build Cache Statistics")
		fmt.Println("======================")
		fmt.Println()
	}

	buildCacheStats := stats.(*cache.BuildCacheStats)

	fmt.Printf("Location:     %s\n", buildCacheStats.Location)
	fmt.Printf("Size:         %s\n", cache.FormatBytes(buildCacheStats.Size))
	fmt.Printf("Entries:      %s\n", cache.FormatCount(buildCacheStats.EntryCount))
	if !buildCacheStats.OldestEntry.IsZero() {
		fmt.Printf("Oldest Entry: %s\n", buildCacheStats.OldestEntry.Format("2006-01-02 15:04:05"))
		fmt.Printf("Newest Entry: %s\n", buildCacheStats.NewestEntry.Format("2006-01-02 15:04:05"))
	}

	return nil
}

func outputModuleStats(manager *cachemgr.UnifiedManager) error {
	stats, err := manager.GetStatsByType("module")
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Println("Module Cache Statistics")
		fmt.Println("=======================")
		fmt.Println()
	}

	moduleCacheStats := stats.(*cache.ModCacheStats)
	fmt.Printf("Location:     %s\n", moduleCacheStats.Location)
	fmt.Printf("Size:         %s\n", cache.FormatBytes(moduleCacheStats.Size))
	fmt.Printf("Modules:      %s\n", cache.FormatCount(moduleCacheStats.ModuleCount))

	if verbose && len(moduleCacheStats.TopModules) > 0 {
		fmt.Println()
		fmt.Println("Top Modules by Size:")
		for i, mod := range moduleCacheStats.TopModules {
			fmt.Printf("   %d. %s (%s)\n", i+1, mod.Path, cache.FormatBytes(mod.Size))
		}
	}

	return nil
}

func outputTestStats(manager *cachemgr.UnifiedManager) error {
	stats, err := manager.GetStatsByType("test")
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Println("Test Cache Statistics")
		fmt.Println("=====================")
		fmt.Println()
	}
	testCacheStats := stats.(*cache.TestCacheStats)
	fmt.Printf("Location:     %s\n", testCacheStats.Location)
	fmt.Printf("Size:         %s\n", cache.FormatBytes(testCacheStats.Size))
	fmt.Printf("Entries:      %s\n", cache.FormatCount(testCacheStats.EntryCount))
	if !testCacheStats.OldestEntry.IsZero() {
		fmt.Printf("Oldest Entry: %s\n", testCacheStats.OldestEntry.Format("2006-01-02 15:04:05"))
		fmt.Printf("Newest Entry: %s\n", testCacheStats.NewestEntry.Format("2006-01-02 15:04:05"))
	}

	return nil
}
