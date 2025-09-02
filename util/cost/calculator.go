package cost

import (
	"fmt"

	"github.com/openai/openai-go"
)

// TokenUsage は OpenAI API のトークン使用量を表す構造体
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// CalculateCost は指定されたモデルとトークン使用量に基づいてコストを計算します
func CalculateCost(model string, promptPricePerM float64, completionPricePerM float64, usage openai.CompletionUsage) (string, error) {
	promptCost := float64(usage.PromptTokens) / 1000_000.0 * promptPricePerM
	completionCost := float64(usage.CompletionTokens) / 1000_000.0 * completionPricePerM
	totalCost := promptCost + completionCost

	logStr := fmt.Sprintf("$%.5f (入力: $%.5f, 出力: $%.5f)", totalCost, promptCost, completionCost)
	return logStr, nil
}
