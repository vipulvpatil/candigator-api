package storage

import "github.com/vipulvpatil/candidate-tracker-go/internal/model"

type FileUploadAccessorConfigurableMock struct {
	GetFileUploadInternal                    func(id string) (*model.FileUpload, error)
	GetFileUploadsForTeamInteral             func(team *model.Team) ([]*model.FileUpload, error)
	CreateFileUploadForTeamInteral           func(name string, team *model.Team) (*model.FileUpload, error)
	UpdateFileUploadWithPresignedUrlInternal func(id, presignedUrl string) error
	UpdateFileUploadWithStatusInternal       func(id, status string) error
}

func (f *FileUploadAccessorConfigurableMock) GetFileUpload(id string) (*model.FileUpload, error) {
	return f.GetFileUploadInternal(id)
}

func (f *FileUploadAccessorConfigurableMock) GetFileUploadsForTeam(team *model.Team) ([]*model.FileUpload, error) {
	return f.GetFileUploadsForTeamInteral(team)
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
