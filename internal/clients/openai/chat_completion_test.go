package openai

import (
	"testing"

	openaigo "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func Test_GetChatCompletionMessages(t *testing.T) {
	t.Run("returs chat completion request messages in appropriate format ", func(t *testing.T) {
		input := ChatCompletionRequest{
			Messages: []ChatCompletionMessage{
				{
					Role:    "system",
					Content: "message for system",
				},
				{
					Role:    "random_user",
					Content: "does not get included in output",
				},
				{
					Role:    "user",
					Content: "message from user",
				},
				{
					Role:    "assistant",
					Content: "message from AI",
				},
			},
		}

		expected := []openaigo.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "message for system",
			},
			{
				Role:    "user",
				Content: "message from user",
			},
			{
				Role:    "assistant",
				Content: "message from AI",
			},
		}

		output := input.GetChatCompletionMessages()
		assert.Equal(t, expected, output)
	})
}
