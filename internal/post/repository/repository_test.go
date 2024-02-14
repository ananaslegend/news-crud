package repository

import (
	"context"
	"database/sql"
	"github.com/ananaslegend/news-crud/internal/post/model"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func TestNewsStorage(t *testing.T) {
	ctx := context.Background()

	// 1. Start the postgres container and run any migrations on it
	container, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Run any migrations on the database
	_, _, err = container.Exec(ctx, []string{"psql", "-U", "test", "-d", "test", "-c", `create table if not exists posts (
  id serial primary key,
  title text not null,
  content text not null,
  author_id integer not null,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now()
);`})
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	dbURL, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	conn, err := sql.Open("postgres", dbURL+"sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	repo := NewPostRepository(conn)

	t.Run("Test inserting post", func(t *testing.T) {
		t.Cleanup(func() {
			_, err := conn.Exec("delete from posts where true")
			if err != nil {
				t.Fatal(err)
			}
		})

		postToInsert := model.Post{
			Title:     "test",
			Content:   "test",
			AuthorID:  1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		postID, err := repo.CreatePost(ctx, postToInsert)

		require.NoError(t, err)

		postInserted, err := repo.GetPostByID(ctx, postID)
		require.NoError(t, err)

		require.Equal(t, postToInsert.Title, postInserted.Title)
		require.Equal(t, postToInsert.Content, postInserted.Content)
		require.Equal(t, postToInsert.AuthorID, postInserted.AuthorID)
	})

	t.Run("Test Get post (fail case)", func(t *testing.T) {
		invalidPostID := 228

		_, err := repo.GetPostByID(ctx, invalidPostID)

		require.ErrorIs(t, err, ErrNoPostWasFound)
	})

}
