package user

import (
	"context"
	"net/http"

	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
)

var keyUser = &ctxKey{"user"}

type ctxKey struct {
	name string
}

func AuthMiddleware(repo entities.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				next.ServeHTTP(w, r)
			}

			tu, err := entities.ParseToken(authHeader)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			u, err := repo.FindByID(tu.ID)
			if err != nil || u == nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), keyUser, u)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func UserByCtx(ctx context.Context) *entities.User {
	return ctx.Value(keyUser).(*entities.User)
}
