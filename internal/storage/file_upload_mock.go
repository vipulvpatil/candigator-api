package storage

import "github.com/vipulvpatil/candidate-tracker-go/internal/model"

type FileUploadAccessorConfigurableMock struct {
	GetFileUploadInternal                               func(id string) (*model.FileUpload, error)
	GetFileUploadUsingTxInternal                        func(id string, tx DatabaseTransaction) (*model.FileUpload, error)
	GetFileUploadsForTeamInteral                        func(team *model.Team) ([]*model.FileUpload, error)
	GetUnprocessedFileUploadsCountForTeamInternal       func(team *model.Team) (int, error)
	GetAllProcessingNotStartedFileUploadIdsInternal     func() ([]string, error)
	CreateFileUploadForTeamInteral                      func(name string, team *model.Team) (*model.FileUpload, error)
	UpdateFileUploadWithPresignedUrlInternal            func(id, presignedUrl string) error
	UpdateFileUploadWithStatusInternal                  func(id, status string) error
	UpdateFileUploadWithProcessingStatusInternal        func(id, processingStatus string) error
	UpdateFileUploadWithProcessingStatusUsingTxInternal func(id, processingStatus string, tx DatabaseTransaction) error
}

func (f *FileUploadAccessorConfigurableMock) GetFileUpload(id string) (*model.FileUpload, error) {
	return f.GetFileUploadInternal(id)
}

func (f *FileUploadAccessorConfigurableMock) GetFileUploadUsingTx(id string, tx DatabaseTransaction) (*model.FileUpload, error) {
	return f.GetFileUploadUsingTxInternal(id, tx)
}

func (f *FileUploadAccessorConfigurableMock) GetFileUploadsForTeam(team *model.Team) ([]*model.FileUpload, error) {
	return f.GetFileUploadsForTeamInteral(team)
}

func (f *FileUploadAccessorConfigurableMock) GetUnprocessedFileUploadsCountForTeam(team *model.Team) (int, error) {
	return f.GetUnprocessedFileUploadsCountForTeamInternal(team)
}

func (f *FileUploadAccessorConfigurableMock) GetAllProcessingNotStartedFileUploadIds() ([]string, error) {
	return f.GetAllProcessingNotStartedFileUploadIdsInternal()
}

func (f *FileUploadAccessorConfigurableMock) CreateFileUploadForTeam(name string, team *model.Team) (*model.FileUpload, error) {
	return f.CreateFileUploadForTeamInteral(name, team)
}

func (f *FileUploadAccessorConfigurableMock) UpdateFileUploadWithPresignedUrl(id, presignedUrl string) error {
	return f.UpdateFileUploadWithPresignedUrlInternal(id, presignedUrl)
}

func (f *FileUploadAccessorConfigurableMock) UpdateFileUploadWithStatus(id, status string) error {
	return f.UpdateFileUploadWithStatusInternal(id, status)
}

func (f *FileUploadAccessorConfigurableMock) UpdateFileUploadWithProcessingStatus(id, processingStatus string) error {
	return f.UpdateFileUploadWithProcessingStatusInternal(id, processingStatus)
}

func (f *FileUploadAccessorConfigurableMock) UpdateFileUploadWithProcessingStatusUsingTx(id, processingStatus string, tx DatabaseTransaction) error {
	return f.UpdateFileUploadWithProcessingStatusUsingTxInternal(id, processingStatus, tx)
}
