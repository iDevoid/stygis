package user

import (
	"context"

	"github.com/iDevoid/stygis/pkg/model"
)

// Persistence initiator includes the functions from storage psql
type Persistence interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, userID int64) (user *model.User, err error)
	Find(ctx context.Context, email, password string) (user *model.User, err error)
	ChangePassword(ctx context.Context, newPassword string, user *model.User) error
	Delete(ctx context.Context, user *model.User) error
}

// Caching initiator contains functions to get data from redis
type Caching interface {
}

type service struct {
	userRepository Persistence
}

// Usecase would be use to contain the business logic functions
type Usecase interface{}

// InitUsecase is the function to initiate the business logic with services that'll be used by business logic
func InitUsecase(repo Persistence) Usecase {
	return &service{
		userRepository: repo,
	}
}
