package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/iDevoid/stygis/internal/module/user"
	"github.com/iDevoid/stygis/internal/constants/model"
)

type userCache struct {
	connection *redis.Client
}

const (
	keyRedisUser   = "user:%d"
	expireRedisKey = 60 * 60 * 24
)

// UserInit is the function to init the user caching
func UserInit(conn *redis.Client) user.Caching {
	return &userCache{
		connection: conn,
	}
}

func (uc *userCache) Save(ctx context.Context, user *model.User) error {
	data, _ := json.Marshal(user)
	key := fmt.Sprintf(keyRedisUser, user.ID)

	err := uc.connection.Set(key, data, time.Duration(expireRedisKey)).Err()
	if err != nil {
		return err
	}

	return err
}

func (uc *userCache) Get(ctx context.Context, userID int64) (*model.User, error) {
	key := fmt.Sprintf(keyRedisUser, userID)

	res, err := uc.connection.Get(key).Result()
	if err != nil {
		return nil, err
	}

	var user *model.User
	err = json.Unmarshal([]byte(res), user)
	return user, err
}

func (uc *userCache) Delete(ctx context.Context, userID int64) error {
	key := fmt.Sprintf(keyRedisUser, userID)
	err := uc.connection.Del(key).Err()
	return err
}
