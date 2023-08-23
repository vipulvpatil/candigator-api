package storage

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type CandidateAccessor interface {
	CreateCandidateWithAiGeneratedPersonaForTeamUsingTx(persona *model.Persona, team *model.Team, tx DatabaseTransaction) error
	GetCandidatesForTeam(team *model.Team) ([]*model.Candidate, error)
	UpdateCandidateForTeam(id string, persona *model.Persona, team *model.Team) error
}

func (s *Storage) CreateCandidateWithAiGeneratedPersonaForTeamUsingTx(persona *model.Persona, team *model.Team, tx DatabaseTransaction) error {
	id := s.IdGenerator.Generate()
	if !persona.IsValid() {
		return errors.New("cannot create Candidate without a valid persona")
	}

	_, err := model.NewCandidate(model.CandidateOptions{
		Id:                 id,
		Team:               team,
		AiGeneratedPersona: persona,
		FileUploadId:       persona.FileUploadId,
	})
	if err != nil {
		return err
	}

	result, err := s.db.Exec(
		`INSERT INTO public."candidates"
		("id", "ai_generated_persona", "team_id", "file_upload_id")
		VALUES
		($1, $2, $3, $4)`,
		id, persona, team.Id(), persona.FileUploadId,
	)
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while inserting Candidate: %s", id))
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while inserting Candidate and changing db: %s", id))
	}
	if rowsAffected != 1 {
		return utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when inserting Candidate in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}

	return nil
}

func (s *Storage) GetCandidatesForTeam(team *model.Team) ([]*model.Candidate, error) {
	if team == nil || utilities.IsBlank(team.Id()) {
		return nil, errors.New("team cannot be blank")
	}

	rows, err := s.db.Query(
		`SELECT id, ai_generated_persona, manually_created_persona, file_upload_id
		FROM public."candidates"
		WHERE team_id = $1 ORDER BY created_at ASC, id ASC`,
		team.Id(),
	)
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to select candidates")
	}
	defer rows.Close()

	candidates := []*model.Candidate{}

	for rows.Next() {
		var id string
		var aiGeneratedPersona, manuallyCreatedPersona model.Persona
		var fileUploadId sql.NullString
		err := rows.Scan(&id, &aiGeneratedPersona, &manuallyCreatedPersona, &fileUploadId)

		if err != nil {
			return nil, utilities.WrapBadError(err, "failed while scanning rows")
		}

		var fileUploadIdString string
		if fileUploadId.Valid {
			fileUploadIdString = fileUploadId.String
		}

		candidate, err := model.NewCandidate(model.CandidateOptions{
			Id:                     id,
			AiGeneratedPersona:     &aiGeneratedPersona,
			ManuallyCreatedPersona: &manuallyCreatedPersona,
			Team:                   team,
			FileUploadId:           fileUploadIdString,
		})

		if err != nil {
			// TODO: Log this error?
			continue
		}

		candidates = append(candidates, candidate)
	}

	err = rows.Err()
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to correctly go through candidates rows")
	}
	return candidates, nil
}

func (s *Storage) UpdateCandidateForTeam(id string, persona *model.Persona, team *model.Team) error {
	return nil
}
