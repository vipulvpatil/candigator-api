package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"google.golang.org/grpc/metadata"
)

func Test_contextWithUserData(t *testing.T) {
	tests := []struct {
		name              string
		input             context.Context
		expectedOutput    metadata.MD
		userRetrieverMock storage.UserRetriever
		errorExpected     bool
		errorString       string
	}{
		{
			name: "populates the context with retrieved user data",
			input: metadata.NewIncomingContext(
				context.Background(),
				metadata.Pairs(requestingUserEmailCtxKey, "test@example.com"),
			),
			expectedOutput: metadata.Pairs(
				requestingUserIdCtxKey, "1",
				requestingUserEmailCtxKey, "test@example.com",
			),
			userRetrieverMock: &storage.UserRetrieverMockSuccess{
				Id:    "1",
				Email: "test@example.com",
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name:              "errors when metadata is not present in incoming context",
			input:             context.Background(),
			expectedOutput:    nil,
			userRetrieverMock: nil, //&storage.UserRetrieverMockEmpty{},
			errorExpected:     true,
			errorString:       "retrieving metadata has failed",
		},
		{
			name:              "errors when requesting_user_email is not present in incoming context metadata",
			input:             metadata.NewIncomingContext(context.Background(), metadata.Pairs()),
			expectedOutput:    nil,
			userRetrieverMock: nil, //&storage.UserRetrieverMockEmpty{},
			errorExpected:     true,
			errorString:       "requesting_user_email is not supplied",
		},
		{
			name: "errors when requesting user email is not retrieved from storage",
			input: metadata.NewIncomingContext(
				context.Background(),
				metadata.Pairs(requestingUserEmailCtxKey, "test@example.com"),
			),
			expectedOutput:    nil,
			userRetrieverMock: &storage.UserRetrieverMockFailure{},
			errorExpected:     true,
			errorString:       "cannot find user by email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedCtx, err := contextWithUserData(tt.input, tt.userRetrieverMock)
			if !tt.errorExpected {
				assert.NoError(t, err)
				md, ok := metadata.FromIncomingContext(updatedCtx)
				assert.Equal(t, ok, true)
				assert.Equal(t, tt.expectedOutput, md)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.ErrorContains(t, err, tt.errorString)
			}
		})
	}
}

func Test_getUserFromContext(t *testing.T) {
	returnUser, _ := model.NewUser(
		model.UserOptions{
			Id:    "1",
			Email: "test@example.com",
		},
	)

	tests := []struct {
		name           string
		input          context.Context
		expectedOutput *model.User
		errorExpected  bool
		errorString    string
	}{
		{
			name: "returns user from context metadata",
			input: metadata.NewIncomingContext(
				context.Background(),
				metadata.Pairs(
					requestingUserIdCtxKey, "1",
					requestingUserEmailCtxKey, "test@example.com",
				),
			),
			expectedOutput: returnUser,
			errorExpected:  false,
			errorString:    "",
		},
		{
			name:           "errors when metadata is not present in incoming context",
			input:          context.Background(),
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "retrieving user data failed",
		},
		{
			name: "Could not retrieve requesting_user_email from context",
			input: metadata.NewIncomingContext(
				context.Background(),
				metadata.Pairs(),
			),
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "could not retrieve requesting_user_email from context",
		},
		{
			name: "Could not retrieve requesting_user_id from context",
			input: metadata.NewIncomingContext(
				context.Background(),
				metadata.Pairs(
					requestingUserEmailCtxKey, "test@example.com",
				),
			),
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "could not retrieve requesting_user_id from context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := getUserFromContext(tt.input)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, user)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.ErrorContains(t, err, tt.errorString)
			}
		})
	}
}
