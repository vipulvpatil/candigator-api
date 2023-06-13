package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
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

func (s *CandidateTrackerGoService) CompleteFileUploads(ctx context.Context, req *pb.CompleteFileUploadsRequest) (*pb.CompleteFileUploadsResponse, error) {
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

	fileUploadUpdates := req.GetFileUploadUpdates()
	for _, fileUploadUpdate := range fileUploadUpdates {
		responseData = append(responseData, s.getUpdatedFileUploadForTeam(fileUploadUpdate, team))
	}

	return &pb.CompleteFileUploadsResponse{
		FileUploads: responseData,
	}, nil
}

func (s *CandidateTrackerGoService) getUpdatedFileUploadForTeam(fileUploadUpdate *pb.FileUploadUpdate, team *model.Team) *pb.FileUpload {
	fileUploadResponse := pb.FileUpload{
		Id: fileUploadUpdate.GetId(),
	}

	fileUpload, err := s.storage.GetFileUpload(fileUploadUpdate.GetId())
	if err != nil {
		return fileUploadResponseWithError(&fileUploadResponse, errors.Wrap(err, "unable to get fileUpload"))
	}

	if fileUpload == nil {
		return fileUploadResponseWithError(&fileUploadResponse, errors.New("unable to get fileUpload"))
	}

	fileUploadResponse.Name = fileUpload.Name()
	fileUploadResponse.PresignedUrl = fileUpload.PresignedUrl()

	if fileUpload.Completed() {
		return fileUploadResponseWithError(&fileUploadResponse, errors.New("unable to update fileUpload"))
	}

	if !fileUpload.BelongsToTeam(team) {
		return fileUploadResponseWithError(&fileUploadResponse, utilities.NewBadError("unauthorized fileUpload attempted"))
	}

	updateStatus := model.FileUploadStatus(fileUploadUpdate.GetStatus())
	if updateStatus.Valid() {
		return fileUploadResponseWithError(&fileUploadResponse, errors.New("invalid update status"))
	}

	err = s.storage.UpdateFileUploadWithStatus(fileUploadResponse.GetId(), fileUploadUpdate.GetStatus())
	if err != nil {
		return fileUploadResponseWithError(&fileUploadResponse, utilities.WrapBadError(err, "unable to update fileUpload"))
	}

	fileUploadResponse.Status = updateStatus.String()

	return &fileUploadResponse
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
