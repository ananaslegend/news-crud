package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ananaslegend/news-crud/internal/post/model"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (pr PostRepository) CreatePost(ctx context.Context, post model.Post) (int, error) {
	const op = "news-crud.internal.post.create.repository.CreatePost"

	res := pr.db.QueryRowContext(ctx, `
		insert into 
		    posts (title, content, created_at, updated_at, author_id)
		values ($1, $2, $3, $4, $5)
		returning id
`, post.Title, post.Content, post.CreatedAt, post.UpdatedAt, post.AuthorID,
	)

	var postID int

	err := res.Scan(&postID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return postID, nil
}

func (pr PostRepository) GetPostByID(ctx context.Context, id int) (model.Post, error) {
	const op = "news-crud.internal.post.get_by_id.repository.GetPostByID"

	var post model.Post
	err := pr.db.QueryRowContext(ctx, `
		select id, title, content, created_at, updated_at, author_id
		from posts
		where id = $1
	`, id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.AuthorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Post{}, ErrNoPostWasFound
		}

		return model.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	return post, nil
}

func (pr PostRepository) GetPostByFilter(ctx context.Context, filter model.Filter) ([]model.Post, error) {
	const op = "news-crud.internal.post.get_by_filter.repository.GetPostByFilter"

	rows, err := pr.db.QueryContext(ctx, `
		select id, title, content, created_at, updated_at 
		from posts
		where created_at >= $1 and created_at <= $2
	`, filter.DateFrom, filter.DateTo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	posts := make([]model.Post, 0)

	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if len(posts) == 0 {
		return nil, ErrNoPostWasFound
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nil
}

func (pr PostRepository) UpdatePost(ctx context.Context, post model.Post) error {
	const op = "news-crud.internal.post.update.repository.postgre.UpdatePost"

	_, err := pr.db.ExecContext(ctx, `
		update posts
		set title = $1, content = $2, updated_at = $3
		where id = $4
	`, post.Title, post.Content, post.UpdatedAt, post.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoPostWasFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pr PostRepository) DeletePost(ctx context.Context, id int) error {
	const op = "news-crud.internal.post.delete.repository.DeletePost"

	if _, err := pr.db.ExecContext(ctx, `
		delete from posts
		where id = $1
	`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoPostWasFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
