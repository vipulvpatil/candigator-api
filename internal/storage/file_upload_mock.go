package storage

import "github.com/vipulvpatil/candidate-tracker-go/internal/model"

type FileUploadAccessorConfigurableMock struct {
	CreateFileUploadForTeamInteral           func(name string, team *model.Team) (*model.FileUpload, error)
	UpdateFileUploadWithPresignedUrlInternal func(id, presignedUrl string) error
}

func (f *FileUploadAccessorConfigurableMock) CreateFileUploadForTeam(name string, team *model.Team) (*model.FileUpload, error) {
	return f.CreateFileUploadForTeamInteral(name, team)
}

func (f *FileUploadAccessorConfigurableMock) UpdateFileUploadWithPresignedUrl(id, presignedUrl string) error {
	return f.UpdateFileUploadWithPresignedUrlInternal(id, presignedUrl)
}
