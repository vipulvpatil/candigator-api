package workers

import (
	"fmt"

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
		err = fmt.Errorf("fileUpload is in incorrect processing state: %s", fileUpload.Id())
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

	logger.LogMessageln(fileUpload.StoragePath())
	logger.LogMessageln(fileUpload.Name())

	_, err = fileStorer.GetLocalFilePath(fileUpload.StoragePath(), fileUpload.Name())
	if err != nil {
		logger.LogError(err)
		return err
	}

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
