package chat

// Model represents an available LLM model
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Provider    string `json:"provider"`
}

// AvailableModels returns the list of available LLM models
var AvailableModels = []Model{
	// Gemini models
	{
		ID:          "gemini-2.0-flash",
		Name:        "Gemini 2.0 Flash",
		Description: "Fast and efficient",
		Provider:    "gemini",
	},
	{
		ID:          "gemini-2.5-flash",
		Name:        "Gemini 2.5 Flash",
		Description: "Latest flash model",
		Provider:    "gemini",
	},
	{
		ID:          "gemini-2.5-pro",
		Name:        "Gemini 2.5 Pro",
		Description: "Most capable Gemini model",
		Provider:    "gemini",
	},
	// OpenAI models
	{
		ID:          "gpt-4o",
		Name:        "GPT-4o",
		Description: "OpenAI's most capable model",
		Provider:    "openai",
	},
	{
		ID:          "gpt-4o-mini",
		Name:        "GPT-4o Mini",
		Description: "Fast and affordable OpenAI model",
		Provider:    "openai",
	},
	// Anthropic models
	{
		ID:          "claude-opus-4-6",
		Name:        "Claude Opus 4.6",
		Description: "Anthropic's most capable model",
		Provider:    "anthropic",
	},
	{
		ID:          "claude-sonnet-4-6",
		Name:        "Claude Sonnet 4.6",
		Description: "Balanced performance and speed",
		Provider:    "anthropic",
	},
	{
		ID:          "claude-haiku-4-5-20251001",
		Name:        "Claude Haiku 4.5",
		Description: "Fast and efficient Anthropic model",
		Provider:    "anthropic",
	},
}

// IsValidModel checks if the given model ID is valid
func IsValidModel(modelID string) bool {
	for _, m := range AvailableModels {
		if m.ID == modelID {
			return true
		}
	}
	return false
}

// GetDefaultModel returns the default model
func GetDefaultModel() Model {
	return AvailableModels[0]
}

// ProviderForModel returns the provider name for a given model ID.
// Returns empty string if the model is not found.
func ProviderForModel(modelID string) string {
	for _, m := range AvailableModels {
		if m.ID == modelID {
			return m.Provider
		}
	}
	return ""
}

// ModelsForProvider returns all models belonging to the given provider.
func ModelsForProvider(provider string) []Model {
	var result []Model
	for _, m := range AvailableModels {
		if m.Provider == provider {
			result = append(result, m)
		}
	}
	return result
}
