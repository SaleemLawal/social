package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound      = errors.New("Resource not found")
	QUERY_TIMEOUT_DURATION = 5 * time.Second
	ErrConflict            = errors.New("Resource already exists")
	ErrDuplicateUsername   = errors.New("Username already exists")
	ErrDuplicateEmail      = errors.New("Email already exists")
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
		Create(context.Context, *sql.Tx, *User) error
		GetById(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		Follow(context.Context, int64, int64) error
		Unfollow(context.Context, int64, int64) error
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
	}

	Comments interface {
		GetByPostId(context.Context, int64) ([]*Comment, error)
		Create(context.Context, *Comment) error
	}

	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
		Roles:    &RoleStore{db},
	}
}

// withTx is a helper function to execute a function within a transaction
func withTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}
