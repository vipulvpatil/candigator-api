package model

type fileUploadProcessingStatus int64

const (
	undefinedFileUploadProcessingStatus fileUploadProcessingStatus = iota
	not_started
	ongoing
	completed
	failed
)

func FileUploadProcessingStatus(str string) fileUploadProcessingStatus {
	switch str {
	case "NOT STARTED":
		return not_started
	case "ONGOING":
		return ongoing
	case "COMPLETED":
		return completed
	case "FAILED":
		return failed
	default:
		return undefinedFileUploadProcessingStatus
	}
}

func (b fileUploadProcessingStatus) String() string {
	switch b {
	case not_started:
		return "NOT STARTED"
	case ongoing:
		return "ONGOING"
	case completed:
		return "COMPLETED"
	case failed:
		return "FAILED"
	default:
		return "UNDEFINED"
	}
}

func (b fileUploadProcessingStatus) Valid() bool {
	return b.String() != "UNDEFINED"
}
