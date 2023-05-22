package storage

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type UserRetriever interface {
	UserByEmail(email string) (*model.User, error)
}

func (s *Storage) UserByEmail(email string) (*model.User, error) {
	if utilities.IsBlank(email) {
		return nil, errors.New("cannot search by blank email")
	}

	userOptions := model.UserOptions{}
	row := s.db.QueryRow(`SELECT id, email FROM public."users" WHERE email = $1`, email)
	err := row.Scan(&userOptions.Id, &userOptions.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Errorf("UserByEmail %s: no such user", email)
		}
		return nil, errors.Errorf("UserByEmail %s: %v", email, err)
	}
	return model.NewUser(userOptions)
}
