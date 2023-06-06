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
			name:           "creates WAITING_FOR_FILE file upload status",
			input:          "WAITING_FOR_FILE",
			expectedOutput: waitingForFile,
		},
		{
			name:           "creates UPLOADING_FILE file upload status",
			input:          "UPLOADING_FILE",
			expectedOutput: uploadingFile,
		},
		{
			name:           "creates FILE_READY file upload status",
			input:          "FILE_READY",
			expectedOutput: fileReady,
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
			name:           "gets WAITING_FOR_FILE from ai game state",
			input:          waitingForFile,
			expectedOutput: "WAITING_FOR_FILE",
		},
		{
			name:           "gets UPLOADING_FILE from human game state",
			input:          uploadingFile,
			expectedOutput: "UPLOADING_FILE",
		},
		{
			name:           "gets FILE_READY from human game state",
			input:          fileReady,
			expectedOutput: "FILE_READY",
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
		assert.True(t, fileReady.Valid())
	})

	t.Run("returns false for a invalid file upload status", func(t *testing.T) {
		assert.False(t, undefinedFileUploadStatus.Valid())
	})
}
