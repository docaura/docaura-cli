package docgen

import (
	"context"
	"github.com/docaura/docaura-cli/pkg/analyzer"
	"github.com/tmc/langchaingo/prompts"
	"strings"
)

// generatePackageExample generates a package-level usage example.
func (g *Generator) generatePackageExample(ctx context.Context, pkg *analyzer.PackageInfo) (string, error) {
	template := prompts.NewPromptTemplate(`
Create a realistic Go code example showing how to use this package:

Package: {{.name}}
Description: {{.description}}
Key Functions: {{range .functions}}{{if .is_exported}}{{.name}}, {{end}}{{end}}
Key Types: {{range .types}}{{if .is_exported}}{{.name}}, {{end}}{{end}}

Write a complete, runnable example that shows:
1. Import statement
2. Basic usage
3. Error handling
4. Realistic use case

Return only the Go code, no explanations.`,
		[]string{"name", "description", "functions", "types"})

	prompt, err := template.Format(map[string]any{
		"name":        pkg.Name,
		"description": pkg.Description,
		"functions":   pkg.Functions,
		"types":       pkg.Types,
	})
	if err != nil {
		return "", err
	}

	response, err := g.generateContent(ctx, prompt)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}

// generateFunctionExample generates an example for a specific function.
func (g *Generator) generateFunctionExample(ctx context.Context, fn *analyzer.FunctionInfo, pkg *analyzer.PackageInfo) (string, error) {
	template := prompts.NewPromptTemplate(`
Create a Go code example for this function:

Function: {{.name}}
Signature: {{.signature}}
Package: {{.package}}
{{if .parameters}}Parameters: {{range .parameters}}{{.name}} {{.type}}, {{end}}{{end}}

Write a realistic example showing how to call this function.
Include proper error handling if needed.
Return only the Go code snippet.`,
		[]string{"name", "signature", "package", "parameters"})

	prompt, err := template.Format(map[string]any{
		"name":       fn.Name,
		"signature":  fn.Signature,
		"package":    pkg.Name,
		"parameters": fn.Parameters,
	})
	if err != nil {
		return "", err
	}

	response, err := g.generateContent(ctx, prompt)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}
