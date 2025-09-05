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
	"github.com/openai/openai-go/option"
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
		cfg := config.GetInstance().GetConfig()
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: cfg.Default.GetLogLevel().Level(),
		}))
		if err := answer(cfg, args, logger); err != nil {
			logger.Error("fail in calling answer", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(answerCmd)
}

func answer(cfg config.Config, args []string, logger *slog.Logger) error {
	client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	controller := controller.NewOpenAIController(&client, cfg.Default.Model, logger)

	if len(args) == 0 {
		return fmt.Errorf("path is required")
	}
	path := args[0]

	// Get system message from configuration
	sysMsg := cfg.Answer.SystemMessage

	// Load file content
	content, err := loadContent(path)
	if err != nil {
		return fmt.Errorf("fail in loading content: %v", err)
	}

	// Get last quote and other content
	lastQuote, otherContents, err := file.LoadLastQuote(content)
	if err != nil {
		return fmt.Errorf("fail in loading last quote: %v", err)
	}

	// Get user message from configuration and apply template processing
	userMsg, err := cfg.Answer.UserMessage.Apply(map[string]string{
		"Question": lastQuote,
		"Context":  otherContents,
	})
	if err != nil {
		return fmt.Errorf("fail in creating user message: %v", err)
	}

	// Get quality settings from configuration
	maxTokens := cfg.Default.Quality.MaxTokens
	temperature := cfg.Default.Quality.Temperature

	// Add character count instruction to system message
	if cfg.Answer.TargetLength > 0 {
		sysMsg += fmt.Sprintf("\n\n**Answer Length Guidance**: Please provide an answer of approximately %d characters.", cfg.Answer.TargetLength)
	}

	// Log configuration values
	logger.Info("using configuration",
		"maxTokens", maxTokens,
		"temperature", temperature,
		"targetLength", cfg.Answer.TargetLength)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	defer func() { _ = f.Close() }()

	f.WriteString("\n\n")
	controller.ControlStreaming(sysMsg, userMsg, cfg.Default.Quality, func(chunk openai.ChatCompletionChunk) error {
		answer := chunk.Choices[0].Delta.Content
		if err != nil {
			return err
		}
		appendText := answer
		if _, err := f.WriteString(appendText); err != nil {
			return err
		}
		return nil
	})
	return nil
}

func loadContent(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
