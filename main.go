package main

import (
	"github.com/docaura/docaura-cli/cmd"
)

// version information set by build system
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Set version information in the cmd package
	cmd.SetVersionInfo(version, commit, date)

	// Execute the root command
	cmd.Execute()
}
