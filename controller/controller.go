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
)

type OpenAIController struct {
	client  *openai.Client
	modelID string
	logger  *slog.Logger
}

func NewOpenAIController(client *openai.Client, modelID string, logger *slog.Logger) *OpenAIController {
	return &OpenAIController{
		client:  client,
		modelID: modelID,
		logger:  logger,
	}
}

func (c *OpenAIController) Control(sysMsg, usrMsg string, quality config.QualityConfig, completionFunc func(res *openai.ChatCompletion) error) error {
	// Use default values if configuration values are 0
	maxTokens := quality.MaxTokens
	temperature := quality.Temperature
	if maxTokens == 0 {
		maxTokens = models.DefaultMaxTokens
	}
	if temperature == 0.0 {
		temperature = models.DefaultTemperature
	}

	completion, err := c.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: c.modelID,
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

	costInfo, err := models.CalculateCostString(c.modelID, completion.Usage)
	if err != nil {
		return fmt.Errorf("cost calculation error: %v", err)
	}
	c.logger.Info("cost information", "costInfo", costInfo)
	// completion := resp.Choices[0].Message.Content

	return completionFunc(completion)
}

func (c *OpenAIController) ControlStreaming(sysMsg, usrMsg string, quality config.QualityConfig, completionFunc func(res openai.ChatCompletionChunk) error) error {
	stream := c.client.Chat.Completions.NewStreaming(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(sysMsg),
			openai.UserMessage(usrMsg),
		},
		Seed:  openai.Int(0),
		Model: c.modelID,
	})

	// optionally, an accumulator helper can be used
	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if content, ok := acc.JustFinishedContent(); ok {
			c.logger.Debug("Content stream finished:", "content", content)
			costInfo, err := models.CalculateCostString(c.modelID, acc.Usage)
			if err != nil {
				return fmt.Errorf("cost calculation error: %v", err)
			}
			c.logger.Info("cost information", "costInfo", costInfo)
		}

		// if using tool calls
		if tool, ok := acc.JustFinishedToolCall(); ok {
			c.logger.Debug("Tool call stream finished:", "index", tool.Index, "name", tool.Name, "arguments", tool.Arguments)
		}

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			c.logger.Debug("Refusal stream finished:", "refusal", refusal)
		}

		// it's best to use chunks after handling JustFinished events
		if len(chunk.Choices) > 0 {
			completionFunc(chunk)
		}
	}

	if stream.Err() != nil {
		return stream.Err()
	}
	return nil
}
