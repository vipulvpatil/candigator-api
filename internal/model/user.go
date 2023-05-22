package model

import (
	"github.com/pkg/errors"

	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type User struct {
	id    string
	email string
}

type UserOptions struct {
	Id    string
	Email string
}

func NewUser(opts UserOptions) (*User, error) {
	if utilities.IsBlank(opts.Id) {
		return nil, errors.New("cannot create user with a empty id")
	}

	if utilities.IsBlank(opts.Email) {
		return nil, errors.New("cannot create user with a empty email")
	}

	return &User{
		id:    opts.Id,
		email: opts.Email,
	}, nil
}

func (u *User) GetId() string {
	return u.id
}
