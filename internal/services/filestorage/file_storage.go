package filestorage

import (
	"github.com/vipulvpatil/candidate-tracker-go/internal/clients/s3"
)

type FileStorer interface {
	GetPresignedUrl(teamId, fileId, fileName string) (string, error)
}

type fileStorage struct {
	s3Client s3.Client
}

func NewFileStorage(s3client s3.Client) (*fileStorage, error) {
	return &fileStorage{
		s3Client: s3client,
	}, nil
}

func (f *fileStorage) GetPresignedUrl(teamId, fileId, fileName string) (string, error) {
	path := teamId + "/" + fileId
	return f.s3Client.GetPresignedUploadUrl(path, fileName)
}
