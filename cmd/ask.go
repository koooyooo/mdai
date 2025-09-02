/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/koooyooo/mdai/models"
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
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		if err := ask(args, logger); err != nil {
			logger.Error("fail in calling ask", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(askCmd)
}

const DefaultSystemMessage = `You are a helpful and detailed assistant. When answering questions based on the given context, please follow these guidelines:

1. Answer in the same language as the question
2. Make full use of the context information
3. Provide specific and practical information
4. Add examples and explanations when necessary
5. Ensure answers are appropriately long and content-rich
6. Provide insights that deepen the questioner's understanding
7. Prefer rich markdown formatting
`

func ask(args []string, logger *slog.Logger) error {
	path := ""
	if len(args) != 0 {
		path = args[0]
	}
	if path == "" {
		return fmt.Errorf("file path is required")
	}
	content, err := loadContent(path)
	if err != nil {
		return fmt.Errorf("fail in loading content: %v", err)
	}
	lastQuote, otherContents, err := file.LoadLastQuote(content)
	if err != nil {
		return err
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("environment variable OPENAI_API_KEY is not set")
	}

	var modelID = openai.ChatModelGPT4oMini
	client := openai.NewClient(option.WithAPIKey(apiKey))
	message := fmt.Sprintf("Context: %s\n\nQuestion: %s", otherContents, lastQuote)
	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: modelID,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(DefaultSystemMessage),
			openai.UserMessage(message),
		},
		MaxTokens:   openai.Int(models.DefaultMaxTokens),
		Temperature: openai.Float(models.DefaultTemperature),
	})
	if err != nil {
		return fmt.Errorf("OpenAI API error: %v", err)
	}
	if len(resp.Choices) == 0 {
		return fmt.Errorf("no response from OpenAI API")
	}

	costInfo, err := models.CalculateCostString(modelID, resp.Usage)
	if err != nil {
		return fmt.Errorf("cost calculation error: %v", err)
	}
	logger.Info("cost information", "costInfo", costInfo)
	answer := resp.Choices[0].Message.Content

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
}

func loadContent(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
