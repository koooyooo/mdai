package models

import (
	"fmt"
	"log"
)

// ExampleSimpleUsage はシンプルな使用例を示します
func ExampleSimpleUsage() {
	fmt.Println("=== 利用可能なモデル ===")
	models := ListModels()
	for _, model := range models {
		fmt.Printf("- %s (%s)\n", model.Name, model.Provider)
	}

	fmt.Println("\n=== コスト計算の例 ===")

	// GPT-4o-miniで1000入力トークン、500出力トークンを使用した場合
	cost, err := CalculateCost("gpt-4o-mini", 1000, 500)
	if err != nil {
		log.Printf("Cost calculation error: %v", err)
	} else {
		fmt.Printf("GPT-4o-mini (1000入力 + 500出力): $%.6f\n", cost)
	}

	// GPT-4oで1000入力トークン、500出力トークンを使用した場合
	cost, err = CalculateCost("gpt-4o", 1000, 500)
	if err != nil {
		log.Printf("Cost calculation error: %v", err)
	} else {
		fmt.Printf("GPT-4o (1000入力 + 500出力): $%.6f\n", cost)
	}

	// Claude 3 Haikuで1000入力トークン、500出力トークンを使用した場合
	cost, err = CalculateCost("claude-3-haiku-20240307", 1000, 500)
	if err != nil {
		log.Printf("Cost calculation error: %v", err)
	} else {
		fmt.Printf("Claude 3 Haiku (1000入力 + 500出力): $%.6f\n", cost)
	}

	fmt.Println("\n=== 価格情報 ===")
	model, err := GetModelByID("gpt-4o-mini")
	if err == nil {
		fmt.Printf("GPT-4o-mini 入力価格: $%.2f per 1M tokens\n", model.PromptPricePer1M)
		fmt.Printf("GPT-4o-mini 出力価格: $%.2f per 1M tokens\n", model.CompletionPricePer1M)
		fmt.Printf("通貨: %s\n", model.Currency)
	}
}
