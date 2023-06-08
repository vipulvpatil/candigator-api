package storage

import (
	"errors"
	"fmt"

	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type FileUploadAccessor interface {
	CreateFileUploadForTeam(name string, team *model.Team) (*model.FileUpload, error)
	UpdateFileUploadWithPresignedUrl(id, presignedUrl string) (*model.FileUpload, error)
}

func (s *Storage) CreateFileUploadForTeam(name string, team *model.Team) (*model.FileUpload, error) {
	id := s.IdGenerator.Generate()
	initialFileUploadStatus := "WAITING_FOR_FILE"

	newFileUpload, err := model.NewFileUpload(model.FileUploadOptions{
		Id:     id,
		Name:   name,
		Status: initialFileUploadStatus,
		Team:   team,
	})
	if err != nil {
		return nil, err
	}

	result, err := s.db.Exec(
		`INSERT INTO public."file_uploads"
		("id", "name", "presigned_url", "status", "team_id")
		VALUES
		($1, $2, '', $3, $4)`,
		id, name, initialFileUploadStatus, team.Id(),
	)
	if err != nil {
		return nil, utilities.WrapBadError(err, fmt.Sprintf("dbError while inserting file_upload: %s", id))
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, utilities.WrapBadError(err, fmt.Sprintf("dbError while inserting file_upload and changing db: %s", id))
	}
	if rowsAffected != 1 {
		return nil, utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when inserting file_upload in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}

	return newFileUpload, nil
}

func (s *Storage) UpdateFileUploadWithPresignedUrl(id, presignedUrl string) error {
	if utilities.IsBlank(id) {
		return errors.New("id cannot be blank")
	}

	if utilities.IsBlank(presignedUrl) {
		return errors.New("presignedUrl cannot be blank")
	}

	result, err := s.db.Exec(`UPDATE public."file_uploads" SET "presigned_url" = $2 WHERE id = $1`, id, presignedUrl)
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while updating fileUpload: %s %s", id, presignedUrl))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while checking affected row while updating fileUpload: %s %s", id, presignedUrl))
	}

	if rowsAffected != 1 {
		return utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when inserting file_upload in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}
	return nil
}
