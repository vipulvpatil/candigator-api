package server

import (
	"context"

	"github.com/vipulvpatil/candidate-tracker-go/internal/clients/openai"
	"github.com/vipulvpatil/candidate-tracker-go/internal/config"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
	pb "github.com/vipulvpatil/candidate-tracker-go/protos"
)

type CandidateTrackerGoService struct {
	pb.UnsafeCandidateTrackerGoServer
	storage      storage.StorageAccessor
	openAiClient openai.Client
	config       *config.Config
	logger       utilities.Logger
}

type ServerDependencies struct {
	Storage      storage.StorageAccessor
	OpenAiClient openai.Client
	Config       *config.Config
	Logger       utilities.Logger
}

func NewServer(deps ServerDependencies) (*CandidateTrackerGoService, error) {
	return &CandidateTrackerGoService{
		storage:      deps.Storage,
		openAiClient: deps.OpenAiClient,
		config:       deps.Config,
		logger:       deps.Logger,
	}, nil
}

func (s *CandidateTrackerGoService) CheckConnection(ctx context.Context, req *pb.CheckConnectionRequest) (*pb.CheckConnectionResponse, error) {
	return &pb.CheckConnectionResponse{
		ConnectionStatus: "okay",
	}, nil
}
