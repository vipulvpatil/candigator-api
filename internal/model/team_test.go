package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewTeam(t *testing.T) {
	currentFileCount := 1
	tests := []struct {
		name           string
		input          TeamOptions
		expectedOutput *Team
		errorExpected  bool
		errorString    string
	}{
		{
			name:           "id is empty",
			input:          TeamOptions{},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create team with an empty id",
		},
		{
			name: "name is empty",
			input: TeamOptions{
				Id: "123",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create team with an empty name",
		},
		{
			name: "explicit currentFileCount not provided",
			input: TeamOptions{
				Id:   "123",
				Name: "test",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create team without explicit current file count",
		},
		{
			name: "fileCountLimit is 0",
			input: TeamOptions{
				Id:               "123",
				Name:             "test",
				CurrentFileCount: &currentFileCount,
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create team with 0 file count limit",
		},
		{
			name: "Team gets created successfully",
			input: TeamOptions{
				Id:               "123",
				Name:             "test",
				CurrentFileCount: &currentFileCount,
				FileCountLimit:   100,
			},
			expectedOutput: &Team{
				id:               "123",
				name:             "test",
				currentFileCount: 1,
				fileCountLimit:   100,
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewTeam(tt.input)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
