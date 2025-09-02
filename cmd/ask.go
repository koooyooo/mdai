/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"log/slog"
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

func ask(args []string, logger *slog.Logger) error {
	path := ""
	if len(args) != 0 {
		path = args[0]
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
		return fmt.Errorf("OpenAI API error: %v", err)
	}
	if len(resp.Choices) == 0 {
		return fmt.Errorf("no response from OpenAI API")
	}

	costInfo, err := cost.CalculateCost(string(openai.ChatModelGPT4oMini), 0.6, 2.4, resp.Usage)
	if err != nil {
		return fmt.Errorf("cost calculation error: %v", err)
	}
	logger.Info("cost information", "costInfo", costInfo)
	answer := resp.Choices[0].Message.Content

	if path != "" {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		appendText := fmt.Sprintf("\n\n%s\n", answer)
		if _, err := f.WriteString(appendText); err != nil {
			return err
		}
	} else {
		fmt.Println(answer)
	}
	return nil
}

func loadContent(path string) (string, error) {
	if path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(b), nil
	} else {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("error reading from standard input: %v", err)
		}
		return string(b), nil
	}
}
