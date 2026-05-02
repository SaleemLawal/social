package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/saleemlawal/social/internal/store"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
	// Posts interface {
	// 	Get(context.Context, int64) (*store.Post, error)
	// 	Set(context.Context, int64, *store.Post) error
	// }
	// Comments interface {
	// 	Get(context.Context, int64) (*store.Comment, error)
	// 	Set(context.Context, int64, *store.Comment) error
	// }
}

func NewRedisStorage(redisClient *redis.Client) Storage {
	return Storage{
		Users: &UserStore{redisClient},
	}
}
