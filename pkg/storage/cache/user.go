package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"

	"github.com/iDevoid/stygis/pkg/internal/user"
	"github.com/iDevoid/stygis/pkg/model"
)

type userCache struct {
	pool *redis.Pool
}

const (
	keyRedisUser   = "user:%d"
	expireRedisKey = 60 * 60 * 24
)

// InitUserCache is the function to init the user caching
func InitUserCache(pool *redis.Pool) user.Caching {
	return &userCache{
		pool: pool,
	}
}

func (uc *userCache) Save(ctx context.Context, user *model.User) error {
	conn, err := uc.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, _ := json.Marshal(user)
	key := fmt.Sprintf(keyRedisUser, user.ID)
	_, err = conn.Do("SET", key, data, "EX", expireRedisKey)
	return err
}

func (uc *userCache) Get(ctx context.Context, userID int64) (*model.User, error) {
	conn, err := uc.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	key := fmt.Sprintf(keyRedisUser, userID)
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	var user *model.User
	err = json.Unmarshal(data, user)
	return user, err
}

func (uc *userCache) Delete(ctx context.Context, userID int64) error {
	conn, err := uc.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	key := fmt.Sprintf(keyRedisUser, userID)
	_, err = conn.Do("DEL", key)
	return err
}
