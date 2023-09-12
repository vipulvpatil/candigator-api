package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewCandidate(t *testing.T) {
	currentFileCount := 1
	team, _ := NewTeam(TeamOptions{
		Id:               "team_id1",
		Name:             "test@example.com",
		CurrentFileCount: &currentFileCount,
		FileCountLimit:   100,
	})
	tests := []struct {
		name           string
		input          CandidateOptions
		expectedOutput *Candidate
		errorExpected  bool
		errorString    string
	}{
		{
			name:           "id is empty",
			input:          CandidateOptions{},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create Candidate with an empty id",
		},
		{
			name: "has nil Team",
			input: CandidateOptions{
				Id: "123",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create Candidate with a nil Team",
		},
		{
			name: "has no associated persoa",
			input: CandidateOptions{
				Id:   "123",
				Team: team,
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create Candidate without a valid persona",
		},
		{
			name: "Candidate gets created successfully with AI Generated Persona",
			input: CandidateOptions{
				Id:                 "123",
				Team:               team,
				AiGeneratedPersona: &Persona{Name: "persona name", BuiltBy: "AI", FileUploadId: "fp_id1"},
			},
			expectedOutput: &Candidate{
				id:                 "123",
				team:               team,
				aiGeneratedPersona: &Persona{Name: "persona name", BuiltBy: "AI", FileUploadId: "fp_id1"},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "Candidate gets created successfully with Manual Persona",
			input: CandidateOptions{
				Id:                     "123",
				Team:                   team,
				ManuallyCreatedPersona: &Persona{Name: "persona name"},
			},
			expectedOutput: &Candidate{
				id:                     "123",
				team:                   team,
				manuallyCreatedPersona: &Persona{Name: "persona name"},
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewCandidate(tt.input)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_Candidate_Id(t *testing.T) {
	currentFileCount := 1
	team, _ := NewTeam(TeamOptions{
		Id:               "team_id1",
		Name:             "test@example.com",
		CurrentFileCount: &currentFileCount,
		FileCountLimit:   100,
	})
	t.Run("returns Id", func(t *testing.T) {
		candidate := &Candidate{
			id: "c_id1",
			aiGeneratedPersona: &Persona{
				Name: "candidate_1",
			},
			team:         team,
			fileUploadId: "fp_id1",
		}
		assert.Equal(t, "c_id1", candidate.Id())
	})
}
func Test_Candidate_FileUploadId(t *testing.T) {
	currentFileCount := 1
	team, _ := NewTeam(TeamOptions{
		Id:               "team_id1",
		Name:             "test@example.com",
		CurrentFileCount: &currentFileCount,
		FileCountLimit:   100,
	})
	t.Run("returns FileUploadId", func(t *testing.T) {
		candidate := &Candidate{
			id: "c_id1",
			aiGeneratedPersona: &Persona{
				Name: "candidate_1",
			},
			team:         team,
			fileUploadId: "fp_id1",
		}
		assert.Equal(t, "fp_id1", candidate.FileUploadId())
	})
}

func Test_Candidate_AiGeneratedPersonaAsJsonString(t *testing.T) {
	currentFileCount := 1
	team, _ := NewTeam(TeamOptions{
		Id:               "team_id1",
		Name:             "test@example.com",
		CurrentFileCount: &currentFileCount,
		FileCountLimit:   100,
	})
	t.Run("returns AiGeneratedPersonaAsJsonString", func(t *testing.T) {
		candidate := &Candidate{
			id: "c_id1",
			aiGeneratedPersona: &Persona{
				Name: "candidate_1",
			},
			team:         team,
			fileUploadId: "fp_id1",
		}
		assert.Equal(t, "{\"Name\":\"candidate_1\"}", candidate.AiGeneratedPersonaAsJsonString())
	})

	t.Run("returns AiGeneratedPersonaAsJsonString", func(t *testing.T) {
		candidate := &Candidate{
			id:           "c_id1",
			team:         team,
			fileUploadId: "fp_id1",
		}
		assert.Equal(t, "", candidate.AiGeneratedPersonaAsJsonString())
	})
}

func Test_Candidate_ManuallyCreatedPersonaAsJsonString(t *testing.T) {
	currentFileCount := 1
	team, _ := NewTeam(TeamOptions{
		Id:               "team_id1",
		Name:             "test@example.com",
		CurrentFileCount: &currentFileCount,
		FileCountLimit:   100,
	})
	t.Run("returns ManuallyCreatedPersonaAsJsonString", func(t *testing.T) {
		candidate := &Candidate{
			id: "c_id1",
			manuallyCreatedPersona: &Persona{
				Name: "candidate_1",
			},
			team:         team,
			fileUploadId: "fp_id1",
		}
		assert.Equal(t, "{\"Name\":\"candidate_1\"}", candidate.ManuallyCreatedPersonaAsJsonString())
	})

	t.Run("returns ManuallyCreatedPersonaAsJsonString", func(t *testing.T) {
		candidate := &Candidate{
			id:           "c_id1",
			team:         team,
			fileUploadId: "fp_id1",
		}
		assert.Equal(t, "", candidate.ManuallyCreatedPersonaAsJsonString())
	})
}
