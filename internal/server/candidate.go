package server

import (
	"context"

	pb "github.com/vipulvpatil/candidate-tracker-go/protos"
)

func (s *CandidateTrackerGoService) GetCandidates(ctx context.Context, req *pb.GetCandidatesRequest) (*pb.GetCandidatesResponse, error) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userWithTeam, err := s.storage.HydrateTeam(user)
	if err != nil {
		return nil, err
	}

	team := userWithTeam.Team()

	responseData := []*pb.Candidate{}
	candidates, err := s.storage.GetCandidatesForTeam(team)
	if err != nil {
		return nil, err
	}

	for _, candidate := range candidates {
		candidateResponse := pb.Candidate{
			Id:                     candidate.Id(),
			AiGeneratedPersona:     candidate.AiGeneratedPersonaAsJsonString(),
			ManuallyCreatedPersona: candidate.ManuallyCreatedPersonaAsJsonString(),
			FileUploadId:           candidate.FileUploadId(),
		}
		responseData = append(responseData, &candidateResponse)
	}

	return &pb.GetCandidatesResponse{
		Candidates: responseData,
	}, nil
}
