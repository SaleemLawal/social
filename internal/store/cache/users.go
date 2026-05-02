package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/saleemlawal/social/internal/store"
)

const userExpirationTime = time.Minute

type UserStore struct {
	redisClient *redis.Client
}

func (s *UserStore) Get(ctx context.Context, id int64) (*store.User, error) {
	key := fmt.Sprintf("user:%v", id)
	val, err := s.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if val == "" {
		return nil, store.ErrRecordNotFound
	}

	user := &store.User{}
	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	key := fmt.Sprintf("user:%v", user.ID)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, key, json, userExpirationTime).Err()
}
