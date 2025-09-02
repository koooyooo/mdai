/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/koooyooo/mdai/util/cost"
	"github.com/koooyooo/mdai/util/file"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/spf13/cobra"
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ask(args)
	},
}

func init() {
	rootCmd.AddCommand(askCmd)
}

func ask(args []string) error {
	path := args[0]
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// 1. マークダウンファイルから最後の引用部分を抽出
	lastQuote, otherContents, err := file.LoadLastQuote(string(b))
	if err != nil {
		return err
	}
	// 2. OpenAI API キー取得
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("環境変数 OPENAI_API_KEY が設定されていません")
	}

	// 3. OpenAI API へ問い合わせ（公式SDKを使用）
	client := openai.NewClient(option.WithAPIKey(apiKey))

	message := fmt.Sprintf("Context: %s\n\nQuestion: %s", otherContents, lastQuote)

	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a helpful assistant. You are given a question and a context. You need to answer the question based on the context. You need to answer the question in the same language as the question."),
			openai.UserMessage(message),
		},
	})
	if err != nil {
		return fmt.Errorf("OpenAI API エラー: %v", err)
	}
	if len(resp.Choices) == 0 {
		return fmt.Errorf("OpenAI API から応答がありません")
	}
	// GPT-4-1106-preview (gpt-4-1-mini) の料金を使用
	// 2024年6月時点: 入力 1K tokens あたり $0.6, 出力 1K tokens あたり $2.4
	costInfo, err := cost.CalculateCost(string(openai.ChatModelGPT4oMini), 0.6, 2.4, resp.Usage)
	if err != nil {
		return fmt.Errorf("コスト計算エラー: %v", err)
	}
	fmt.Println(costInfo)
	answer := resp.Choices[0].Message.Content

	// 4. マークダウンファイルの末尾に結果を追記
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	appendText := fmt.Sprintf("\n\n%s\n", answer)
	if _, err := f.WriteString(appendText); err != nil {
		return err
	}

	fmt.Println("OpenAI の応答をファイル末尾に追記しました。")

	return nil
}
