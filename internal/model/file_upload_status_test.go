package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileUploadStatus(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput fileUploadStatus
	}{
		{
			name:           "creates INITIATED file upload status",
			input:          "INITIATED",
			expectedOutput: initiated,
		},
		{
			name:           "creates SUCCESS file upload status",
			input:          "SUCCESS",
			expectedOutput: success,
		},
		{
			name:           "creates FAILURE file upload status",
			input:          "FAILURE",
			expectedOutput: failure,
		},
		{
			name:           "handles unknown file upload status",
			input:          "unknown",
			expectedOutput: undefinedFileUploadStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := FileUploadStatus(tt.input)
			assert.Equal(t, state, tt.expectedOutput)
		})
	}
}

func Test_FileUploadStatus_String(t *testing.T) {
	tests := []struct {
		name           string
		input          fileUploadStatus
		expectedOutput string
	}{
		{
			name:           "gets INITIATED from ai game state",
			input:          initiated,
			expectedOutput: "INITIATED",
		},
		{
			name:           "gets SUCCESS from human game state",
			input:          success,
			expectedOutput: "SUCCESS",
		},
		{
			name:           "gets FAILURE from human game state",
			input:          failure,
			expectedOutput: "FAILURE",
		},
		{
			name:           "gets unknown from undefinedFileUploadStatus game state",
			input:          undefinedFileUploadStatus,
			expectedOutput: "UNDEFINED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileUploadStatusString := tt.input.String()
			assert.Equal(t, fileUploadStatusString, tt.expectedOutput)
		})
	}
}

func Test_BotType_Valid(t *testing.T) {
	t.Run("returns true for a valid file upload status", func(t *testing.T) {
		assert.True(t, failure.Valid())
	})

	t.Run("returns false for a invalid file upload status", func(t *testing.T) {
		assert.False(t, undefinedFileUploadStatus.Valid())
	})
}
