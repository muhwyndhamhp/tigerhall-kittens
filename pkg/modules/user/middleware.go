package user

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
)

var KeyUser = &ctxKey{"user"}

type ctxKey struct {
	name string
}

func AuthMiddleware(ur entities.UserRepository, tr entities.TokenHistoryRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			u, err := ExtractUserFromJWT(c.Request().Context(), ur, tr, authHeader)
			if err != nil {
				log.Error(err)
				return next(c)
			}

			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), KeyUser, u)))

			return next(c)
		}
	}
}

func ExtractUserFromJWT(ctx context.Context, ur entities.UserRepository, tr entities.TokenHistoryRepository, authHeader string) (*entities.User, error) {
	if authHeader == "" {
		return nil, entities.ErrUserByCtxNotFound
	}

	t, _ := tr.FindByToken(ctx, authHeader)
	if t != nil {
		return nil, entities.ErrTokenAlreadyInvalidated
	}

	tu, err := entities.ParseToken(authHeader)
	if err != nil {
		return nil, err
	}

	u, err := ur.FindByID(ctx, tu.ID)
	if err != nil || u == nil {
		return nil, entities.ErrUserByCtxNotFound
	}

	return u, nil
}

func UserByCtx(ctx context.Context) (*entities.User, error) {
	i := ctx.Value(KeyUser)
	if i == nil {
		return nil, entities.ErrUserByCtxNotFound
	}

	v := i.(*entities.User)
	return v, nil
}
