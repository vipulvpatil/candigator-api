package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileUploadProcessingStatus(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput fileUploadProcessingStatus
	}{
		{
			name:           "creates NOT STARTED file upload status",
			input:          "NOT STARTED",
			expectedOutput: not_started,
		},
		{
			name:           "creates ONGOING file upload status",
			input:          "ONGOING",
			expectedOutput: ongoing,
		},
		{
			name:           "creates COMPLETED file upload status",
			input:          "COMPLETED",
			expectedOutput: completed,
		},
		{
			name:           "creates FAILED file upload status",
			input:          "FAILED",
			expectedOutput: failed,
		},
		{
			name:           "handles unknown file upload processing status",
			input:          "unknown",
			expectedOutput: undefinedFileUploadProcessingStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := FileUploadProcessingStatus(tt.input)
			assert.Equal(t, state, tt.expectedOutput)
		})
	}
}

func Test_FileUploadProcessingStatus_String(t *testing.T) {
	tests := []struct {
		name           string
		input          fileUploadProcessingStatus
		expectedOutput string
	}{
		{
			name:           "gets NOT STARTED from not_started file upload processing state",
			input:          not_started,
			expectedOutput: "NOT STARTED",
		},
		{
			name:           "gets ONGOING from ongoing file upload processing state",
			input:          ongoing,
			expectedOutput: "ONGOING",
		},
		{
			name:           "gets COMPLETED from completed file upload processing state",
			input:          completed,
			expectedOutput: "COMPLETED",
		},
		{
			name:           "gets FAILED from failed file upload processing state",
			input:          failed,
			expectedOutput: "FAILED",
		},
		{
			name:           "gets unknown from undefinedFileUploadProcessingStatus file upload processing state",
			input:          undefinedFileUploadProcessingStatus,
			expectedOutput: "UNDEFINED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileUploadProcessingStatusString := tt.input.String()
			assert.Equal(t, fileUploadProcessingStatusString, tt.expectedOutput)
		})
	}
}

func Test_FileUploadProcessingStatus_Valid(t *testing.T) {
	t.Run("returns true for a valid file upload processing status", func(t *testing.T) {
		assert.True(t, completed.Valid())
	})

	t.Run("returns false for a invalid file upload processing status", func(t *testing.T) {
		assert.False(t, undefinedFileUploadProcessingStatus.Valid())
	})
}
