package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewFileUpload(t *testing.T) {
	tests := []struct {
		name           string
		input          FileUploadOptions
		expectedOutput *FileUpload
		errorExpected  bool
		errorString    string
	}{
		{
			name:           "id is empty",
			input:          FileUploadOptions{},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create FileUpload with an empty id",
		},
		{
			name: "name is empty",
			input: FileUploadOptions{
				Id: "123",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create FileUpload with an empty name",
		},
		{
			name: "presigned Url is empty",
			input: FileUploadOptions{
				Id:   "123",
				Name: "test",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create FileUpload with an empty presignedUrl",
		},
		{
			name: "user is nil",
			input: FileUploadOptions{
				Id:           "123",
				Name:         "test",
				PresignedUrl: "some_url",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create FileUpload with a nil team",
		},
		{
			name: "FileUpload gets created successfully",
			input: FileUploadOptions{
				Id:           "123",
				Name:         "test",
				PresignedUrl: "some_url",
				Status:       "FILE_READY",
				Team: &Team{
					id:   "team_id1",
					name: "test",
				},
			},
			expectedOutput: &FileUpload{
				id:           "123",
				name:         "test",
				presignedUrl: "some_url",
				status:       fileReady,
				team: &Team{
					id:   "team_id1",
					name: "test",
				},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "FileUpload gets created successfully with default status",
			input: FileUploadOptions{
				Id:           "123",
				Name:         "test",
				PresignedUrl: "some_url",
				Team: &Team{
					id:   "team_id1",
					name: "test",
				},
			},
			expectedOutput: &FileUpload{
				id:           "123",
				name:         "test",
				presignedUrl: "some_url",
				status:       waitingForFile,
				team: &Team{
					id:   "team_id1",
					name: "test",
				},
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewFileUpload(tt.input)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
