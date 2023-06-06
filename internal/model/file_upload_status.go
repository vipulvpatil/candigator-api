package model

type fileUploadStatus int64

const (
	undefinedFileUploadStatus fileUploadStatus = iota
	waitingForFile
	uploadingFile
	fileReady
)

func FileUploadStatus(str string) fileUploadStatus {
	switch str {
	case "WAITING_FOR_FILE":
		return waitingForFile
	case "UPLOADING_FILE":
		return uploadingFile
	case "FILE_READY":
		return fileReady
	default:
		return undefinedFileUploadStatus
	}
}

func (b fileUploadStatus) String() string {
	switch b {
	case waitingForFile:
		return "WAITING_FOR_FILE"
	case uploadingFile:
		return "UPLOADING_FILE"
	case fileReady:
		return "FILE_READY"
	default:
		return "UNDEFINED"
	}
}

func (b fileUploadStatus) Valid() bool {
	return b.String() != "UNDEFINED"
}
