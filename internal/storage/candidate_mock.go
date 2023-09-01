package storage

import "github.com/vipulvpatil/candidate-tracker-go/internal/model"

type CandidateAccessorConfigurableMock struct {
	CreateCandidateWithAiGeneratedPersonaForTeamUsingTxInternal func(persona *model.Persona, team *model.Team, tx DatabaseTransaction) error
	GetCandidatesForTeamInternal                                func(team *model.Team) ([]*model.Candidate, error)
	GetCandidateForTeamInternal                                 func(id string, team *model.Team) (*model.Candidate, error)
	UpdateCandidateWithManuallyCreatedPersonaForTeamInternal    func(id string, persona *model.Persona, team *model.Team) (string, error)
}

func (c *CandidateAccessorConfigurableMock) CreateCandidateWithAiGeneratedPersonaForTeamUsingTx(persona *model.Persona, team *model.Team, tx DatabaseTransaction) error {
	return c.CreateCandidateWithAiGeneratedPersonaForTeamUsingTxInternal(persona, team, tx)
}

func (c *CandidateAccessorConfigurableMock) GetCandidatesForTeam(team *model.Team) ([]*model.Candidate, error) {
	return c.GetCandidatesForTeamInternal(team)
}

func (c *CandidateAccessorConfigurableMock) GetCandidateForTeam(id string, team *model.Team) (*model.Candidate, error) {
	return c.GetCandidateForTeamInternal(id, team)
}

func (c *CandidateAccessorConfigurableMock) UpdateCandidateWithManuallyCreatedPersonaForTeam(id string, persona *model.Persona, team *model.Team) (string, error) {
	return c.UpdateCandidateWithManuallyCreatedPersonaForTeamInternal(id, persona, team)
}
