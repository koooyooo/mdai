/*
Copyright © 2025 koooyooo
*/
package controller

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/koooyooo/mdai/config"
	"github.com/koooyooo/mdai/util/file"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// AppendConfig holds configuration for a specific append operation
type AppendConfig struct {
	SystemMessage string
	UserMessage   config.UserMessageTemplate
	ExtraArgs     []string
}

// Append performs an append operation on a markdown file
func Append(cfg config.Config, operation string, path string, extraArgs []string, logger *slog.Logger) error {
	// Get operation configuration dynamically
	opConfig, err := getAppendOperationConfig(cfg, operation)
	if err != nil {
		return err
	}

	// Create append configuration
	appendConfig := &AppendConfig{
		SystemMessage: opConfig.SystemMessage,
		UserMessage:   opConfig.UserMessage,
		ExtraArgs:     extraArgs,
	}

	// Execute append operation
	return executeAppend(cfg, appendConfig, path, extraArgs, logger)
}

func getAppendOperationConfig(cfg config.Config, operation string) (config.OperationConfig, error) {
	// Get operation configuration from map
	opConfig, exists := cfg.Append.Operations[operation]
	if !exists {
		return config.OperationConfig{}, fmt.Errorf("unsupported append operation: %s", operation)
	}

	return opConfig, nil
}

func executeAppend(cfg config.Config, appendConfig *AppendConfig, path string, extraArgs []string, logger *slog.Logger) error {
	// Validate file
	if err := validateAppendFile(path); err != nil {
		return err
	}

	// Load file content
	content, err := file.LoadContent(path)
	if err != nil {
		return fmt.Errorf("fail in loading content: %v", err)
	}

	// Prepare messages based on operation type
	sysMsg, userMsg, err := prepareAppendMessages(cfg, appendConfig, content, extraArgs)
	if err != nil {
		return err
	}

	// Execute append operation
	client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	openAIController := NewOpenAIController(&client, cfg.Default.Model, logger)

	// Log configuration values
	logger.Info("using configuration",
		"maxTokens", cfg.Default.Quality.MaxTokens,
		"temperature", cfg.Default.Quality.Temperature)

	// Open file for appending
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer func() { _ = f.Close() }()

	if _, err := f.WriteString("\n\n"); err != nil {
		return fmt.Errorf("failed to write newlines: %v", err)
	}

	// Check if streaming should be disabled
	nonStream := false // This could be passed as a parameter or from config
	if nonStream || cfg.Default.DisableStream {
		// Non-streaming mode with cost calculation
		return openAIController.Control(sysMsg, userMsg, cfg.Default.Quality, func(completion *openai.ChatCompletion) error {
			answer := completion.Choices[0].Message.Content
			if _, err := f.WriteString(answer); err != nil {
				return fmt.Errorf("failed to write answer: %v", err)
			}
			return nil
		})
	}

	// Streaming mode
	return openAIController.ControlStreaming(sysMsg, userMsg, cfg.Default.Quality, func(chunk openai.ChatCompletionChunk) error {
		answer := chunk.Choices[0].Delta.Content
		if _, err := f.WriteString(answer); err != nil {
			return fmt.Errorf("failed to write chunk: %v", err)
		}
		return nil
	})
}

func validateAppendFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", path)
	}

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(path), ".md") {
		return fmt.Errorf("file must have .md extension: %s", path)
	}

	return nil
}

func prepareAppendMessages(cfg config.Config, appendConfig *AppendConfig, content string, extraArgs []string) (string, string, error) {
	sysMsg := appendConfig.SystemMessage

	// Prepare template variables
	templateVars := map[string]string{
		"Content": content,
	}

	// Add operation-specific template variables based on extraArgs
	// For answer operation, extract last quote and other content
	if len(extraArgs) == 0 { // answer operation doesn't have extra args
		lastQuote, otherContents, err := file.LoadLastQuote(content)
		if err != nil {
			return "", "", fmt.Errorf("fail in loading last quote: %v", err)
		}
		templateVars["Question"] = lastQuote
		templateVars["Context"] = otherContents
	}

	// Apply template processing
	userMsg, err := appendConfig.UserMessage.Apply(templateVars)
	if err != nil {
		return "", "", fmt.Errorf("fail in creating user message: %v", err)
	}

	return sysMsg, userMsg, nil
}
