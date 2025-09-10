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

// summarizeCmd represents the summarize command
var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Summarize the content of a markdown file",
	Long: `Summarize the content of a markdown file using AI.
The summarized content will be saved to a new file with "_sum" suffix.
For example, if the input file is "document.md", the output will be "document_sum.md".`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetInstance().GetConfig()
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: cfg.Default.GetLogLevel().Level(),
		}))
		if err := summarize(cfg, args, logger); err != nil {
			logger.Error("fail in calling summarize", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(summarizeCmd)
}

func summarize(cfg config.Config, args []string, logger *slog.Logger) error {
	if len(args) == 0 {
		return fmt.Errorf("path is required")
	}

	path := args[0]
	extraArgs := []string{}

	// Call transform controller directly
	return controller.Transform(cfg, "summarize", path, extraArgs, logger)
}
