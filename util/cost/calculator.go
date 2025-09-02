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
func CalculateCost(model string, promptPricePerM float64, completionPricePerM float64, usage openai.Usage) (string, error) {
	promptCost := float64(usage.PromptTokens) / 1000_000.0 * promptPricePerM
	completionCost := float64(usage.CompletionTokens) / 1000_000.0 * completionPricePerM
	totalCost := promptCost + completionCost

	result := fmt.Sprintf("---- 費用情報 ----\n")
	result += fmt.Sprintf("モデル: %s\n", model)
	result += fmt.Sprintf("Prompt tokens: %d\n", usage.PromptTokens)
	result += fmt.Sprintf("Completion tokens: %d\n", usage.CompletionTokens)
	result += fmt.Sprintf("Total tokens: %d\n", usage.TotalTokens)
	result += fmt.Sprintf("推定費用: $%.5f (入力: $%.5f, 出力: $%.5f)\n", totalCost, promptCost, completionCost)
	result += "-----------------"

	return result, nil
}
