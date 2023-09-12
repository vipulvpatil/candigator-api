package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
)

func Test_UserByEmail(t *testing.T) {
	currentFileCount := 1
	user, _ := model.NewUser(model.UserOptions{
		Id:    "user_id1",
		Email: "test@example.com",
	})
	team, _ := model.NewTeam(model.TeamOptions{
		Id:               "team_id1",
		Name:             "test@example.com",
		CurrentFileCount: &currentFileCount,
		FileCountLimit:   100,
	})
	userWithTeam, _ := model.NewUser(model.UserOptions{
		Id:    "user_id1",
		Email: "test@example.com",
		Team:  team,
	})
	tests := []struct {
		name            string
		input           string
		output          *model.User
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name:            "errors when email is blank",
			input:           "",
			output:          nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "cannot search by blank email",
		},
		{
			name:            "errors if no such user",
			input:           "absent@user.com",
			output:          nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "UserByEmail absent@user.com: no such user",
		},
		{
			name:   "returns user without team if not associated",
			input:  "test@example.com",
			output: user,
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."users" (
						"id", "email"
					)
					VALUES (
						'user_id1', 'test@example.com'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."users" WHERE id = 'user_id1'`},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name:   "returns user with team if associated",
			input:  "test@example.com",
			output: userWithTeam,
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."teams" (
						"id", "name"
					)
					VALUES (
						'team_id1', 'test@example.com'
					)`,
				},
				{
					Query: `INSERT INTO public."users" (
						"id", "email", "team_id"
					)
					VALUES (
						'user_id1', 'test@example.com', 'team_id1'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."teams" WHERE id = 'team_id1'`},
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewDbStorage(
				StorageOptions{
					Db: testDb,
				},
			)

			runSqlOnDb(t, s.db, tt.setupSqlStmts)
			defer runSqlOnDb(t, s.db, tt.cleanupSqlStmts)
			user, err := s.UserByEmail(tt.input)
			assert.Equal(t, tt.output, user)
			if !tt.errorExpected {
				assert.NoError(t, err)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}
