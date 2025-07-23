package cmd

import (
	"fmt"
	"github.com/docaura/docaura-cli/internal/app"
	"github.com/spf13/cobra"
)

var (
	// Generate command flags
	projectDir  string
	outputDir   string
	packageName string
	watch       bool
	style       string
	examples    bool
	private     bool
)

var generateCmd = &cobra.Command{
	Use:   "generate [flags]",
	Short: "Generate documentation for Go packages",
	Long: `Analyze Go source code and generate enhanced documentation with AI-powered
descriptions and examples. Supports multiple output formats and can watch
for file changes to automatically regenerate documentation.`,
	Aliases: []string{"gen", "g"},
	RunE:    runGenerate,
	Example: `  # Generate docs for current directory
  docaura generate

  # Generate docs with custom output directory
  docaura generate --output ./my-docs

  # Generate docs for specific package
  docaura generate --package ./internal/utils

  # Watch for changes and auto-regenerate
  docaura generate --watch

  # Generate HTML documentation
  docaura generate --style html

  # Include private symbols and disable examples
  docaura generate --private --examples=false`,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Command-specific flags
	generateCmd.Flags().StringVarP(&projectDir, "dir", "d", ".", "project directory to analyze")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "./docs", "output directory for documentation")
	generateCmd.Flags().StringVarP(&packageName, "package", "p", "", "specific package to document (relative to project dir)")
	generateCmd.Flags().BoolVarP(&watch, "watch", "w", false, "watch for file changes and regenerate documentation")
	generateCmd.Flags().StringVarP(&style, "style", "s", "markdown", "documentation style (markdown, godoc, html)")
	generateCmd.Flags().BoolVar(&examples, "examples", true, "generate AI-enhanced examples")
	generateCmd.Flags().BoolVar(&private, "private", false, "include private (unexported) symbols")

	// Mark commonly used flags
	generateCmd.Flags().Lookup("dir").Usage = "project directory to analyze"
	generateCmd.Flags().Lookup("output").Usage = "output directory for documentation"
	generateCmd.Flags().Lookup("style").Usage = "documentation style: markdown, godoc, or html"
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Create application configuration from global and command flags
	config := GetGlobalConfig()

	// Override with command-specific flags
	config.ProjectDir = projectDir
	config.OutputDir = outputDir
	config.PackageName = packageName
	config.Watch = watch
	config.Style = style
	config.Examples = examples
	config.Private = private

	// Create and run application
	application, err := app.New(config)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	if err := application.Run(); err != nil {
		return fmt.Errorf("application error: %w", err)
	}

	return nil
}
