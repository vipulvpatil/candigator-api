package candidatebuilder

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/vipulvpatil/candidate-tracker-go/internal/clients/openai"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

const NOT_A_RESUME = "NOT A RESUME"

func Build(resumeText string, openAiClient openai.Client) (string, error) {
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
				Role:    "assistant",
				Content: "Here is the persona consisting of Name, Email, Phone, City, State, Country, Years of experience, Top 5 technical skills present in this profile as Tech Skills, Top 5 soft skills present in this profile as Soft Skills, Top 3 recommended Job positions.",
			},
		},
	}

	response, err := openAiClient.CallChatCompletionApi(&openAiChatCompletionRequest)
	if err != nil {
		return "", err
	}

	fmt.Println("response-----")
	fmt.Println(response)

	candidateData, err := getCandidateDataFromOpenAiResponse(response)
	if err != nil {
		return "", err
	}

	fmt.Println("candidateData-----")
	fmt.Println(candidateData)

	return response, nil
}

type Candidate struct {
	Name            string
	Email           string
	Phone           string
	City            string
	State           string
	Country         string
	YoE             string
	TechSkills      []string
	SoftSkills      []string
	RecommendedJobs []string
}

func getCandidateDataFromOpenAiResponse(response string) (*Candidate, error) {
	fmt.Println(response)
	if strings.Contains(response, NOT_A_RESUME) {
		return nil, errors.New("needs a valid resume to parse")
	}

	lines := linesFrom(response)
	keyValuePairs := []keyValuePair{}

	for i, line := range lines {
		fmt.Printf("Line %d:\n", i)
		fmt.Printf("%s\n", line)
		keyValuePairs = append(keyValuePairs, splitIntoKeyValues(line))
	}

	createCandidateMapFromKeyValuePairs(keyValuePairs)

	fmt.Println(keyValuePairs)

	return &Candidate{}, nil
}

func linesFrom(response string) []string {
	allLines := strings.Split(response, "\n")
	nonEmptyLines := []string{}
	for _, line := range allLines {
		fmt.Printf("\nwhat:%s:not", line)
		if !utilities.IsBlank(line) {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	return nonEmptyLines
}

func createCandidateMapFromKeyValuePairs(kvps []keyValuePair) {
	startingWithNumberRegex := regexp.MustCompile(`^\d*. `)
	candidateMap := map[string][]string{}
	currentAttribute := ""
	collectedValues := []string{}
	for _, keyValue := range kvps {
		if keyValue.hasKey() {
			if !utilities.IsBlank(currentAttribute) {
				candidateMap[currentAttribute] = collectedValues
				collectedValues = []string{}
			}
			currentAttribute = keyValue.Key
		}
		if keyValue.hasValue() {
			values := strings.Split(keyValue.Value, ",")
			for _, v := range values {
				trimmedValue := strings.Trim(v, " ")
				valueAfterRemovingNumber := string(startingWithNumberRegex.ReplaceAll([]byte(trimmedValue), []byte("")))
				if !utilities.IsBlank(valueAfterRemovingNumber) {
					collectedValues = append(collectedValues, valueAfterRemovingNumber)
				}
			}
		}
	}
	candidateMap[currentAttribute] = collectedValues

	fmt.Println("**************")
	for k, v := range candidateMap {
		fmt.Println("attribute--")
		fmt.Println(k)
		fmt.Println("values--")
		for _, va := range v {
			fmt.Printf("v: %s\n", va)
		}
	}
	fmt.Println("**************")
}

type keyValuePair struct {
	Key   string
	Value string
}

func (kvp keyValuePair) hasKey() bool {
	return !utilities.IsBlank(kvp.Key)
}

func (kvp keyValuePair) hasValue() bool {
	return !utilities.IsBlank(kvp.Value)
}

func splitIntoKeyValues(line string) keyValuePair {
	return keyValuePair{
		Key:   getKeyInLine(line),
		Value: getValueInLine(line),
	}
}

func getKeyInLine(line string) string {
	splits := strings.Split(line, ":")
	if len(splits) > 1 {
		return splits[0]
	}
	return ""
}

func getValueInLine(line string) string {
	splits := strings.Split(line, ":")
	if len(splits) > 1 {
		return splits[1]
	}
	if len(splits) > 0 {
		return splits[0]
	}
	return ""
}
