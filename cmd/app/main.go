package main

import (
	"database/sql"
	"github.com/ananaslegend/news-crud/internal/config"
	"github.com/ananaslegend/news-crud/internal/middleware"
	permissionRepository "github.com/ananaslegend/news-crud/internal/permission/repository"
	permissionService "github.com/ananaslegend/news-crud/internal/permission/service"
	postHandler "github.com/ananaslegend/news-crud/internal/post/handler"
	postRepository "github.com/ananaslegend/news-crud/internal/post/repository"
	postService "github.com/ananaslegend/news-crud/internal/post/service"
	"github.com/ananaslegend/news-crud/pkg/logs"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("cant get config", logs.Err(err))
	}

	logger := logs.SetUpLogger(*cfg)

	db, err := sql.Open("postgres", cfg.DBConn)
	if err != nil {
		logger.Error("cant connect to database", logs.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	// todo graceful shutdown

	permissionRepo := permissionRepository.NewPermissionRepository(db)
	permissionSrv := permissionService.NewPermissionService(permissionRepo)

	postRepo := postRepository.NewPostRepository(db)
	postSrv := postService.NewPostService(
		logger,
		postRepo,
		postRepo,
		postRepo,
		postRepo,
		postRepo,
		permissionSrv,
		permissionSrv,
	)
	postHdl := postHandler.NewPostHandler(
		logger,
		postSrv,
		postSrv,
		postSrv,
		postSrv,
		postSrv,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /posts", middleware.Auth(cfg.Secret, postHdl.CreatePost))
	mux.HandleFunc("GET /posts", postHdl.GetPostByFilter)
	mux.HandleFunc("GET /posts/{id}", postHdl.GetPostByID)
	mux.HandleFunc("PUT /posts/{id}", middleware.Auth(cfg.Secret, postHdl.UpdatePostByID))
	mux.HandleFunc("DELETE /posts/{id}", middleware.Auth(cfg.Secret, postHdl.DeletePost))

	s := http.Server{
		Addr:    cfg.HttpPort,
		Handler: mux, // todo recover middleware
	}

	s.ListenAndServe()
}
