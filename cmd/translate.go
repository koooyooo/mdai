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
	controller := controller.NewOpenAIController(os.Getenv("OPENAI_API_KEY"), logger)

	if len(args) < 2 {
		return fmt.Errorf("both filepath and language are required")
	}
	path := args[0]
	language := args[1]

	// ファイルの存在確認
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", path)
	}

	// ファイル拡張子の確認
	if !strings.HasSuffix(strings.ToLower(path), ".md") {
		return fmt.Errorf("file must have .md extension: %s", path)
	}

	// 言語コードの検証
	if !isValidLanguageCode(language) {
		return fmt.Errorf("invalid language code: %s. Please use standard language codes like 'en', 'ja', 'zh', etc.", language)
	}

	// 出力ファイル名の生成
	outputPath := generateTranslateOutputPath(path, language)

	// ファイル内容の読み込み
	content, err := loadContent(path)
	if err != nil {
		return fmt.Errorf("fail in loading content: %v", err)
	}

	// 設定ファイルを読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Warn("failed to load config, using defaults", "error", err)
		// エラーが発生した場合はデフォルト設定を使用
		cfg = config.GetDefaultConfig()
	}

	// 設定ファイルからシステムメッセージを取得
	sysMsg := cfg.Translate.SystemMessage
	if sysMsg == "" {
		sysMsg = createTranslateSystemMessage()
	}

	// 言語固有の指示を追加
	sysMsg += fmt.Sprintf("\n\n**Target Language**: %s", getLanguageName(language))

	userMsg := createTranslateUserMessage(content, language)

	// 設定ファイルから品質設定を取得
	maxTokens := cfg.Default.Quality.MaxTokens
	temperature := cfg.Default.Quality.Temperature

	// 設定値をログに出力
	logger.Info("using configuration",
		"maxTokens", maxTokens,
		"temperature", temperature,
		"targetLanguage", language)

	controller.Control(sysMsg, userMsg, cfg.Default.Quality, func(completion *openai.ChatCompletion) error {
		translatedContent := completion.Choices[0].Message.Content

		// 翻訳結果をファイルに保存
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
	// 一般的な言語コードのリスト
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

func createTranslateSystemMessage() string {
	return `You are a professional translator specialized in translating markdown documents. When translating content, please follow these guidelines:

1. Translate the content to the specified target language accurately and naturally
2. Maintain the original markdown formatting and structure
3. Preserve all headings, lists, code blocks, and formatting elements
4. Keep the same tone and style as the original document
5. Ensure technical terms are translated appropriately for the target language
6. Maintain the document's readability and flow in the target language
7. Preserve any links, references, or citations
8. Keep the same level of detail and information as the original
9. Use appropriate language conventions for the target language
10. Ensure the translation sounds natural to native speakers of the target language`
}

func createTranslateUserMessage(content, language string) string {
	return fmt.Sprintf(`Please translate the following markdown content to %s:

%s

Please maintain all markdown formatting, structure, and ensure the translation is natural and accurate in the target language.`, getLanguageName(language), content)
}

func saveTranslation(outputPath, translatedContent, originalPath, language string) error {
	// ファイルに書き込み
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// 翻訳内容を直接書き込み
	if _, err := f.WriteString(translatedContent); err != nil {
		return err
	}

	return nil
}
