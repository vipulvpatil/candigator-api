package storage

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type TeamAccessor interface {
	TeamForUserId(userId string) (*model.Team, error)
	CreateTeamForUserId(userId string) (*model.Team, error)
}

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
	row := tx.QueryRow(`
		SELECT users.id, users.email, teams.id, teams.name
		FROM public."users"
		LEFT JOIN public."teams" ON teams.id = users.team_id
		WHERE users.id = $1
		ORDER BY teams.created_at ASC, teams.id
		FOR UPDATE OF users
		LIMIT 1
	`, user.GetId())
	err = row.Scan(&userOpts.Id, &userOpts.Email, &teamId, &teamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Errorf("HydrateTeam %s: no such user", user.GetId())
		}
		return nil, utilities.WrapBadError(err, fmt.Sprintf("HydrateTeam %s", user.GetId()))
	}
	if teamId.Valid && teamName.Valid {
		teamOpts = model.TeamOptions{
			Id:   teamId.String,
			Name: teamName.String,
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

		teamOpts = model.TeamOptions{
			Id:   id,
			Name: userOpts.Email,
		}
	}

	team, err := model.NewTeam(teamOpts)
	if err != nil {
		return nil, utilities.WrapBadError(err, "invalid team options")
	}

	err = tx.Commit()
	if err != nil {
		return nil, utilities.WrapBadError(err, "dbError while hydrating user with team tx")
	}
	userOpts.Team = team

	return model.NewUser(userOpts)
}
