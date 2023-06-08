package storage

import (
	"errors"

	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
)

type TeamHydratorMockSuccess struct {
	User *model.User
}

func (t *TeamHydratorMockSuccess) HydrateTeam(user *model.User) (*model.User, error) {
	return t.User, nil
}

type TeamHydratorMockFailure struct {
}

func (t *TeamHydratorMockFailure) HydrateTeam(user *model.User) (*model.User, error) {
	return nil, errors.New("unable to hydrate team")
}
