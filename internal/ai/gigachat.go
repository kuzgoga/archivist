package ai

import (
	"errors"
	"fmt"
	"github.com/paulrzcz/go-gigachat"
)

type GigaChat struct {
	client      *gigachat.Client
	model       string
	temperature float64
}

func NewGigaChat(clientId string, clientSecret string, model string) (*GigaChat, error) {
	client, err := gigachat.NewInsecureClient(clientId, clientSecret)
	if err != nil {
		return nil, err
	}
	err = client.Auth()
	if err != nil {
		return nil, err
	}

	return &GigaChat{
		client:      client,
		model:       model,
		temperature: 0.7,
	}, nil
}

func (g *GigaChat) Ask(request string) (ChatResponse, error) {
	var n int64 = 1
	var maxTokens int64 = 120
	var repetitionPenalty = 1.1
	const normalFinishingReason string = "stop"

	chat, err := g.client.Chat(&gigachat.ChatRequest{
		Model: g.model,
		Messages: []gigachat.Message{
			{
				Role:    gigachat.UserRole,
				Content: request,
			},
		},
		N:                 &n,
		Stream:            new(bool),
		MaxTokens:         &maxTokens,
		RepetitionPenalty: &repetitionPenalty,
	})
	if err != nil {
		return ChatResponse{}, err
	}

	if chat.Choices[0].FinishReason != normalFinishingReason {
		fmt.Println(chat.Choices[0].FinishReason)
		return ChatResponse{
			Answer:     chat.Choices[0].Message.Content,
			Successful: false,
		}, errors.New("unexpected finishing reason")
	} else {
		return ChatResponse{
			Answer:     chat.Choices[0].Message.Content,
			Successful: false,
		}, nil
	}
}
