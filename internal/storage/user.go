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
	var teamFileCountLimit, teamCurrentFileCount sql.NullInt64
	row := s.db.QueryRow(`
		SELECT users.id, users.email, t.id, t.name, t.file_count_limit, t.current_file_count
		FROM public."users"
		LEFT JOIN (
			SELECT
				teams.id,
				teams.name,
				teams.file_count_limit,
				teams.created_at,
				count(file_uploads.id) AS current_file_count
			FROM public."teams"
			LEFT JOIN
			public."file_uploads"
			ON teams.id = file_uploads.team_id
			GROUP BY teams.id
		) t
		ON t.id = users.team_id
		WHERE users.email = $1
		ORDER BY t.created_at ASC, t.id
	`, email)
	err := row.Scan(&userOptions.Id, &userOptions.Email, &teamId, &teamName, &teamFileCountLimit, &teamCurrentFileCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Errorf("UserByEmail %s: no such user", email)
		}
		return nil, errors.Errorf("UserByEmail %s: %v", email, err)
	}

	if teamId.Valid &&
		teamName.Valid &&
		teamCurrentFileCount.Valid &&
		teamFileCountLimit.Valid {
		currentFileCount := int(teamCurrentFileCount.Int64)
		team, err := model.NewTeam(model.TeamOptions{
			Id:               teamId.String,
			Name:             teamName.String,
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
