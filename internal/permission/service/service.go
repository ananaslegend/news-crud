package service

import (
	"context"
	"errors"
	"github.com/ananaslegend/news-crud/internal/permission/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type GetAuthorPostByID interface {
	GetAuthorPostByID(ctx context.Context, id int) (int, error)
}

type PermissionService struct {
	postByID GetAuthorPostByID
}

func NewPermissionService(getPostByID GetAuthorPostByID) *PermissionService {
	return &PermissionService{
		postByID: getPostByID,
	}
}

func (s PermissionService) UserCanUpdatePost(ctx context.Context, userID, postID int) bool {
	authorID, err := s.postByID.GetAuthorPostByID(ctx, postID)
	if err != nil {
		if errors.Is(err, repository.ErrNoPostWasFound) {
			return false
		}

		return false
	}

	if authorID != userID {
		return false
	}

	return true
}

func (s PermissionService) UserCanDeletePost(ctx context.Context, userID, postID int) bool {
	authorID, err := s.postByID.GetAuthorPostByID(ctx, postID)
	if err != nil {
		if errors.Is(err, repository.ErrNoPostWasFound) {
			return false
		}

		return false
	}

	if authorID != userID {
		return false
	}

	return true
}
