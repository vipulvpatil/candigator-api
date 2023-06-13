package server

import (
	"context"

	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	pb "github.com/vipulvpatil/candidate-tracker-go/protos"
)

func (s *CandidateTrackerGoService) UploadFiles(ctx context.Context, req *pb.UploadFilesRequest) (*pb.UploadFilesResponse, error) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userWithTeam, err := s.storage.HydrateTeam(user)
	if err != nil {
		return nil, err
	}

	team := userWithTeam.Team()
	responseData := []*pb.FileUpload{}

	files := req.GetFiles()
	for _, file := range files {
		fileName := file.Name
		responseData = append(responseData, s.newFileUploadForTeam(fileName, team))
	}

	return &pb.UploadFilesResponse{
		FileUploads: responseData,
	}, nil
}

func (s *CandidateTrackerGoService) CompleteFileUpload(ctx context.Context, req *pb.CompleteFileUploadRequest) (*pb.CompleteFileUploadResponse, error) {
	return nil, nil
}

func (s *CandidateTrackerGoService) newFileUploadForTeam(fileName string, team *model.Team) *pb.FileUpload {
	fileUploadResponse := pb.FileUpload{
		Name: fileName,
	}

	fileUpload, err := s.storage.CreateFileUploadForTeam(fileName, team)
	if err != nil {
		return fileUploadResponseWithError(&fileUploadResponse, err)
	}
	fileUploadId := fileUpload.Id()
	fileUploadResponse.Id = fileUploadId
	fileUploadResponse.Status = fileUpload.Status()

	presignedUrl, err := s.fileStorer.GetPresignedUrl(team.Id(), fileUploadId, fileName)
	if err != nil {
		return fileUploadResponseWithError(&fileUploadResponse, err)
	}
	fileUploadResponse.PresignedUrl = presignedUrl

	err = s.storage.UpdateFileUploadWithPresignedUrl(fileUploadId, presignedUrl)
	if err != nil {
		return fileUploadResponseWithError(&fileUploadResponse, err)
	}
	return &fileUploadResponse
}

func fileUploadResponseWithError(fileUploadResponse *pb.FileUpload, err error) *pb.FileUpload {
	fileUploadResponse.Error = err.Error()
	return fileUploadResponse
}
