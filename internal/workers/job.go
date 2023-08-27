package workers

import (
	"fmt"

	"github.com/gocraft/work"
	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/lib/parser"
	"github.com/vipulvpatil/candidate-tracker-go/internal/lib/parser/personabuilder"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type jobContext struct{}

func (j *jobContext) processFileUpload(job *work.Job) error {
	fileUploadId := job.ArgString("fileUploadId")

	fileUpload, err := updateFileUploadToProcessing(fileUploadId)
	if err != nil {
		logger.LogError(err)
		return err
	}

	err = processFileUploadUsingAi(fileUpload)
	if err != nil {
		logger.LogError(err)
		skippedErr := workerStorage.UpdateFileUploadWithProcessingStatus(fileUpload.Id(), "FAILED")
		if skippedErr != nil {
			logger.LogError(skippedErr)
		}
		return err
	}

	return nil
}

func updateFileUploadToProcessing(fileUploadId string) (*model.FileUpload, error) {
	if utilities.IsBlank(fileUploadId) {
		err := errors.New("fileUploadId is required")
		logger.LogError(err)
		return nil, err
	}

	tx, err := workerStorage.BeginTransaction()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	defer tx.Rollback()

	fileUpload, err := workerStorage.GetFileUploadUsingTx(fileUploadId, tx)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	if fileUpload.ProcessingOngoing() || fileUpload.ProcessingFinised() {
		err = fmt.Errorf("fileUpload is in incorrect processing state: %s", fileUpload.Id())
		logger.LogError(err)
		return nil, err
	}

	err = workerStorage.UpdateFileUploadWithProcessingStatusUsingTx(fileUploadId, "ONGOING", tx)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return fileUpload, tx.Commit()
}

func processFileUploadUsingAi(fileUpload *model.FileUpload) error {
	if fileUpload == nil {
		err := errors.New("fileUpload is required")
		logger.LogError(err)
		return err
	}

	tx, err := workerStorage.BeginTransaction()
	if err != nil {
		logger.LogError(err)
		return err
	}
	defer tx.Rollback()

	localFilePath, err := fileStorer.GetLocalFilePath(fileUpload.StoragePath(), fileUpload.Name())
	if err != nil {
		logger.LogError(err)
		return err
	}

	text, err := parser.GetTextFromPdf(localFilePath)
	if err != nil {
		logger.LogError(err)
		return err
	}

	logger.LogMessageln(text)

	persona, err := personabuilder.Build(text, openAiClient)
	if err != nil {
		logger.LogError(err)
		return err
	}

	persona.FileUploadId = fileUpload.Id()

	err = workerStorage.CreateCandidateWithAiGeneratedPersonaForTeamUsingTx(persona, fileUpload.Team(), tx)
	if err != nil {
		logger.LogError(err)
		return err
	}

	err = workerStorage.UpdateFileUploadWithProcessingStatusUsingTx(fileUpload.Id(), "COMPLETED", tx)
	if err != nil {
		logger.LogError(err)
		return err
	}

	return tx.Commit()
}
