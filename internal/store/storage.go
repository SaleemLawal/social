package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("resource not found")
	QUERY_TIMEOUT_DURATION = 5 * time.Second
)

type Storage struct {
	Posts interface{
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
	}
	Users interface{
		Create(context.Context, *User) error
	}
	Comments interface{
		GetByPostId(context.Context, int64) ([]*Comment, error)
		Create(context.Context, *Comment) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db},
		Users: &UserStore{db},
		Comments: &CommentStore{db},
	}
}
