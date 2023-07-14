package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewCandidate(t *testing.T) {
	tests := []struct {
		name           string
		input          CandidateOptions
		expectedOutput *Candidate
		errorExpected  bool
		errorString    string
	}{
		{
			name:           "id is empty",
			input:          CandidateOptions{},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create candidate with an empty id",
		},
		{
			name: "has no associated persoa",
			input: CandidateOptions{
				Id: "123",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create candidate without a valid persona",
		},
		{
			name: "Candidate gets created successfully with AI Generated Persona",
			input: CandidateOptions{
				Id:                 "123",
				AiGeneratedPersona: &Persona{Name: "persona name"},
			},
			expectedOutput: &Candidate{
				id:                 "123",
				aiGeneratedPersona: &Persona{Name: "persona name"},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "Candidate gets created successfully with Manual Persona",
			input: CandidateOptions{
				Id:                     "123",
				ManuallyCreatedPersona: &Persona{Name: "persona name"},
			},
			expectedOutput: &Candidate{
				id:                     "123",
				manuallyCreatedPersona: &Persona{Name: "persona name"},
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewCandidate(tt.input)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
