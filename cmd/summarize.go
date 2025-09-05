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
	"github.com/openai/openai-go/option"
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
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Warn("failed to load config, using defaults", "error", err)
		// Use default configuration if error occurs
		cfg = config.GetDefaultConfig()
	}
	client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	controller := controller.NewOpenAIController(&client, cfg.Default.Model, logger)

	if len(args) == 0 {
		return fmt.Errorf("path is required")
	}
	path := args[0]

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", path)
	}

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(path), ".md") {
		return fmt.Errorf("file must have .md extension: %s", path)
	}

	// Generate output filename
	outputPath := generateOutputPath(path)

	// Load file content
	content, err := loadContent(path)
	if err != nil {
		return fmt.Errorf("fail in loading content: %v", err)
	}

	// Load configuration file
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Warn("failed to load config, using defaults", "error", err)
		// Use default configuration if error occurs
		cfg = config.GetDefaultConfig()
	}

	// Get system message from configuration
	sysMsg := cfg.Summarize.SystemMessage

	// Get user message from configuration and apply template processing
	userMsg, err := cfg.Summarize.UserMessage.Apply(map[string]string{
		"Content": content,
	})
	if err != nil {
		return fmt.Errorf("fail in creating user message: %v", err)
	}

	// Get quality settings from configuration
	maxTokens := cfg.Default.Quality.MaxTokens
	temperature := cfg.Default.Quality.Temperature

	// Add character count instruction to system message
	if cfg.Summarize.TargetLength > 0 {
		sysMsg += fmt.Sprintf("\n\n**Summary Length Guidance**: Please provide a summary of approximately %d characters.", cfg.Summarize.TargetLength)
	}

	// Log configuration values
	logger.Info("using configuration",
		"maxTokens", maxTokens,
		"temperature", temperature,
		"targetLength", cfg.Summarize.TargetLength)

	controller.Control(sysMsg, userMsg, cfg.Default.Quality, func(completion *openai.ChatCompletion) error {
		summary := completion.Choices[0].Message.Content

		// Save summary result to file
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

func saveSummary(outputPath, summary, originalPath string) error {
	// Write to file
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write summary content directly
	if _, err := f.WriteString(summary); err != nil {
		return err
	}

	return nil
}
