package model

type fileUploadStatus int64

const (
	undefinedFileUploadStatus fileUploadStatus = iota
	initiated
	success
	failure
)

func FileUploadStatus(str string) fileUploadStatus {
	switch str {
	case "INITIATED":
		return initiated
	case "SUCCESS":
		return success
	case "FAILURE":
		return failure
	default:
		return undefinedFileUploadStatus
	}
}

func (b fileUploadStatus) String() string {
	switch b {
	case initiated:
		return "INITIATED"
	case success:
		return "SUCCESS"
	case failure:
		return "FAILURE"
	default:
		return "UNDEFINED"
	}
}

func (b fileUploadStatus) Valid() bool {
	return b.String() != "UNDEFINED"
}
