package storage

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

func Test_HydrateTeam(t *testing.T) {
	zeroFileCount := 0
	currentFileCount := 1
	inputUserWithoutTeam, _ := model.NewUser(model.UserOptions{
		Id:    "user_id1",
		Email: "test@example.com",
	})
	team, _ := model.NewTeam(model.TeamOptions{
		Id:               "team_id1",
		Name:             "test@example.com",
		CurrentFileCount: &currentFileCount,
		FileCountLimit:   100,
	})
	inputUserWithTeam, _ := model.NewUser(model.UserOptions{
		Id:    "user_id1",
		Email: "test@example.com",
		Team:  team,
	})
	teamWithZeroFiles, _ := model.NewTeam(model.TeamOptions{
		Id:               "team_id1",
		Name:             "test@example.com",
		CurrentFileCount: &zeroFileCount,
		FileCountLimit:   100,
	})
	inputUserWithZeroFilesTeam, _ := model.NewUser(model.UserOptions{
		Id:    "user_id1",
		Email: "test@example.com",
		Team:  teamWithZeroFiles,
	})
	tests := []struct {
		name            string
		input           *model.User
		output          *model.User
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		dbUpdateCheck   func(*sql.DB) bool
		errorExpected   bool
		errorString     string
	}{
		{
			name:            "errors when user is nil",
			input:           nil,
			output:          nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "cannot hydrate a nil user",
		},
		{
			name:            "returns user if it has team without checking database",
			input:           inputUserWithTeam,
			output:          inputUserWithTeam,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   false,
			errorString:     "",
		},
		{
			name:   "hydrates and returns user if there is an associated team in database.",
			input:  inputUserWithoutTeam,
			output: inputUserWithTeam,
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
					Query: `INSERT INTO public."file_uploads" (
						"id", "name", "presigned_url", "status", "processing_status", "team_id"
					)
					VALUES (
						'fp_id1', 'file1.pdf', 'https://presigned_url1', 'INITIATED', 'NOT STARTED', 'team_id1'
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
			dbUpdateCheck: nil,
			errorExpected: false,
			errorString:   "",
		},
		{
			name:   "errors if user not in database.",
			input:  inputUserWithoutTeam,
			output: nil,
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."teams" (
						"id", "name"
					)
					VALUES (
						'team_id1', 'test@example.com'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."teams" WHERE id = 'team_id1'`},
			},
			dbUpdateCheck: nil,
			errorExpected: true,
			errorString:   "HydrateTeam user_id1: no such user",
		},
		{
			name:   "hydrates and returns user by creating a new team in database.",
			input:  inputUserWithoutTeam,
			output: inputUserWithZeroFilesTeam,
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
				{Query: `DELETE FROM public."teams" WHERE id = 'team_id1'`},
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var id string
				row := db.QueryRow(
					`SELECT id FROM public."teams" WHERE id = 'team_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&id)
				assert.NoError(t, err)
				assert.Equal(t, "team_id1", id)

				var teamId sql.NullString
				row = db.QueryRow(
					`SELECT team_id FROM public."users" WHERE id = 'user_id1'`,
				)
				assert.NoError(t, row.Err())
				err = row.Scan(&teamId)
				assert.NoError(t, err)
				assert.True(t, teamId.Valid)
				assert.Equal(t, "team_id1", teamId.String)

				return true
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name:   "errors if creating a new team in database errors.",
			input:  inputUserWithoutTeam,
			output: nil,
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."users" (
						"id", "email"
					)
					VALUES (
						'user_id1', 'test@example.com'
					)`,
				},
				{
					Query: `INSERT INTO public."teams" (
						"id", "name"
					)
					VALUES (
						'team_id1', 'test@example.com'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."users" WHERE id = 'user_id1'`},
				{Query: `DELETE FROM public."teams" WHERE id = 'team_id1'`},
			},
			dbUpdateCheck: nil,
			errorExpected: true,
			errorString:   "THIS IS BAD: dbError while inserting team: team_id1 test@example.com: pq: duplicate key value violates unique constraint \"teams_pkey\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewDbStorage(
				StorageOptions{
					Db:          testDb,
					IdGenerator: &utilities.IdGeneratorMockConstant{Id: "team_id1"},
				},
			)

			runSqlOnDb(t, s.db, tt.setupSqlStmts)
			defer runSqlOnDb(t, s.db, tt.cleanupSqlStmts)
			user, err := s.HydrateTeam(tt.input)
			assert.Equal(t, tt.output, user)
			if tt.output != nil {
				assert.Equal(t, tt.output.Team(), user.Team())
			}
			if !tt.errorExpected {
				assert.NoError(t, err)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
			if tt.dbUpdateCheck != nil {
				assert.True(t, tt.dbUpdateCheck(s.db))
			}
		})
	}
}
