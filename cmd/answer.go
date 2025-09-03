/*
Copyright © 2025 koooyooo
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/koooyooo/mdai/config"
	"github.com/koooyooo/mdai/controller"
	"github.com/koooyooo/mdai/util/file"
	"github.com/openai/openai-go"
	"github.com/spf13/cobra"
)

// answerCmd represents the answer command
var answerCmd = &cobra.Command{
	Use:   "answer",
	Short: "Answer the question based on the content of a markdown file",
	Long: `Answer the question based on the content of a markdown file.
	The question will be extracted from the last quote in the file.
	The answer will be appended to the end of the file.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		if err := answer(args, logger); err != nil {
			logger.Error("fail in calling answer", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(answerCmd)
}

func answer(args []string, logger *slog.Logger) error {
	// 設定ファイルを読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Warn("failed to load config, using defaults", "error", err)
		// エラーが発生した場合はデフォルト設定を使用
		cfg = config.GetDefaultConfig()
	}

	controller := controller.NewOpenAIController(os.Getenv("OPENAI_API_KEY"), logger)

	if len(args) == 0 {
		return fmt.Errorf("path is required")
	}
	path := args[0]

	// 設定ファイルからシステムメッセージを取得
	sysMsg := cfg.Answer.SystemMessage

	userMsg, err := createUserMessage(path)
	if err != nil {
		return fmt.Errorf("fail in creating user message: %v", err)
	}

	// 設定ファイルから品質設定を取得
	maxTokens := cfg.Default.Quality.MaxTokens
	temperature := cfg.Default.Quality.Temperature

	// システムメッセージに文字数の指示を追加
	if cfg.Answer.TargetChars > 0 {
		sysMsg += fmt.Sprintf("\n\n**Answer Length Guidance**: Please provide an answer of approximately %d characters.", cfg.Answer.TargetChars)
	}

	// 設定値をログに出力
	logger.Info("using configuration",
		"maxTokens", maxTokens,
		"temperature", temperature,
		"targetChars", cfg.Answer.TargetChars)

	controller.Control(sysMsg, userMsg, cfg.Default.Quality, func(completion *openai.ChatCompletion) error {
		answer := completion.Choices[0].Message.Content
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		appendText := fmt.Sprintf("\n\n%s\n", answer)
		if _, err := f.WriteString(appendText); err != nil {
			return err
		}
		return nil
	})
	return nil
}

func createUserMessage(path string) (string, error) {
	content, err := loadContent(path)
	if err != nil {
		return "", fmt.Errorf("fail in loading content: %v", err)
	}
	lastQuote, otherContents, err := file.LoadLastQuote(content)
	if err != nil {
		return "", err
	}
	userMsg := fmt.Sprintf("Context: %s\n\nQuestion: %s", otherContents, lastQuote)
	return userMsg, nil
}

func loadContent(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
