package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewUser(t *testing.T) {
	tests := []struct {
		name           string
		input          UserOptions
		expectedOutput *User
		errorExpected  bool
		errorString    string
	}{
		{
			name: "id is empty",
			input: UserOptions{
				Email: "test@example.com",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create user with a empty id",
		},
		{
			name: "email is empty",
			input: UserOptions{
				Id: "123",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create user with a empty email",
		},
		{
			name: "User gets created successfully",
			input: UserOptions{
				Id:    "123",
				Email: "test@example.com",
			},
			expectedOutput: &User{
				id:    "123",
				email: "test@example.com",
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewUser(tt.input)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_UserGetId(t *testing.T) {
	t.Run("GetId returns user's id", func(t *testing.T) {
		user := &User{
			id:    "123",
			email: "test@example.com",
		}
		assert.Equal(t, "123", user.GetId())
	})
}
