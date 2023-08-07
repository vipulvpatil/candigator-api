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

func Test_GetCandidatesForTeam(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "Team1",
	})
	persona1 := model.Persona{
		Name:       "ai persona 1",
		Email:      "email_1",
		Phone:      "phone_1",
		City:       "city_1",
		State:      "state_1",
		Country:    "country_1",
		YoE:        5,
		TechSkills: []string{"tech skill 1", "tech skill 2", "tech skill 3"},
	}
	persona2 := model.Persona{
		Name:       "ai persona 2",
		Email:      "email_2",
		Phone:      "phone_2",
		City:       "city_2",
		State:      "state_2",
		Country:    "country_2",
		YoE:        2,
		TechSkills: []string{"tech skill 10", "tech skill 20", "tech skill 30"},
	}
	persona3 := model.Persona{
		Name:       "manual persona 1",
		Email:      "email_3",
		Phone:      "phone_3",
		City:       "city_3",
		State:      "state_3",
		Country:    "country_3",
		YoE:        7,
		TechSkills: []string{"tech skill 11", "tech skill 21", "tech skill 31"},
	}
	persona4 := model.Persona{
		Name:       "manual persona 2",
		Email:      "email_4",
		Phone:      "phone_4",
		City:       "city_4",
		State:      "state_4",
		Country:    "country_4",
		YoE:        51,
		TechSkills: []string{"tech skill 13", "tech skill 23", "tech skill 33"},
	}
	candidate1, _ := model.NewCandidate(model.CandidateOptions{
		Id:                 "c_id1",
		AiGeneratedPersona: &persona1,
		Team:               team,
		FileUploadId:       "fp_id1",
	})
	candidate2, _ := model.NewCandidate(model.CandidateOptions{
		Id:                     "c_id2",
		AiGeneratedPersona:     &persona2,
		ManuallyCreatedPersona: &persona3,
		Team:                   team,
		FileUploadId:           "fp_id2",
	})
	candidate3, _ := model.NewCandidate(model.CandidateOptions{
		Id:                     "c_id3",
		ManuallyCreatedPersona: &persona4,
		Team:                   team,
	})
	tests := []struct {
		name            string
		input           *model.Team
		output          []*model.Candidate
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name:            "errors when team is empty",
			input:           nil,
			output:          nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "team cannot be blank",
		},
		{
			name:            "errors when team id is empty",
			input:           &model.Team{},
			output:          nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "team cannot be blank",
		},
		{
			name:   "successfully gets candidates",
			input:  team,
			output: []*model.Candidate{candidate1, candidate2, candidate3},
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
				{
					Query: `INSERT INTO public."file_uploads" (
								"id", "name", "presigned_url", "status", "processing_status", "team_id"
							)
							VALUES (
								'fp_id2', 'file2.pdf', 'https://presigned_url2', 'INITIATED', 'NOT STARTED', 'team_id1'
							)`,
				},
				{
					Query: `INSERT INTO public."candidates" (
								"id", "ai_generated_persona", "manually_created_persona","team_id", "file_upload_id"
							)
							VALUES (
								'c_id1', $1, $2, 'team_id1', 'fp_id1'
							),(
								'c_id2', $3, $4, 'team_id1', 'fp_id2'
							),(
								'c_id3', $5, $6, 'team_id1', NULL
							)`,
					Args: []any{
						&persona1, nil,
						&persona2, &persona3,
						nil, &persona4,
					},
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
			candidates, err := s.GetCandidatesForTeam(tt.input)
			assert.Equal(t, tt.output, candidates)
			if !tt.errorExpected {
				assert.NoError(t, err)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}
