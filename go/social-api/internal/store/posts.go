package store

import (
	"context"
	"database/sql"
	"social-api/internal/models"

	"github.com/lib/pq"
)

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *models.Post) error {
	query := `INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	return s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserId, pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}
