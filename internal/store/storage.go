package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound      = errors.New("resource not found")
	QUERY_TIMEOUT_DURATION = 5 * time.Second
	ErrConflict            = errors.New("resource already exists")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
		GetFeeds(context.Context, int64, *PaginationFeedsQuery) ([]*Feed, error)
	}

	Users interface {
		Create(context.Context, *User) error
		GetById(context.Context, int64) (*User, error)
		Follow(context.Context, int64, int64) error
		Unfollow(context.Context, int64, int64) error
	}

	Comments interface {
		GetByPostId(context.Context, int64) ([]*Comment, error)
		Create(context.Context, *Comment) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
