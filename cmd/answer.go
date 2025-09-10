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

var nonStream bool

func init() {
	answerCmd.Flags().BoolVarP(&nonStream, "nonstream", "n", false, "Disable streaming output (enables cost calculation)")
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
	answerConfig := cfg.GetAnswerConfig("default")
	sysMsg := answerConfig.SystemMessage

	// Load file content
	content, err := file.LoadContent(path)
	if err != nil {
		return fmt.Errorf("fail in loading content: %v", err)
	}

	// Get last quote and other content
	lastQuote, otherContents, err := file.LoadLastQuote(content)
	if err != nil {
		return fmt.Errorf("fail in loading last quote: %v", err)
	}

	// Get user message from configuration and apply template processing
	userMsg, err := answerConfig.UserMessage.Apply(map[string]string{
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
	if answerConfig.TargetLength > 0 {
		sysMsg += fmt.Sprintf("\n\n**Answer Length Guidance**: Please provide an answer of approximately %d characters.", answerConfig.TargetLength)
	}

	// Log configuration values
	logger.Info("using configuration",
		"maxTokens", maxTokens,
		"temperature", temperature,
		"targetLength", answerConfig.TargetLength)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer func() { _ = f.Close() }()

	if _, err := f.WriteString("\n\n"); err != nil {
		return fmt.Errorf("failed to write newlines: %v", err)
	}

	// Disable streaming if either the nonstream flag is set or it's disabled in config
	if nonStream || cfg.Default.DisableStream {
		// Non-streaming mode with cost calculation
		return controller.Control(sysMsg, userMsg, cfg.Default.Quality, func(completion *openai.ChatCompletion) error {
			answer := completion.Choices[0].Message.Content
			if _, err := f.WriteString(answer); err != nil {
				return fmt.Errorf("failed to write answer: %v", err)
			}
			return nil
		})
	}

	// Streaming mode
	return controller.ControlStreaming(sysMsg, userMsg, cfg.Default.Quality, func(chunk openai.ChatCompletionChunk) error {
		answer := chunk.Choices[0].Delta.Content
		if _, err := f.WriteString(answer); err != nil {
			return fmt.Errorf("failed to write chunk: %v", err)
		}
		return nil
	})
}
