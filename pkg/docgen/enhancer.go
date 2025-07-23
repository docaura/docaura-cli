package docgen

import (
	"context"
	"fmt"
	"github.com/docaura/docaura-cli/pkg/analyzer"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	"strings"
)

// enhancePackageDescription generates an enhanced description for a package.
func (g *Generator) enhancePackageDescription(ctx context.Context, pkg *analyzer.PackageInfo) (string, error) {
	template := prompts.NewPromptTemplate(`
Analyze this Go package and write a clear, concise description (2-3 sentences):

Package: {{.name}}
Path: {{.path}}

Functions: {{range .functions}}{{.name}}, {{end}}
Types: {{range .types}}{{.name}}, {{end}}

Write a professional description that explains:
1. What this package does
2. Who would use it
3. Key capabilities

Keep it under 200 words and avoid marketing language.`,
		[]string{"name", "path", "functions", "types"})

	prompt, err := template.Format(map[string]any{
		"name":      pkg.Name,
		"path":      pkg.Path,
		"functions": pkg.Functions,
		"types":     pkg.Types,
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

// enhanceFunctionDescription generates an enhanced description for a function.
func (g *Generator) enhanceFunctionDescription(ctx context.Context, fn *analyzer.FunctionInfo) (string, error) {
	template := prompts.NewPromptTemplate(`
Write a clear description for this Go function:

Function: {{.name}}
Signature: {{.signature}}
{{if .parameters}}Parameters: {{range .parameters}}{{.name}} {{.type}}, {{end}}{{end}}
{{if .returns}}Returns: {{range .returns}}{{.type}}, {{end}}{{end}}

Describe what it does, when to use it, and any important behavior.
Keep it concise (1-2 sentences).`,
		[]string{"name", "signature", "parameters", "returns"})

	prompt, err := template.Format(map[string]any{
		"name":       fn.Name,
		"signature":  fn.Signature,
		"parameters": fn.Parameters,
		"returns":    fn.Returns,
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

// enhanceTypeDescription generates an enhanced description for a type.
func (g *Generator) enhanceTypeDescription(ctx context.Context, typ *analyzer.TypeInfo) (string, error) {
	template := prompts.NewPromptTemplate(`
Write a clear description for this Go type:

Type: {{.name}} ({{.kind}})
{{if .fields}}Fields: {{range .fields}}{{.name}} {{.type}}, {{end}}{{end}}
{{if .methods}}Methods: {{range .methods}}{{.}}, {{end}}{{end}}

Describe what it represents and how it's used.
Keep it concise (1-2 sentences).`,
		[]string{"name", "kind", "fields", "methods"})

	prompt, err := template.Format(map[string]any{
		"name":    typ.Name,
		"kind":    typ.Kind,
		"fields":  typ.Fields,
		"methods": typ.Methods,
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

// generateContent is a helper method to generate content using the LLM.
func (g *Generator) generateContent(ctx context.Context, prompt string) (string, error) {
	response, err := g.llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	})
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return response.Choices[0].Content, nil
}
