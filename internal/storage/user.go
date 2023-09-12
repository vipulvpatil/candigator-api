package storage

import (
	"database/sql"
	"fmt"

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
	var teamId, teamName sql.NullString
	row := s.db.QueryRow(`
		SELECT users.id, users.email, teams.id, teams.name
		FROM public."users"
		LEFT JOIN public."teams" ON teams.id = users.team_id
		WHERE users.email = $1
	`, email)
	err := row.Scan(&userOptions.Id, &userOptions.Email, &teamId, &teamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Errorf("UserByEmail %s: no such user", email)
		}
		return nil, errors.Errorf("UserByEmail %s: %v", email, err)
	}

	if teamId.Valid && teamName.Valid {
		currentFileCount := 1
		team, err := model.NewTeam(model.TeamOptions{
			Id:   teamId.String,
			Name: teamName.String,
			// TODO: Verify this
			CurrentFileCount: &currentFileCount,
			FileCountLimit:   100,
		})
		if err != nil {
			return nil, utilities.WrapBadError(err, fmt.Sprintf("UserByEmail %s: invalid team options", email))
		}
		userOptions.Team = team
	}
	return model.NewUser(userOptions)
}
