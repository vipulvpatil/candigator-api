package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/lib/parser/personabuilder"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
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

func (s *CandidateTrackerGoService) GetCandidate(ctx context.Context, req *pb.GetCandidateRequest) (*pb.GetCandidateResponse, error) {
	id := req.GetId()
	if utilities.IsBlank(id) {
		return nil, errors.New("id cannot be blank")
	}

	user, err := getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userWithTeam, err := s.storage.HydrateTeam(user)
	if err != nil {
		return nil, err
	}

	team := userWithTeam.Team()

	candidate, err := s.storage.GetCandidateForTeam(id, team)
	if err != nil {
		return nil, err
	}

	return &pb.GetCandidateResponse{
		Candidate: &pb.Candidate{
			Id:                     candidate.Id(),
			AiGeneratedPersona:     candidate.AiGeneratedPersonaAsJsonString(),
			ManuallyCreatedPersona: candidate.ManuallyCreatedPersonaAsJsonString(),
			FileUploadId:           candidate.FileUploadId(),
		},
	}, nil
}

func (s *CandidateTrackerGoService) UpdateCandidate(ctx context.Context, req *pb.UpdateCandidateRequest) (*pb.UpdateCandidateResponse, error) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userWithTeam, err := s.storage.HydrateTeam(user)
	if err != nil {
		return nil, err
	}

	team := userWithTeam.Team()
	candidateId := req.GetId()
	personaJson := req.GetManuallyCreatedPersona()

	persona, err := personabuilder.ParsePersonaFromJson(personaJson)
	if err != nil {
		return nil, err
	}
	persona.BuiltBy = "HUMAN"

	err = s.storage.UpdateCandidateWithManuallyCreatedPersonaForTeam(candidateId, persona, team)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateCandidateResponse{}, nil
}
