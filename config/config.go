package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config は設定ファイルの構造を表します
type Config struct {
	Default   DefaultConfig   `yaml:"default"`
	Answer    AnswerConfig    `yaml:"answer"`
	Summarize SummarizeConfig `yaml:"summarize"`
}

// DefaultConfig はデフォルト設定を表します
type DefaultConfig struct {
	Model    string        `yaml:"model"`
	Quality  QualityConfig `yaml:"quality"`
	LogLevel string        `yaml:"log_level"`
}

// QualityConfig は品質設定を表します
type QualityConfig struct {
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
}

// AnswerConfig はanswerコマンドの設定を表します
type AnswerConfig struct {
	SystemMessage string `yaml:"system_message"`
	TargetChars   int    `yaml:"target_chars"`
}

// SummarizeConfig はsummarizeコマンドの設定を表します
type SummarizeConfig struct {
	SystemMessage string `yaml:"system_message"`
	TargetChars   int    `yaml:"target_chars"`
}

// LoadConfig は設定ファイルを読み込みます
func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	configPath := filepath.Join(homeDir, ".mdai", "config.yml")

	// 設定ファイルが存在しない場合はデフォルト設定を返す
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return GetDefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// GetDefaultConfig はデフォルト設定を返します
func GetDefaultConfig() *Config {
	return &Config{
		Default: DefaultConfig{
			Model: "gpt-4o-mini",
			Quality: QualityConfig{
				MaxTokens:   2000,
				Temperature: 0.7,
			},
			LogLevel: "info",
		},
		Answer: AnswerConfig{
			SystemMessage: `You are a helpful and detailed assistant. When answering questions based on the given context, please follow these guidelines:

1. Answer in the same language as the question
2. Make full use of the context information
3. Provide specific and practical information
4. Add examples and explanations when necessary
5. Ensure answers are appropriately long and content-rich
6. Provide insights that deepen the questioner's understanding
7. Prefer rich markdown formatting`,
			TargetChars: 500,
		},
		Summarize: SummarizeConfig{
			SystemMessage: `You are a helpful and detailed assistant specialized in summarizing markdown documents. When summarizing content, please follow these guidelines:

1. Provide a comprehensive yet concise summary of the main content
2. Maintain the key points and important information
3. Use clear and organized structure with markdown formatting
4. Include main headings and subheadings when relevant
5. Preserve important details, examples, and references
6. Make the summary easy to read and understand
7. Use appropriate markdown elements (headers, lists, emphasis, etc.)
8. Keep the summary appropriately long - not too brief, not too verbose
9. Focus on the most valuable and actionable information
10. Maintain the original tone and style when appropriate`,
			TargetChars: 800,
		},
	}
}

// GetLogLevel はログレベルを取得します
func (c *Config) GetLogLevel() string {
	if c.Default.LogLevel == "" {
		return "info"
	}
	return c.Default.LogLevel
}

// GetModel はデフォルトモデルを取得します
func (c *Config) GetModel() string {
	if c.Default.Model == "" {
		return "gpt-4o-mini"
	}
	return c.Default.Model
}

// GetMaxTokens はデフォルトの最大トークン数を取得します
func (c *Config) GetMaxTokens() int {
	if c.Default.Quality.MaxTokens == 0 {
		return 2000
	}
	return c.Default.Quality.MaxTokens
}

// GetTemperature はデフォルトの温度設定を取得します
func (c *Config) GetTemperature() float64 {
	if c.Default.Quality.Temperature == 0 {
		return 0.7
	}
	return c.Default.Quality.Temperature
}
