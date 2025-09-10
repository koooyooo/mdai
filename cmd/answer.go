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

var watchFlag bool

// answerCmd represents the answer command
var answerCmd = &cobra.Command{
	Use:   "answer [file]",
	Short: "Answer the question based on the content of a markdown file",
	Long: `Answer the question based on the content of a markdown file.
	The question will be extracted from the last quote in the file.
	The answer will be appended to the end of the file.
	
	Use --watch or -w to enable file watching mode, which will automatically
	answer questions when the file is modified.`,
	Args: cobra.ExactArgs(1),
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

	// Add watch flag
	answerCmd.Flags().BoolVarP(&watchFlag, "watch", "w", false, "Watch file for changes and auto-answer")
}

func answer(cfg config.Config, args []string, logger *slog.Logger) error {
	if len(args) == 0 {
		return fmt.Errorf("path is required")
	}

	path := args[0]
	extraArgs := []string{}

	if watchFlag {
		// Watch mode: use WatchAndAppend
		watchConfig := &controller.WatchConfig{
			FilePath:   path,
			Operation:  "answer",
			ExtraArgs:  extraArgs,
			DebounceMs: 500, // 500ms debounce to avoid rapid successive calls
		}
		return controller.WatchAndAppend(cfg, watchConfig, logger)
	} else {
		// Normal mode: single execution
		return controller.Append(cfg, "answer", path, extraArgs, logger)
	}
}
