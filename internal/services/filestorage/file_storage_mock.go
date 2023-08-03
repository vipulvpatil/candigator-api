package filestorage

import (
	"errors"

	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type FileStorerMock struct {
	PresignedUrl  string
	LocalFilePath string
}

func (f *FileStorerMock) GetPresignedUrl(path, fileName string) (string, error) {
	if utilities.IsBlank(f.PresignedUrl) {
		return "", errors.New("unable to get PresignedUrl")
	}
	return f.PresignedUrl, nil
}

func (f *FileStorerMock) GetLocalFilePath(path, fileName string) (string, error) {
	if utilities.IsBlank(f.LocalFilePath) {
		return "", errors.New("unable to get LocalFilePath")
	}
	return f.LocalFilePath, nil
}
