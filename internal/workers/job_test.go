package workers

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/clients/openai"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/services/filestorage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

func Test_updateFileUploadToProcessing(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "test@example.com",
	})
	fileUpload, _ := model.NewFileUpload(model.FileUploadOptions{
		Id:               "fp_id1",
		Name:             "file1.pdf",
		PresignedUrl:     "https://presigned_url1",
		Status:           "INITIATED",
		ProcessingStatus: "NOT STARTED",
		Team:             team,
	})

	tests := []struct {
		name                   string
		input                  string
		output                 *model.FileUpload
		fileUploadAccessorMock storage.FileUploadAccessor
		txMock                 *storage.DatabaseTransactionMock
		txShouldCommit         bool
		errorExpected          bool
		errorString            string
	}{
		{
			name:                   "errors if fileUploadId is blank",
			input:                  "",
			output:                 nil,
			fileUploadAccessorMock: nil,
			txMock:                 nil,
			txShouldCommit:         false,
			errorExpected:          true,
			errorString:            "fileUploadId is required",
		},
		{
			name:                   "errors if unable to get transaction",
			input:                  "fp_id1",
			output:                 nil,
			fileUploadAccessorMock: nil,
			txMock:                 nil,
			txShouldCommit:         false,
			errorExpected:          true,
			errorString:            "unable to begin a db transaction",
		},
		{
			name:   "errors if fileUpload is not in db",
			input:  "fp_id1",
			output: nil,
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetFileUploadUsingTxInternal: func(string, storage.DatabaseTransaction) (*model.FileUpload, error) {
					return nil, errors.New("fileUpload not in db")
				},
			},
			txMock:         &storage.DatabaseTransactionMock{},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "fileUpload not in db",
		},
		{
			name:   "errors if fileUpload is in wrong state",
			input:  "fp_id1",
			output: nil,
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetFileUploadUsingTxInternal: func(string, storage.DatabaseTransaction) (*model.FileUpload, error) {
					fileUpload, _ := model.NewFileUpload(model.FileUploadOptions{
						Id:               "fp_id1",
						Name:             "file1.pdf",
						PresignedUrl:     "https://presigned_url1",
						Status:           "INITIATED",
						ProcessingStatus: "ONGOING",
						Team:             team,
					})
					return fileUpload, nil
				},
			},
			txMock:         &storage.DatabaseTransactionMock{},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "fileUpload is in incorrect processing state: fp_id1",
		},
		{
			name:   "errors if unable to update fileUpload",
			input:  "fp_id1",
			output: nil,
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetFileUploadUsingTxInternal: func(string, storage.DatabaseTransaction) (*model.FileUpload, error) {
					fileUpload, _ := model.NewFileUpload(model.FileUploadOptions{
						Id:               "fp_id1",
						Name:             "file1.pdf",
						PresignedUrl:     "https://presigned_url1",
						Status:           "INITIATED",
						ProcessingStatus: "NOT STARTED",
						Team:             team,
					})
					return fileUpload, nil
				},
				UpdateFileUploadWithProcessingStatusUsingTxInternal: func(id, processingStatus string, tx storage.DatabaseTransaction) error {
					return errors.New("unable to update file upload")
				},
			},
			txMock:         &storage.DatabaseTransactionMock{},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "unable to update file upload",
		},
		{
			name:   "successfully returns a file upload",
			input:  "fp_id1",
			output: fileUpload,
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetFileUploadUsingTxInternal: func(string, storage.DatabaseTransaction) (*model.FileUpload, error) {
					fileUpload, _ := model.NewFileUpload(model.FileUploadOptions{
						Id:               "fp_id1",
						Name:             "file1.pdf",
						PresignedUrl:     "https://presigned_url1",
						Status:           "INITIATED",
						ProcessingStatus: "NOT STARTED",
						Team:             team,
					})
					return fileUpload, nil
				},
				UpdateFileUploadWithProcessingStatusUsingTxInternal: func(id, processingStatus string, tx storage.DatabaseTransaction) error {
					return nil
				},
			},
			txMock:         &storage.DatabaseTransactionMock{},
			txShouldCommit: true,
			errorExpected:  false,
			errorString:    "",
		},
	}

	for _, tt := range tests {
		logger = &utilities.NullLogger{}
		workerStorage = storage.NewStorageAccessorMock(
			storage.WithDatabaseTransactionProviderMock(&storage.DatabaseTransactionProviderMock{
				Transaction: tt.txMock,
			}),
			storage.WithFileUploadAccessorMock(tt.fileUploadAccessorMock),
		)

		t.Run(tt.name, func(t *testing.T) {
			fileUpload, err := updateFileUploadToProcessing(tt.input)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.output, fileUpload)
			}

			if tt.txMock != nil {
				if tt.txShouldCommit {
					assert.True(t, tt.txMock.Committed, "transaction should have committed")
				} else {
					assert.True(t, tt.txMock.Rolledback, "transaction should have rolledback")
					assert.False(t, tt.txMock.Committed, "transaction should not have committed")
				}
			}
		})
	}
}

