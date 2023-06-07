package storage

import (
	"fmt"

	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type FileUploadAccessor interface {
	CreateFileUploadForTeam(user *model.User) (*model.FileUpload, error)
}

func (s *Storage) CreateFileUploadForTeam(name, presignedUrl string, team *model.Team) (*model.FileUpload, error) {
	id := s.IdGenerator.Generate()
	initialFileUploadStatus := "WAITING_FOR_FILE"

	newFileUpload, err := model.NewFileUpload(model.FileUploadOptions{
		Id:           id,
		Name:         name,
		PresignedUrl: presignedUrl,
		Status:       initialFileUploadStatus,
		Team:         team,
	})
	if err != nil {
		return nil, err
	}

	result, err := s.db.Exec(
		`INSERT INTO public."file_uploads"
		("id", "name", "presigned_url", "status", "team_id")
		VALUES
		($1, $2, $3, $4, $5)`,
		id, name, presignedUrl, initialFileUploadStatus, team.Id(),
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
