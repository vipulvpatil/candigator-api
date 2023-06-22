package filestorage

import (
	"errors"

	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type FileStorerMock struct {
	PresignedUrl string
	FileData     string
}

func (f *FileStorerMock) GetPresignedUrl(path, fileName string) (string, error) {
	if utilities.IsBlank(f.PresignedUrl) {
		return "", errors.New("unable to get PresignedUrl")
	}
	return f.PresignedUrl, nil
}

func (f *FileStorerMock) GetFileData(path, fileName string) (string, error) {
	if utilities.IsBlank(f.FileData) {
		return "", errors.New("unable to get FileData")
	}
	return f.FileData, nil
}
