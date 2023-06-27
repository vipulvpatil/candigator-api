package openai

import (
	"github.com/pkg/errors"
	openaigo "github.com/sashabaranov/go-openai"
)

type chatCompletionRequest interface {
	GetChatCompletionMessages() []openaigo.ChatCompletionMessage
}

type ChatCompletionRequest struct {
	Messages []ChatCompletionMessage
}

type ChatCompletionMessage struct {
	Role    string
	Content string
}

func (m *ChatCompletionMessage) getChatCompletionMessageRole() (string, error) {
	switch m.Role {
	case "system", "user", "assistant":
		return m.Role, nil
	default:
		return "", errors.Errorf("unknown role: %s", m.Role)
	}
}

func (c *ChatCompletionRequest) GetChatCompletionMessages() []openaigo.ChatCompletionMessage {
	messages := []openaigo.ChatCompletionMessage{}
	for _, m := range c.Messages {
		role, err := m.getChatCompletionMessageRole()
		if err == nil {
			messages = append(messages, openaigo.ChatCompletionMessage{
				Role:    role,
				Content: m.Content,
			})
		}
	}
	return messages
}
