package repository

import (
	"context"

	"github.com/iDevoid/stygis/internal/module/user"
	"github.com/iDevoid/stygis/internal/constants/model"
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

func (ur *userRepo) DataProfile(ctx context.Context, userID int64) (*model.User, error) {
	user, err := ur.cache.Get(ctx, userID)
	if err == nil {
		return user, nil
	}
	user, err = ur.persistence.FindByID(ctx, userID)
	return user, err
}
