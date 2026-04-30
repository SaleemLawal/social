package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type password struct {
	text *string
	hash []byte
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
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

func (p *password) Set(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.hash = hash
	p.text = &password
	return nil
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	var query string = `
		INSERT INTO users (username, password, email) VALUES($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, user.Username, user.Password.hash, user.Email).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
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

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, expiresAt time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createUserInvitation(ctx, tx, token, user.ID, expiresAt); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, userId int64, expiresAt time.Duration) error {
	query := `INSERT INTO user_invitations (token, user_id, expires_at) VALUES($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userId, time.Now().Add(expiresAt))
	if err != nil {
		return err
	}

	return nil
}
