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
	GetFileUploadUsingTx(id string, tx DatabaseTransaction) (*model.FileUpload, error)
	GetFileUploadsForTeam(team *model.Team) ([]*model.FileUpload, error)
	GetUnprocessedFileUploadsCountForTeam(team *model.Team) (int, error)
	GetAllProcessingNotStartedFileUploadIds() ([]string, error)
	CreateFileUploadForTeam(name string, team *model.Team) (*model.FileUpload, error)
	UpdateFileUploadWithPresignedUrl(id, presignedUrl string) error
	UpdateFileUploadWithStatus(id, status string) error
	UpdateFileUploadWithProcessingStatus(id, processingStatus string) error
	UpdateFileUploadWithProcessingStatusUsingTx(id, processingStatus string, tx DatabaseTransaction) error
}

func (s *Storage) GetFileUpload(id string) (*model.FileUpload, error) {
	return getFileUploadUsingCustomDbHandler(s.db, id, false)
}

func (s *Storage) GetFileUploadUsingTx(id string, tx DatabaseTransaction) (*model.FileUpload, error) {
	return getFileUploadUsingCustomDbHandler(tx, id, true)
}

func getFileUploadUsingCustomDbHandler(customDb customDbHandler, id string, exclusiveLock bool) (*model.FileUpload, error) {
	if utilities.IsBlank(id) {
		return nil, errors.New("id cannot be blank")
	}

	var name, status, presignedUrl, teamId, teamName, processingStatus string

	queryWithoutLock := `SELECT
	f.name, f.status, f.presigned_url, f.processing_status, teams.id, teams.name
	FROM public."file_uploads" AS f
	JOIN public."teams" ON f.team_id = teams.id
	WHERE f.id = $1 ORDER BY f.created_at ASC LIMIT 1`

	queryWithLock := `SELECT
	f.name, f.status, f.presigned_url, f.processing_status, teams.id, teams.name
	FROM public."file_uploads" AS f
	JOIN public."teams" ON f.team_id = teams.id
	WHERE f.id = $1 ORDER BY f.created_at ASC LIMIT 1
	FOR UPDATE`

	var query string

	if exclusiveLock {
		query = queryWithLock
	} else {
		query = queryWithoutLock
	}

	row := customDb.QueryRow(
		query, id,
	)
	err := row.Scan(&name, &status, &presignedUrl, &processingStatus, &teamId, &teamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
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
		Id:               id,
		Name:             name,
		PresignedUrl:     presignedUrl,
		ProcessingStatus: processingStatus,
		Status:           status,
		Team:             team,
	})
}

func (s *Storage) GetFileUploadsForTeam(team *model.Team) ([]*model.FileUpload, error) {
	if team == nil || utilities.IsBlank(team.Id()) {
		return nil, errors.New("team cannot be blank")
	}

	rows, err := s.db.Query(
		`SELECT f.id, f.name, f.status, f.presigned_url, f.processing_status
		FROM public."file_uploads" AS f
		WHERE f.team_id = $1 ORDER BY f.created_at ASC, f.id ASC`,
		team.Id(),
	)
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to select file_uploads")
	}
	defer rows.Close()

	fileUploads := []*model.FileUpload{}

	for rows.Next() {
		var id, name, status, presignedUrl, processingStatus string
		err := rows.Scan(&id, &name, &status, &presignedUrl, &processingStatus)

		if err != nil {
			return nil, utilities.WrapBadError(err, "failed while scanning rows")
		}

		fileUpload, err := model.NewFileUpload(model.FileUploadOptions{
			Id:               id,
			Name:             name,
			PresignedUrl:     presignedUrl,
			ProcessingStatus: processingStatus,
			Status:           status,
			Team:             team,
		})

		if err != nil {
			// TODO: Log this error?
			continue
		}

		fileUploads = append(fileUploads, fileUpload)
	}

	err = rows.Err()
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to correctly go through file_upload rows")
	}
	return fileUploads, nil
}

func (s *Storage) GetUnprocessedFileUploadsCountForTeam(team *model.Team) (int, error) {
	if team == nil || utilities.IsBlank(team.Id()) {
		return 0, errors.New("team cannot be blank")
	}

	var count int
	row := s.db.QueryRow(
		`SELECT count(id)
		FROM public."file_uploads"
		WHERE team_id = $1
		AND processing_status <> 'COMPLETED'`,
		team.Id(),
	)
	err := row.Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, errors.Errorf("getting unprocessed file upload count for team %s: %v", team.Id(), err)
	}

	return count, nil
}

func (s *Storage) GetAllProcessingNotStartedFileUploadIds() ([]string, error) {
	rows, err := s.db.Query(
		`SELECT id
		FROM public."file_uploads"
		WHERE status = 'SUCCESS'
		AND processing_status = 'NOT STARTED'
		ORDER BY created_at ASC, id ASC`,
	)
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to select file_upload ids")
	}
	defer rows.Close()

	fileUploadIds := []string{}

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, utilities.WrapBadError(err, "failed while scanning rows")
		}

		fileUploadIds = append(fileUploadIds, id)
	}

	err = rows.Err()
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to correctly go through file_upload id rows")
	}
	return fileUploadIds, nil
}

func (s *Storage) CreateFileUploadForTeam(name string, team *model.Team) (*model.FileUpload, error) {
	id := s.IdGenerator.Generate()
	initialFileUploadStatus := "INITIATED"
	initialFileUploadProcessingStatus := "NOT STARTED"

	newFileUpload, err := model.NewFileUpload(model.FileUploadOptions{
		Id:               id,
		Name:             name,
		Status:           initialFileUploadStatus,
		ProcessingStatus: initialFileUploadProcessingStatus,
		Team:             team,
	})
	if err != nil {
		return nil, err
	}

	result, err := s.db.Exec(
		`INSERT INTO public."file_uploads"
		("id", "name", "presigned_url", "status", "processing_status", "team_id")
		VALUES
		($1, $2, '', $3, $4, $5)`,
		id, name, initialFileUploadStatus, initialFileUploadProcessingStatus, team.Id(),
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

func (s *Storage) UpdateFileUploadWithProcessingStatus(id, processingStatus string) error {
	return updateFileUploadWithProcessingStatusUsingCustomDbHandler(s.db, id, processingStatus)
}

func (s *Storage) UpdateFileUploadWithProcessingStatusUsingTx(id, processingStatus string, tx DatabaseTransaction) error {
	return updateFileUploadWithProcessingStatusUsingCustomDbHandler(tx, id, processingStatus)
}

func updateFileUploadWithProcessingStatusUsingCustomDbHandler(customDb customDbHandler, id, processingStatus string) error {
	if utilities.IsBlank(id) {
		return errors.New("id cannot be blank")
	}

	if !model.FileUploadProcessingStatus(processingStatus).Valid() {
		return errors.New("processing status should be valid")
	}

	result, err := customDb.Exec(`UPDATE public."file_uploads" SET "processing_status" = $2 WHERE id = $1`, id, processingStatus)
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while updating fileUpload: %s %s", id, processingStatus))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while checking affected row while updating fileUpload: %s %s", id, processingStatus))
	}

	if rowsAffected != 1 {
		return utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when inserting file_upload in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}
	return nil
}
