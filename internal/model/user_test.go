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
			name: "User gets created successfully with nil Team",
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
		{
			name: "User gets created successfully with non-nil Team",
			input: UserOptions{
				Id:    "123",
				Email: "test@example.com",
				Team: &Team{
					id:   "team_id1",
					name: "team_name1",
				},
			},
			expectedOutput: &User{
				id:    "123",
				email: "test@example.com",
				team: &Team{
					id:   "team_id1",
					name: "team_name1",
				},
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

func Test_User_GetId(t *testing.T) {
	t.Run("GetId returns user's id", func(t *testing.T) {
		user := &User{
			id:    "123",
			email: "test@example.com",
		}
		assert.Equal(t, "123", user.GetId())
	})
}

func Test_User_Team(t *testing.T) {
	t.Run("Team returns user's team", func(t *testing.T) {
		user := &User{
			id:    "123",
			email: "test@example.com",
			team: &Team{
				id:   "team_id1",
				name: "team_name1",
			},
		}
		team := user.Team()
		assert.NotNil(t, team)
		assert.Equal(t, "team_id1", team.id)
		assert.Equal(t, "team_name1", team.name)
	})

	t.Run("Team returns nil if user has no team", func(t *testing.T) {
		user := &User{
			id:    "123",
			email: "test@example.com",
		}
		team := user.Team()
		assert.Nil(t, team)
	})
}
