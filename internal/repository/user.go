package repository

import (
	"context"

	"github.com/iDevoid/stygis/internal/module/user"
	"github.com/iDevoid/stygis/internal/constants/model"
)

type userRepo struct {
	cache       user.Caching
	persistance user.Persistence
}

// UserInit to initiate the repository of user domain
func UserInit(cache user.Caching, persistance user.Persistence) user.Repository {
	return &userRepo{
		cache:       cache,
		persistance: persistance,
	}
}

func (ur *userRepo) DataProfile(ctx context.Context, userID int64) (*model.User, error) {
	user, err := ur.cache.Get(ctx, userID)
	if err == nil {
		return user, nil
	}
	user, err = ur.persistance.FindByID(ctx, userID)
	return user, err
}
