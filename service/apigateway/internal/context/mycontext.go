package mycontext

import (
	"context"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

const ApiKeyContextKey contextKey = "apiKey"

const TokenContextKey contextKey = "token"

func ApiKeyFromContext(ctx context.Context) (string, bool) {
	apiKey, ok := ctx.Value(ApiKeyContextKey).(string)
	return apiKey, ok
}

func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, TokenContextKey, token)
}

func TokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(TokenContextKey).(string)
	return token, ok
}

func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}

func UserForContext(ctx context.Context) (int, bool) {
	uid := ctx.Value(UserIDContextKey)
	if uid == nil {
		return 0, false
	}

	id, ok := uid.(int)
	return id, ok
}
