package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Follower struct {
	UserId     int64     `json:"user_id"`
	FollowerId int64     `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	var query string = `
		INSERT INTO users (username, password, email) VALUES($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	if err := s.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&user.ID, &user.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	var query string = `
		SELECT id, username, email, created_at FROM users WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	var user = &User{}
	if err := s.db.QueryRowContext(ctx, query, userId).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UserStore) Follow(ctx context.Context, followerId, userId int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id) VALUES($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, followerId)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		case errors.As(err, &pqErr) && pqErr.Code == "23505":
			return ErrConflict
		default:
			return err
		}
	}
	return nil
}

func (s *UserStore) Unfollow(ctx context.Context, followerId, userId int64) error {
	query := `
		DELETE FROM followers WHERE user_id = $1 AND follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, followerId)

	return err
}
