/*
Copyright © 2025 koooyooo
*/
package cmd

import (
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed config.sample.yml
var configSample []byte

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
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	// Path to ~/.mdai directory
	configDir := filepath.Join(homeDir, ".mdai")
	configPath := filepath.Join(configDir, "config.yml")

	// Create config directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		logger.Info("creating config directory", "path", configDir)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		logger.Info("config file already exists", "path", configPath)
		logger.Info("if you want to overwrite, please remove the existing file first")
		return nil
	}

	// Write embedded config sample to config file
	if err := os.WriteFile(configPath, configSample, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	logger.Info("configuration initialized successfully", "path", configPath)
	logger.Info("you can now customize the configuration by editing this file")

	return nil
}
