package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/muhammadali7768/gocachectl/internal/cache"
	"github.com/muhammadali7768/gocachectl/internal/cachemgr"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show cache location and environment information",
	Long: `Display information about Go cache locations and environment:
- GOCACHE location
- GOMODCACHE location
- Go version
- Cache availability status`,
	Example: `  gocachectl info
  gocachectl info --json`,
	RunE: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	// Create unified manager
	manager, err := cachemgr.NewUnifiedManager()
	if err != nil {
		return fmt.Errorf("failed to initialize cache manager: %w", err)
	}

	// Get cache info
	info, err := manager.GetCacheInfo()
	if err != nil {
		return fmt.Errorf("failed to get cache info: %w", err)
	}

	if jsonOutput {
		return outputInfoJSON(cmd, info)
	}

	return outputInfoHuman(info)
}

func outputInfoJSON(cmd *cobra.Command, info *cache.CacheInfo) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(info)
}

func outputInfoHuman(info *cache.CacheInfo) error {
	if !quiet {
		fmt.Println("Go Cache Information")
		fmt.Println("====================")
		fmt.Println()
	}

	// Go Version
	fmt.Printf("Go Version:       %s\n", info.GoVersion)
	fmt.Println()

	// Build Cache
	fmt.Println("Build Cache (GOCACHE):")
	fmt.Printf("   Location:      %s\n", info.GOCACHE)
	if info.BuildCacheOK {
		fmt.Printf("   Status:        ✓ Available\n")
	} else {
		fmt.Printf("   Status:        ✗ Not available\n")
	}
	fmt.Println()

	// Module Cache
	fmt.Println("Module Cache (GOMODCACHE):")
	fmt.Printf("   Location:      %s\n", info.GOMODCACHE)
	if info.ModCacheOK {
		fmt.Printf("   Status:        ✓ Available\n")
	} else {
		fmt.Printf("   Status:        ✗ Not available\n")
	}

	if verbose {
		fmt.Println()
		fmt.Println("Note: Test cache is part of the build cache.")
		fmt.Println("Use 'gocachectl stats' to see cache sizes.")
	}

	return nil
}
