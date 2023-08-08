package main

import (
	"fmt"
	"os"

	"github.com/vipulvpatil/candidate-tracker-go/internal/clients/openai"
	"github.com/vipulvpatil/candidate-tracker-go/internal/lib/parser"
	"github.com/vipulvpatil/candidate-tracker-go/internal/lib/parser/personabuilder"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("no filePath provided. correct usage includes one filePath or more filePaths")
		return
	}

	for i := 1; i < len(os.Args); i++ {
		filePath := os.Args[i]

		_, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("File reading error", err)
			return
		}

		fmt.Println("-------")
		fmt.Println(filePath)

		text, err := parser.GetTextFromPdf(filePath)
		if err != nil {
			fmt.Println("unable to parse given file")
			fmt.Println(err)
			return
		}

		openaiApiKey, ok := os.LookupEnv("OPENAI_API_KEY")
		if !ok {
			fmt.Println("OPENAI_API_KEY needed in ENV vars")
			return
		}

		openAiClient := openai.NewClient(openai.ClientOptions{ApiKey: openaiApiKey}, &utilities.StdoutLogger{})

		response, err := personabuilder.OpenAiResponseForResumeText(text, openAiClient)
		if err != nil {
			fmt.Println("openai error")
			fmt.Println(err)
			return
		}

		fmt.Println(response)
	}
}
