package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"google.golang.org/grpc/metadata"
)

type mockRequestWithUserEmail struct{}

func (m *mockRequestWithUserEmail) GetUserEmail() string {
	return "test@example.com"
}

type mockRequestWithEmptyUserEmail struct{}

func (m *mockRequestWithEmptyUserEmail) GetUserEmail() string {
	return ""
}

func Test_contextWithUserData(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			ctx context.Context
			req RequestWithUserEmail
		}
		expectedOutput    metadata.MD
		userRetrieverMock storage.UserRetriever
		errorExpected     bool
		errorString       string
	}{
		{
			name: "populates the context with retrieved user data based on email in context",
			input: struct {
				ctx context.Context
				req RequestWithUserEmail
			}{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.Pairs(requestingUserEmailCtxKey, "test@example.com"),
				),
				req: &mockRequestWithUserEmail{},
			},
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
			name: "populates the context with retrieved user data based on email in request",
			input: struct {
				ctx context.Context
				req RequestWithUserEmail
			}{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs()),
				req: &mockRequestWithUserEmail{},
			},
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
			name: "errors when metadata is not present in incoming context",
			input: struct {
				ctx context.Context
				req RequestWithUserEmail
			}{
				ctx: context.Background(),
				req: nil,
			},
			expectedOutput:    nil,
			userRetrieverMock: nil,
			errorExpected:     true,
			errorString:       "retrieving metadata has failed",
		},
		{
			name: "errors when requesting_user_email is not present in incoming context metadata or request",
			input: struct {
				ctx context.Context
				req RequestWithUserEmail
			}{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs()),
				req: &mockRequestWithEmptyUserEmail{},
			},
			expectedOutput:    nil,
			userRetrieverMock: nil,
			errorExpected:     true,
			errorString:       "unable to determine requesting user",
		},
		{
			name: "errors when requesting user email is not retrieved from storage",
			input: struct {
				ctx context.Context
				req RequestWithUserEmail
			}{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.Pairs(requestingUserEmailCtxKey, "test@example.com"),
				),
				req: &mockRequestWithUserEmail{},
			},
			expectedOutput:    nil,
			userRetrieverMock: &storage.UserRetrieverMockFailure{},
			errorExpected:     true,
			errorString:       "cannot find user by email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedCtx, err := contextWithUserData(tt.input.ctx, tt.input.req, tt.userRetrieverMock)
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
