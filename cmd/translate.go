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

// translateCmd represents the translate command
var translateCmd = &cobra.Command{
	Use:   "translate [filepath] [language]",
	Short: "Translate markdown file to specified language",
	Long: `Translate a markdown file to the specified language using AI.
The translated content will be saved to a new file with "_[language]" suffix.
For example, if the input file is "document.md" and language is "ja", 
the output will be "document_ja.md".

Supported language codes: "en", "ja", "zh", "ko", "es", "fr", "de", etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		if err := translate(args, logger); err != nil {
			logger.Error("fail in calling translate", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)
}

func translate(args []string, logger *slog.Logger) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Warn("failed to load config, using defaults", "error", err)
		// Use default configuration if error occurs
		cfg = config.GetDefaultConfig()
	}
	client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	controller := controller.NewOpenAIController(&client, cfg.Default.Model, logger)

	if len(args) < 2 {
		return fmt.Errorf("both filepath and language are required")
	}
	path := args[0]
	language := args[1]

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", path)
	}

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(path), ".md") {
		return fmt.Errorf("file must have .md extension: %s", path)
	}

	// Validate language code
	if !isValidLanguageCode(language) {
		return fmt.Errorf("invalid language code: %s. Please use standard language codes like 'en', 'ja', 'zh', etc.", language)
	}

	// Generate output filename
	outputPath := generateTranslateOutputPath(path, language)

	// Load file content
	content, err := loadContent(path)
	if err != nil {
		return fmt.Errorf("fail in loading content: %v", err)
	}

	// Get system message from configuration
	sysMsg := cfg.Translate.SystemMessage

	// Add language-specific instruction
	sysMsg += fmt.Sprintf("\n\n**Target Language**: %s", getLanguageName(language))

	// Get user message from configuration and apply template processing
	userMsg, err := cfg.Translate.UserMessage.Apply(map[string]string{
		"Content":        content,
		"TargetLanguage": getLanguageName(language),
	})
	if err != nil {
		return fmt.Errorf("fail in creating user message: %v", err)
	}

	// Get quality settings from configuration
	maxTokens := cfg.Default.Quality.MaxTokens
	temperature := cfg.Default.Quality.Temperature

	// Log configuration values
	logger.Info("using configuration",
		"maxTokens", maxTokens,
		"temperature", temperature,
		"targetLanguage", language)

	controller.Control(sysMsg, userMsg, cfg.Default.Quality, func(completion *openai.ChatCompletion) error {
		translatedContent := completion.Choices[0].Message.Content

		// Save translation result to file
		if err := saveTranslation(outputPath, translatedContent, path, language); err != nil {
			return fmt.Errorf("fail in saving translation: %v", err)
		}

		logger.Info("translation created successfully", "input", path, "output", outputPath, "language", language)
		return nil
	})

	return nil
}

func generateTranslateOutputPath(inputPath, language string) string {
	dir := filepath.Dir(inputPath)
	filename := filepath.Base(inputPath)
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	return filepath.Join(dir, nameWithoutExt+"_"+language+".md")
}

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

func saveTranslation(outputPath, translatedContent, originalPath, language string) error {
	// Write to file
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write translation content directly
	if _, err := f.WriteString(translatedContent); err != nil {
		return err
	}

	return nil
}
