package cache

import (
	"context"

	"github.com/gomodule/redigo/redis"

	"github.com/iDevoid/stygis/pkg/internal/user"
	"github.com/iDevoid/stygis/pkg/model"
)

type userCache struct {
	pool *redis.Pool
}

// InitUserCache is the function to init the user caching
func InitUserCache(pool *redis.Pool) user.Caching {
	return &userCache{
		pool: pool,
	}
}

func (uc *userCache) Save(context context.Context, user *model.User) error {
	return nil
}

func (uc *userCache) Get(context context.Context, userID int64) (*model.User, error) {
	return nil, nil
}

func (uc *userCache) Delete(context context.Context, userID int64) error {
	return nil
}
