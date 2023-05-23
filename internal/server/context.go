package server

import (
	"context"

	"github.com/vipulvpatil/candidate-tracker-go/internal/model"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const requestingUserIdCtxKey = "requesting_user_id"
const requestingUserEmailCtxKey = "requesting_user_email"

func contextWithUserData(ctx context.Context, req RequestWithUserEmail, userRetriever storage.UserRetriever) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.InvalidArgument, "retrieving metadata has failed")
	}

	requestingUserEmail := ""
	requestingUserEmails := md.Get(requestingUserEmailCtxKey)
	if len(requestingUserEmails) > 0 {
		requestingUserEmail = requestingUserEmails[0]
	} else {
		// TODO: This is unclean code.
		// I found a bug in the GRPC JS library where in the metadata passed in, gets dropped intermittently.
		// Until it is debugged further, we will pass the requesting user in the request alongside the metadata.
		// So if it is not found in the metadata we can get it from the request.
		// If found in metadata, it is populated in context for futher use.
		userEmailInRequest := req.GetUserEmail()
		if utilities.IsBlank(userEmailInRequest) {
			return ctx, status.Errorf(codes.Unauthenticated, "unable to determine requesting user")
		} else {
			requestingUserEmail = userEmailInRequest
		}
		md.Append(requestingUserEmailCtxKey, requestingUserEmail)
	}

	user, err := userRetriever.UserByEmail(requestingUserEmail)
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
