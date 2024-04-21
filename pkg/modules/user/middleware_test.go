package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestMiddleware_AuthMiddleware(t *testing.T) {
	testCase := []struct {
		name        string
		authHeader  string
		mockRepo    *entities.User
		mockRepoErr error
		want        *entities.User
		wantErr     error
	}{
		{
			name:       "success get user from context",
			authHeader: GenerateJWT(nil),
			mockRepo: &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "user-1",
				Email: "email-1@example.com",
			},
			mockRepoErr: nil,
			want: &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "user-1",
				Email: "email-1@example.com",
			},
			wantErr: nil,
		},
		{
			name:       "failed get user from context",
			authHeader: "failed-token",
			mockRepo: &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "user-1",
				Email: "email-1@example.com",
			},
			mockRepoErr: nil,
			want:        nil,
			wantErr:     nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ur := mocks.NewUserRepository(t)
			tr := mocks.NewTokenHistoryRepository(t)

			ur.
				On("FindByID", mock.Anything, uint(1)).
				Return(tc.mockRepo, tc.mockRepoErr).
				Maybe()

			tr.
				On("FindByToken", mock.Anything, tc.authHeader).
				Return(nil, nil).
				Maybe()

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/query", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Request().Header.Add("Authorization", tc.authHeader)

			mw := AuthMiddleware(ur, tr)

			next := echo.HandlerFunc(func(c echo.Context) error {
				return nil
			})

			err := mw(next)(c)

			u, _ := UserByCtx(c.Request().Context())

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, u)
		})
	}
}

func TestMiddleware_ExtractUserFromJWT(t *testing.T) {
	token := GenerateJWT(nil)
	testCase := []struct {
		name        string
		authHeader  string
		mockToken   *entities.TokenHistory
		mockRepo    *entities.User
		mockRepoErr error
		expected    *entities.User
		expectedErr error
	}{
		{
			name:       "success extract user from jwt",
			authHeader: token,
			mockRepo: &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "user-1",
				Email: "mail-1@example.com",
			},
			mockRepoErr: nil,
			expected: &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "user-1",
				Email: "mail-1@example.com",
			},
			expectedErr: nil,
		},
		{
			name:        "failed extract user from jwt given user record not found",
			authHeader:  token,
			mockRepo:    nil,
			mockRepoErr: entities.ErrUserByCtxNotFound,
			expected:    nil,
			expectedErr: jwt.ValidationError{Inner: jwt.ErrSignatureInvalid},
		},
		{
			name:        "failed extract user from jwt given empty auth header",
			authHeader:  "",
			mockRepo:    nil,
			mockRepoErr: entities.ErrUserByCtxNotFound,
			expected:    nil,
			expectedErr: entities.ErrUserByCtxNotFound,
		},
		{
			name:        "failed extract user from jwt given invalid token",
			authHeader:  "invalid-token",
			mockRepo:    nil,
			mockRepoErr: entities.ErrUserByCtxNotFound,
			expected:    nil,
			expectedErr: jwt.ValidationError{Inner: jwt.ErrSignatureInvalid},
		},
		{
			name:       "failed extract user from jwt given token already invalidated",
			authHeader: token,
			mockToken: &entities.TokenHistory{
				Token:     token,
				RevokedAt: time.Now(),
			},
			mockRepo: &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "user-1",
				Email: "email-1@example.com",
			},
			mockRepoErr: nil,
			expected:    nil,
			expectedErr: entities.ErrTokenAlreadyInvalidated,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ur := mocks.NewUserRepository(t)
			tr := mocks.NewTokenHistoryRepository(t)

			ur.
				On("FindByID", mock.Anything, uint(1)).
				Return(tc.mockRepo, tc.mockRepoErr).
				Maybe()

			tr.
				On("FindByToken", mock.Anything, tc.authHeader).
				Return(tc.mockToken, nil).
				Maybe()

			u, err := ExtractUserFromJWT(context.Background(), ur, tr, tc.authHeader)
			assert.Equal(t, tc.expected, u)
			if tc.expectedErr == entities.ErrUserByCtxNotFound {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.Equal(t, tc.expectedErr != nil, err != nil)
			}
		})
	}
}

func TestMiddleware_UserByCtx(t *testing.T) {
	testCase := []struct {
		name        string
		ctx         context.Context
		expected    *entities.User
		expectedErr error
	}{
		{
			name: "success get user from context",
			ctx: context.WithValue(context.Background(), KeyUser, &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "user-1",
				Email: "mail-1@example.com",
			}),
			expected: &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "user-1",
				Email: "mail-1@example.com",
			},
			expectedErr: nil,
		},
		{
			name:        "failed get user from context",
			ctx:         context.Background(),
			expected:    nil,
			expectedErr: entities.ErrUserByCtxNotFound,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			u, err := UserByCtx(tc.ctx)
			assert.Equal(t, tc.expected, u)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
