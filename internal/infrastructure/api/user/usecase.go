package user

import (
	"context"
	"errors"

	"github.com/qkitzero/event-service/internal/application/user"
	userv1 "github.com/qkitzero/user-service/gen/go/user/v1"
	"google.golang.org/grpc/metadata"
)

type userUsecase struct {
	client userv1.UserServiceClient
}

func NewUserUsecase(client userv1.UserServiceClient) user.UserUsecase {
	return &userUsecase{client: client}
}

func (s *userUsecase) GetUser(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata is missing")
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	getUserRequest := &userv1.GetUserRequest{}

	getUserResponse, err := s.client.GetUser(ctx, getUserRequest)
	if err != nil {
		return "", err
	}

	return getUserResponse.GetUserId(), nil
}
