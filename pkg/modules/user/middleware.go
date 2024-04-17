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
			u, err := ExtractUserFromJWT(c.Request().Context(), repo, authHeader)
			if err != nil {
				return next(c)
			}

			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), keyUser, u)))

			return next(c)
		}
	}
}

func ExtractUserFromJWT(ctx context.Context, repo entities.UserRepository, authHeader string) (*entities.User, error) {
	if authHeader == "" {
		return nil, entities.ErrUserByCtxNotFound
	}

	tu, err := entities.ParseToken(authHeader)
	if err != nil {
		return nil, err
	}

	u, err := repo.FindByID(ctx, tu.ID)
	if err != nil || u == nil {
		return nil, entities.ErrUserByCtxNotFound
	}

	return u, nil
}

func UserByCtx(ctx context.Context) (*entities.User, error) {
	i := ctx.Value(keyUser)
	if i == nil {
		return nil, entities.ErrUserByCtxNotFound
	}

	v := i.(*entities.User)
	return v, nil
}
