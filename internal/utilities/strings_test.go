package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsBlank(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput bool
	}{
		{
			name:           "Non-empty string returns false",
			input:          "Non empty string",
			expectedOutput: false,
		},
		{
			name:           "Non-empty string with trailing white spaces returns false",
			input:          "Non empty string   ",
			expectedOutput: false,
		},
		{
			name:           "Non-empty string with leading white spaces returns false",
			input:          "   Non empty string",
			expectedOutput: false,
		},
		{
			name:           "Empty string returns true",
			input:          "",
			expectedOutput: true,
		},
		{
			name:           "Empty white spaces returns false",
			input:          "   ",
			expectedOutput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBlank(tt.input)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
