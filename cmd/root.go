package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	verbose    bool
	jsonOutput bool
	quiet      bool

	// Version info
	version string
	commit  string
	date    string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "gocachectl",
	Short: "Universal Go cache management tool",
	Long: `gocachectl is a comprehensive CLI tool to manage all Go caches.

Manage build cache, module cache, and test cache from a single interface.
View statistics, clear caches, and optimize disk usage.

Examples:
  gocachectl stats               Show all cache statistics
  gocachectl stats --modules     Show only module cache stats
  gocachectl clear --all         Clear all caches
  gocachectl info                Show cache locations`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gocachectl.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "minimal output")

	// Version command
	rootCmd.AddCommand(versionCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
			}
			return
		}

		// Search config in home directory with name ".gocachectl" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gocachectl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// SetVersionInfo sets version information
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gocachectl %s\n", version)
		if verbose {
			fmt.Printf("commit: %s\n", commit)
			fmt.Printf("built: %s\n", date)
		}
	},
}
