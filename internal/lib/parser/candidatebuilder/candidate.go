package candidatebuilder

import "github.com/vipulvpatil/candidate-tracker-go/internal/clients/openai"

func Build(resumeText string, openAiClient openai.Client) (string, error) {
	openAiChatCompletionRequest := openai.ChatCompletionRequest{
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are Resume analyser. You read a resume and build a persona based on the given criteria. If the provided resume does not seem like a resume, you respond with \"NOT A RESUME\"",
			},
			{
				Role:    "assistant",
				Content: "Please share your resume.",
			},
			{
				Role:    "user",
				Content: resumeText,
			},
			{
				Role:    "assistant",
				Content: "Here is the persona consisting of Name, Email, Phone, City, State, Country, Years of experience, Top 5 technical skills present in this profile as Tech Skills, Top 5 soft skills present in this profile as Soft Skills, Top 3 recommended Job positions.",
			},
		},
	}

	return openAiClient.CallChatCompletionApi(&openAiChatCompletionRequest)
}
