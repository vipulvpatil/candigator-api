package model

import (
	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type FileUpload struct {
	id               string
	name             string
	presignedUrl     string
	processingStatus fileUploadProcessingStatus
	status           fileUploadStatus
	team             *Team
}

type FileUploadOptions struct {
	Id               string
	Name             string
	PresignedUrl     string
	ProcessingStatus string
	Status           string
	Team             *Team
}

func NewFileUpload(opts FileUploadOptions) (*FileUpload, error) {
	if utilities.IsBlank(opts.Id) {
		return nil, errors.New("cannot create FileUpload with an empty id")
	}

	if utilities.IsBlank(opts.Name) {
		return nil, errors.New("cannot create FileUpload with an empty name")
	}

	status := FileUploadStatus(opts.Status)
	if !status.Valid() {
		status = initiated
	}

	processingStatus := FileUploadProcessingStatus(opts.ProcessingStatus)
	if !processingStatus.Valid() {
		return nil, errors.New("cannot create FileUpload with an invalid processing status")
	}

	if opts.Team == nil {
		return nil, errors.New("cannot create FileUpload with a nil team")
	}

	return &FileUpload{
		id:               opts.Id,
		name:             opts.Name,
		presignedUrl:     opts.PresignedUrl,
		processingStatus: processingStatus,
		status:           status,
		team:             opts.Team,
	}, nil
}

func (f *FileUpload) Id() string {
	return f.id
}

func (f *FileUpload) Name() string {
	return f.name
}

func (f *FileUpload) PresignedUrl() string {
	return f.presignedUrl
}

func (f *FileUpload) Status() string {
	return f.status.String()
}

func (f *FileUpload) ProcessingStatus() string {
	return f.processingStatus.String()
}

func (f *FileUpload) Completed() bool {
	return f.status == success || f.status == failure
}

func (f *FileUpload) BelongsToTeam(t *Team) bool {
	return f.team.id == t.id
}
