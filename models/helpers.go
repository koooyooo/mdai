package models

import (
	"fmt"

	"github.com/openai/openai-go"
)

// GetModelByID はモデルIDからモデル情報を取得します
func GetModelByID(modelID string) (*AIModel, error) {
	switch modelID {
	case "gpt-4o-mini":
		return GPT4oMini, nil
	case "gpt-4o":
		return GPT4o, nil
	case "gpt-3.5-turbo":
		return GPT35Turbo, nil
	case "claude-3-haiku-20240307":
		return Claude3Haiku, nil
	case "claude-3-sonnet-20240229":
		return Claude3Sonnet, nil
	case "claude-3-opus-20240229":
		return Claude3Opus, nil
	default:
		return nil, fmt.Errorf("model not found: %s", modelID)
	}
}

// CalculateCost は指定されたモデルのトークン使用量に基づいてコストを計算します
func CalculateCost(modelID string, promptTokens, completionTokens int) (float64, error) {
	model, err := GetModelByID(modelID)
	if err != nil {
		return 0, err
	}
	return model.CalculateTotalCost(promptTokens, completionTokens), nil
}

// CalculateCostString は既存のutil/costパッケージと互換性のある形式でコストを返します
func CalculateCostString(modelID string, usage openai.CompletionUsage) (string, error) {
	model, err := GetModelByID(modelID)
	if err != nil {
		return "", err
	}

	promptCost := model.CalculatePromptCost(int(usage.PromptTokens))
	completionCost := model.CalculateCompletionCost(int(usage.CompletionTokens))
	totalCost := promptCost + completionCost

	return fmt.Sprintf("[%s] $%.5f (入力: $%.5f, 出力: $%.5f)", modelID, totalCost, promptCost, completionCost), nil
}

// ListModels は利用可能なモデルの一覧を返します
func ListModels() []*AIModel {
	return []*AIModel{
		GPT4oMini,
		GPT4o,
		GPT35Turbo,
		Claude3Haiku,
		Claude3Sonnet,
		Claude3Opus,
	}
}
