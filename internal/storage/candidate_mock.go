package storage

import "github.com/vipulvpatil/candidate-tracker-go/internal/model"

type CandidateAccessorConfigurableMock struct {
	CreateCandidateWithAiGeneratedPersonaForTeamUsingTxInternal func(persona *model.Persona, team *model.Team, tx DatabaseTransaction) error
	GetCandidatesForTeamInternal                                func(team *model.Team) ([]*model.Candidate, error)
}

func (c *CandidateAccessorConfigurableMock) CreateCandidateWithAiGeneratedPersonaForTeamUsingTx(persona *model.Persona, team *model.Team, tx DatabaseTransaction) error {
	return c.CreateCandidateWithAiGeneratedPersonaForTeamUsingTxInternal(persona, team, tx)
}

func (c *CandidateAccessorConfigurableMock) GetCandidatesForTeam(team *model.Team) ([]*model.Candidate, error) {
	return c.GetCandidatesForTeamInternal(team)
}
