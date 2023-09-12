package model

import (
	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type Team struct {
	id               string
	name             string
	currentFileCount int
	fileCountLimit   int
}

type TeamOptions struct {
	Id               string
	Name             string
	CurrentFileCount *int
	FileCountLimit   int
}

func NewTeam(opts TeamOptions) (*Team, error) {
	if utilities.IsBlank(opts.Id) {
		return nil, errors.New("cannot create team with an empty id")
	}

	if utilities.IsBlank(opts.Name) {
		return nil, errors.New("cannot create team with an empty name")
	}

	if opts.CurrentFileCount == nil {
		return nil, errors.New("cannot create team without explicit current file count")
	}

	if opts.FileCountLimit == 0 {
		return nil, errors.New("cannot create team with 0 file count limit")
	}

	return &Team{
		id:               opts.Id,
		name:             opts.Name,
		currentFileCount: *opts.CurrentFileCount,
		fileCountLimit:   opts.FileCountLimit,
	}, nil
}

func (t *Team) Id() string {
	return t.id
}
