package user

import (
	"context"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities/mocks"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMiddleware_ExtractUserFromJWT(t *testing.T) {
	testCase := []struct {
		name        string
		authHeader  string
		mockRepo    *entities.User
		mockRepoErr error
		expected    *entities.User
		expectedErr error
	}{
		{
			name:       "success extract user from jwt",
			authHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImthemVhbS5wbHVzQGdtYWlsLmNvbSIsImV4cCI6MTcxMzQwNzkzMiwiaWQiOjEsInVzZXJuYW1lIjoibXVod3luZGhhbSJ9.gYOSnOkJhPU-nZXP8g3fwD72kmCNA65RB_TJkXWwzPA",
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
			authHeader:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImthemVhbS5wbHVzQGdtYWlsLmNvbSIsImV4cCI6MTcxMzQwNzkzMiwiaWQiOjEsInVzZXJuYW1lIjoibXVod3luZGhhbSJ9.gYOSnOkJhPU-nZXP8g3fwD72kmCNA65RB_TJkXWwzPA",
			mockRepo:    nil,
			mockRepoErr: entities.ErrUserByCtxNotFound,
			expected:    nil,
			expectedErr: entities.ErrUserByCtxNotFound,
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
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)

			repo.
				On("FindByID", context.Background(), uint(1)).
				Return(tc.mockRepo, tc.mockRepoErr).
				Maybe()

			u, err := ExtractUserFromJWT(context.Background(), repo, tc.authHeader)
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
			ctx: context.WithValue(context.Background(), keyUser, &entities.User{
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
