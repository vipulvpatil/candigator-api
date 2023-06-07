package model

import (
	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type Team struct {
	id   string
	name string
}

type TeamOptions struct {
	Id   string
	Name string
}

func NewTeam(opts TeamOptions) (*Team, error) {
	if utilities.IsBlank(opts.Id) {
		return nil, errors.New("cannot create team with an empty id")
	}

	if utilities.IsBlank(opts.Name) {
		return nil, errors.New("cannot create team with an empty name")
	}

	return &Team{
		id:   opts.Id,
		name: opts.Name,
	}, nil
}

func (t *Team) Id() string {
	return t.id
}
