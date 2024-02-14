package model

import "time"

type Post struct {
	ID        int
	Title     string
	Content   string
	AuthorID  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPost(title, content string, authorID int) Post {
	return Post{
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
