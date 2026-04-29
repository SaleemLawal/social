package store

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user"`
	Likes     int       `json:"likes"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetByPostId(ctx context.Context, postId int64) ([]*Comment, error) {
	query := `
		SELECT c.id, c.content, c.post_id, c.user_id, c.created_at, u.username, u.id, u.email FROM comments c
		JOIN users u ON c.user_id = u.id
		where c.post_id = $1
		ORDER BY c.created_at DESC;
	`
	var comments []*Comment
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Comment
		c.User = User{}
		if err := rows.Scan(
			&c.ID,
			&c.Content,
			&c.PostID,
			&c.UserID,
			&c.CreatedAt,
			&c.User.Username,
			&c.User.ID,
			&c.User.Email,
		); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		WITH inserted AS (
			INSERT INTO comments (content, post_id, user_id, likes)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, user_id
		)
		SELECT i.id, i.created_at, u.id, u.username, u.email, u.created_at
		FROM inserted i
		JOIN users u ON u.id = i.user_id
	`
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	if err := s.db.QueryRowContext(ctx, query, comment.Content, comment.PostID, comment.UserID, comment.Likes).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UserID,
		&comment.User.Username,
		&comment.User.Email,
		&comment.User.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}
