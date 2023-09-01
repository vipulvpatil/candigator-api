package server

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
	pb "github.com/vipulvpatil/candidate-tracker-go/protos"
	"google.golang.org/grpc/metadata"
)

func Test_GetCandidates(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "test@example.com",
	})
	userWithTeam, _ := model.NewUser(model.UserOptions{
		Id:    "user_id1",
		Email: "test@example.com",
		Team:  team,
	})

	candidate1, _ := model.NewCandidate(model.CandidateOptions{
		Id: "c_id1",
		AiGeneratedPersona: &model.Persona{
			Name:       "ai persona 1",
			Email:      "email_1",
			Phone:      "phone_1",
			City:       "city_1",
			State:      "state_1",
			Country:    "country_1",
			YoE:        5,
			TechSkills: []string{"tech skill 1", "tech skill 2", "tech skill 3"},
		},
		Team:         team,
		FileUploadId: "fp_id1",
	})
	candidate2, _ := model.NewCandidate(model.CandidateOptions{
		Id: "c_id2",
		AiGeneratedPersona: &model.Persona{
			Name:       "ai persona 1",
			Email:      "email_1",
			Phone:      "phone_1",
			City:       "city_1",
			State:      "state_1",
			Country:    "country_1",
			YoE:        5,
			TechSkills: []string{"tech skill 1", "tech skill 2", "tech skill 3"},
		},
		ManuallyCreatedPersona: &model.Persona{
			Name:       "manual persona 1",
			Email:      "email_1",
			Phone:      "phone_1",
			City:       "city_1",
			State:      "state_1",
			Country:    "country_1",
			YoE:        5,
			TechSkills: []string{"tech skill 1", "tech skill 2", "tech skill 3"},
		},
		Team:         team,
		FileUploadId: "fp_id2",
	})
	candidate3, _ := model.NewCandidate(model.CandidateOptions{
		Id: "c_id3",
		ManuallyCreatedPersona: &model.Persona{
			Name:       "manual persona 1",
			Email:      "email_1",
			Phone:      "phone_1",
			City:       "city_1",
			State:      "state_1",
			Country:    "country_1",
			YoE:        5,
			TechSkills: []string{"tech skill 1", "tech skill 2", "tech skill 3"},
		},
		Team: team,
	})

	tests := []struct {
		name                  string
		ctx                   context.Context
		input                 *pb.GetCandidatesRequest
		output                *pb.GetCandidatesResponse
		teamHydratorMock      storage.TeamHydrator
		candidateAccessorMock storage.CandidateAccessor
		errorExpected         bool
		errorString           string
	}{
		{
			name:                  "errors if no user in context",
			ctx:                   context.Background(),
			input:                 &pb.GetCandidatesRequest{},
			output:                &pb.GetCandidatesResponse{},
			teamHydratorMock:      nil,
			candidateAccessorMock: nil,
			errorExpected:         true,
			errorString:           "rpc error: code = Unauthenticated desc = retrieving user data failed",
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
			input:                 &pb.GetCandidatesRequest{},
			output:                &pb.GetCandidatesResponse{},
			teamHydratorMock:      &storage.TeamHydratorMockFailure{},
			candidateAccessorMock: nil,
			errorExpected:         true,
			errorString:           "unable to hydrate team",
		},
		{
			name: "returns error if database errors when getting candidates",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input:            &pb.GetCandidatesRequest{},
			output:           nil,
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			candidateAccessorMock: &storage.CandidateAccessorConfigurableMock{
				GetCandidatesForTeamInternal: func(team *model.Team) ([]*model.Candidate, error) {
					return nil, errors.New("dbError when querying")
				},
			},
			errorExpected: true,
			errorString:   "dbError when querying",
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
			input: &pb.GetCandidatesRequest{},
			output: &pb.GetCandidatesResponse{
				Candidates: []*pb.Candidate{
					{
						Id:                 "c_id1",
						AiGeneratedPersona: "{\"Name\":\"ai persona 1\",\"Email\":\"email_1\",\"Phone\":\"phone_1\",\"City\":\"city_1\",\"State\":\"state_1\",\"Country\":\"country_1\",\"YoE\":5,\"Tech Skills\":[\"tech skill 1\",\"tech skill 2\",\"tech skill 3\"]}",
						FileUploadId:       "fp_id1",
					},
					{
						Id:                     "c_id2",
						AiGeneratedPersona:     "{\"Name\":\"ai persona 1\",\"Email\":\"email_1\",\"Phone\":\"phone_1\",\"City\":\"city_1\",\"State\":\"state_1\",\"Country\":\"country_1\",\"YoE\":5,\"Tech Skills\":[\"tech skill 1\",\"tech skill 2\",\"tech skill 3\"]}",
						ManuallyCreatedPersona: "{\"Name\":\"manual persona 1\",\"Email\":\"email_1\",\"Phone\":\"phone_1\",\"City\":\"city_1\",\"State\":\"state_1\",\"Country\":\"country_1\",\"YoE\":5,\"Tech Skills\":[\"tech skill 1\",\"tech skill 2\",\"tech skill 3\"]}",
						FileUploadId:           "fp_id2",
					},
					{
						Id:                     "c_id3",
						ManuallyCreatedPersona: "{\"Name\":\"manual persona 1\",\"Email\":\"email_1\",\"Phone\":\"phone_1\",\"City\":\"city_1\",\"State\":\"state_1\",\"Country\":\"country_1\",\"YoE\":5,\"Tech Skills\":[\"tech skill 1\",\"tech skill 2\",\"tech skill 3\"]}",
					},
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			candidateAccessorMock: &storage.CandidateAccessorConfigurableMock{
				GetCandidatesForTeamInternal: func(team *model.Team) ([]*model.Candidate, error) {
					return []*model.Candidate{
						candidate1, candidate2, candidate3,
					}, nil
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
					storage.WithCandidateAccessorMock(tt.candidateAccessorMock),
				),
				Logger: &utilities.NullLogger{},
			})

			response, err := server.GetCandidates(
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

func Test_GetCandidate(t *testing.T) {
	team, _ := model.NewTeam(model.TeamOptions{
		Id:   "team_id1",
		Name: "test@example.com",
	})
	userWithTeam, _ := model.NewUser(model.UserOptions{
		Id:    "user_id1",
		Email: "test@example.com",
		Team:  team,
	})

	candidate1, _ := model.NewCandidate(model.CandidateOptions{
		Id: "c_id1",
		AiGeneratedPersona: &model.Persona{
			Name:       "ai persona 1",
			Email:      "email_1",
			Phone:      "phone_1",
			City:       "city_1",
			State:      "state_1",
			Country:    "country_1",
			YoE:        5,
			TechSkills: []string{"tech skill 1", "tech skill 2", "tech skill 3"},
		},
		ManuallyCreatedPersona: &model.Persona{
			Name:       "manual persona 1",
			Email:      "email_1",
			Phone:      "phone_1",
			City:       "city_1",
			State:      "state_1",
			Country:    "country_1",
			YoE:        5,
			TechSkills: []string{"tech skill 1", "tech skill 2", "tech skill 3"},
		},
		Team:         team,
		FileUploadId: "fp_id1",
	})

	tests := []struct {
		name                  string
		ctx                   context.Context
		input                 *pb.GetCandidateRequest
		output                *pb.GetCandidateResponse
		teamHydratorMock      storage.TeamHydrator
		candidateAccessorMock storage.CandidateAccessor
		errorExpected         bool
		errorString           string
	}{
		{
			name:                  "errors if id is blank",
			ctx:                   context.Background(),
			input:                 &pb.GetCandidateRequest{},
			output:                &pb.GetCandidateResponse{},
			teamHydratorMock:      nil,
			candidateAccessorMock: nil,
			errorExpected:         true,
			errorString:           "id cannot be blank",
		},
		{
			name: "errors if no user in context",
			ctx:  context.Background(),
			input: &pb.GetCandidateRequest{
				Id: "c_id1",
			},
			output:                &pb.GetCandidateResponse{},
			teamHydratorMock:      nil,
			candidateAccessorMock: nil,
			errorExpected:         true,
			errorString:           "rpc error: code = Unauthenticated desc = retrieving user data failed",
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
			input: &pb.GetCandidateRequest{
				Id: "c_id1",
			},
			output:                &pb.GetCandidateResponse{},
			teamHydratorMock:      &storage.TeamHydratorMockFailure{},
			candidateAccessorMock: nil,
			errorExpected:         true,
			errorString:           "unable to hydrate team",
		},
		{
			name: "returns error if database errors when getting candidates",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.GetCandidateRequest{
				Id: "c_id1",
			},
			output:           nil,
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			candidateAccessorMock: &storage.CandidateAccessorConfigurableMock{
				GetCandidateForTeamInternal: func(id string, team *model.Team) (*model.Candidate, error) {
					return nil, errors.New("dbError when querying")
				},
			},
			errorExpected: true,
			errorString:   "dbError when querying",
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
			input: &pb.GetCandidateRequest{
				Id: "c_id1",
			},
			output: &pb.GetCandidateResponse{
				Candidate: &pb.Candidate{
					Id:                     "c_id1",
					AiGeneratedPersona:     "{\"Name\":\"ai persona 1\",\"Email\":\"email_1\",\"Phone\":\"phone_1\",\"City\":\"city_1\",\"State\":\"state_1\",\"Country\":\"country_1\",\"YoE\":5,\"Tech Skills\":[\"tech skill 1\",\"tech skill 2\",\"tech skill 3\"]}",
					ManuallyCreatedPersona: "{\"Name\":\"manual persona 1\",\"Email\":\"email_1\",\"Phone\":\"phone_1\",\"City\":\"city_1\",\"State\":\"state_1\",\"Country\":\"country_1\",\"YoE\":5,\"Tech Skills\":[\"tech skill 1\",\"tech skill 2\",\"tech skill 3\"]}",
					FileUploadId:           "fp_id1",
				},
			},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			candidateAccessorMock: &storage.CandidateAccessorConfigurableMock{
				GetCandidateForTeamInternal: func(id string, team *model.Team) (*model.Candidate, error) {
					return candidate1, nil
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
					storage.WithCandidateAccessorMock(tt.candidateAccessorMock),
				),
				Logger: &utilities.NullLogger{},
			})

			response, err := server.GetCandidate(
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

func Test_UpdateCandidate(t *testing.T) {
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
		name                  string
		ctx                   context.Context
		input                 *pb.UpdateCandidateRequest
		output                *pb.UpdateCandidateResponse
		teamHydratorMock      storage.TeamHydrator
		candidateAccessorMock storage.CandidateAccessor
		errorExpected         bool
		errorString           string
	}{
		{
			name:                  "errors if no user in context",
			ctx:                   context.Background(),
			input:                 &pb.UpdateCandidateRequest{},
			output:                &pb.UpdateCandidateResponse{},
			teamHydratorMock:      nil,
			candidateAccessorMock: nil,
			errorExpected:         true,
			errorString:           "rpc error: code = Unauthenticated desc = retrieving user data failed",
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
			input:                 &pb.UpdateCandidateRequest{},
			output:                &pb.UpdateCandidateResponse{},
			teamHydratorMock:      &storage.TeamHydratorMockFailure{},
			candidateAccessorMock: nil,
			errorExpected:         true,
			errorString:           "unable to hydrate team",
		},
		{
			name: "returns error if unable to parse persona",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input:                 &pb.UpdateCandidateRequest{},
			output:                nil,
			teamHydratorMock:      &storage.TeamHydratorMockSuccess{User: userWithTeam},
			candidateAccessorMock: &storage.CandidateAccessorConfigurableMock{},
			errorExpected:         true,
			errorString:           "unable to parse persona json: unexpected end of JSON input",
		},
		{
			name: "returns error if database errors when updating candidate",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.New(
					map[string]string{
						requestingUserIdCtxKey:    "user_id1",
						requestingUserEmailCtxKey: "user@example.com",
					},
				),
			),
			input: &pb.UpdateCandidateRequest{
				ManuallyCreatedPersona: "{\"Name\":\"manual persona 1\",\"Email\":\"email_1\",\"Phone\":\"phone_1\",\"City\":\"city_1\",\"State\":\"state_1\",\"Country\":\"country_1\",\"YoE\":5,\"Tech Skills\":[\"tech skill 1\",\"tech skill 2\",\"tech skill 3\"]}",
			},
			output:           nil,
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			candidateAccessorMock: &storage.CandidateAccessorConfigurableMock{
				UpdateCandidateWithManuallyCreatedPersonaForTeamInternal: func(id string, persona *model.Persona, team *model.Team) (string, error) {
					return "", errors.New("dbError when updating")
				},
			},
			errorExpected: true,
			errorString:   "dbError when updating",
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
			input: &pb.UpdateCandidateRequest{
				ManuallyCreatedPersona: "{\"Name\":\"manual persona 1\",\"Email\":\"email_1\",\"Phone\":\"phone_1\",\"City\":\"city_1\",\"State\":\"state_1\",\"Country\":\"country_1\",\"YoE\":5,\"Tech Skills\":[\"tech skill 1\",\"tech skill 2\",\"tech skill 3\"]}",
			},
			output:           &pb.UpdateCandidateResponse{Id: "new_id1"},
			teamHydratorMock: &storage.TeamHydratorMockSuccess{User: userWithTeam},
			candidateAccessorMock: &storage.CandidateAccessorConfigurableMock{
				UpdateCandidateWithManuallyCreatedPersonaForTeamInternal: func(id string, persona *model.Persona, team *model.Team) (string, error) {
					return "new_id1", nil
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
					storage.WithCandidateAccessorMock(tt.candidateAccessorMock),
				),
				Logger: &utilities.NullLogger{},
			})

			response, err := server.UpdateCandidate(
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
