package server

import (
	"context"

	"google.golang.org/grpc"
)

// The calls to this service are authenticated using mutual TLS.
// This following interceptor adds a valid user if one exists
// on whose behalf the current request has been made.
func (s *CandidateTrackerGoService) RequestingUserInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	updatedCtx, err := contextWithUserData(ctx, s.storage)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	} else {
		return handler(updatedCtx, req)
	}
}
