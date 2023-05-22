package chatgpt

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

type OpenaiClient struct {
	client *openai.Client
}

func NewOpenAIClient() *OpenaiClient {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, defaulting to system envs")
	}

	return &OpenaiClient{
		client: openai.NewClient(os.Getenv("OPENAI_API_TOKEN")),
	}
}

func (oaic *OpenaiClient) SendChatGPT(msg string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ChatGPTTimeout)
	defer cancel()
	resp, err := oaic.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
