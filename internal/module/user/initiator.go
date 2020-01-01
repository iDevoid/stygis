package user

import (
	"context"

	"github.com/iDevoid/stygis/internal/constants/model"
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
	Save(ctx context.Context, user *model.User) error
	Get(ctx context.Context, userID int64) (*model.User, error)
	Delete(ctx context.Context, userID int64) error
}

// Repository is the data logic of the flow of getting or storing data
type Repository interface {
	DataProfile(ctx context.Context, userID int64) (*model.User, error)
}

type service struct {
	userPersistence Persistence
	userRepository  Repository
}

// Usecase would be use to contain the business logic functions
type Usecase interface {
	NewAccount(ctx context.Context, user *model.User) error
}

// InitializeDomain is the function to initiate the business logic with services that'll be used by business logic
func InitializeDomain(persistance Persistence, repository Repository) Usecase {
	return &service{
		userPersistence: persistance,
		userRepository:  repository,
	}
}
