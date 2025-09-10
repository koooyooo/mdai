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
		cfg := config.GetInstance().GetConfig()
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: cfg.Default.GetLogLevel().Level(),
		}))
		if err := translate(cfg, args, logger); err != nil {
			logger.Error("fail in calling translate", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)
}

func translate(cfg config.Config, args []string, logger *slog.Logger) error {
	if len(args) < 2 {
		return fmt.Errorf("both filepath and language are required")
	}

	path := args[0]
	language := args[1]
	extraArgs := []string{language}

	// Call transform controller directly
	return controller.Transform(cfg, "translate", path, extraArgs, logger)
}
