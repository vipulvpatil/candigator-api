package storage

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

func Test_CreateFileUploadForTeam(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "Team1",
	})
	fileUpload, _ := model.NewFileUpload(model.FileUploadOptions{
		Id:           "fp_id1",
		Name:         "file1.pdf",
		PresignedUrl: "",
		Status:       "WAITING_FOR_FILE",
		Team:         team,
	})
	tests := []struct {
		name  string
		input struct {
			name string
			team *model.Team
		}
		output          *model.FileUpload
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		dbUpdateCheck   func(db *sql.DB) bool
		errorExpected   bool
		errorString     string
	}{
		{
			name: "errors when name is empty",
			input: struct {
				name string
				team *model.Team
			}{},
			output:          nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "cannot create FileUpload with an empty name",
		},
		{
			name: "errors when team is nil",
			input: struct {
				name string
				team *model.Team
			}{
				name: "file1.pdf",
			},
			output:          nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "cannot create FileUpload with a nil team",
		},
		{
			name: "errors when team does not exist in Database",
			input: struct {
				name string
				team *model.Team
			}{
				name: "file1.pdf",
				team: team,
			},
			output:          nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "THIS IS BAD: dbError while inserting file_upload: fp_id1: pq: insert or update on table \"file_uploads\" violates foreign key constraint \"file_uploads_team_id_fkey\"",
		},
		{
			name: "successfully creates a new file upload",
			input: struct {
				name string
				team *model.Team
			}{
				name: "file1.pdf",
				team: team,
			},
			output: fileUpload,
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."teams" (
								"id", "name"
							)
							VALUES (
								'team_id1', 'Team1'
							)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."teams" WHERE id = 'team_id1'`},
				{Query: `DELETE FROM public."file_uploads" WHERE id = 'fp_id1'`},
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var id, name, presignedUrl, status, createdAt string
				row := db.QueryRow(
					`SELECT id, name, presigned_url, status, created_at FROM public."file_uploads" WHERE team_id = 'team_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&id, &name, &presignedUrl, &status, &createdAt)
				assert.NoError(t, err)
				assert.Equal(t, "fp_id1", id)
				assert.Equal(t, "file1.pdf", name)
				assert.Equal(t, "", presignedUrl)
				assert.Equal(t, model.FileUploadStatus("WAITING_FOR_FILE").String(), status)
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
					IdGenerator: &utilities.IdGeneratorMockConstant{Id: "fp_id1"},
				},
			)

			runSqlOnDb(t, s.db, tt.setupSqlStmts)
			defer runSqlOnDb(t, s.db, tt.cleanupSqlStmts)
			fileUpload, err := s.CreateFileUploadForTeam(tt.input.name, tt.input.team)
			assert.Equal(t, tt.output, fileUpload)
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

func Test_UpdateFileUploadWithPresignedUrl(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			id           string
			presignedUrl string
		}
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		dbUpdateCheck   func(db *sql.DB) bool
		errorExpected   bool
		errorString     string
	}{
		{
			name: "errors when id is empty",
			input: struct {
				id           string
				presignedUrl string
			}{},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "id cannot be blank",
		},
		{
			name: "errors when presignedUrl is empty",
			input: struct {
				id           string
				presignedUrl string
			}{
				id: "fp_id1",
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "presignedUrl cannot be blank",
		},
		{
			name: "errors when fileUpload does not exist in database",
			input: struct {
				id           string
				presignedUrl string
			}{
				id:           "fp_id1",
				presignedUrl: "http://presigned_url1",
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "THIS IS BAD: Very few or too many rows were affected when inserting file_upload in db. This is highly unexpected. rowsAffected: 0",
		},
		{
			name: "successfully updates file upload",
			input: struct {
				id           string
				presignedUrl string
			}{
				id:           "fp_id1",
				presignedUrl: "http://presigned_url1",
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
								"id", "name", "presigned_url", "status", "team_id"
							)
							VALUES (
								'fp_id1', 'file1.pdf', '', 'WAITING_FOR_FILE', 'team_id1'
							)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."teams" WHERE id = 'team_id1'`},
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var id, name, presignedUrl, status, createdAt string
				row := db.QueryRow(
					`SELECT id, name, presigned_url, status, created_at FROM public."file_uploads" WHERE team_id = 'team_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&id, &name, &presignedUrl, &status, &createdAt)
				assert.NoError(t, err)
				assert.Equal(t, "fp_id1", id)
				assert.Equal(t, "file1.pdf", name)
				assert.Equal(t, "http://presigned_url1", presignedUrl)
				assert.Equal(t, model.FileUploadStatus("WAITING_FOR_FILE").String(), status)
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
					IdGenerator: &utilities.IdGeneratorMockConstant{Id: "fp_id1"},
				},
			)

			runSqlOnDb(t, s.db, tt.setupSqlStmts)
			defer runSqlOnDb(t, s.db, tt.cleanupSqlStmts)
			err := s.UpdateFileUploadWithPresignedUrl(tt.input.id, tt.input.presignedUrl)
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
