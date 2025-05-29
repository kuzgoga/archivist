package copilot

import (
	"bismark/internal/ai"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	copilotApi "github.com/stong1994/github-copilot-api"
	"math/rand"
	"time"
)

type Client struct {
	client *copilotApi.Copilot
}

type copilotResponse struct {
	Answer string
}

func NewClient() (*Client, error) {
	client, err := copilotApi.NewCopilot()
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

func (c *Client) Ask(request string) (ai.ChatResponse, error) {
	prompt := fmt.Sprintf(`give me json code with key 'Answer' that value will be a full answer to the question "%s"`, request)
	response, err := c.client.CreateCompletion(context.Background(), &copilotApi.CompletionRequest{
		Messages: []copilotApi.Message{
			{
				Role:    "system",
				Content: "You are powerful JSON programmer, but silent. You don't know Markdown and you use only json in answer. Don't use ```json!",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		StreamingFunc: func(_ context.Context, _ []byte) error {
			return nil
		},
	})
	if err != nil {
		return ai.ChatResponse{Successful: false}, err
	}
	if response.Choices[0].FinishReason != "stop" {
		return ai.ChatResponse{Successful: false}, errors.New("unexpected finishing reason")
	}
	answer := copilotResponse{}
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &answer)
	if err != nil {
		return ai.ChatResponse{Successful: false}, err
	}

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))

	return ai.ChatResponse{
		Answer:     answer.Answer,
		Successful: true,
	}, nil
}
