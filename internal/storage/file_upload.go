package storage

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type FileUploadAccessor interface {
	GetFileUpload(id string) (*model.FileUpload, error)
	CreateFileUploadForTeam(name string, team *model.Team) (*model.FileUpload, error)
	UpdateFileUploadWithPresignedUrl(id, presignedUrl string) error
	UpdateFileUploadWithStatus(id, status string) error
}

func (s *Storage) GetFileUpload(id string) (*model.FileUpload, error) {

	if utilities.IsBlank(id) {
		return nil, errors.New("id cannot be blank")
	}

	var name, status, presignedUrl, teamId, teamName string

	row := s.db.QueryRow(
		`SELECT f.name, f.status, f.presigned_url, teams.id, teams.name
		FROM public."file_uploads" AS f
		JOIN public."teams" ON f.team_id = teams.id
		WHERE f.id = $1 ORDER BY f.created_at ASC LIMIT 1`, id,
	)
	err := row.Scan(&name, &status, &presignedUrl, &teamId, &teamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Errorf("no file upload with id %s", id)
		}
		return nil, errors.Errorf("getting file upload for id %s: %v", id, err)
	}

	team, err := model.NewTeam(model.TeamOptions{
		Id:   teamId,
		Name: teamName,
	})

	if err != nil {
		return nil, err
	}

	return model.NewFileUpload(model.FileUploadOptions{
		Id:     id,
		Name:   name,
		Status: status,
		Team:   team,
	})
}

func (s *Storage) CreateFileUploadForTeam(name string, team *model.Team) (*model.FileUpload, error) {
	id := s.IdGenerator.Generate()
	initialFileUploadStatus := "INITIATED"

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

func (s *Storage) UpdateFileUploadWithStatus(id, status string) error {
	if utilities.IsBlank(id) {
		return errors.New("id cannot be blank")
	}

	if !model.FileUploadStatus(status).Valid() {
		return errors.New("status should be valid")
	}

	result, err := s.db.Exec(`UPDATE public."file_uploads" SET "status" = $2 WHERE id = $1`, id, status)
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while updating fileUpload: %s %s", id, status))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while checking affected row while updating fileUpload: %s %s", id, status))
	}

	if rowsAffected != 1 {
		return utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when inserting file_upload in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}
	return nil
}
