package model

import (
	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type Team struct {
	id   string
	name string
	user *User
}

type TeamOptions struct {
	Id   string
	Name string
	User *User
}

func NewTeam(opts TeamOptions) (*Team, error) {
	if utilities.IsBlank(opts.Id) {
		return nil, errors.New("cannot create team with an empty id")
	}

	if utilities.IsBlank(opts.Name) {
		return nil, errors.New("cannot create team with an empty name")
	}

	if opts.User == nil {
		return nil, errors.New("cannot create team with a nil user")
	}

	return &Team{
		id:   opts.Id,
		name: opts.Name,
		user: opts.User,
	}, nil
}
