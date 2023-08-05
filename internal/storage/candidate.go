package storage

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type CandidateAccessor interface {
	CreateCandidateWithAiGeneratedPersonaForTeamUsingTx(persona *model.Persona, team *model.Team, tx DatabaseTransaction) error
	GetCandidatesForTeam(team *model.Team) ([]*model.Candidate, error)
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
	return nil, nil
}
