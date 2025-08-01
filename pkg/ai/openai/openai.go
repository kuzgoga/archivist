package openai

import (
	"archivist/pkg/ai"
	"context"
	"errors"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

type Client struct {
	client *openai.Client
	model  shared.ChatModel
}

func NewClient(apiKey string, model shared.ChatModel) *Client {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &Client{
		client: &client,
		model:  model,
	}
}

func (c *Client) Ask(request string) (ai.ChatResponse, error) {
	chatCompletion, err := c.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(request),
		},
		Model: c.model,
	})
	if err != nil {
		return ai.ChatResponse{}, err
	}
	if chatCompletion.Choices[0].FinishReason != "stop" {
		return ai.ChatResponse{
			Successful: false,
			Answer:     chatCompletion.Choices[0].Message.Content,
		}, errors.New("unexpected finishing reason")
	} else {
		return ai.ChatResponse{
			Successful: true,
			Answer:     chatCompletion.Choices[0].Message.Content,
		}, nil
	}
}
