package docgen

import (
	"context"
	"fmt"
	"github.com/docaura/docaura-cli/pkg/analyzer"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"os"
)

// Generator generates documentation from Go package information using AI enhancement.
type Generator struct {
	llm       llms.Model
	templates *TemplateManager
}

// New creates a new documentation generator instance.
func New() (*Generator, error) {
	llm, err := openai.New(
		openai.WithModel("llama3-8b-8192"),
		openai.WithToken(os.Getenv("GROQ_API_KEY")),
		openai.WithBaseURL("https://api.groq.com/openai/v1"))
	if err != nil {
		return nil, fmt.Errorf("create LLM client: %w", err)
	}

	templates, err := NewTemplateManager()
	if err != nil {
		return nil, fmt.Errorf("create template manager: %w", err)
	}

	return &Generator{
		llm:       llm,
		templates: templates,
	}, nil
}

// NewWithLLM creates a new documentation generator with a custom LLM.
func NewWithLLM(llm llms.Model) (*Generator, error) {
	templates, err := NewTemplateManager()
	if err != nil {
		return nil, fmt.Errorf("create template manager: %w", err)
	}

	return &Generator{
		llm:       llm,
		templates: templates,
	}, nil
}

// GeneratePackageDoc generates documentation for a Go package.
func (g *Generator) GeneratePackageDoc(ctx context.Context, pkg *analyzer.PackageInfo, config Config) (string, error) {
	if err := config.Validate(); err != nil {
		return "", fmt.Errorf("invalid config: %w", err)
	}

	// Create a copy to avoid modifying the original
	enhancedPkg := *pkg

	// Enhance descriptions with AI
	if err := g.enhanceDescriptions(ctx, &enhancedPkg); err != nil {
		return "", fmt.Errorf("enhance descriptions: %w", err)
	}

	// Generate usage examples if requested
	if config.GenerateExamples {
		if err := g.generateExamples(ctx, &enhancedPkg); err != nil {
			return "", fmt.Errorf("generate examples: %w", err)
		}
	}

	// Apply template based on style
	result, err := g.templates.Execute(config.Style, &enhancedPkg)
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return result, nil
}

// enhanceDescriptions enhances package descriptions using AI.
func (g *Generator) enhanceDescriptions(ctx context.Context, pkg *analyzer.PackageInfo) error {
	// Enhance package description if empty or too brief
	if len(pkg.Description) < minDescriptionLength {
		if enhanced, err := g.enhancePackageDescription(ctx, pkg); err == nil && enhanced != "" {
			pkg.Description = enhanced
		}
	}

	// Enhance function descriptions
	for i := range pkg.Functions {
		if len(pkg.Functions[i].Description) < minDescriptionLength {
			if enhanced, err := g.enhanceFunctionDescription(ctx, &pkg.Functions[i]); err == nil && enhanced != "" {
				pkg.Functions[i].Description = enhanced
			}
		}
	}

	// Enhance type descriptions
	for i := range pkg.Types {
		if len(pkg.Types[i].Description) < minDescriptionLength {
			if enhanced, err := g.enhanceTypeDescription(ctx, &pkg.Types[i]); err == nil && enhanced != "" {
				pkg.Types[i].Description = enhanced
			}
		}
	}

	return nil
}

// generateExamples generates code examples using AI.
func (g *Generator) generateExamples(ctx context.Context, pkg *analyzer.PackageInfo) error {
	// Generate package-level usage example
	if len(pkg.Examples) == 0 {
		if example, err := g.generatePackageExample(ctx, pkg); err == nil && example != "" {
			pkg.Examples = append(pkg.Examples, analyzer.ExampleInfo{
				Name: "Basic Usage",
				Code: example,
				Doc:  "Basic usage example",
			})
		}
	}

	// Generate function examples
	for i := range pkg.Functions {
		fn := &pkg.Functions[i]
		if len(fn.Examples) == 0 && fn.IsExported {
			if example, err := g.generateFunctionExample(ctx, fn, pkg); err == nil && example != "" {
				fn.Examples = append(fn.Examples, example)
			}
		}
	}

	return nil
}

// Constants for description enhancement
const (
	minDescriptionLength = 20
	maxDescriptionLength = 500
)
