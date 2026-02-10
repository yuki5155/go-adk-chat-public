package chat

// Model represents an available LLM model
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AvailableModels returns the list of available LLM models
var AvailableModels = []Model{
	{
		ID:          "gemini-2.0-flash",
		Name:        "Gemini 2.0 Flash",
		Description: "Fast and efficient",
	},
	{
		ID:          "gemini-2.5-flash",
		Name:        "Gemini 2.5 Flash",
		Description: "Latest flash model",
	},
	{
		ID:          "gemini-2.5-pro",
		Name:        "Gemini 2.5 Pro",
		Description: "Most capable",
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
