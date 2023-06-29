package openai

import (
	"context"

	"github.com/pkg/errors"
	openaigo "github.com/sashabaranov/go-openai"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type client struct {
	apiKey string
	logger utilities.Logger
}

type Client interface {
	CallCompletionApi(prompt string) (string, error)
	CallChatCompletionApi(request chatCompletionRequest) (string, error)
}

type ClientOptions struct {
	ApiKey string
}

func NewClient(opts ClientOptions, logger utilities.Logger) Client {
	return &client{
		apiKey: opts.ApiKey,
		logger: logger,
	}
}

func (c *client) CallCompletionApi(prompt string) (string, error) {
	c.logger.LogMessageln(prompt)
	openAiGoClient := openaigo.NewClient(c.apiKey)
	ctx := context.Background()

	req := openaigo.CompletionRequest{
		Model:     openaigo.GPT3TextDavinci003,
		MaxTokens: 50,
		Prompt:    prompt,
	}
	resp, err := openAiGoClient.CreateCompletion(ctx, req)
	if err != nil {
		c.logger.LogError(err)
		return "", errors.Wrap(err, "Open Ai error")
	}
	return resp.Choices[0].Text, nil
}

func (c *client) CallChatCompletionApi(request chatCompletionRequest) (string, error) {
	messages := request.GetChatCompletionMessages()
	if len(messages) == 0 {
		return "", errors.New("no messages provided to be sent to OpenAI")
	}

	openAiGoClient := openaigo.NewClient(c.apiKey)
	ctx := context.Background()

	req := openaigo.ChatCompletionRequest{
		Model:     openaigo.GPT3Dot5Turbo,
		MaxTokens: 1000,
		Messages:  messages,
	}

	resp, err := openAiGoClient.CreateChatCompletion(ctx, req)
	if err != nil {
		c.logger.LogError(err)
		return "", errors.Wrap(err, "Open Ai error")
	}

	return resp.Choices[0].Message.Content, nil
}
