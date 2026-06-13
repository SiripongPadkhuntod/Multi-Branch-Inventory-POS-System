package domain

import "context"

type contextKey string

const currentUserKey contextKey = "current_user"

func ContextWithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, currentUserKey, user)
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(currentUserKey).(User)
	return user, ok
}
