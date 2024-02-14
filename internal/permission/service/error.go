package service

import "errors"

var (
	ErrUserHasNoPermission = errors.New("user has no permission")
	ErrNoPostWasFound      = errors.New("no post was found")
)
