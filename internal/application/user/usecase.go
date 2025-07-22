package user

import "context"

type UserUsecase interface {
	GetUser(ctx context.Context) (string, error)
}
