package docgen

import (
	"fmt"
	"strings"
	"text/template"
)

// TemplateManager manages documentation templates.
type TemplateManager struct {
	templates map[string]*template.Template
}

// NewTemplateManager creates a new template manager with default templates.
func NewTemplateManager() (*TemplateManager, error) {
	tm := &TemplateManager{
		templates: make(map[string]*template.Template),
	}

	if err := tm.loadDefaultTemplates(); err != nil {
		return nil, fmt.Errorf("load default templates: %w", err)
	}

	return tm, nil
}

// Execute executes the template for the given style.
func (tm *TemplateManager) Execute(style string, data interface{}) (string, error) {
	tmpl, exists := tm.templates[style]
	if !exists {
		return "", fmt.Errorf("template for style %q not found", style)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return result.String(), nil
}

// loadDefaultTemplates loads the default templates.
func (tm *TemplateManager) loadDefaultTemplates() error {
	// Markdown template
	markdownTemplate := `# {{.Name}}

{{.Description}}

## Installation

` + "```bash" + `
go get {{.Path}}
` + "```" + `

## Usage

{{if .Examples}}
{{range .Examples}}
` + "```go" + `
{{.Code}}
` + "```" + `
{{end}}
{{end}}

## API Reference

{{if .Functions}}
### Functions

{{range .Functions}}
{{if .IsExported}}
#### {{.Name}}

` + "```go" + `
{{.Signature}}
` + "```" + `

{{.Description}}

{{if .Parameters}}
**Parameters:**
{{range .Parameters}}
- ` + "`{{.Name}}`" + ` ({{.Type}})
{{end}}
{{end}}

{{if .Returns}}
**Returns:**
{{range .Returns}}
- {{.Type}}{{if .Description}} - {{.Description}}{{end}}
{{end}}
{{end}}

{{if .Examples}}
**Example:**
{{range .Examples}}
` + "```go" + `
{{.}}
` + "```" + `
{{end}}
{{end}}

{{end}}
{{end}}
{{end}}

{{if .Types}}
### Types

{{range .Types}}
{{if .IsExported}}
#### {{.Name}}

` + "```go" + `
type {{.Name}} {{.Kind}}
` + "```" + `

{{.Description}}

{{if .Fields}}
**Fields:**
{{range .Fields}}
- ` + "`{{.Name}}`" + ` {{.Type}}{{if .Description}} - {{.Description}}{{end}}
{{end}}
{{end}}

{{if .Methods}}
**Methods:**
{{range .Methods}}
- [{{.}}](#{{. | lower}})
{{end}}
{{end}}

{{end}}
{{end}}
{{end}}
`

	if err := tm.addTemplate("markdown", markdownTemplate); err != nil {
		return fmt.Errorf("add markdown template: %w", err)
	}

	// Add other templates (godoc, html) here as needed

	return nil
}

// addTemplate adds a template with the given name and content.
func (tm *TemplateManager) addTemplate(name, content string) error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	tmpl, err := template.New(name).Funcs(funcMap).Parse(content)
	if err != nil {
		return fmt.Errorf("parse template %q: %w", name, err)
	}

	tm.templates[name] = tmpl
	return nil
}
