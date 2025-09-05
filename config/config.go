package config

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of the configuration file
type Config struct {
	Default   DefaultConfig   `yaml:"default"`
	Answer    AnswerConfig    `yaml:"answer"`
	Summarize SummarizeConfig `yaml:"summarize"`
	Translate TranslateConfig `yaml:"translate"`
}

// DefaultConfig represents the default configuration
type DefaultConfig struct {
	Model         string        `yaml:"model"`
	Quality       QualityConfig `yaml:"quality"`
	LogLevel      string        `yaml:"log_level"`
	DisableStream bool          `yaml:"disable_stream"`
}

func (c DefaultConfig) GetLogLevel() slog.Level {
	switch c.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}
	return slog.LevelInfo
}

// QualityConfig represents quality settings
type QualityConfig struct {
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
}

// AnswerConfig represents the configuration for the answer command
type AnswerConfig struct {
	SystemMessage string              `yaml:"system_message"`
	UserMessage   UserMessageTemplate `yaml:"user_message"`
	TargetLength  int                 `yaml:"target_length"`
}

// UserMessageTemplate represents the template for user messages
type UserMessageTemplate struct {
	Template string
}

// Apply generates a message by applying variables to the template
func (t *UserMessageTemplate) Apply(vars map[string]string) (string, error) {
	if t.Template == "" {
		return "", fmt.Errorf("template is empty")
	}

	// Create template
	tmpl, err := template.New("userMessage").Parse(t.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return buf.String(), nil
}

// SummarizeConfig represents the configuration for the summarize command
type SummarizeConfig struct {
	SystemMessage string              `yaml:"system_message"`
	UserMessage   UserMessageTemplate `yaml:"user_message"`
	TargetLength  int                 `yaml:"target_length"`
}

// TranslateConfig represents the configuration for the translate command
type TranslateConfig struct {
	SystemMessage string              `yaml:"system_message"`
	UserMessage   UserMessageTemplate `yaml:"user_message"`
}

// LoadConfig loads the configuration file
func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	configPath := filepath.Join(homeDir, ".mdai", "config.yml")

	// Return default configuration if config file doesn't exist
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

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() *Config {
	return &Config{
		Default: DefaultConfig{
			Model: "gpt-4o-mini",
			Quality: QualityConfig{
				MaxTokens:   2000,
				Temperature: 0.7,
			},
			LogLevel:      "info",
			DisableStream: false,
		},
		Answer: AnswerConfig{
			SystemMessage: `You are a helpful and detailed assistant. When answering questions based on the given context, please follow these guidelines:

1. Answer in the same language as the question
2. Make full use of the context information
3. Add examples and explanations when necessary
4. Ensure answers are appropriately long and content-rich
5. Provide insights that deepen the questioner's understanding
6. Prefer rich markdown formatting`,
			UserMessage: UserMessageTemplate{
				Template: `Context: {{.Context}}

Question: {{.Question}}`,
			},
			TargetLength: 500,
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
			UserMessage: UserMessageTemplate{
				Template: `Please provide a comprehensive summary of the following markdown content:

{{.Content}}

Please create a well-structured summary that captures the essence and key points of this content.`,
			},
			TargetLength: 800,
		},
		Translate: TranslateConfig{
			SystemMessage: `You are a professional translator specialized in translating markdown documents. When translating content, please follow these guidelines:

1. Translate the content to the specified target language accurately and naturally
2. Maintain the original markdown formatting and structure
3. Preserve all headings, lists, code blocks, and formatting elements
4. Keep the same tone and style as the original document
5. Ensure technical terms are translated appropriately for the target language
6. Maintain the document's readability and flow in the target language
7. Preserve any links, references, or citations
8. Keep the same level of detail and information as the original
9. Use appropriate language conventions for the target language
10. Ensure the translation sounds natural to native speakers of the target language`,
			UserMessage: UserMessageTemplate{
				Template: `Please translate the following content to {{.TargetLanguage}}:

{{.Content}}

Please maintain the original markdown formatting and structure while ensuring the translation is accurate and natural.`,
			},
		},
	}
}

// GetLogLevel gets the log level
func (c *Config) GetLogLevel() string {
	if c.Default.LogLevel == "" {
		return "info"
	}
	return c.Default.LogLevel
}

// GetModel gets the default model
func (c *Config) GetModel() string {
	if c.Default.Model == "" {
		return "gpt-4o-mini"
	}
	return c.Default.Model
}

// GetMaxTokens gets the default maximum token count
func (c *Config) GetMaxTokens() int {
	if c.Default.Quality.MaxTokens == 0 {
		return 2000
	}
	return c.Default.Quality.MaxTokens
}

// GetTemperature gets the default temperature setting
func (c *Config) GetTemperature() float64 {
	if c.Default.Quality.Temperature == 0 {
		return 0.7
	}
	return c.Default.Quality.Temperature
}
