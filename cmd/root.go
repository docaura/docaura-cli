package cmd

import (
	"fmt"
	"github.com/docaura/docaura-cli/internal/app"
	"github.com/spf13/cobra"
	"os"
)

var (
	// Global flags
	verbose    bool
	configFile string

	// Version info (set by build system)
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "docaura",
	Short: "Docaura is an AI-powered documentation generator for Go projects",
	Long: `Docaura analyzes Go source code and generates enhanced documentation
with AI-powered descriptions and examples. It supports multiple output
formats including Markdown, HTML, and Godoc-style documentation.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to configuration file (JSON)")

	// Set version info
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)

	// Add version template
	rootCmd.SetVersionTemplate(`{{printf "%s version %s\n" .Name .Version}}`)
}

// SetVersionInfo sets the version information (called from main.go)
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

// GetGlobalConfig returns a base config with global flags set
func GetGlobalConfig() app.Config {
	config := app.DefaultConfig()
	config.ConfigFile = configFile
	config.Verbose = verbose
	return config
}