func Test_processFileUploadUsingAi(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "test@example.com",
	})
	fileUpload, _ := model.NewFileUpload(model.FileUploadOptions{
		Id:               "fp_id1",
		Name:             "file1.pdf",
		PresignedUrl:     "https://presigned_url1",
		Status:           "INITIATED",
		ProcessingStatus: "ONGOING",
		Team:             team,
	})
	personaJson := `{
		"Name": "Person",
		"Email": "someemail@example.com",
		"Phone": "+91-1234567890",
		"City": "",
		"State": "",
		"Country": "",
		"YoE": 9,
		"Tech Skills": [
			"Go",
			"NodeJS",
			"Ruby on Rails",
			"Xcode",
			"ObjC"
		],
		"Soft Skills": [],
		"Recommended Roles": [
			"Senior Software Engineer",
			"Backend Developer",
			"iOS Developer"
		],
		"Education": [
			{
				"Qualification": "B.E., Computer Engineering",
				"CompletionYear": "2008",
				"Institute": "Ramrao Adik Institute of Technology"
			},
			{
				"Qualification": "HSC, Science",
				"CompletionYear": "2004",
				"Institute": "Ramnivas Ruia Junior College"
			},
			{
				"Qualification": "SSC",
				"CompletionYear": "2002",
				"Institute": "St. John The Baptist High School"
			}
		]
	}`
	tests := []struct {
		name                   string
		input                  *model.FileUpload
		fileUploadAccessorMock storage.FileUploadAccessor
		candidateAccessorMock  storage.CandidateAccessor
		fileStorerMock         filestorage.FileStorer
		openAiClientMock       openai.Client
		txMock                 *storage.DatabaseTransactionMock
		txShouldCommit         bool
		errorExpected          bool
		errorString            string
	}{
		{
			name:                   "errors if fileUpload is nil",
			input:                  nil,
			fileUploadAccessorMock: nil,
			txMock:                 nil,
			txShouldCommit:         false,
			errorExpected:          true,
			errorString:            "fileUpload is required",
		},
		{
			name:                   "errors if unable to get transaction",
			input:                  fileUpload,
			fileUploadAccessorMock: nil,
			txMock:                 nil,
			txShouldCommit:         false,
			errorExpected:          true,
			errorString:            "unable to begin a db transaction",
		},
		{
			name:                   "errors if unable to get fileUpload local path",
			input:                  fileUpload,
			fileUploadAccessorMock: nil,
			fileStorerMock:         &filestorage.FileStorerMock{},
			txMock:                 &storage.DatabaseTransactionMock{},
			txShouldCommit:         false,
			errorExpected:          true,
			errorString:            "unable to get LocalFilePath",
		},
		{
			name:                   "errors if unable to get text from pdf file",
			input:                  fileUpload,
			fileUploadAccessorMock: nil,
			fileStorerMock: &filestorage.FileStorerMock{
				LocalFilePath: "invalid_path.pdf",
			},
			txMock:         &storage.DatabaseTransactionMock{},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "exit status 1",
		},
		{
			name:                   "errors if unable to build persona",
			input:                  fileUpload,
			fileUploadAccessorMock: nil,
			fileStorerMock: &filestorage.FileStorerMock{
				LocalFilePath: "test_fixtures/test-resume.pdf",
			},
			openAiClientMock: &openai.MockClientSuccess{
				Text: "what",
			},
			txMock:         &storage.DatabaseTransactionMock{},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "unable to parse response: invalid character 'w' looking for beginning of value",
		},
		{
			name:                   "errors if unable to create candidate",
			input:                  fileUpload,
			fileUploadAccessorMock: nil,
			candidateAccessorMock: &storage.CandidateAccessorConfigurableMock{
				CreateCandidateWithAiGeneratedPersonaForTeamUsingTxInternal: func(persona *model.Persona, team *model.Team, tx storage.DatabaseTransaction) error {
					return errors.New("unable to create candidate")
				},
			},
			fileStorerMock: &filestorage.FileStorerMock{
				LocalFilePath: "test_fixtures/test-resume.pdf",
			},
			openAiClientMock: &openai.MockClientSuccess{
				Text: personaJson,
			},
			txMock:         &storage.DatabaseTransactionMock{},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "unable to create candidate",
		},
		{
			name:  "errors if unable to update file upload",
			input: fileUpload,
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				UpdateFileUploadWithProcessingStatusUsingTxInternal: func(id, processingStatus string, tx storage.DatabaseTransaction) error {
					return errors.New("unable to update file upload")
				},
			},
			candidateAccessorMock: &storage.CandidateAccessorConfigurableMock{
				CreateCandidateWithAiGeneratedPersonaForTeamUsingTxInternal: func(persona *model.Persona, team *model.Team, tx storage.DatabaseTransaction) error {
					return nil
				},
			},
			fileStorerMock: &filestorage.FileStorerMock{
				LocalFilePath: "test_fixtures/test-resume.pdf",
			},
			openAiClientMock: &openai.MockClientSuccess{
				Text: personaJson,
			},
			txMock:         &storage.DatabaseTransactionMock{},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "unable to update file upload",
		},
		{
			name:  "success",
			input: fileUpload,
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				UpdateFileUploadWithProcessingStatusUsingTxInternal: func(id, processingStatus string, tx storage.DatabaseTransaction) error {
					return nil
				},
			},
			candidateAccessorMock: &storage.CandidateAccessorConfigurableMock{
				CreateCandidateWithAiGeneratedPersonaForTeamUsingTxInternal: func(persona *model.Persona, team *model.Team, tx storage.DatabaseTransaction) error {
					return nil
				},
			},
			fileStorerMock: &filestorage.FileStorerMock{
				LocalFilePath: "test_fixtures/test-resume.pdf",
			},
			openAiClientMock: &openai.MockClientSuccess{
				Text: personaJson,
			},
			txMock:         &storage.DatabaseTransactionMock{},
			txShouldCommit: true,
			errorExpected:  false,
			errorString:    "",
		},
	}

	for _, tt := range tests {
		logger = &utilities.NullLogger{}
		workerStorage = storage.NewStorageAccessorMock(
			storage.WithDatabaseTransactionProviderMock(&storage.DatabaseTransactionProviderMock{
				Transaction: tt.txMock,
			}),
			storage.WithFileUploadAccessorMock(tt.fileUploadAccessorMock),
			storage.WithCandidateAccessorMock(tt.candidateAccessorMock),
		)
		fileStorer = tt.fileStorerMock
		openAiClient = tt.openAiClientMock

		t.Run(tt.name, func(t *testing.T) {
			err := processFileUploadUsingAi(tt.input)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}

			if tt.txMock != nil {
				if tt.txShouldCommit {
					assert.True(t, tt.txMock.Committed, "transaction should have committed")
				} else {
					assert.True(t, tt.txMock.Rolledback, "transaction should have rolledback")
					assert.False(t, tt.txMock.Committed, "transaction should not have committed")
				}
			}
		})
	}
}
