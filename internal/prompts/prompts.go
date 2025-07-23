package prompts

// PromptConfig represents configuration for prompt generation.
type PromptConfig struct {
	MaxLength    int
	Style        string
	IncludeTypes bool
	IncludeUsage bool
}

// DefaultConfig returns a default prompt configuration.
func DefaultConfig() PromptConfig {
	return PromptConfig{
		MaxLength:    200,
		Style:        "professional",
		IncludeTypes: true,
		IncludeUsage: true,
	}
}
