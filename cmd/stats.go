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
		stats, err := manager.GetBuildStats()
		if err != nil {
			return err
		}
		data = stats
	} else if showModules {
		stats, err := manager.GetModuleStats()
		if err != nil {
			return err
		}
		data = stats
	} else if showTest {
		stats, err := manager.GetTestStats()
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
	stats, err := manager.GetAllStats()
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Println("Go Cache Statistics")
		fmt.Println("===================")
		fmt.Println()
	}

	// Build Cache
	fmt.Println("📦 Build Cache")
	fmt.Printf("   Location:     %s\n", stats.BuildCache.Location)
	fmt.Printf("   Size:         %s\n", cache.FormatBytes(stats.BuildCache.Size))
	fmt.Printf("   Entries:      %s\n", cache.FormatCount(stats.BuildCache.EntryCount))
	if !stats.BuildCache.OldestEntry.IsZero() {
		fmt.Printf("   Oldest:       %s\n", stats.BuildCache.OldestEntry.Format("2006-01-02 15:04:05"))
		fmt.Printf("   Newest:       %s\n", stats.BuildCache.NewestEntry.Format("2006-01-02 15:04:05"))
	}
	fmt.Println()

	// Module Cache
	fmt.Println("📚 Module Cache")
	fmt.Printf("   Location:     %s\n", stats.ModCache.Location)
	fmt.Printf("   Size:         %s\n", cache.FormatBytes(stats.ModCache.Size))
	fmt.Printf("   Modules:      %s\n", cache.FormatCount(stats.ModCache.ModuleCount))
	fmt.Println()

	// Test Cache
	fmt.Println("🧪 Test Cache")
	fmt.Printf("   Location:     %s\n", stats.TestCache.Location)
	fmt.Printf("   Size:         %s\n", cache.FormatBytes(stats.TestCache.Size))
	fmt.Printf("   Entries:      %s\n", cache.FormatCount(stats.TestCache.EntryCount))
	if !stats.TestCache.OldestEntry.IsZero() {
		fmt.Printf("   Oldest:       %s\n", stats.TestCache.OldestEntry.Format("2006-01-02 15:04:05"))
		fmt.Printf("   Newest:       %s\n", stats.TestCache.NewestEntry.Format("2006-01-02 15:04:05"))
	}
	fmt.Println()

	// Total
	fmt.Println("📊 Total")
	fmt.Printf("   Total Size:   %s\n", cache.FormatBytes(stats.TotalSize))
	fmt.Printf("   Total Items:  %s\n", cache.FormatCount(stats.TotalCount))

	// Verbose output
	if verbose {
		fmt.Println()
		fmt.Println("Size Distribution (Build Cache):")
		fmt.Printf("   Small (<1MB):    %d entries (%s)\n",
			stats.BuildCache.Distribution.Small, cache.FormatBytes(stats.BuildCache.Distribution.SmallSize))
		fmt.Printf("   Medium (1-10MB): %d entries (%s)\n",
			stats.BuildCache.Distribution.Medium, cache.FormatBytes(stats.BuildCache.Distribution.MediumSize))
		fmt.Printf("   Large (>10MB):   %d entries (%s)\n",
			stats.BuildCache.Distribution.Large, cache.FormatBytes(stats.BuildCache.Distribution.LargeSize))

		if len(stats.ModCache.TopModules) > 0 {
			fmt.Println()
			fmt.Println("Top Modules by Size:")
			for i, mod := range stats.ModCache.TopModules {
				if i >= 5 {
					break
				}
				fmt.Printf("   %d. %s (%s)\n", i+1, mod.Path, cache.FormatBytes(mod.Size))
			}
		}
	}

	return nil
}

func outputBuildStats(manager *cachemgr.UnifiedManager) error {
	stats, err := manager.GetBuildStats()
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Println("Build Cache Statistics")
		fmt.Println("======================")
		fmt.Println()
	}

	fmt.Printf("Location:     %s\n", stats.Location)
	fmt.Printf("Size:         %s\n", cache.FormatBytes(stats.Size))
	fmt.Printf("Entries:      %s\n", cache.FormatCount(stats.EntryCount))
	if !stats.OldestEntry.IsZero() {
		fmt.Printf("Oldest Entry: %s\n", stats.OldestEntry.Format("2006-01-02 15:04:05"))
		fmt.Printf("Newest Entry: %s\n", stats.NewestEntry.Format("2006-01-02 15:04:05"))
	}

	return nil
}

func outputModuleStats(manager *cachemgr.UnifiedManager) error {
	stats, err := manager.GetModuleStats()
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Println("Module Cache Statistics")
		fmt.Println("=======================")
		fmt.Println()
	}

	fmt.Printf("Location:     %s\n", stats.Location)
	fmt.Printf("Size:         %s\n", cache.FormatBytes(stats.Size))
	fmt.Printf("Modules:      %s\n", cache.FormatCount(stats.ModuleCount))

	if verbose && len(stats.TopModules) > 0 {
		fmt.Println()
		fmt.Println("Top Modules by Size:")
		for i, mod := range stats.TopModules {
			fmt.Printf("   %d. %s (%s)\n", i+1, mod.Path, cache.FormatBytes(mod.Size))
		}
	}

	return nil
}

func outputTestStats(manager *cachemgr.UnifiedManager) error {
	stats, err := manager.GetTestStats()
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Println("Test Cache Statistics")
		fmt.Println("=====================")
		fmt.Println()
	}

	fmt.Printf("Location:     %s\n", stats.Location)
	fmt.Printf("Size:         %s\n", cache.FormatBytes(stats.Size))
	fmt.Printf("Entries:      %s\n", cache.FormatCount(stats.EntryCount))
	if !stats.OldestEntry.IsZero() {
		fmt.Printf("Oldest Entry: %s\n", stats.OldestEntry.Format("2006-01-02 15:04:05"))
		fmt.Printf("Newest Entry: %s\n", stats.NewestEntry.Format("2006-01-02 15:04:05"))
	}

	return nil
}
