package workers

import (
	"math/rand"
	"testing"

	"github.com/gocraft/work"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

func Test_processFileUploads(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "test@example.com",
	})
	tests := []struct {
		name                   string
		input                  map[string]interface{}
		fileUploadAccessorMock storage.FileUploadAccessor
		txMock                 *storage.DatabaseTransactionMock
		txShouldCommit         bool
		errorExpected          bool
		errorString            string
	}{
		{
			name: "errors if fileUploadId is blank",
			input: map[string]interface{}{
				"fileUploadId": "",
			},
			fileUploadAccessorMock: nil,
			txMock:                 nil,
			txShouldCommit:         false,
			errorExpected:          true,
			errorString:            "fileUploadId is required",
		},
		{
			name: "errors if unable to get transaction",
			input: map[string]interface{}{
				"fileUploadId": "fp_id1",
			},
			fileUploadAccessorMock: nil,
			txMock:                 nil,
			txShouldCommit:         false,
			errorExpected:          true,
			errorString:            "unable to begin a db transaction",
		},
		{
			name: "errors if fileUpload is not in db",
			input: map[string]interface{}{
				"fileUploadId": "fp_id1",
			},
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
			name: "errors if fileUpload is in wrong state",
			input: map[string]interface{}{
				"fileUploadId": "fp_id1",
			},
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
			errorString:    "fileUpload is in incorrect processing state",
		},
		{
			name: "errors if unable to update fileUpload",
			input: map[string]interface{}{
				"fileUploadId": "fp_id1",
			},
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
			rand.Seed(0)
			jc := jobContext{}
			err := jc.processFileUploads(&work.Job{
				Args: tt.input,
			})
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
