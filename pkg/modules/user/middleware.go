package user

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
)

var keyUser = &ctxKey{"user"}

type ctxKey struct {
	name string
}

func AuthMiddleware(repo entities.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				return next(c)
			}

			tu, err := entities.ParseToken(authHeader)
			if err != nil {
				return next(c)
			}

			if tu == nil {
				return next(c)
			}

			u, err := repo.FindByID(tu.ID)
			if err != nil || u == nil {
				return next(c)
			}

			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), keyUser, u)))

			return next(c)
		}
	}
}

//
// func AuthMiddleware(repo entities.UserRepository) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			authHeader := r.Header.Get("Authorization")
//
// 			if authHeader == "" {
// 				next.ServeHTTP(w, r)
// 			}
//
// 			tu, err := entities.ParseToken(authHeader)
// 			if err != nil {
// 				next.ServeHTTP(w, r)
// 			}
//
// 			if tu == nil {
// 				next.ServeHTTP(w, r)
// 			}
//
// 			u, err := repo.FindByID(tu.ID)
// 			if err != nil || u == nil {
// 				next.ServeHTTP(w, r)
// 			}
//
// 			ctx := context.WithValue(r.Context(), keyUser, u)
// 			r = r.WithContext(ctx)
//
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

func UserByCtx(ctx context.Context) *entities.User {
	return ctx.Value(keyUser).(*entities.User)
}
