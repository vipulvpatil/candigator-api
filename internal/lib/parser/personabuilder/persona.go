package personabuilder

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/clients/openai"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
)

const NOT_A_RESUME = "NOT A RESUME"
const BUILDER_VERSION = "1.0.0"

func Build(resumeText string, openAiClient openai.Client) (*model.Persona, error) {
	response, err := OpenAiResponseForResumeText(resumeText, openAiClient)
	if err != nil {
		return nil, err
	}
	personaData, err := getPersonaDataFromOpenAiResponse(response)
	if err != nil {
		return nil, err
	}

	return personaData, nil
}

func OpenAiResponseForResumeText(resumeText string, openAiClient openai.Client) (string, error) {
	openAiChatCompletionRequest := openai.ChatCompletionRequest{
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: fmt.Sprintf("You are Resume analyser. You read a resume and build a persona based on the given criteria. If the provided resume does not seem like a resume, your response should start with \"%s\"", NOT_A_RESUME),
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
				Role:    "user",
				Content: "Given the above resume, please build a persona in JSON with the following attributes (return empty if data is unavailable). Name, Email, Phone, City, State, Country, Years of experience as \"YoE\" (type int), Top 5 technical skills present in this profile as \"Tech Skills\" (type array of string), Top 5 soft skills present in this profile as \"Soft Skills\" (type array of string), Top 3 recommended job positions as \"Recommended Roles\" (type array of string), Certifications (type array of string), Institutes attended as \"Education\" (type array) including \"Qualification\" and \"CompletionYear\" (type string), Jobs held as \"Experience\" (type array) including \"Title\", \"Company Name\", \"Starting Year\" (type string), \"Ending Year\" (type string), \"Ongoing\" (type boolean).",
			},
		},
	}

	return openAiClient.CallChatCompletionApi(&openAiChatCompletionRequest)
}

func getPersonaDataFromOpenAiResponse(response string) (*model.Persona, error) {
	if strings.Contains(response, NOT_A_RESUME) {
		return nil, errors.New("needs a valid resume to parse")
	}

	var persona model.Persona

	err := json.Unmarshal([]byte(response), &persona)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse response")
	}

	persona.BuilderVersion = BUILDER_VERSION
	persona.BuiltBy = "AI"

	return &persona, nil
}
