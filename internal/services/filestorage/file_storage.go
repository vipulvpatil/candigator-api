package filestorage

import "fmt"

type FileStorer interface {
	GetPresignedUrl(fileId string, teamId string) (string, error)
}

type fileStorage struct{}

func NewFileStorage() (*fileStorage, error) {
	return &fileStorage{}, nil
}

// TODO: This is a placeholder implementation. Please use S3 client to actually get Presigned Url
func (f *fileStorage) GetPresignedUrl(fileId string, teamId string) (string, error) {
	return fmt.Sprintf("http://%s/%s", teamId, fileId), nil
}
