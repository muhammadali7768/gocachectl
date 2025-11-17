package main

import (
	"fmt"
	"os"

	"github.com/muhammadali7768/gocachectl/cmd"
)

var (
	// Set during build time using ldflags
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version info
	cmd.SetVersionInfo(version, commit, date)

	// Execute root command
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
