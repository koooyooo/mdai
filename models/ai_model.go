package models

// ModelType はAIモデルのタイプを表す定数
type ModelType string

const (
	ModelTypeChat       ModelType = "chat"
	ModelTypeCompletion ModelType = "completion"
	ModelTypeEmbedding  ModelType = "embedding"
)

// Provider はAIプロバイダーを表す定数
type Provider string

const (
	ProviderOpenAI    Provider = "OpenAI"
	ProviderAnthropic Provider = "Anthropic"
	ProviderGoogle    Provider = "Google"
)

// AIModel はAIモデルの基本情報と価格設定を表す構造体
type AIModel struct {
	ID                   string    // モデルの一意識別子
	Name                 string    // モデル名
	Provider             Provider  // プロバイダー
	ModelType            ModelType // モデルタイプ
	ContextSize          int       // コンテキストサイズ
	MaxTokens            int       // 最大トークン数
	PromptPricePer1M     float64   // 入力トークン100万個あたりの価格
	CompletionPricePer1M float64   // 出力トークン100万個あたりの価格
	EmbeddingPricePer1M  float64   // 埋め込みトークン100万個あたりの価格
	Currency             string    // 通貨（USD、JPY等）
}

// String はモデルの文字列表現を返します
func (m *AIModel) String() string {
	return m.Name
}

// CalculatePromptCost は入力トークン数に基づいてコストを計算します
func (m *AIModel) CalculatePromptCost(tokenCount int) float64 {
	return float64(tokenCount) / 1_000_000.0 * m.PromptPricePer1M
}

// CalculateCompletionCost は出力トークン数に基づいてコストを計算します
func (m *AIModel) CalculateCompletionCost(tokenCount int) float64 {
	return float64(tokenCount) / 1_000_000.0 * m.CompletionPricePer1M
}

// CalculateEmbeddingCost は埋め込みトークン数に基づいてコストを計算します
func (m *AIModel) CalculateEmbeddingCost(tokenCount int) float64 {
	return float64(tokenCount) / 1_000_000.0 * m.EmbeddingPricePer1M
}

// CalculateTotalCost は入力・出力トークン数に基づいて総コストを計算します
func (m *AIModel) CalculateTotalCost(promptTokens, completionTokens int) float64 {
	return m.CalculatePromptCost(promptTokens) + m.CalculateCompletionCost(completionTokens)
}
