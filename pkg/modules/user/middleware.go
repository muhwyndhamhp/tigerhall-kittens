package user

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
)

var keyUser = &ctxKey{"user"}

var ErrUserByCtxNotFound = errors.New("ErrUserByCtxNotFound: user not found in context")

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

func UserByCtx(ctx context.Context) (*entities.User, error) {
	i := ctx.Value(keyUser)
	if i == nil {
		return nil, ErrUserByCtxNotFound
	}

	v := i.(*entities.User)
	if v == nil {
		return nil, ErrUserByCtxNotFound
	}
	return v, nil
}
