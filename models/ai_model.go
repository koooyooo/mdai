/*
Copyright © 2025 koooyooo
*/
package models

// ModelType represents the type of AI model
type ModelType string

const (
	ModelTypeChat       ModelType = "chat"
	ModelTypeCompletion ModelType = "completion"
	ModelTypeEmbedding  ModelType = "embedding"
)

// Provider represents the AI provider
type Provider string

const (
	ProviderOpenAI    Provider = "OpenAI"
	ProviderAnthropic Provider = "Anthropic"
	ProviderGoogle    Provider = "Google"
)

// AIModel represents the basic information and pricing of an AI model
type AIModel struct {
	ID                   string    // Unique identifier for the model
	Name                 string    // Model name
	Provider             Provider  // Provider
	ModelType            ModelType // Model type
	ContextSize          int       // Context size
	MaxTokens            int       // Maximum token count
	PromptPricePer1M     float64   // Price per 1M input tokens
	CompletionPricePer1M float64   // Price per 1M output tokens
	EmbeddingPricePer1M  float64   // Price per 1M embedding tokens
	Currency             string    // Currency (USD, JPY, etc.)
}

// String returns the string representation of the model
func (m *AIModel) String() string {
	return m.Name
}

// CalculatePromptCost calculates cost based on input token count
func (m *AIModel) CalculatePromptCost(tokenCount int) float64 {
	return float64(tokenCount) / 1_000_000.0 * m.PromptPricePer1M
}

// CalculateCompletionCost calculates cost based on output token count
func (m *AIModel) CalculateCompletionCost(tokenCount int) float64 {
	return float64(tokenCount) / 1_000_000.0 * m.CompletionPricePer1M
}

// CalculateEmbeddingCost calculates cost based on embedding token count
func (m *AIModel) CalculateEmbeddingCost(tokenCount int) float64 {
	return float64(tokenCount) / 1_000_000.0 * m.EmbeddingPricePer1M
}

// CalculateTotalCost calculates total cost based on input and output token counts
func (m *AIModel) CalculateTotalCost(promptTokens, completionTokens int) float64 {
	return m.CalculatePromptCost(promptTokens) + m.CalculateCompletionCost(completionTokens)
}
