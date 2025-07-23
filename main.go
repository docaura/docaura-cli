package main

import (
	"flag"
	"github.com/docaura/docaura-cli/internal/app"
	"log"
	"os"
)

// version information set by build system
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	var (
		projectDir  = flag.String("dir", ".", "project directory to analyze")
		outputDir   = flag.String("output", "./docs", "output directory for documentation")
		configFile  = flag.String("config", "", "path to configuration file (JSON)")
		watch       = flag.Bool("watch", false, "watch for file changes and regenerate documentation")
		packageName = flag.String("package", "", "specific package to document (relative to project dir)")
		style       = flag.String("style", "markdown", "documentation style (markdown, godoc, html)")
		examples    = flag.Bool("examples", true, "generate AI-enhanced examples")
		private     = flag.Bool("private", false, "include private (unexported) symbols")
		verbose     = flag.Bool("v", false, "verbose output")
		showVersion = flag.Bool("version", false, "show version information")
	)
	flag.Parse()

	// Show version information
	if *showVersion {
		showVersionInfo()
		os.Exit(0)
	}

	// Create application configuration
	config := app.Config{
		ProjectDir:  *projectDir,
		OutputDir:   *outputDir,
		ConfigFile:  *configFile,
		PackageName: *packageName,
		Watch:       *watch,
		Style:       *style,
		Examples:    *examples,
		Private:     *private,
		Verbose:     *verbose,
	}

	// Create and run application
	application, err := app.New(config)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func showVersionInfo() {
	log.Printf("gendocs version %s (commit: %s, built: %s)", version, commit, date)
}
