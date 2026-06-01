package ctxkeys

import (
	"context"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
)

func WithUser(ctx context.Context, user repo.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func UserFrom(ctx context.Context) (repo.User, bool) {
	user, ok := ctx.Value(UserKey).(repo.User)
	return user, ok
}
