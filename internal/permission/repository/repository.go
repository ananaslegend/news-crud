package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	DefaultAuthorID = 0
)

type PermissionRepository struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (pr PermissionRepository) GetAuthorPostByID(ctx context.Context, id int) (int, error) {
	const op = "news-crud.internal.post.get_by_id.repository.GetPostByID"

	stmt, err := pr.db.Prepare(`
		select author_id 
		from post
		where id == ?
	`)
	if err != nil {
		return DefaultAuthorID, err
	}

	var authorID int
	if err = stmt.QueryRowContext(ctx, id).Scan(&authorID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DefaultAuthorID, ErrNoPostWasFound
		}
		return DefaultAuthorID, fmt.Errorf("%s: %w", op, err)
	}

	return authorID, nil
}
