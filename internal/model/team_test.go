package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewTeam(t *testing.T) {
	tests := []struct {
		name           string
		input          TeamOptions
		expectedOutput *Team
		errorExpected  bool
		errorString    string
	}{
		{
			name:           "id is empty",
			input:          TeamOptions{},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create team with a empty id",
		},
		{
			name: "name is empty",
			input: TeamOptions{
				Id: "123",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create team with a empty name",
		},
		{
			name: "user is nil",
			input: TeamOptions{
				Id:   "123",
				Name: "test",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create team with a nil user",
		},
		{
			name: "User gets created successfully",
			input: TeamOptions{
				Id:   "123",
				Name: "test",
				User: &User{
					id:    "user_id1",
					email: "test@example.com",
				},
			},
			expectedOutput: &Team{
				id:   "123",
				name: "test",
				user: &User{
					id:    "user_id1",
					email: "test@example.com",
				},
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewTeam(tt.input)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
