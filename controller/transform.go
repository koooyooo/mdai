/*
Copyright © 2025 koooyooo
*/
package controller

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/koooyooo/mdai/config"
	"github.com/koooyooo/mdai/util/file"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// TransformOperation represents the type of transformation to perform
type TransformOperation string

const (
	OperationSummarize TransformOperation = "summarize"
	OperationTranslate TransformOperation = "translate"
)

// TransformConfig holds configuration for a specific transformation
type TransformConfig struct {
	SystemMessage  string
	UserMessage    config.UserMessageTemplate
	SuffixTemplate config.UserMessageTemplate
	ExtraArgs      []string
}

// Transform performs a transformation operation on a markdown file
func Transform(cfg config.Config, operation TransformOperation, path string, extraArgs []string, logger *slog.Logger) error {
	// Get operation configuration dynamically
	opConfig, err := getOperationConfig(cfg, operation)
	if err != nil {
		return err
	}

	// Validate arguments using configuration
	if err := validateArgs(extraArgs, opConfig.Args); err != nil {
		return err
	}

	// Create transform configuration
	transformConfig := &TransformConfig{
		SystemMessage:  opConfig.SystemMessage,
		UserMessage:    opConfig.UserMessage,
		SuffixTemplate: opConfig.Suffix,
		ExtraArgs:      extraArgs,
	}

	// Execute transformation
	return executeTransform(cfg, transformConfig, path, extraArgs, logger)
}

func getOperationConfig(cfg config.Config, operation TransformOperation) (config.OperationConfig, error) {
	// Get operation configuration from map
	opConfig, exists := cfg.Transform.Operations[string(operation)]
	if !exists {
		return config.OperationConfig{}, fmt.Errorf("unsupported operation: %s", operation)
	}

	return opConfig, nil
}

func validateArgs(extraArgs []string, argsConfig config.ArgsConfig) error {
	argCount := len(extraArgs)

	if argCount < argsConfig.MinCount {
		return fmt.Errorf("operation requires at least %d arguments, got %d", argsConfig.MinCount, argCount)
	}

	if argsConfig.MaxCount > 0 && argCount > argsConfig.MaxCount {
		return fmt.Errorf("operation accepts at most %d arguments, got %d", argsConfig.MaxCount, argCount)
	}

	return nil
}

func executeTransform(cfg config.Config, transformConfig *TransformConfig, path string, extraArgs []string, logger *slog.Logger) error {
	// Validate file
	if err := validateFile(path); err != nil {
		return err
	}

	// Generate output filename using suffix template
	suffix, err := generateSuffix(transformConfig.SuffixTemplate, extraArgs)
	if err != nil {
		return fmt.Errorf("fail in generating output suffix: %v", err)
	}
	outputPath := generateOutputPath(path, suffix)

	// Load file content
	content, err := file.LoadContent(path)
	if err != nil {
		return fmt.Errorf("fail in loading content: %v", err)
	}

	// Prepare messages
	sysMsg, userMsg, err := prepareMessages(cfg, transformConfig, content, extraArgs)
	if err != nil {
		return err
	}

	// Execute transformation
	client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	openAIController := NewOpenAIController(&client, cfg.Default.Model, logger)

	// Log configuration values
	logger.Info("using configuration",
		"maxTokens", cfg.Default.Quality.MaxTokens,
		"temperature", cfg.Default.Quality.Temperature)

	openAIController.Control(sysMsg, userMsg, cfg.Default.Quality, func(completion *openai.ChatCompletion) error {
		result := completion.Choices[0].Message.Content

		// Save result to file
		if err := saveResult(outputPath, result, path, extraArgs); err != nil {
			return fmt.Errorf("fail in saving result: %v", err)
		}

		logger.Info("transformation completed successfully",
			"input", path,
			"output", outputPath)
		return nil
	})

	return nil
}

func validateFile(path string) error {
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

func generateSuffix(suffixTemplate config.UserMessageTemplate, extraArgs []string) (string, error) {
	// Prepare template variables
	templateVars := map[string]string{}

	// Add extra arguments as Arg0, Arg1, etc.
	for i, arg := range extraArgs {
		templateVars[fmt.Sprintf("Arg%d", i)] = arg
	}

	// Apply template processing
	suffix, err := suffixTemplate.Apply(templateVars)
	if err != nil {
		return "", fmt.Errorf("fail in applying suffix template: %v", err)
	}

	return suffix, nil
}

func generateOutputPath(inputPath, suffix string) string {
	dir := filepath.Dir(inputPath)
	filename := filepath.Base(inputPath)
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	return filepath.Join(dir, nameWithoutExt+suffix+".md")
}

func prepareMessages(cfg config.Config, transformConfig *TransformConfig, content string, extraArgs []string) (string, string, error) {
	sysMsg := transformConfig.SystemMessage

	// Prepare template variables
	templateVars := map[string]string{
		"Content": content,
	}

	// Add operation-specific template variables based on extraArgs
	// This allows each operation to define its own template variables
	if len(extraArgs) > 0 {
		// For translate operation, add TargetLanguage variable
		if len(extraArgs) == 1 && isValidLanguageCode(extraArgs[0]) {
			templateVars["TargetLanguage"] = getLanguageName(extraArgs[0])
		}
	}

	// Apply template processing
	userMsg, err := transformConfig.UserMessage.Apply(templateVars)
	if err != nil {
		return "", "", fmt.Errorf("fail in creating user message: %v", err)
	}

	return sysMsg, userMsg, nil
}

func saveResult(outputPath, result, originalPath string, extraArgs []string) error {
	// Write to file
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write result content directly
	if _, err := f.WriteString(result); err != nil {
		return err
	}

	return nil
}

// Language-related utility functions (used by translate operation)
func isValidLanguageCode(language string) bool {
	// List of common language codes
	validLanguages := map[string]bool{
		"en": true, "ja": true, "zh": true, "ko": true, "es": true,
		"fr": true, "de": true, "it": true, "pt": true, "ru": true,
		"ar": true, "hi": true, "th": true, "vi": true, "nl": true,
		"sv": true, "no": true, "da": true, "fi": true, "pl": true,
		"tr": true, "he": true, "id": true, "ms": true, "ca": true,
	}
	return validLanguages[strings.ToLower(language)]
}

func getLanguageName(languageCode string) string {
	languageNames := map[string]string{
		"en": "English",
		"ja": "Japanese (日本語)",
		"zh": "Chinese (中文)",
		"ko": "Korean (한국어)",
		"es": "Spanish (Español)",
		"fr": "French (Français)",
		"de": "German (Deutsch)",
		"it": "Italian (Italiano)",
		"pt": "Portuguese (Português)",
		"ru": "Russian (Русский)",
		"ar": "Arabic (العربية)",
		"hi": "Hindi (हिन्दी)",
		"th": "Thai (ไทย)",
		"vi": "Vietnamese (Tiếng Việt)",
		"nl": "Dutch (Nederlands)",
		"sv": "Swedish (Svenska)",
		"no": "Norwegian (Norsk)",
		"da": "Danish (Dansk)",
		"fi": "Finnish (Suomi)",
		"pl": "Polish (Polski)",
		"tr": "Turkish (Türkçe)",
		"he": "Hebrew (עברית)",
		"id": "Indonesian (Bahasa Indonesia)",
		"ms": "Malay (Bahasa Melayu)",
		"ca": "Catalan (Català)",
	}

	if name, exists := languageNames[strings.ToLower(languageCode)]; exists {
		return name
	}
	return languageCode
}
