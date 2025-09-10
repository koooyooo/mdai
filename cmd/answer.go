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

// answerCmd represents the answer command
var answerCmd = &cobra.Command{
	Use:   "answer",
	Short: "Answer the question based on the content of a markdown file",
	Long: `Answer the question based on the content of a markdown file.
	The question will be extracted from the last quote in the file.
	The answer will be appended to the end of the file.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetInstance().GetConfig()
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: cfg.Default.GetLogLevel().Level(),
		}))
		if err := answer(cfg, args, logger); err != nil {
			logger.Error("fail in calling answer", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(answerCmd)
}

func answer(cfg config.Config, args []string, logger *slog.Logger) error {
	if len(args) == 0 {
		return fmt.Errorf("path is required")
	}

	path := args[0]
	extraArgs := []string{}

	// Call append controller directly
	return controller.Append(cfg, "answer", path, extraArgs, logger)
}
