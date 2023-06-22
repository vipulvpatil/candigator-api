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
			name: "processingStatus is invalid",
			input: FileUploadOptions{
				Id:               "123",
				Name:             "test",
				PresignedUrl:     "some_url",
				ProcessingStatus: "",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create FileUpload with an invalid processing status",
		},
		{
			name: "team is nil",
			input: FileUploadOptions{
				Id:               "123",
				Name:             "test",
				PresignedUrl:     "some_url",
				ProcessingStatus: "NOT STARTED",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create FileUpload with a nil team",
		},
		{
			name: "FileUpload gets created successfully",
			input: FileUploadOptions{
				Id:               "123",
				Name:             "test",
				PresignedUrl:     "some_url",
				ProcessingStatus: "NOT STARTED",
				Status:           "FAILURE",
				Team: &Team{
					id:   "team_id1",
					name: "test",
				},
			},
			expectedOutput: &FileUpload{
				id:               "123",
				name:             "test",
				presignedUrl:     "some_url",
				processingStatus: not_started,
				status:           failure,
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
				Id:               "123",
				Name:             "test",
				PresignedUrl:     "some_url",
				ProcessingStatus: "NOT STARTED",
				Team: &Team{
					id:   "team_id1",
					name: "test",
				},
			},
			expectedOutput: &FileUpload{
				id:               "123",
				name:             "test",
				presignedUrl:     "some_url",
				processingStatus: not_started,
				status:           initiated,
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

func Test_FileUpload_Id(t *testing.T) {
	t.Run("Id returns fileUpload's id", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: not_started,
			status:           initiated,
		}
		assert.Equal(t, "fp_id1", fileUpload.Id())
	})
}

func Test_FileUpload_Name(t *testing.T) {
	t.Run("Id returns fileUpload's id", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: not_started,
			status:           initiated,
		}
		assert.Equal(t, "file1.pdf", fileUpload.Name())
	})
}

func Test_FileUpload_PresignedUrl(t *testing.T) {
	t.Run("Id returns fileUpload's id", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: not_started,
			status:           initiated,
		}
		assert.Equal(t, "http://presignedUrl1", fileUpload.PresignedUrl())
	})
}

func Test_FileUpload_Status(t *testing.T) {
	t.Run("Status returns fileUpload's status if valid", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: not_started,
			status:           initiated,
		}
		assert.Equal(t, "INITIATED", fileUpload.Status())
	})

	t.Run("Status returns undefined if status is invalid", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: not_started,
		}
		assert.Equal(t, "UNDEFINED", fileUpload.Status())
	})
}

func Test_FileUpload_ProcessingStatus(t *testing.T) {
	t.Run("ProcessingStatus returns fileUpload's processing status if valid", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: ongoing,
			status:           initiated,
		}
		assert.Equal(t, "ONGOING", fileUpload.ProcessingStatus())
	})

	t.Run("ProcessingStatus returns undefined if processing status is invalid", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:           "fp_id1",
			name:         "file1.pdf",
			presignedUrl: "http://presignedUrl1",
		}
		assert.Equal(t, "UNDEFINED", fileUpload.ProcessingStatus())
	})
}

func Test_FileUpload_Completed(t *testing.T) {
	t.Run("Completed returns true if status is success", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: not_started,
			status:           success,
		}
		assert.True(t, fileUpload.Completed())
	})

	t.Run("Completed returns true if status is failure", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: not_started,
			status:           failure,
		}
		assert.True(t, fileUpload.Completed())
	})

	t.Run("Completed returns false if status is initiated", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: not_started,
			status:           initiated,
		}
		assert.False(t, fileUpload.Completed())
	})
}

func Test_FileUpload_BelongsToTeam(t *testing.T) {
	t.Run("BelongsToTeam returns true", func(t *testing.T) {
		team := &Team{
			id:   "team_id1",
			name: "team1",
		}

		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: not_started,
			status:           initiated,
			team: &Team{
				id:   "team_id1",
				name: "team1",
			},
		}
		assert.True(t, fileUpload.BelongsToTeam(team))
	})

	t.Run("BelongsToTeam returns false", func(t *testing.T) {
		team := &Team{
			id:   "team_id1",
			name: "team1",
		}

		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: not_started,
			status:           initiated,
			team: &Team{
				id:   "team_id2",
				name: "team2",
			},
		}
		assert.False(t, fileUpload.BelongsToTeam(team))
	})
}

func Test_FileUpload_ProcessingOngoing(t *testing.T) {
	t.Run("ProcessingOngoing returns false if processingStatus is non_started", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: not_started,
			status:           success,
		}
		assert.False(t, fileUpload.ProcessingOngoing())
	})

	t.Run("ProcessingOngoing returns true if processingStatus is ongoing", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: ongoing,
			status:           failure,
		}
		assert.True(t, fileUpload.ProcessingOngoing())
	})

	t.Run("ProcessingOngoing returns false if processingStatus is completed", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: completed,
			status:           initiated,
		}
		assert.False(t, fileUpload.ProcessingOngoing())
	})

	t.Run("ProcessingOngoing returns false if processingStatus is failed", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: failed,
			status:           initiated,
		}
		assert.False(t, fileUpload.ProcessingOngoing())
	})
}

func Test_FileUpload_ProcessingFinised(t *testing.T) {
	t.Run("ProcessingFinised returns false if processingStatus is non_started", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: not_started,
			status:           success,
		}
		assert.False(t, fileUpload.ProcessingFinised())
	})

	t.Run("ProcessingFinised returns false if processingStatus is ongoing", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: ongoing,
			status:           failure,
		}
		assert.False(t, fileUpload.ProcessingFinised())
	})

	t.Run("ProcessingFinised returns true if processingStatus is completed", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: completed,
			status:           initiated,
		}
		assert.True(t, fileUpload.ProcessingFinised())
	})

	t.Run("ProcessingFinised returns true if processingStatus is failed", func(t *testing.T) {
		fileUpload := &FileUpload{
			id:               "fp_id1",
			name:             "file1.pdf",
			presignedUrl:     "http://presignedUrl1",
			processingStatus: failed,
			status:           initiated,
		}
		assert.True(t, fileUpload.ProcessingFinised())
	})
}
