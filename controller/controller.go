/*
Copyright © 2025 koooyooo
*/
package controller

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/koooyooo/mdai/config"
	"github.com/koooyooo/mdai/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIController struct {
	apiKey string
	logger *slog.Logger
}

func NewOpenAIController(apiKey string, logger *slog.Logger) *OpenAIController {
	return &OpenAIController{
		apiKey: apiKey,
		logger: logger,
	}
}

func (c *OpenAIController) Control(sysMsg, usrMsg string, quality config.QualityConfig, completionFunc func(res *openai.ChatCompletion) error) error {
	var modelID = openai.ChatModelGPT4oMini
	client := openai.NewClient(option.WithAPIKey(c.apiKey))

	// Use default values if configuration values are 0
	maxTokens := quality.MaxTokens
	temperature := quality.Temperature
	if maxTokens == 0 {
		maxTokens = models.DefaultMaxTokens
	}
	if temperature == 0.0 {
		temperature = models.DefaultTemperature
	}

	completion, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: modelID,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(sysMsg),
			openai.UserMessage(usrMsg),
		},
		MaxTokens:   openai.Int(int64(maxTokens)),
		Temperature: openai.Float(temperature),
	})
	if err != nil {
		return fmt.Errorf("OpenAI API error: %v", err)
	}
	if len(completion.Choices) == 0 {
		return fmt.Errorf("no response from OpenAI API")
	}

	costInfo, err := models.CalculateCostString(modelID, completion.Usage)
	if err != nil {
		return fmt.Errorf("cost calculation error: %v", err)
	}
	c.logger.Info("cost information", "costInfo", costInfo)
	// completion := resp.Choices[0].Message.Content

	return completionFunc(completion)
}
