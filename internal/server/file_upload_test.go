package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/services/filestorage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
	pb "github.com/vipulvpatil/candidate-tracker-go/protos"
	"google.golang.org/grpc/metadata"
)

func Test_UploadFiles(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "test@example.com",
	})
	userWithTeam, _ := model.NewUser(model.UserOptions{
		Id:    "user_id1",
		Email: "test@example.com",
		Team:  team,
	})
	tests := []struct {
		name                   string
		ctx                    context.Context
		input                  *pb.UploadFilesRequest
		output                 *pb.UploadFilesResponse
		teamHydratorMock       storage.TeamHydrator
		fileUploadAccessorMock storage.FileUploadAccessor
		fileStorerMock         filestorage.FileStorer
		errorExpected          bool
		errorString            string
	}{
		{
			name:                   "errors if no user in context",
			ctx:                    context.Background(),
			input:                  &pb.UploadFilesRequest{},
			output:                 &pb.UploadFilesResponse{},
			teamHydratorMock:       nil,
			fileUploadAccessorMock: nil,
			fileStorerMock:         nil,
			errorExpected:          true,
			errorString:            "rpc error: code = Unauthenticated desc = retrieving user data failed",
		},
		{
			name: "errors if unable to hydrate team",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input:                  &pb.UploadFilesRequest{},
			output:                 &pb.UploadFilesResponse{},
			teamHydratorMock:       &storage.TeamHydratorMockFailure{},
			fileUploadAccessorMock: nil,
			fileStorerMock:         nil,
			errorExpected:          true,
			errorString:            "unable to hydrate team",
		},
		{
			name: "returns response with error if database errors on creation",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.UploadFilesRequest{
				Files: []*pb.UploadFile{
					{
						Name: "file1.pdf",
					},
				},
			},
			output: &pb.UploadFilesResponse{
				FileUploads: []*pb.FileUpload{
					{
						Id:           "",
						Name:         "file1.pdf",
						PresignedUrl: "",
						Status:       "",
						Error:        "unable to create fileUpload",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				CreateFileUploadForTeamInteral: func(name string, team *model.Team) (*model.FileUpload, error) {
					return nil, errors.New("unable to create fileUpload")
				},
			},
			fileStorerMock: nil,
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "returns response with error if unable to get presigedUrl",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.UploadFilesRequest{
				Files: []*pb.UploadFile{
					{
						Name: "file1.pdf",
					},
				},
			},
			output: &pb.UploadFilesResponse{
				FileUploads: []*pb.FileUpload{
					{
						Id:               "fp_id1",
						Name:             "file1.pdf",
						PresignedUrl:     "",
						Status:           "INITIATED",
						ProcessingStatus: "NOT STARTED",
						Error:            "unable to get presignedUrl",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				CreateFileUploadForTeamInteral: func(name string, team *model.Team) (*model.FileUpload, error) {
					return model.NewFileUpload(model.FileUploadOptions{
						Id:               "fp_id1",
						Name:             name,
						Status:           "INITIATED",
						ProcessingStatus: "NOT STARTED",
						Team:             team,
					})
				},
			},
			fileStorerMock: &filestorage.FileStorerMockFailure{},
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "returns response with error if database cannot update FileUpload",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.UploadFilesRequest{
				Files: []*pb.UploadFile{
					{
						Name: "file1.pdf",
					},
				},
			},
			output: &pb.UploadFilesResponse{
				FileUploads: []*pb.FileUpload{
					{
						Id:               "fp_id1",
						Name:             "file1.pdf",
						PresignedUrl:     "http://presigned_url1",
						Status:           "INITIATED",
						ProcessingStatus: "NOT STARTED",
						Error:            "unable to update fileUpload",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				CreateFileUploadForTeamInteral: func(name string, team *model.Team) (*model.FileUpload, error) {
					return model.NewFileUpload(model.FileUploadOptions{
						Id:               "fp_id1",
						Name:             name,
						Status:           "INITIATED",
						ProcessingStatus: "NOT STARTED",
						Team:             team,
					})
				},
				UpdateFileUploadWithPresignedUrlInternal: func(id, presignedUrl string) error {
					return errors.New("unable to update fileUpload")
				},
			},
			fileStorerMock: &filestorage.FileStorerMockSuccess{PresignedUrl: "http://presigned_url1"},
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "returns response without error if nothing fails",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.UploadFilesRequest{
				Files: []*pb.UploadFile{
					{
						Name: "file1.pdf",
					},
				},
			},
			output: &pb.UploadFilesResponse{
				FileUploads: []*pb.FileUpload{
					{
						Id:               "fp_id1",
						Name:             "file1.pdf",
						PresignedUrl:     "http://presigned_url1",
						Status:           "INITIATED",
						ProcessingStatus: "NOT STARTED",
						Error:            "",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				CreateFileUploadForTeamInteral: func(name string, team *model.Team) (*model.FileUpload, error) {
					return model.NewFileUpload(model.FileUploadOptions{
						Id:               "fp_id1",
						Name:             name,
						Status:           "INITIATED",
						ProcessingStatus: "NOT STARTED",
						Team:             team,
					})
				},
				UpdateFileUploadWithPresignedUrlInternal: func(id, presignedUrl string) error {
					return nil
				},
			},
			fileStorerMock: &filestorage.FileStorerMockSuccess{PresignedUrl: "http://presigned_url1"},
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "returns response with some errors if some files fail",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.UploadFilesRequest{
				Files: []*pb.UploadFile{
					{
						Name: "file1.pdf",
					},
					{
						Name: "file2.pdf",
					},
					{
						Name: "file3.pdf",
					},
				},
			},
			output: &pb.UploadFilesResponse{
				FileUploads: []*pb.FileUpload{
					{
						Id:               "fp_id1",
						Name:             "file1.pdf",
						PresignedUrl:     "http://presigned_url1",
						Status:           "INITIATED",
						ProcessingStatus: "NOT STARTED",
						Error:            "",
					},
					{
						Id:               "fp_id2",
						Name:             "file2.pdf",
						PresignedUrl:     "http://presigned_url1",
						Status:           "INITIATED",
						ProcessingStatus: "NOT STARTED",
						Error:            "unable to upload fileUpload",
					},
					{
						Id:               "",
						Name:             "file3.pdf",
						PresignedUrl:     "",
						Status:           "",
						ProcessingStatus: "",
						Error:            "unable to create fileUpload",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				CreateFileUploadForTeamInteral: func(name string, team *model.Team) (*model.FileUpload, error) {
					if name == "file1.pdf" {
						return model.NewFileUpload(model.FileUploadOptions{
							Id:               "fp_id1",
							Name:             name,
							Status:           "INITIATED",
							ProcessingStatus: "NOT STARTED",
							Team:             team,
						})
					} else if name == "file2.pdf" {
						return model.NewFileUpload(model.FileUploadOptions{
							Id:               "fp_id2",
							Name:             name,
							Status:           "INITIATED",
							ProcessingStatus: "NOT STARTED",
							Team:             team,
						})
					} else {
						return nil, errors.New("unable to create fileUpload")
					}
				},
				UpdateFileUploadWithPresignedUrlInternal: func(id, presignedUrl string) error {
					if id == "fp_id2" {
						return errors.New("unable to upload fileUpload")
					}
					return nil
				},
			},
			fileStorerMock: &filestorage.FileStorerMockSuccess{PresignedUrl: "http://presigned_url1"},
			errorExpected:  false,
			errorString:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithTeamHydratorMock(tt.teamHydratorMock),
					storage.WithFileUploadAccessorMock(tt.fileUploadAccessorMock),
				),
				Logger:     &utilities.NullLogger{},
				FileStorer: tt.fileStorerMock,
			})

			response, err := server.UploadFiles(
				tt.ctx,
				tt.input,
			)
			if !tt.errorExpected {
				assert.Empty(t, tt.errorString)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.output, response)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}

func Test_CompleteFileUploads(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "test@example.com",
	})
	otherTeam, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id2",
		Name: "test2@example.com",
	})
	userWithTeam, _ := model.NewUser(model.UserOptions{
		Id:    "user_id1",
		Email: "test@example.com",
		Team:  team,
	})
	initiatedFileUpload, _ := model.NewFileUpload(model.FileUploadOptions{
		Id:               "fp_id1",
		Name:             "file1.pdf",
		PresignedUrl:     "https://presigned_url1",
		Status:           "INITIATED",
		ProcessingStatus: "NOT STARTED",
		Team:             team,
	})
	initiatedFileUpload2, _ := model.NewFileUpload(model.FileUploadOptions{
		Id:               "fp_id2",
		Name:             "file2.pdf",
		PresignedUrl:     "https://presigned_url2",
		Status:           "INITIATED",
		ProcessingStatus: "NOT STARTED",
		Team:             team,
	})
	completedFileUpload, _ := model.NewFileUpload(model.FileUploadOptions{
		Id:               "fp_id1",
		Name:             "file1.pdf",
		PresignedUrl:     "https://presigned_url1",
		Status:           "SUCCESS",
		ProcessingStatus: "NOT STARTED",
		Team:             team,
	})
	anotherFileUpload, _ := model.NewFileUpload(model.FileUploadOptions{
		Id:               "fp_id1",
		Name:             "file1.pdf",
		PresignedUrl:     "https://presigned_url1",
		Status:           "INITIATED",
		ProcessingStatus: "NOT STARTED",
		Team:             otherTeam,
	})

	tests := []struct {
		name                   string
		ctx                    context.Context
		input                  *pb.CompleteFileUploadsRequest
		output                 *pb.CompleteFileUploadsResponse
		teamHydratorMock       storage.TeamHydrator
		fileUploadAccessorMock storage.FileUploadAccessor
		fileStorerMock         filestorage.FileStorer
		errorExpected          bool
		errorString            string
	}{
		{
			name:                   "errors if no user in context",
			ctx:                    context.Background(),
			input:                  &pb.CompleteFileUploadsRequest{},
			output:                 &pb.CompleteFileUploadsResponse{},
			teamHydratorMock:       nil,
			fileUploadAccessorMock: nil,
			fileStorerMock:         nil,
			errorExpected:          true,
			errorString:            "rpc error: code = Unauthenticated desc = retrieving user data failed",
		},
		{
			name: "errors if unable to hydrate team",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input:                  &pb.CompleteFileUploadsRequest{},
			output:                 &pb.CompleteFileUploadsResponse{},
			teamHydratorMock:       &storage.TeamHydratorMockFailure{},
			fileUploadAccessorMock: nil,
			fileStorerMock:         nil,
			errorExpected:          true,
			errorString:            "unable to hydrate team",
		},
		{
			name: "returns response with error if database errors when getting fileUpload",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.CompleteFileUploadsRequest{
				FileUploadUpdates: []*pb.FileUploadUpdate{
					{
						Id:     "fp_id1",
						Status: "SUCCESS",
					},
				},
			},
			output: &pb.CompleteFileUploadsResponse{
				FileUploads: []*pb.FileUpload{
					{
						Id:               "fp_id1",
						Name:             "",
						PresignedUrl:     "",
						Status:           "",
						ProcessingStatus: "",
						Error:            "unable to get fileUpload: dbError when querying",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetFileUploadInternal: func(id string) (*model.FileUpload, error) {
					return nil, errors.New("dbError when querying")
				},
			},
			fileStorerMock: nil,
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "returns response with error if fileUpload not in database",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.CompleteFileUploadsRequest{
				FileUploadUpdates: []*pb.FileUploadUpdate{
					{
						Id:     "fp_id1",
						Status: "SUCCESS",
					},
				},
			},
			output: &pb.CompleteFileUploadsResponse{
				FileUploads: []*pb.FileUpload{
					{
						Id:               "fp_id1",
						Name:             "",
						PresignedUrl:     "",
						Status:           "",
						ProcessingStatus: "",
						Error:            "unable to get fileUpload",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetFileUploadInternal: func(id string) (*model.FileUpload, error) {
					return nil, nil
				},
			},
			fileStorerMock: nil,
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "returns response with error if fileUpload in database is already completed",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.CompleteFileUploadsRequest{
				FileUploadUpdates: []*pb.FileUploadUpdate{
					{
						Id:     "fp_id1",
						Status: "SUCCESS",
					},
				},
			},
			output: &pb.CompleteFileUploadsResponse{
				FileUploads: []*pb.FileUpload{
					{
						Id:               "fp_id1",
						Name:             "file1.pdf",
						PresignedUrl:     "https://presigned_url1",
						Status:           "",
						ProcessingStatus: "NOT STARTED",
						Error:            "unable to update fileUpload",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetFileUploadInternal: func(id string) (*model.FileUpload, error) {
					return completedFileUpload, nil
				},
			},
			fileStorerMock: nil,
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "returns response with error if fileUpload belongs to different team",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.CompleteFileUploadsRequest{
				FileUploadUpdates: []*pb.FileUploadUpdate{
					{
						Id:     "fp_id1",
						Status: "SUCCESS",
					},
				},
			},
			output: &pb.CompleteFileUploadsResponse{
				FileUploads: []*pb.FileUpload{
					{
						Id:               "fp_id1",
						Name:             "file1.pdf",
						PresignedUrl:     "https://presigned_url1",
						Status:           "",
						ProcessingStatus: "NOT STARTED",
						Error:            "THIS IS BAD: unauthorized update of fileUpload attempted",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetFileUploadInternal: func(id string) (*model.FileUpload, error) {
					return anotherFileUpload, nil
				},
			},
			fileStorerMock: nil,
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "returns response with error if new status is invalid",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.CompleteFileUploadsRequest{
				FileUploadUpdates: []*pb.FileUploadUpdate{
					{
						Id:     "fp_id1",
						Status: "DONE",
					},
				},
			},
			output: &pb.CompleteFileUploadsResponse{
				FileUploads: []*pb.FileUpload{
					{
						Id:               "fp_id1",
						Name:             "file1.pdf",
						PresignedUrl:     "https://presigned_url1",
						Status:           "",
						ProcessingStatus: "NOT STARTED",
						Error:            "invalid update status",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetFileUploadInternal: func(id string) (*model.FileUpload, error) {
					return initiatedFileUpload, nil
				},
			},
			fileStorerMock: nil,
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "returns response with error if updating fails in database",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.CompleteFileUploadsRequest{
				FileUploadUpdates: []*pb.FileUploadUpdate{
					{
						Id:     "fp_id1",
						Status: "SUCCESS",
					},
				},
			},
			output: &pb.CompleteFileUploadsResponse{
				FileUploads: []*pb.FileUpload{
					{
						Id:               "fp_id1",
						Name:             "file1.pdf",
						PresignedUrl:     "https://presigned_url1",
						Status:           "",
						ProcessingStatus: "NOT STARTED",
						Error:            "THIS IS BAD: unable to update fileUpload: dbError while updating",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetFileUploadInternal: func(id string) (*model.FileUpload, error) {
					return initiatedFileUpload, nil
				},
				UpdateFileUploadWithStatusInternal: func(id, status string) error {
					return errors.New("dbError while updating")
				},
			},
			fileStorerMock: nil,
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "runs successfully with  partial errors",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.CompleteFileUploadsRequest{
				FileUploadUpdates: []*pb.FileUploadUpdate{
					{
						Id:     "fp_id1",
						Status: "SUCCESS",
					},
					{
						Id:     "fp_id2",
						Status: "SUCCESS",
					},
				},
			},
			output: &pb.CompleteFileUploadsResponse{
				FileUploads: []*pb.FileUpload{
					{
						Id:               "fp_id1",
						Name:             "file1.pdf",
						PresignedUrl:     "https://presigned_url1",
						Status:           "",
						ProcessingStatus: "NOT STARTED",
						Error:            "THIS IS BAD: unable to update fileUpload: dbError while updating",
					},
					{
						Id:               "fp_id2",
						Name:             "file2.pdf",
						PresignedUrl:     "https://presigned_url2",
						Status:           "SUCCESS",
						ProcessingStatus: "NOT STARTED",
						Error:            "",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetFileUploadInternal: func(id string) (*model.FileUpload, error) {
					if id == "fp_id1" {
						return initiatedFileUpload, nil
					} else {
						return initiatedFileUpload2, nil
					}
				},
				UpdateFileUploadWithStatusInternal: func(id, status string) error {
					if id == "fp_id1" {
						return errors.New("dbError while updating")
					} else {
						return nil
					}
				},
			},
			fileStorerMock: nil,
			errorExpected:  false,
			errorString:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithTeamHydratorMock(tt.teamHydratorMock),
					storage.WithFileUploadAccessorMock(tt.fileUploadAccessorMock),
				),
				Logger:     &utilities.NullLogger{},
				FileStorer: tt.fileStorerMock,
			})

			response, err := server.CompleteFileUploads(
				tt.ctx,
				tt.input,
			)
			if !tt.errorExpected {
				assert.Empty(t, tt.errorString)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.output, response)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}
