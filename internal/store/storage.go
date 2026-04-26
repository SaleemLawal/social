package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("resource not found")
)

type Storage struct {
	Posts
	Users
}

type Posts interface {
	Create(context.Context, *Post) error
	GetById(context.Context, int64) (*Post, error)
}

type Users interface {
	Create(context.Context, *User) error
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db},
		Users: &UsersStore{db},
	}
}
