package service

import "errors"

var (
	ErrCantCretePost       = errors.New("post can not be created")
	ErrNoPostWasFound      = errors.New("post not found")
	ErrUserHasNoPermission = errors.New("user has no permission")
	ErrUserIsNotAuthor     = errors.New("user is not author")
)
