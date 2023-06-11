package filestorage

import (
	"errors"
)

type FileStorerMockSuccess struct {
	PresignedUrl string
}

func (f *FileStorerMockSuccess) GetPresignedUrl(teamId, fileId, fileName string) (string, error) {
	return f.PresignedUrl, nil
}

type FileStorerMockFailure struct{}

func (f *FileStorerMockFailure) GetPresignedUrl(teamId, fileId, fileName string) (string, error) {
	return "", errors.New("unable to get presignedUrl")
}
