package personabuilder

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/clients/openai"
)

const NOT_A_RESUME = "NOT A RESUME"
const BUILDER_VERSION = "1.0.0"

func Build(resumeText string, openAiClient openai.Client) (*Persona, error) {
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
				Content: "Given the above resume, please build a persona in JSON with the following attributes. Name, Email, Phone, City, State, Country, Years of experience as \"YoE\" (type int), Top 5 technical skills present in this profile as \"Tech Skills\", Top 5 soft skills present in this profile as \"Soft Skills\", Top 3 recommended job positions as \"Recommended Roles\", Certifications, Institutes attended as \"Education\" including \"Qualification\" and \"CompletionYear\" (type string).",
			},
		},
	}

	return openAiClient.CallChatCompletionApi(&openAiChatCompletionRequest)
}

type Education struct {
	Institute      string `json:"Institute"`
	Qualification  string `json:"Qualification"`
	CompletionYear string `json:"CompletionYear"`
}
type Persona struct {
	Name             string      `json:"Name"`
	Email            string      `json:"Email"`
	Phone            string      `json:"Phone"`
	City             string      `json:"City"`
	State            string      `json:"State"`
	Country          string      `json:"Country"`
	YoE              int         `json:"YoE"`
	TechSkills       []string    `json:"Tech Skills"`
	SoftSkills       []string    `json:"Soft Skills"`
	RecommendedRoles []string    `json:"Recommended Roles"`
	Education        []Education `json:"Education"`
	Certifications   []string    `json:"Certifications"`
	BuilderVersion   string
}

func getPersonaDataFromOpenAiResponse(response string) (*Persona, error) {
	if strings.Contains(response, NOT_A_RESUME) {
		return nil, errors.New("needs a valid resume to parse")
	}

	var persona Persona

	err := json.Unmarshal([]byte(response), &persona)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse response")
	}

	persona.BuilderVersion = BUILDER_VERSION

	return &persona, nil
}
