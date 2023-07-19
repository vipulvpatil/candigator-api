package storage

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

func Test_CreateCandidateWithAiGeneratedPersonaForTeamUsingTx(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "Team1",
	})
	tests := []struct {
		name  string
		input struct {
			persona *model.Persona
			team    *model.Team
		}
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		dbUpdateCheck   func(db *sql.DB) bool
		errorExpected   bool
		errorString     string
	}{
		{
			name: "errors when team is nil",
			input: struct {
				persona *model.Persona
				team    *model.Team
			}{
				persona: &model.Persona{Name: "user_1"},
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "cannot create Candidate with a nil Team",
		},
		{
			name: "errors when persona is invalid",
			input: struct {
				persona *model.Persona
				team    *model.Team
			}{
				team: team,
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "cannot create Candidate without a valid persona",
		},
		{
			name: "errors when fileUpload does not exist in Database",
			input: struct {
				persona *model.Persona
				team    *model.Team
			}{
				persona: &model.Persona{Name: "user_1", BuiltBy: "AI", FileUploadId: "fp_id1"},
				team:    team,
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "THIS IS BAD: dbError while inserting Candidate: can_id1: pq: insert or update on table \"candidates\" violates foreign key constraint \"candidates_file_upload_id_fkey\"",
		},
		{
			name: "errors when team does not exist in Database",
			input: struct {
				persona *model.Persona
				team    *model.Team
			}{
				persona: &model.Persona{Name: "user_1", BuiltBy: "AI", FileUploadId: "fp_id1"},
				team:    team,
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."teams" (
						"id", "name"
					)
					VALUES (
						'team_id2', 'Team2'
					)`,
				},
				{
					Query: `INSERT INTO public."file_uploads" (
						"id", "name", "presigned_url", "status", "processing_status", "team_id"
					)
					VALUES (
						'fp_id1', 'file1.pdf', 'https://presigned_url1', 'INITIATED', 'NOT STARTED', 'team_id2'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."teams" WHERE id = 'team_id2'`},
			},
			dbUpdateCheck: nil,
			errorExpected: true,
			errorString:   "THIS IS BAD: dbError while inserting Candidate: can_id1: pq: insert or update on table \"candidates\" violates foreign key constraint \"candidates_team_id_fkey\"",
		},
		{
			name: "successfully creates a new file upload",
			input: struct {
				persona *model.Persona
				team    *model.Team
			}{
				persona: &model.Persona{Name: "user_1", BuiltBy: "AI", FileUploadId: "fp_id1"},
				team:    team,
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."teams" (
								"id", "name"
							)
							VALUES (
								'team_id1', 'Team1'
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
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."teams" WHERE id = 'team_id1'`},
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var id string
				var persona model.Persona
				var createdAt time.Time
				row := db.QueryRow(
					`SELECT id, ai_generated_persona, created_at FROM public."candidates" WHERE team_id = 'team_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&id, &persona, &createdAt)
				assert.NoError(t, err)
				assert.Equal(t, "can_id1", id)
				assert.Equal(t, "user_1", persona.Name)
				return true
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewDbStorage(
				StorageOptions{
					Db:          testDb,
					IdGenerator: &utilities.IdGeneratorMockConstant{Id: "can_id1"},
				},
			)

			runSqlOnDb(t, s.db, tt.setupSqlStmts)
			defer runSqlOnDb(t, s.db, tt.cleanupSqlStmts)
			tx, err := s.BeginTransaction()
			assert.NoError(t, err)
			err = s.CreateCandidateWithAiGeneratedPersonaForTeamUsingTx(
				tt.input.persona, tt.input.team, tx,
			)
			tx.Commit()
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
