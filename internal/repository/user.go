package repository

import (
	"context"

	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/iDevoid/stygis/internal/module/user"
)

type userRepo struct {
	cache       user.Caching
	persistence user.Persistence
}

// UserInit to initiate the repository of user domain
func UserInit(cache user.Caching, persistence user.Persistence) user.Repository {
	return &userRepo{
		cache:       cache,
		persistence: persistence,
	}
}

// DataProfile gets the data user from cache, if it doesn't exist then it selects the data from persistence db and save it to cache
func (ur *userRepo) DataProfile(ctx context.Context, userID int64) (*model.User, error) {
	user, err := ur.cache.Get(ctx, userID)
	if err == nil {
		return user, nil
	}
	user, err = ur.persistence.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	err = ur.cache.Save(ctx, user)
	return user, err
}
