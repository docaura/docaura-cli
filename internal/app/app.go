package app

import (
	"context"
	"fmt"
	"github.com/docaura/docaura-cli/internal/fileutils"
	"github.com/docaura/docaura-cli/pkg/analyzer"
	"github.com/docaura/docaura-cli/pkg/docgen"
	"log"
	"os"
	"path/filepath"
)

// App represents the main application.
type App struct {
	config    Config
	analyzer  *analyzer.Analyzer
	generator *docgen.Generator
	watcher   *Watcher
}

// New creates a new application instance.
func New(config Config) (*App, error) {
	// Load config from file if specified
	if err := config.LoadFromFile(config.ConfigFile); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Create analyzer
	analyzer := analyzer.New()

	// Create generator
	generator, err := docgen.New()
	if err != nil {
		return nil, fmt.Errorf("create generator: %w", err)
	}

	app := &App{
		config:    config,
		analyzer:  analyzer,
		generator: generator,
	}

	// Create watcher if needed
	if config.Watch {
		watcher, err := NewWatcher(config)
		if err != nil {
			return nil, fmt.Errorf("create watcher: %w", err)
		}
		app.watcher = watcher
	}

	return app, nil
}

// Run runs the application.
func (a *App) Run() error {
	if a.config.Verbose {
		log.Printf("Starting documentation generation for %s", a.config.ProjectDir)
	}

	if a.config.Watch {
		return a.runWatcher()
	}

	return a.generateOnce()
}

// generateOnce generates documentation once and exits.
func (a *App) generateOnce() error {
	if a.config.PackageName != "" {
		return a.generateSinglePackage()
	}

	return a.generateAllPackages()
}

// generateSinglePackage generates documentation for a specific package.
func (a *App) generateSinglePackage() error {
	packagePath := filepath.Join(a.config.ProjectDir, a.config.PackageName)
	return a.generatePackageDocs(packagePath)
}

// generateAllPackages generates documentation for all packages in the project.
func (a *App) generateAllPackages() error {
	packages, err := fileutils.FindGoPackages(a.config.ProjectDir, a.config.ExcludeDirs)
	if err != nil {
		return fmt.Errorf("find Go packages: %w", err)
	}

	if len(packages) == 0 {
		log.Println("No Go packages found in project directory")
		return nil
	}

	if a.config.Verbose {
		log.Printf("Found %d packages to document", len(packages))
	}

	var errors []error
	for _, packagePath := range packages {
		if err := a.generatePackageDocs(packagePath); err != nil {
			if a.config.Verbose {
				log.Printf("Error documenting package %s: %v", packagePath, err)
			}
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to document %d packages", len(errors))
	}

	return nil
}

// generatePackageDocs generates documentation for a single package.
func (a *App) generatePackageDocs(packagePath string) error {
	if a.config.Verbose {
		log.Printf("Analyzing package: %s", packagePath)
	}

	// Analyze package
	pkg, err := a.analyzer.AnalyzePackage(packagePath)
	if err != nil {
		return fmt.Errorf("analyze package %q: %w", packagePath, err)
	}

	// Generate documentation
	ctx := context.Background()
	docgenConfig := a.config.ToDocgenConfig()

	doc, err := a.generator.GeneratePackageDoc(ctx, pkg, docgenConfig)
	if err != nil {
		return fmt.Errorf("generate documentation: %w", err)
	}

	// Write to file
	outputPath := a.getOutputPath(pkg.Name)
	if err := a.writeDocumentation(outputPath, doc); err != nil {
		return fmt.Errorf("write documentation: %w", err)
	}

	if a.config.Verbose {
		log.Printf("Generated documentation: %s", outputPath)
	}

	return nil
}

// getOutputPath returns the output path for a package's documentation.
func (a *App) getOutputPath(packageName string) string {
	var filename string
	switch a.config.Style {
	case "markdown":
		filename = packageName + ".md"
	case "html":
		filename = packageName + ".html"
	default:
		filename = packageName + ".md"
	}

	return filepath.Join(a.config.OutputDir, filename)
}

// writeDocumentation writes documentation to a file.
func (a *App) writeDocumentation(outputPath, content string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// Write documentation file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("write file %q: %w", outputPath, err)
	}

	return nil
}

// runWatcher runs the file watcher.
func (a *App) runWatcher() error {
	if a.watcher == nil {
		return fmt.Errorf("watcher not initialized")
	}

	// Generate initial documentation
	if err := a.generateOnce(); err != nil {
		log.Printf("Initial generation failed: %v", err)
	}

	// Start watching
	return a.watcher.Watch(func() error {
		if a.config.Verbose {
			log.Println("Regenerating documentation due to file changes...")
		}
		return a.generateAllPackages()
	})
}
