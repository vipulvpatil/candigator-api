package filestorage

import (
	"github.com/vipulvpatil/candidate-tracker-go/internal/clients/s3"
)

type FileStorer interface {
	GetPresignedUrl(path, fileName string) (string, error)
	GetFileData(path, fileName string) (string, error)
}

type fileStorage struct {
	s3Client s3.Client
}

func NewFileStorage(s3client s3.Client) (*fileStorage, error) {
	return &fileStorage{
		s3Client: s3client,
	}, nil
}

func (f *fileStorage) GetPresignedUrl(path, fileName string) (string, error) {
	return f.s3Client.GetPresignedUploadUrl(path, fileName)
}

func (f *fileStorage) GetFileData(path, fileName string) (string, error) {
	return f.s3Client.GetFileData(path, fileName)
}
