package model

import (
	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type FileUpload struct {
	id           string
	name         string
	presignedUrl string
	status       fileUploadStatus
	team         *Team
}

type FileUploadOptions struct {
	Id           string
	Name         string
	PresignedUrl string
	Status       string
	Team         *Team
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
		status = waitingForFile
	}

	if opts.Team == nil {
		return nil, errors.New("cannot create FileUpload with a nil team")
	}

	return &FileUpload{
		id:           opts.Id,
		name:         opts.Name,
		presignedUrl: opts.PresignedUrl,
		status:       status,
		team:         opts.Team,
	}, nil
}

func (f *FileUpload) Id() string {
	return f.id
}

func (f *FileUpload) Status() string {
	return f.status.String()
}
