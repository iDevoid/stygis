package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/iDevoid/stygis/internal/module/user"
)

type userCache struct {
	connection *redis.Client
}

const (
	keyRedisUser   = "user:%d"
	expireRedisKey = time.Second * 60 * 60 * 24
)

// UserInit is the function to init the user caching
func UserInit(conn *redis.Client) user.Caching {
	return &userCache{
		connection: conn,
	}
}

// Save takes the struct of user data, and save it to redis as json string
func (uc *userCache) Save(ctx context.Context, user *model.User) error {
	data, _ := json.Marshal(user)
	key := fmt.Sprintf(keyRedisUser, user.ID)
	user.Password = "" // no sensitive data allowed to be saved in cache

	err := uc.connection.Set(key, data, time.Duration(expireRedisKey)).Err()
	return err
}

// Get returns the cached user data in redis
func (uc *userCache) Get(ctx context.Context, userID int64) (*model.User, error) {
	key := fmt.Sprintf(keyRedisUser, userID)

	res, err := uc.connection.Get(key).Result()
	if err != nil {
		return nil, err
	}

	user := new(model.User)
	err = json.Unmarshal([]byte(res), user)
	return user, err
}

// Delete removes the cache by userID
func (uc *userCache) Delete(ctx context.Context, userID int64) error {
	key := fmt.Sprintf(keyRedisUser, userID)
	err := uc.connection.Del(key).Err()
	return err
}
