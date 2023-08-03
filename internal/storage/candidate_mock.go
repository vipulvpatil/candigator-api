package storage

import "github.com/vipulvpatil/candidate-tracker-go/internal/model"

type CandidateAccessorConfigurableMock struct {
	CreateCandidateWithAiGeneratedPersonaForTeamUsingTxInternal func(persona *model.Persona, team *model.Team, tx DatabaseTransaction) error
}

func (c *CandidateAccessorConfigurableMock) CreateCandidateWithAiGeneratedPersonaForTeamUsingTx(persona *model.Persona, team *model.Team, tx DatabaseTransaction) error {
	return c.CreateCandidateWithAiGeneratedPersonaForTeamUsingTxInternal(persona, team, tx)
}
