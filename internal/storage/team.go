package storage

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type TeamHydrator interface {
	HydrateTeam(user *model.User) (*model.User, error)
}

func (s *Storage) HydrateTeam(user *model.User) (*model.User, error) {
	if user == nil {
		return nil, errors.Errorf("cannot hydrate a nil user")
	}

	if user.Team() != nil {
		return user, nil
	}

	tx, err := s.BeginTransaction()
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to start db transaction")
	}
	defer tx.Rollback()

	userOpts := model.UserOptions{}
	var teamOpts model.TeamOptions
	var teamId, teamName sql.NullString
	var teamFileCountLimit, teamCurrentFileCount sql.NullInt64
	row := tx.QueryRow(`
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
		WHERE users.id = $1
		ORDER BY t.created_at ASC, t.id
		FOR UPDATE OF users
		LIMIT 1
	`, user.GetId())
	err = row.Scan(&userOpts.Id, &userOpts.Email, &teamId, &teamName, &teamFileCountLimit, &teamCurrentFileCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Errorf("HydrateTeam %s: no such user", user.GetId())
		}
		return nil, utilities.WrapBadError(err, fmt.Sprintf("HydrateTeam %s", user.GetId()))
	}
	if teamId.Valid &&
		teamName.Valid &&
		teamCurrentFileCount.Valid &&
		teamFileCountLimit.Valid {
		currentFileCount := int(teamCurrentFileCount.Int64)
		teamOpts = model.TeamOptions{
			Id:               teamId.String,
			Name:             teamName.String,
			CurrentFileCount: &currentFileCount,
			FileCountLimit:   int(teamFileCountLimit.Int64),
		}
	} else {
		id := s.IdGenerator.Generate()
		result, err := tx.Exec(
			`INSERT INTO public."teams" ("id", "name") VALUES ($1, $2)`, id, userOpts.Email,
		)
		if err != nil {
			return nil, utilities.WrapBadError(err, fmt.Sprintf("dbError while inserting team: %s %s", id, userOpts.Email))
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return nil, utilities.WrapBadError(err, fmt.Sprintf("dbError while inserting team and changing db: %s %s", id, userOpts.Email))
		}
		if rowsAffected != 1 {
			return nil, utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when inserting team in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
		}

		result, err = tx.Exec(
			`UPDATE public."users" SET "team_id" = $1 WHERE id = $2`, id, userOpts.Id,
		)
		if err != nil {
			return nil, utilities.WrapBadError(err, fmt.Sprintf("dbError while connecting team to user: %s %s", id, userOpts.Id))
		}
		rowsAffected, err = result.RowsAffected()
		if err != nil {
			return nil, utilities.WrapBadError(err, fmt.Sprintf("dbError while checking affected row after connecting team to user: %s %s", id, userOpts.Id))
		}
		if rowsAffected != 1 {
			return nil, utilities.NewBadError(fmt.Sprintf("Very few or too many rows rows were affected when team was connected to user. This is highly unexpected. rowsAffected: %d", rowsAffected))
		}

		newFileCount := 0
		teamOpts = model.TeamOptions{
			Id:               id,
			Name:             userOpts.Email,
			CurrentFileCount: &newFileCount,
			FileCountLimit:   100,
		}
	}

	team, err := model.NewTeam(teamOpts)
	if err != nil {
		return nil, utilities.WrapBadError(err, fmt.Sprintf("invalid team options error while hydrating user: %s", userOpts.Id))
	}

	err = tx.Commit()
	if err != nil {
		return nil, utilities.WrapBadError(err, "dbError while hydrating user with team tx")
	}
	userOpts.Team = team

	return model.NewUser(userOpts)
}
