package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/docaura/docaura-cli/internal/app"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	// Init command flags
	initForce bool
	initName  string
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new docaura configuration file",
	Long: `Create a new docaura configuration file (docaura.json) in the specified
directory or current directory. This file can be used to customize
documentation generation settings.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
	Example: `  # Initialize config in current directory
  docaura init

  # Initialize config in specific directory
  docaura init ./my-project

  # Force overwrite existing config
  docaura init --force

  # Initialize with custom project name
  docaura init --name "My Project"`,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Command-specific flags
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "overwrite existing configuration file")
	initCmd.Flags().StringVarP(&initName, "name", "n", "", "project name for the configuration")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// Convert to absolute path
	absDir, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("resolve target directory: %w", err)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(absDir, 0755); err != nil {
		return fmt.Errorf("create target directory: %w", err)
	}

	configPath := filepath.Join(absDir, "docaura.json")

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil && !initForce {
		return fmt.Errorf("configuration file already exists at %s (use --force to overwrite)", configPath)
	}

	// Create default configuration
	config := app.DefaultConfig()
	config.ProjectDir = absDir

	// Set project name if provided or derive from directory
	if initName != "" {
		config.ProjectName = initName
	} else {
		config.ProjectName = filepath.Base(absDir)
	}

	// Set a default description
	config.ProjectDescription = fmt.Sprintf("Documentation for %s", config.ProjectName)

	// Convert config to JSON
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal configuration: %w", err)
	}

	// Write configuration file
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return fmt.Errorf("write configuration file: %w", err)
	}

	fmt.Printf("✓ Created configuration file: %s\n", configPath)
	fmt.Printf("✓ Project name: %s\n", config.ProjectName)
	fmt.Printf("✓ Output directory: %s\n", config.OutputDir)
	fmt.Println("\nYou can now run 'docaura generate' to create documentation.")
	fmt.Println("Edit docaura.json to customize your documentation settings.")

	return nil
}
