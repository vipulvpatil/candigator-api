package server

import (
	"context"

	pb "github.com/vipulvpatil/candidate-tracker-go/protos"
)

func (s *CandidateTrackerGoService) GetUserData(ctx context.Context, req *pb.GetUserDataRequest) (*pb.GetUserDataResponse, error) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userWithTeam, err := s.storage.HydrateTeam(user)
	if err != nil {
		return nil, err
	}

	team := userWithTeam.Team()

	count, err := s.storage.GetUnprocessedFileUploadsCountForTeam(team)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserDataResponse{
		FileCountLimit:       int64(team.FileCountLimit()),
		CurrentFileCount:     int64(team.CurrentFileCount()),
		UnprocessedFileCount: int64(count),
	}, nil
}
