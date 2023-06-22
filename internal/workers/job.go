package workers

import (
	"github.com/gocraft/work"
	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type jobContext struct{}

func (j *jobContext) processFileUpload(job *work.Job) error {
	fileUploadId := job.ArgString("fileUploadId")

	if utilities.IsBlank(fileUploadId) {
		err := errors.New("fileUploadId is required")
		logger.LogError(err)
		return err
	}

	tx, err := workerStorage.BeginTransaction()
	if err != nil {
		logger.LogError(err)
		return err
	}
	defer tx.Rollback()

	fileUpload, err := workerStorage.GetFileUploadUsingTx(fileUploadId, tx)
	if err != nil {
		logger.LogError(err)
		return err
	}

	if fileUpload.ProcessingOngoing() || fileUpload.ProcessingFinised() {
		err = errors.New("fileUpload is in incorrect processing state")
		logger.LogError(err)
		return err
	}

	err = workerStorage.UpdateFileUploadWithProcessingStatusUsingTx(fileUploadId, "ONGOING", tx)
	if err != nil {
		logger.LogError(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		logger.LogError(err)
		return err
	}

	// TODO: Get PDF file from storage
	// TODO: Parse PDF
	// TODO: Make call to Open AI
	// TODO: Create Candidate object
	// TODO: Update FileUpload
	// err = workerStorage.UpdateFileUploadWithProcessingStatus(fileUploadId, "COMPLETED")
	// if err != nil {
	// 	logger.LogError(err)
	// 	return err
	// }

	return nil
}
