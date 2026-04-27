package store

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	ID int64 `json:"id"`
	Content string `json:"content"`
	PostID int64 `json:"post_id"`
	UserID int64 `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	User User `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetByPostId(ctx context.Context, postId int64) ([]*Comment, error) {
	query := `
		SELECT c.id, c.content, c.post_id, c.user_id, c.created_at, u.username, u.id FROM comments c
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
		); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}