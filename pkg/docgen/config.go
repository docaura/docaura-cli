package docgen

import "fmt"

// Config represents configuration for documentation generation.
type Config struct {
	ProjectName      string `json:"project_name"`
	ProjectDesc      string `json:"project_description"`
	OutputDir        string `json:"output_dir"`
	IncludePrivate   bool   `json:"include_private"`
	GenerateExamples bool   `json:"generate_examples"`
	Style            string `json:"style"` // "godoc", "markdown", "html"
}

// Validate validates the configuration and sets defaults.
func (c *Config) Validate() error {
	if c.OutputDir == "" {
		c.OutputDir = "./docs"
	}

	if c.Style == "" {
		c.Style = "markdown"
	}

	validStyles := map[string]bool{
		"godoc":    true,
		"markdown": true,
		"html":     true,
	}

	if !validStyles[c.Style] {
		return fmt.Errorf("invalid style %q: must be one of godoc, markdown, html", c.Style)
	}

	return nil
}
