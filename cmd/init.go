/*
Copyright © 2025 koooyooo
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize mdai configuration",
	Long: `Initialize mdai configuration by creating the config directory and copying the sample config file.
	
This command will:
1. Create ~/.mdai directory if it doesn't exist
2. Copy config.sample.yml to ~/.mdai/config.yml if it doesn't exist
3. Display the path of the created config file`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		if err := initConfig(logger); err != nil {
			logger.Error("fail in calling init", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initConfig(logger *slog.Logger) error {
	// ホームディレクトリを取得
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	// ~/.mdaiディレクトリのパス
	configDir := filepath.Join(homeDir, ".mdai")
	configPath := filepath.Join(configDir, "config.yml")

	// 設定ディレクトリが存在しない場合は作成
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		logger.Info("creating config directory", "path", configDir)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	// 設定ファイルが既に存在する場合は確認
	if _, err := os.Stat(configPath); err == nil {
		logger.Info("config file already exists", "path", configPath)
		logger.Info("if you want to overwrite, please remove the existing file first")
		return nil
	}

	// サンプル設定ファイルのパス（実行ファイルと同じディレクトリ）
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}
	execDir := filepath.Dir(execPath)
	sampleConfigPath := filepath.Join(execDir, "config.sample.yml")

	// サンプル設定ファイルが存在しない場合は、カレントディレクトリを試す
	if _, err := os.Stat(sampleConfigPath); os.IsNotExist(err) {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %v", err)
		}
		sampleConfigPath = filepath.Join(currentDir, "config.sample.yml")
	}

	// サンプル設定ファイルの存在確認
	if _, err := os.Stat(sampleConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("sample config file not found at: %s", sampleConfigPath)
	}

	// サンプル設定ファイルを読み込み
	sampleData, err := os.ReadFile(sampleConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read sample config file: %v", err)
	}

	// 設定ファイルに書き込み
	if err := os.WriteFile(configPath, sampleData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	logger.Info("configuration initialized successfully", "path", configPath)
	logger.Info("you can now customize the configuration by editing this file")

	return nil
}
