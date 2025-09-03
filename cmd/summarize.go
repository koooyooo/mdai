/*
Copyright © 2025 koooyooo
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/koooyooo/mdai/config"
	"github.com/koooyooo/mdai/controller"
	"github.com/openai/openai-go"
	"github.com/spf13/cobra"
)

// summarizeCmd represents the summarize command
var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Summarize the content of a markdown file",
	Long: `Summarize the content of a markdown file using AI.
The summarized content will be saved to a new file with "_sum" suffix.
For example, if the input file is "document.md", the output will be "document_sum.md".`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		if err := summarize(args, logger); err != nil {
			logger.Error("fail in calling summarize", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(summarizeCmd)
}

func summarize(args []string, logger *slog.Logger) error {
	controller := controller.NewOpenAIController(os.Getenv("OPENAI_API_KEY"), logger)

	if len(args) == 0 {
		return fmt.Errorf("path is required")
	}
	path := args[0]

	// ファイルの存在確認
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", path)
	}

	// ファイル拡張子の確認
	if !strings.HasSuffix(strings.ToLower(path), ".md") {
		return fmt.Errorf("file must have .md extension: %s", path)
	}

	// 出力ファイル名の生成
	outputPath := generateOutputPath(path)

	// ファイル内容の読み込み
	content, err := loadContent(path)
	if err != nil {
		return fmt.Errorf("fail in loading content: %v", err)
	}

	// 設定ファイルを読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Warn("failed to load config, using defaults", "error", err)
		// エラーが発生した場合はデフォルト設定を使用
		cfg = config.GetDefaultConfig()
	}

	// 設定ファイルからシステムメッセージを取得
	sysMsg := cfg.Summarize.SystemMessage
	if sysMsg == "" {
		sysMsg = createSummarizeSystemMessage()
	}

	userMsg := createSummarizeUserMessage(content)

	// 設定ファイルから品質設定を取得
	maxTokens := cfg.Default.Quality.MaxTokens
	temperature := cfg.Default.Quality.Temperature

	// システムメッセージに文字数の指示を追加
	if cfg.Summarize.TargetLength > 0 {
		sysMsg += fmt.Sprintf("\n\n**Summary Length Guidance**: Please provide a summary of approximately %d characters.", cfg.Summarize.TargetLength)
	}

	// 設定値をログに出力
	logger.Info("using configuration",
		"maxTokens", maxTokens,
		"temperature", temperature,
		"targetLength", cfg.Summarize.TargetLength)

	controller.Control(sysMsg, userMsg, cfg.Default.Quality, func(completion *openai.ChatCompletion) error {
		summary := completion.Choices[0].Message.Content

		// 要約結果をファイルに保存
		if err := saveSummary(outputPath, summary, path); err != nil {
			return fmt.Errorf("fail in saving summary: %v", err)
		}

		logger.Info("summary created successfully", "input", path, "output", outputPath)
		return nil
	})

	return nil
}

func generateOutputPath(inputPath string) string {
	dir := filepath.Dir(inputPath)
	filename := filepath.Base(inputPath)
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	return filepath.Join(dir, nameWithoutExt+"_sum.md")
}

func createSummarizeSystemMessage() string {
	return `You are a helpful and detailed assistant specialized in summarizing markdown documents. When summarizing content, please follow these guidelines:

1. Provide a comprehensive yet concise summary of the main content
2. Maintain the key points and important information
3. Use clear and organized structure with markdown formatting
4. Include main headings and subheadings when relevant
5. Preserve important details, examples, and references
6. Make the summary easy to read and understand
7. Use appropriate markdown elements (headers, lists, emphasis, etc.)
8. Keep the summary appropriately long - not too brief, not too verbose
9. Focus on the most valuable and actionable information
10. Maintain the original tone and style when appropriate
`
}

func createSummarizeUserMessage(content string) string {
	return fmt.Sprintf(`Please provide a comprehensive summary of the following markdown content:

%s

Please create a well-structured summary that captures the essence and key points of this content.`, content)
}

func saveSummary(outputPath, summary, originalPath string) error {
	// ファイルに書き込み
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// 要約内容を直接書き込み
	if _, err := f.WriteString(summary); err != nil {
		return err
	}

	return nil
}
