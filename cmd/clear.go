package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/muhammadali7768/gocachectl/internal/cache"
	"github.com/muhammadali7768/gocachectl/internal/cachemgr"
	"github.com/spf13/cobra"
)

var (
	clearAll    bool
	clearBuild  bool
	clearMod    bool
	clearTest   bool
	clearForce  bool
	clearDryRun bool
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear cache entries",
	Long: `Clear Go cache entries. You can clear all caches or specific ones.

By default, a confirmation prompt will be shown before deletion.
Use --force to skip the confirmation prompt.
Use --dry-run to see what would be deleted without actually deleting.`,
	Example: `  gocachectl clear --all                 # Clear all caches (with confirmation)
  gocachectl clear --build               # Clear only build cache
  gocachectl clear --modules             # Clear only module cache
  gocachectl clear --test                # Clear only test cache
  gocachectl clear --all --force         # Clear all without confirmation
  gocachectl clear --all --dry-run       # Show what would be deleted`,
	RunE: runClear,
}

func init() {
	rootCmd.AddCommand(clearCmd)

	clearCmd.Flags().BoolVar(&clearAll, "all", false, "clear all caches")
	clearCmd.Flags().BoolVar(&clearBuild, "build", false, "clear build cache")
	clearCmd.Flags().BoolVar(&clearMod, "modules", false, "clear module cache")
	clearCmd.Flags().BoolVar(&clearTest, "test", false, "clear test cache")
	clearCmd.Flags().BoolVarP(&clearForce, "force", "f", false, "skip confirmation prompt")
	clearCmd.Flags().BoolVar(&clearDryRun, "dry-run", false, "show what would be deleted")
}

func runClear(cmd *cobra.Command, args []string) error {
	// Validate flags
	if !clearAll && !clearBuild && !clearMod && !clearTest {
		return fmt.Errorf("must specify at least one cache to clear: --all, --build, --modules, or --test")
	}

	// Create unified manager
	manager, err := cachemgr.NewUnifiedManager()
	if err != nil {
		return fmt.Errorf("failed to initialize cache manager: %w", err)
	}

	// Get current stats before clearing
	var stats *cache.UnifiedStats
	if !quiet {
		stats, err = manager.GetAllStats()
		if err != nil {
			return fmt.Errorf("failed to get cache stats: %w", err)
		}
	}

	// Show what will be cleared
	if !quiet {
		fmt.Println("Cache entries to be cleared:")
		fmt.Println("============================")
		fmt.Println()

		if clearAll || clearBuild {
			fmt.Printf("📦 Build Cache:  %s (%s entries)\n",
				cache.FormatBytes(stats.BuildCache.Size),
				cache.FormatCount(stats.BuildCache.EntryCount))
		}
		if clearAll || clearMod {
			fmt.Printf("📚 Module Cache: %s (%s modules)\n",
				cache.FormatBytes(stats.ModCache.Size),
				cache.FormatCount(stats.ModCache.ModuleCount))
		}
		if clearAll || clearTest {
			fmt.Printf("🧪 Test Cache:   %s (%s entries)\n",
				cache.FormatBytes(stats.TestCache.Size),
				cache.FormatCount(stats.TestCache.EntryCount))
		}

		totalSize := int64(0)
		if clearAll || clearBuild {
			totalSize += stats.BuildCache.Size
		}
		if clearAll || clearMod {
			totalSize += stats.ModCache.Size
		}
		if clearAll || clearTest {
			totalSize += stats.TestCache.Size
		}

		fmt.Println()
		fmt.Printf("Total to be cleared: %s\n", cache.FormatBytes(totalSize))
		fmt.Println()
	}

	// Dry run mode
	if clearDryRun {
		if !quiet {
			fmt.Println("[DRY RUN] No entries were deleted")
		}
		return nil
	}

	// Confirmation prompt
	if !clearForce {
		if !confirm("Are you sure you want to delete these caches?") {
			if !quiet {
				fmt.Println("Operation cancelled")
			}
			return nil
		}
	}

	// Prepare clear options
	opts := cache.ClearOptions{
		Build:   clearBuild,
		Modules: clearMod,
		Test:    clearTest,
		All:     clearAll,
		Force:   clearForce,
		DryRun:  clearDryRun,
	}

	// Perform clearing
	if !quiet {
		fmt.Println("Clearing caches...")
	}

	result, err := manager.Clear(opts)
	if err != nil {
		return fmt.Errorf("failed to clear caches: %w", err)
	}

	// Show results
	if !quiet {
		fmt.Println()
		fmt.Println("Results:")
		fmt.Println("========")
		fmt.Println()

		if clearAll || clearBuild {
			fmt.Printf("📦 Build Cache:  %s entries deleted\n", cache.FormatCount(result.BuildDeleted))
		}
		if clearAll || clearMod {
			fmt.Printf("📚 Module Cache: %s entries deleted\n", cache.FormatCount(result.ModulesDeleted))
		}
		if clearAll || clearTest {
			fmt.Printf("🧪 Test Cache:   %s entries deleted\n", cache.FormatCount(result.TestDeleted))
		}

		fmt.Println()
		fmt.Printf("Total space freed: %s\n", cache.FormatBytes(result.TotalFreed))

		if result.Errors > 0 {
			fmt.Printf("\n⚠️  Warning: %d errors occurred during clearing\n", result.Errors)
		}
	}

	return nil
}

// confirm prompts the user for confirmation
func confirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", message)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
