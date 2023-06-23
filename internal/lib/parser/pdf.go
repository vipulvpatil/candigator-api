package parser

import "github.com/rudolfoborges/pdf2go"

func GetTextFromPdf(filePath string) (string, error) {
	pdf, err := pdf2go.New(filePath, pdf2go.Config{
		LogLevel: pdf2go.LogLevelError,
	})

	if err != nil {
		return "", err
	}

	return pdf.Text()
}
