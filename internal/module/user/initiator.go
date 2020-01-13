package user

import (
	"context"

	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/iDevoid/stygis/platform/routers"
	atreugo "github.com/savsgio/atreugo/v10"
)

// Persistence initiator includes the functions from storage psql
type Persistence interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, userID int64) (*model.User, error)
	Find(ctx context.Context, email, password string) (*model.User, error)
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
	userCaching     Caching
	userRepository  Repository
}

// Usecase would be use to contain the business logic functions
type Usecase interface {
	NewAccount(ctx context.Context, user *model.User) error
}

// InitializeDomain is the function to initiate the business logic with services that'll be used by business logic
func InitializeDomain(persistence Persistence, caching Caching, repository Repository) Usecase {
	return &service{
		userPersistence: persistence,
		userCaching:     caching,
		userRepository:  repository,
	}
}

// Handler contains all the functions for handling http request
type Handler interface {
	Test(ctx *atreugo.RequestCtx) error
	CreateNewAccount(ctx *atreugo.RequestCtx) error
}

// Router contains the functions that will be used for the routing domain user
type Router interface {
	NewRouters() []*routers.Router
}
