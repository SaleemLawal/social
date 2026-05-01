package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Comments  []*Comment `json:"comments"`
	Version   int        `json:"version"`
}

type Feed struct {
	Post
	CommentCount int    `json:"comment_count"`
	Username     string `json:"username"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	if err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags)).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetById(ctx context.Context, postId int64) (*Post, error) {
	query := `
		SELECT id, content, title, user_id, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	var post Post
	if err := s.db.QueryRowContext(ctx, query, postId).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET content = $1, title = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	if err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.ID, post.Version).Scan(
		&post.Version,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *PostStore) Delete(ctx context.Context, postId int64) error {
	query := `
		DELETE FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (s *PostStore) GetFeeds(ctx context.Context, userID int64, fq *PaginationFeedsQuery) ([]*Feed, error) {
	sinceClause := ""
	if fq.Since != "" {
		sinceClause = "AND (p.created_at >= $6)"
	}
	untilClause := ""
	if fq.Until != "" {
		untilClause = "AND (p.created_at <= $7)"
	}

	query := fmt.Sprintf(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username, COUNT(c.id) AS comment_count,
			COALESCE(
				json_agg(
					json_build_object(
						'id', c.id,
						'content', c.content,
						'post_id', c.post_id,
						'user_id', c.user_id,
						'created_at', c.created_at,
						'likes', c.likes,
						'user', json_build_object(
							'id', cu.id,
							'username', cu.username,
							'email', cu.email
						)
					)
				) FILTER (WHERE c.id IS NOT NULL),
				'[]'
			) AS comments
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN users cu ON c.user_id = cu.id
		LEFT JOIN users u ON p.user_id = u.id
		WHERE (p.user_id = $1
		OR p.user_id IN (
			SELECT followers.followed_id FROM followers WHERE follower_id = $1
		))
		AND (p.title ILIKE '%%' || $4 || '%%' OR p.content ILIKE '%%' || $4 || '%%')
		AND (p.tags @> $5 OR $5 = '{}')
		%s
		%s
		GROUP BY p.id, u.username
		ORDER BY p.created_at %s
		LIMIT $2 OFFSET $3
	`, sinceClause, untilClause, fq.Sort)
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	args := []any{userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags)}
	if fq.Since != "" {
		args = append(args, fq.Since)
	}
	if fq.Until != "" {
		args = append(args, fq.Until)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []*Feed
	for rows.Next() {
		var feed Feed
		var commentsJSON []byte
		if err := rows.Scan(
			&feed.ID,
			&feed.UserID,
			&feed.Title,
			&feed.Content,
			&feed.CreatedAt,
			&feed.Version,
			pq.Array(&feed.Tags),
			&feed.Username,
			&feed.CommentCount,
			&commentsJSON,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(commentsJSON, &feed.Comments); err != nil {
			return nil, err
		}
		feeds = append(feeds, &feed)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}
