package storage

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

func Test_GetFileUpload(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "Team1",
	})
	fileUpload, _ := model.NewFileUpload(model.FileUploadOptions{
		Id:               "fp_id1",
		Name:             "file1.pdf",
		PresignedUrl:     "https://presigned_url1",
		Status:           "INITIATED",
		ProcessingStatus: "NOT STARTED",
		Team:             team,
	})
	tests := []struct {
		name            string
		input           string
		output          *model.FileUpload
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name:            "errors when id is empty",
			input:           "",
			output:          nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "id cannot be blank",
		},
		{
			name:            "does not error if fileUpload does not exist in Database",
			input:           "fp_id1",
			output:          nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   false,
			errorString:     "",
		},
		{
			name:   "successfully gets file upload",
			input:  "fp_id1",
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
			fileUpload, err := s.GetFileUpload(tt.input)
			assert.Equal(t, tt.output, fileUpload)
			if !tt.errorExpected {
				assert.NoError(t, err)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}

func Test_GetFileUploadsForTeam(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "Team1",
	})
	fileUpload1, _ := model.NewFileUpload(model.FileUploadOptions{
		Id:               "fp_id1",
		Name:             "file1.pdf",
		PresignedUrl:     "https://presigned_url1",
		Status:           "INITIATED",
		ProcessingStatus: "NOT STARTED",
		Team:             team,
	})
	fileUpload2, _ := model.NewFileUpload(model.FileUploadOptions{
		Id:               "fp_id2",
		Name:             "file2.pdf",
		PresignedUrl:     "https://presigned_url2",
		Status:           "INITIATED",
		ProcessingStatus: "NOT STARTED",
		Team:             team,
	})
	tests := []struct {
		name            string
		input           *model.Team
		output          []*model.FileUpload
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
			name:   "successfully gets file uploads",
			input:  team,
			output: []*model.FileUpload{fileUpload1, fileUpload2},
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
			fileUpload, err := s.GetFileUploadsForTeam(tt.input)
			assert.Equal(t, tt.output, fileUpload)
			if !tt.errorExpected {
				assert.NoError(t, err)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}

func Test_GetUnprocessedFileUploadsCountForTeam(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "Team1",
	})
	tests := []struct {
		name            string
		input           *model.Team
		output          int
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name:            "errors when team is empty",
			input:           nil,
			output:          0,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "team cannot be blank",
		},
		{
			name:            "errors when team id is empty",
			input:           &model.Team{},
			output:          0,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "team cannot be blank",
		},
		{
			name:   "successfully gets correct unprocessed file upload count",
			input:  team,
			output: 2,
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
								'fp_id1', 'file1.pdf', 'https://presigned_url1', 'INITIATED', 'NOT STARTED', 'team_id1'
							)`,
				},
				{
					Query: `INSERT INTO public."file_uploads" (
								"id", "name", "presigned_url", "status", "processing_status", "team_id"
							)
							VALUES (
								'fp_id2', 'file2.pdf', 'https://presigned_url2', 'INITIATED', 'ONGOING', 'team_id1'
							)`,
				},
				{
					Query: `INSERT INTO public."file_uploads" (
								"id", "name", "presigned_url", "status", "processing_status", "team_id"
							)
							VALUES (
								'fp_id3', 'file3.pdf', 'https://presigned_url3', 'INITIATED', 'ONGOING', 'team_id2'
							)`,
				},
				{
					Query: `INSERT INTO public."file_uploads" (
								"id", "name", "presigned_url", "status", "processing_status", "team_id"
							)
							VALUES (
								'fp_id4', 'file4.pdf', 'https://presigned_url4', 'SUCCESS', 'COMPLETED', 'team_id1'
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
			fileUpload, err := s.GetUnprocessedFileUploadsCountForTeam(tt.input)
			assert.Equal(t, tt.output, fileUpload)
			if !tt.errorExpected {
				assert.NoError(t, err)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}

func Test_CreateFileUploadForTeam(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "Team1",
	})
	fileUpload, _ := model.NewFileUpload(model.FileUploadOptions{
		Id:               "fp_id1",
		Name:             "file1.pdf",
		PresignedUrl:     "",
		Status:           "INITIATED",
		ProcessingStatus: "NOT STARTED",
		Team:             team,
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
				assert.Equal(t, model.FileUploadStatus("INITIATED").String(), status)
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
								"id", "name", "presigned_url", "status", "processing_status", "team_id"
							)
							VALUES (
								'fp_id1', 'file1.pdf', '', 'INITIATED', 'NOT STARTED', 'team_id1'
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
				assert.Equal(t, model.FileUploadStatus("INITIATED").String(), status)
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

func Test_UpdateFileUploadWithStatus(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			id     string
			status string
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
				id     string
				status string
			}{},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "id cannot be blank",
		},
		{
			name: "errors when status is not valid",
			input: struct {
				id     string
				status string
			}{
				id: "fp_id1",
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "status should be valid",
		},
		{
			name: "errors when fileUpload does not exist in database",
			input: struct {
				id     string
				status string
			}{
				id:     "fp_id1",
				status: "SUCCESS",
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
				id     string
				status string
			}{
				id:     "fp_id1",
				status: "FAILURE",
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
								'fp_id1', 'file1.pdf', 'http://presigned_url1', 'INITIATED', 'NOT STARTED', 'team_id1'
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
				assert.Equal(t, model.FileUploadStatus("FAILURE").String(), status)
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
			err := s.UpdateFileUploadWithStatus(tt.input.id, tt.input.status)
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
