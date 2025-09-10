/*
Copyright © 2025 koooyooo
*/
package controller

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
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
	success, sysMsg, userMsg, err := prepareAppendMessages(cfg, appendConfig, content, extraArgs)
	if err != nil {
		return err
	}

	if !success {
		return nil
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

func prepareAppendMessages(cfg config.Config, appendConfig *AppendConfig, content string, extraArgs []string) (bool, string, string, error) {
	sysMsg := appendConfig.SystemMessage

	// Prepare template variables
	templateVars := map[string]string{
		"Content": content,
	}

	// Add operation-specific template variables based on extraArgs
	// For answer operation, extract last quote and other content
	if len(extraArgs) == 0 { // answer operation doesn't have extra args
		success, lastQuote, otherContents, err := file.LoadLastQuote(content)
		if err != nil {
			return false, "", "", fmt.Errorf("fail in loading last quote: %v", err)
		}
		if !success {
			// 引用が見つからない場合は処理をキャンセル
			return false, "", "", nil
		}
		templateVars["Question"] = lastQuote
		templateVars["Context"] = otherContents
	}

	// Apply template processing
	userMsg, err := appendConfig.UserMessage.Apply(templateVars)
	if err != nil {
		return false, "", "", fmt.Errorf("fail in creating user message: %v", err)
	}

	return true, sysMsg, userMsg, nil
}

type WatchConfig struct {
	FilePath   string
	Operation  string
	ExtraArgs  []string
	DebounceMs int
}

func WatchAndAppend(cfg config.Config, watchConfig *WatchConfig, logger *slog.Logger) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(watchConfig.FilePath)
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	logger.Info("file watching started", "file", watchConfig.FilePath, "operation", watchConfig.Operation)
	logger.Info("press Ctrl+C to exit")
	logger.Info("")

	working := false
	var mu sync.RWMutex
	cycleCount := 0
	maxCycles := 20

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				mu.RLock()
				if working {
					mu.RUnlock()
					continue
				} else {
					mu.RUnlock()
				}

				// サイクルカウントをチェック
				cycleCount++
				if cycleCount > maxCycles {
					logger.Info("maximum cycles reached", "cycles", maxCycles)
					logger.Info("exiting...")
					return nil
				}

				mu.Lock()
				working = true
				mu.Unlock()

				logger.Info("file changed", "file", event.Name, "cycle", cycleCount)
				time.Sleep(time.Duration(watchConfig.DebounceMs) * time.Millisecond)
				logger.Info("starting append operation")
				err := Append(cfg, watchConfig.Operation, watchConfig.FilePath, watchConfig.ExtraArgs, logger)
				if err != nil {
					logger.Error("append failed", "error", err)
				} else {
					logger.Info("append operation completed")
					logger.Info("")
				}
				mu.Lock()
				working = false
				mu.Unlock()
			}
		case err := <-watcher.Errors:
			logger.Error("watcher error", "error", err)
		case sig := <-sigChan:
			logger.Info("received signal to terminate", "signal", sig)
			logger.Info("exiting...")
			return nil
		}
	}
}
