package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/docaura/docaura-cli/internal/app"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands",
	Long:  `Commands for managing docaura configuration files and settings.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show [config-file]",
	Short: "Show current configuration",
	Long: `Display the current configuration settings. If no config file is specified,
it will look for docaura.json in the current directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runConfigShow,
	Example: `  # Show config from current directory
  docaura config show

  # Show specific config file
  docaura config show ./path/to/docaura.json`,
}

var configValidateCmd = &cobra.Command{
	Use:   "validate [config-file]",
	Short: "Validate configuration file",
	Long: `Validate a docaura configuration file for syntax and required fields.
If no config file is specified, it will validate docaura.json in the current directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runConfigValidate,
	Example: `  # Validate config in current directory
  docaura config validate

  # Validate specific config file
  docaura config validate ./path/to/docaura.json`,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	configPath := getConfigPath(args)

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found: %s", configPath)
	}

	// Load and display configuration
	config := app.DefaultConfig()
	if err := config.LoadFromFile(configPath); err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	// Pretty print the configuration
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal configuration: %w", err)
	}

	fmt.Printf("Configuration from %s:\n\n", configPath)
	fmt.Println(string(configData))

	return nil
}

func runConfigValidate(cmd *cobra.Command, args []string) error {
	configPath := getConfigPath(args)

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found: %s", configPath)
	}

	// Load configuration
	config := app.DefaultConfig()
	if err := config.LoadFromFile(configPath); err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		fmt.Printf("❌ Configuration validation failed: %v\n", err)
		return err
	}

	fmt.Printf("✓ Configuration file %s is valid\n", configPath)
	return nil
}

func getConfigPath(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return filepath.Join(".", "docaura.json")
}
