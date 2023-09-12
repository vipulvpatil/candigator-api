package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
	pb "github.com/vipulvpatil/candidate-tracker-go/protos"
	"google.golang.org/grpc/metadata"
)

func Test_GetUserData(t *testing.T) {
	currentFileCount := 1
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
		name                   string
		ctx                    context.Context
		input                  *pb.GetUserDataRequest
		output                 *pb.GetUserDataResponse
		teamHydratorMock       storage.TeamHydrator
		fileUploadAccessorMock storage.FileUploadAccessor
		errorExpected          bool
		errorString            string
	}{
		{
			name:                   "errors if no user in context",
			ctx:                    context.Background(),
			input:                  &pb.GetUserDataRequest{},
			output:                 &pb.GetUserDataResponse{},
			teamHydratorMock:       nil,
			fileUploadAccessorMock: nil,
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
			input:                  &pb.GetUserDataRequest{},
			output:                 &pb.GetUserDataResponse{},
			teamHydratorMock:       &storage.TeamHydratorMockFailure{},
			fileUploadAccessorMock: nil,
			errorExpected:          true,
			errorString:            "unable to hydrate team",
		},
		{
			name: "returns error if database errors when getting fileUpload count",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input:            &pb.GetUserDataRequest{},
			output:           nil,
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetUnprocessedFileUploadsCountForTeamInternal: func(id *model.Team) (int, error) {
					return 0, errors.New("dbError when querying")
				},
			},
			errorExpected: true,
			errorString:   "dbError when querying",
		},
		{
			name: "returns response with 0 if fileUpload count not in database",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.GetUserDataRequest{},
			output: &pb.GetUserDataResponse{
				FileCountLimit:       100,
				CurrentFileCount:     1,
				UnprocessedFileCount: 0,
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetUnprocessedFileUploadsCountForTeamInternal: func(id *model.Team) (int, error) {
					return 0, nil
				},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "runs successfully",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.GetUserDataRequest{},
			output: &pb.GetUserDataResponse{
				FileCountLimit:       100,
				CurrentFileCount:     1,
				UnprocessedFileCount: 3,
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			fileUploadAccessorMock: &storage.FileUploadAccessorConfigurableMock{
				GetUnprocessedFileUploadsCountForTeamInternal: func(id *model.Team) (int, error) {
					return 3, nil
				},
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithTeamHydratorMock(tt.teamHydratorMock),
					storage.WithFileUploadAccessorMock(tt.fileUploadAccessorMock),
				),
				Logger: &utilities.NullLogger{},
			})

			response, err := server.GetUserData(
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
