package models

// よく使用されるAIモデルの定義（価格情報を含む）
var (
	// OpenAI Models
	GPT4oMini = &AIModel{
		ID:                   "gpt-4o-mini",
		Name:                 "GPT-4o-mini",
		Provider:             ProviderOpenAI,
		ModelType:            ModelTypeChat,
		ContextSize:          128000,
		MaxTokens:            4096,
		PromptPricePer1M:     0.15, // $0.15 per 1M tokens
		CompletionPricePer1M: 0.60, // $0.60 per 1M tokens
		EmbeddingPricePer1M:  0.10, // $0.10 per 1M tokens
		Currency:             "USD",
	}

	GPT4o = &AIModel{
		ID:                   "gpt-4o",
		Name:                 "GPT-4o",
		Provider:             ProviderOpenAI,
		ModelType:            ModelTypeChat,
		ContextSize:          128000,
		MaxTokens:            4096,
		PromptPricePer1M:     2.50,  // $2.50 per 1M tokens
		CompletionPricePer1M: 10.00, // $10.00 per 1M tokens
		EmbeddingPricePer1M:  0.10,  // $0.10 per 1M tokens
		Currency:             "USD",
	}

	GPT4Turbo = &AIModel{
		ID:                   "gpt-4-turbo",
		Name:                 "GPT-4 Turbo",
		Provider:             ProviderOpenAI,
		ModelType:            ModelTypeChat,
		ContextSize:          128000,
		MaxTokens:            4096,
		PromptPricePer1M:     10.00, // $10.00 per 1M tokens
		CompletionPricePer1M: 30.00, // $30.00 per 1M tokens
		EmbeddingPricePer1M:  0.10,  // $0.10 per 1M tokens
		Currency:             "USD",
	}

	GPT35Turbo = &AIModel{
		ID:                   "gpt-3.5-turbo",
		Name:                 "GPT-3.5-turbo",
		Provider:             ProviderOpenAI,
		ModelType:            ModelTypeChat,
		ContextSize:          16385,
		MaxTokens:            4096,
		PromptPricePer1M:     0.50, // $0.50 per 1M tokens
		CompletionPricePer1M: 1.50, // $1.50 per 1M tokens
		EmbeddingPricePer1M:  0.10, // $0.10 per 1M tokens
		Currency:             "USD",
	}

	// Anthropic Models
	Claude3Haiku = &AIModel{
		ID:                   "claude-3-haiku-20240307",
		Name:                 "Claude 3 Haiku",
		Provider:             ProviderAnthropic,
		ModelType:            ModelTypeChat,
		ContextSize:          200000,
		MaxTokens:            4096,
		PromptPricePer1M:     0.25, // $0.25 per 1M tokens
		CompletionPricePer1M: 1.25, // $1.25 per 1M tokens
		EmbeddingPricePer1M:  0.10, // $0.10 per 1M tokens
		Currency:             "USD",
	}

	Claude3Sonnet = &AIModel{
		ID:                   "claude-3-sonnet-20240229",
		Name:                 "Claude 3 Sonnet",
		Provider:             ProviderAnthropic,
		ModelType:            ModelTypeChat,
		ContextSize:          200000,
		MaxTokens:            4096,
		PromptPricePer1M:     3.00,  // $3.00 per 1M tokens
		CompletionPricePer1M: 15.00, // $15.00 per 1M tokens
		EmbeddingPricePer1M:  0.10,  // $0.10 per 1M tokens
		Currency:             "USD",
	}

	Claude3Opus = &AIModel{
		ID:                   "claude-3-opus-20240229",
		Name:                 "Claude 3 Opus",
		Provider:             ProviderAnthropic,
		ModelType:            ModelTypeChat,
		ContextSize:          200000,
		MaxTokens:            4096,
		PromptPricePer1M:     15.00, // $15.00 per 1M tokens
		CompletionPricePer1M: 75.00, // $75.00 per 1M tokens
		EmbeddingPricePer1M:  0.10,  // $0.10 per 1M tokens
		Currency:             "USD",
	}
)

// チャット補完の設定に関する定数
const (
	// デフォルトの最大トークン数
	DefaultMaxTokens = 2000

	// デフォルトの温度設定（創造性の調整）
	DefaultTemperature = 0.7

	// 高品質回答用の設定
	HighQualityMaxTokens   = 4000
	HighQualityTemperature = 0.5

	// 創造的回答用の設定
	CreativeMaxTokens   = 3000
	CreativeTemperature = 0.9
)
