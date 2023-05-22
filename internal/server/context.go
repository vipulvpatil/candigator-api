package server

import (
	"context"

	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const requestingUserIdCtxKey = "requesting_user_id"
const requestingUserEmailCtxKey = "requesting_user_email"

func contextWithUserData(ctx context.Context, userRetriever storage.UserRetriever) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.InvalidArgument, "retrieving metadata has failed")
	}

	requestingUserEmails := md.Get(requestingUserEmailCtxKey)
	if len(requestingUserEmails) < 1 {
		return ctx, status.Errorf(codes.Unauthenticated, "%s is not supplied", requestingUserEmailCtxKey)
	}

	user, err := userRetriever.UserByEmail(requestingUserEmails[0])
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, err.Error())
	}

	md.Append(requestingUserIdCtxKey, user.GetId())

	updatedCtx := metadata.NewIncomingContext(ctx, md)

	return updatedCtx, nil
}

func getUserFromContext(ctx context.Context) (*model.User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "retrieving user data failed")
	}

	requestingUserEmails := md.Get(requestingUserEmailCtxKey)
	if len(requestingUserEmails) < 1 {
		return nil, status.Errorf(codes.Unauthenticated, "could not retrieve %s from context", requestingUserEmailCtxKey)
	}

	requestingUserIds := md.Get(requestingUserIdCtxKey)
	if len(requestingUserIds) < 1 {
		return nil, status.Errorf(codes.Unauthenticated, "could not retrieve %s from context", requestingUserIdCtxKey)
	}

	return model.NewUser(
		model.UserOptions{
			Id:    requestingUserIds[0],
			Email: requestingUserEmails[0],
		},
	)
}
