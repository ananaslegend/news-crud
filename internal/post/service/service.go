package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ananaslegend/news-crud/internal/post/model"
	"github.com/ananaslegend/news-crud/internal/post/repository"
	"github.com/ananaslegend/news-crud/pkg/logs"
	"log/slog"
	"time"
)

type CreatePostRepository interface {
	CreatePost(ctx context.Context, post model.Post) (int, error)
}

type DeletePostRepository interface {
	DeletePost(ctx context.Context, id int) error
}

type GetPostByIDRepository interface {
	GetPostByID(ctx context.Context, id int) (model.Post, error)
}

type GetPostByFilterRepository interface {
	GetPostByFilter(ctx context.Context, filter model.Filter) ([]model.Post, error)
}

type UpdatePostRepository interface {
	UpdatePost(ctx context.Context, post model.Post) error
}

type UserPostUpdatePermissionService interface {
	UserCanUpdatePost(ctx context.Context, userID, postID int) bool
}

type UserPostDeletePermissionService interface {
	UserCanDeletePost(ctx context.Context, userID, postID int) bool
}

type PostService struct {
	logger *slog.Logger

	createPostRepository   CreatePostRepository
	postByIDRepository     GetPostByIDRepository
	postByFilterRepository GetPostByFilterRepository
	updatePostRepository   UpdatePostRepository
	deletePostRepository   DeletePostRepository

	updatePermissionService UserPostUpdatePermissionService
	deletePermissionService UserPostDeletePermissionService
}

func NewPostService(
	logger *slog.Logger,
	createPostRepository CreatePostRepository,
	postByIDRepository GetPostByIDRepository,
	postByFilterRepository GetPostByFilterRepository,
	updatePostRepository UpdatePostRepository,
	deletePostRepository DeletePostRepository,
	updatePermissionService UserPostUpdatePermissionService,
	deletePermissionService UserPostDeletePermissionService,
) *PostService {
	return &PostService{
		logger:                  logger,
		createPostRepository:    createPostRepository,
		postByIDRepository:      postByIDRepository,
		postByFilterRepository:  postByFilterRepository,
		updatePostRepository:    updatePostRepository,
		deletePostRepository:    deletePostRepository,
		updatePermissionService: updatePermissionService,
		deletePermissionService: deletePermissionService,
	}
}

func (ps PostService) CreatePost(ctx context.Context, title, content string, authorID int) (int, error) {
	const op = "news-crud.internal.post.create.service.CreatePost"
	logger := ps.logger.With(slog.String("op", op))

	post := model.NewPost(title, content, authorID)

	postID, err := ps.createPostRepository.CreatePost(ctx, post)
	if err != nil {
		logger.Error("cant create post", logs.Err(err))
		return 0, ErrCantCretePost
	}

	return postID, nil
}

func (ps PostService) GetPostByID(ctx context.Context, id int) (model.Post, error) {
	const op = "news-crud.internal.post.get_by_id.service.GetPostByID"

	post, err := ps.postByIDRepository.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoPostWasFound) {
			return model.Post{}, ErrNoPostWasFound
		}

		return model.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	return post, nil
}

func (ps PostService) GetPostByFilter(ctx context.Context, filter model.Filter) ([]model.Post, error) {
	const op = "news-crud.internal.post.get_by_filter.service.GetPostByFilter"

	posts, err := ps.postByFilterRepository.GetPostByFilter(ctx, filter)
	if err != nil {
		if errors.Is(err, repository.ErrNoPostWasFound) {
			return []model.Post{}, ErrNoPostWasFound
		}

		return []model.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nil
}

func (ps PostService) UpdatePost(ctx context.Context, userID int, post model.Post) error {
	const op = "news-crud.internal.post.update.service.UpdatePost"

	if ok := ps.updatePermissionService.UserCanUpdatePost(ctx, userID, post.ID); !ok {
		return ErrUserHasNoPermission
	}

	post.UpdatedAt = time.Now()

	if err := ps.updatePostRepository.UpdatePost(ctx, post); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (ps PostService) DeletePost(ctx context.Context, userID, postID int) error {
	const op = "news-crud.internal.post.delete.service.DeletePost"

	if ok := ps.deletePermissionService.UserCanDeletePost(ctx, userID, postID); !ok {
		return ErrUserHasNoPermission
	}

	if err := ps.deletePostRepository.DeletePost(ctx, postID); err != nil {
		if errors.Is(err, repository.ErrNoPostWasFound) {
			return ErrNoPostWasFound
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
