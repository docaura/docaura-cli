package app

import (
	"encoding/json"
	"fmt"
	"github.com/docaura/docaura-cli/pkg/docgen"
	"os"
	"path/filepath"
)

// Config represents the application configuration.
type Config struct {
	// CLI flags
	ProjectDir  string `json:"project_dir"`
	OutputDir   string `json:"output_dir"`
	ConfigFile  string `json:"-"` // Not serialized
	PackageName string `json:"package_name"`
	Watch       bool   `json:"watch"`
	Style       string `json:"style"`
	Examples    bool   `json:"examples"`
	Private     bool   `json:"private"`
	Verbose     bool   `json:"verbose"`

	// Additional config file options
	ProjectName        string   `json:"project_name"`
	ProjectDescription string   `json:"project_description"`
	ExcludeDirs        []string `json:"exclude_dirs"`
	WatchInterval      int      `json:"watch_interval_seconds"`
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() Config {
	return Config{
		ProjectDir:    ".",
		OutputDir:     "./docs",
		Style:         "markdown",
		Examples:      true,
		Private:       false,
		Verbose:       false,
		WatchInterval: 5,
		ExcludeDirs: []string{
			"vendor",
			".git",
			"testdata",
			"node_modules",
		},
	}
}

// LoadFromFile loads configuration from a JSON file, merging with existing values.
func (c *Config) LoadFromFile(filename string) error {
	if filename == "" {
		return nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read config file %q: %w", filename, err)
	}

	// Create a temporary config to unmarshal into
	var fileConfig Config
	if err := json.Unmarshal(data, &fileConfig); err != nil {
		return fmt.Errorf("parse config file %q: %w", filename, err)
	}

	// Merge file config with CLI config (CLI takes precedence)
	c.mergeFrom(fileConfig)

	return nil
}

// Validate validates the configuration and sets any missing defaults.
func (c *Config) Validate() error {
	// Ensure project directory exists
	if _, err := os.Stat(c.ProjectDir); os.IsNotExist(err) {
		return fmt.Errorf("project directory %q does not exist", c.ProjectDir)
	}

	// Convert to absolute paths
	var err error
	if c.ProjectDir, err = filepath.Abs(c.ProjectDir); err != nil {
		return fmt.Errorf("resolve project directory: %w", err)
	}

	if c.OutputDir, err = filepath.Abs(c.OutputDir); err != nil {
		return fmt.Errorf("resolve output directory: %w", err)
	}

	// Validate style
	validStyles := map[string]bool{
		"markdown": true,
		"godoc":    true,
		"html":     true,
	}
	if !validStyles[c.Style] {
		return fmt.Errorf("invalid style %q: must be one of markdown, godoc, html", c.Style)
	}

	// Set defaults
	if c.WatchInterval <= 0 {
		c.WatchInterval = 5
	}

	if len(c.ExcludeDirs) == 0 {
		c.ExcludeDirs = DefaultConfig().ExcludeDirs
	}

	return nil
}

// ToDocgenConfig converts the app config to a docgen.Config.
func (c *Config) ToDocgenConfig() docgen.Config {
	return docgen.Config{
		ProjectName:      c.ProjectName,
		ProjectDesc:      c.ProjectDescription,
		OutputDir:        c.OutputDir,
		IncludePrivate:   c.Private,
		GenerateExamples: c.Examples,
		Style:            c.Style,
	}
}

// mergeFrom merges values from another config, keeping existing non-zero values.
func (c *Config) mergeFrom(other Config) {
	if c.ProjectName == "" && other.ProjectName != "" {
		c.ProjectName = other.ProjectName
	}
	if c.ProjectDescription == "" && other.ProjectDescription != "" {
		c.ProjectDescription = other.ProjectDescription
	}
	if len(other.ExcludeDirs) > 0 {
		c.ExcludeDirs = other.ExcludeDirs
	}
	if other.WatchInterval > 0 {
		c.WatchInterval = other.WatchInterval
	}
}
