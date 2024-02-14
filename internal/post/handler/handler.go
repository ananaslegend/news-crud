package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ananaslegend/news-crud/internal/contexts"
	"github.com/ananaslegend/news-crud/internal/post/model"
	"github.com/ananaslegend/news-crud/internal/post/service"
	"github.com/ananaslegend/news-crud/pkg/logs"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type CreatePostService interface {
	CreatePost(ctx context.Context, title, content string, authorID int) (int, error)
}

type GetPostByFilterService interface {
	GetPostByFilter(ctx context.Context, filter model.Filter) ([]model.Post, error)
}

type GetPostByIDService interface {
	GetPostByID(ctx context.Context, id int) (model.Post, error)
}

type UpdatePostService interface {
	UpdatePost(ctx context.Context, userID int, post model.Post) error
}

type DeletePostService interface {
	DeletePost(ctx context.Context, userID, postID int) error
}

type PostHandler struct {
	logger *slog.Logger

	createPostService      CreatePostService
	getPostByFilterService GetPostByFilterService
	getPostByIDService     GetPostByIDService
	updatePostService      UpdatePostService
	deletePostService      DeletePostService
}

func NewPostHandler(
	logger *slog.Logger,
	createPostService CreatePostService,
	getPostByFilterService GetPostByFilterService,
	getPostByIDService GetPostByIDService,
	updatePostService UpdatePostService,
	deletePostService DeletePostService) *PostHandler {
	return &PostHandler{
		logger:                 logger,
		createPostService:      createPostService,
		getPostByFilterService: getPostByFilterService,
		getPostByIDService:     getPostByIDService,
		updatePostService:      updatePostService,
		deletePostService:      deletePostService,
	}
}

type CreatePostRequest struct { // todo
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (p PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	const op = "news-crud.internal.post.create.handler.HandleHTTP"
	logger := p.logger.With(slog.String("op", op))

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("cant decode request", logs.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := validator.New().Struct(req); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := contexts.MustGetUserID(r.Context())

	postID, err := p.createPostService.CreatePost(r.Context(), req.Title, req.Content, userID)
	if err != nil {
		logger.Error("cant create post", logs.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]int{"post_id": postID}); err != nil {
		logger.Error("cant encode response", logs.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func (p PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	const op = "news-crud.internal.post.handler.get_by_id.HandleHTTP"
	logger := p.logger.With(slog.String("op", op))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	post, err := p.getPostByIDService.GetPostByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoPostWasFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			logger.Error("cant get post by id", logs.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	jsonPost, err := json.Marshal(post)
	if err != nil {
		logger.Error("cant marshal post", logs.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonPost)

	return
}

func (p PostHandler) GetPostByFilter(w http.ResponseWriter, r *http.Request) {
	const op = "news-crud.internal.post.get_by_id.handler.HandleHTTP"
	logger := p.logger.With(slog.String("op", op))

	filter := model.Filter{}

	filter.DateFrom, _ = time.Parse(time.RFC3339, r.URL.Query().Get("dateFrom"))
	filter.DateTo, _ = time.Parse(time.RFC3339, r.URL.Query().Get("dateTo"))

	if err := filter.Validation(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	posts, err := p.getPostByFilterService.GetPostByFilter(r.Context(), filter)
	if err != nil {
		if errors.Is(err, service.ErrNoPostWasFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		logger.Error("cant get post by filter", logs.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonPost, err := json.Marshal(posts)
	if err != nil {
		logger.Error("cant marshal post", logs.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonPost)
	w.WriteHeader(http.StatusOK)
}

func (p PostHandler) UpdatePostByID(w http.ResponseWriter, r *http.Request) {
	const op = "news-crud.internal.post.handler.update.HandleHTTP"
	logger := p.logger.With(slog.String("op", op))

	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		logger.Error("cant decode request", logs.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := validator.New().Struct(post); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := contexts.MustGetUserID(r.Context())

	if err := p.updatePostService.UpdatePost(r.Context(), userID, post); err != nil {
		switch {
		case errors.Is(err, service.ErrNoPostWasFound):
			w.WriteHeader(http.StatusNotFound)
			return
		case errors.Is(err, service.ErrUserHasNoPermission):
			w.WriteHeader(http.StatusForbidden)
			return
		default:
			logger.Error("cant update post", logs.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

func (p PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	const op = "news-crud.internal.post.delete.handler.DeletePost"
	logger := p.logger.With(slog.String("op", op))

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || postID < 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	userID := contexts.MustGetUserID(r.Context())

	if err = p.deletePostService.DeletePost(r.Context(), userID, postID); err != nil {
		switch {
		case errors.Is(err, service.ErrNoPostWasFound):
			w.WriteHeader(http.StatusNotFound)
			return
		case errors.Is(err, service.ErrUserHasNoPermission):
			w.WriteHeader(http.StatusForbidden)
			return
		default:
			logger.Error(op, logs.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusOK)
}
