package storage

import (
	"github.com/pkg/errors"

	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
)

type UserRetrieverMockSuccess struct {
	Id    string
	Email string
}

func (u *UserRetrieverMockSuccess) UserByEmail(email string) (*model.User, error) {
	return model.NewUser(model.UserOptions{
		Id:    u.Id,
		Email: u.Email,
	})
}

type UserRetrieverMockFailure struct {
	Id    string
	Email string
}

func (u *UserRetrieverMockFailure) UserByEmail(email string) (*model.User, error) {
	return nil, errors.New("cannot find user by email")
}
