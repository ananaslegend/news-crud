package contexts

import (
	"context"
)

type userIDKey struct{}

func SetUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func MustGetUserID(ctx context.Context) int {
	return ctx.Value(userIDKey{}).(int)
}
