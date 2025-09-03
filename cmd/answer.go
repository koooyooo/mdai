/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/koooyooo/mdai/controller"
	"github.com/koooyooo/mdai/util/file"
	"github.com/openai/openai-go"
	"github.com/spf13/cobra"
)

// answerCmd represents the answer command
var answerCmd = &cobra.Command{
	Use:   "answer",
	Short: "A brief description of your command",
	Long:  ``,
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
	controller := controller.NewOpenAIController(os.Getenv("OPENAI_API_KEY"), logger)

	if len(args) == 0 {
		return fmt.Errorf("path is required")
	}
	path := args[0]
	sysMsg := createSystemMessage()
	userMsg, err := createUserMessage(path)
	if err != nil {
		return fmt.Errorf("fail in creating user message: %v", err)
	}
	controller.Control(sysMsg, userMsg, func(completion *openai.ChatCompletion) error {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		answer := completion.Choices[0].Message.Content
		appendText := fmt.Sprintf("\n\n%s\n", answer)
		if _, err := f.WriteString(appendText); err != nil {
			return err
		}
		return nil
	})
	return nil
}

func createSystemMessage() string {
	return `You are a helpful and detailed assistant. When answering questions based on the given context, please follow these guidelines:

1. Answer in the same language as the question
2. Make full use of the context information
3. Provide specific and practical information
4. Add examples and explanations when necessary
5. Ensure answers are appropriately long and content-rich
6. Provide insights that deepen the questioner's understanding
7. Prefer rich markdown formatting
`
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
