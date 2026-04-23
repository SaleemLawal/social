package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts
	Users
}

type Posts interface {
	Create(context.Context) error
}

type Users interface {
	Create(context.Context) error
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		
	}
}
