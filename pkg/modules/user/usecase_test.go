package user

import (
	"context"
	"testing"

	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUsecase_CreateUser(t *testing.T) {
	token := GenerateJWT(&entities.User{
		Model: gorm.Model{
			ID: 0,
		},
		Name:         "user-1",
		Email:        "email-1@example.com",
		PasswordHash: "inipasswordnya!",
	})
	testCases := []struct {
		name string

		usr *model.NewUser

		findByEmailResp *entities.User
		findByEmailErr  error

		createErr error
		want      string
		wantErr   error
	}{
		{
			name: "should return token and nil error",
			usr: &model.NewUser{
				Name:     "user-1",
				Email:    "email-1@example.com",
				Password: "inipasswordnya!",
			},
			findByEmailErr: gorm.ErrRecordNotFound,
			createErr:      nil,
			want:           token,
			wantErr:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur := mocks.NewUserRepository(t)
			tr := mocks.NewTokenHistoryRepository(t)

			uc := NewUserUsecase(ur, tr)

			ur.
				On("FindByEmail", mock.Anything, tc.usr.Email).
				Return(tc.findByEmailResp, tc.findByEmailErr).
				Once()

			ur.
				On("Create", mock.Anything, mock.Anything).
				Return(tc.createErr).
				Maybe()

			token, err := uc.CreateUser(context.Background(), tc.usr)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, token)
		})
	}
}

func TestUsecase_Login(t *testing.T) {
	token := GenerateJWT(nil)
	pwHash := "$2a$10$MGPcG.T8.KzfqkwgPq9TDuiOGLi45guJQ8PQSM.yXMrjeoRs.Wi2C"
	testCases := []struct {
		name string

		email    string
		password string

		findByEmailErr error
		validateErr    error
		want           string
		wantErr        error
	}{
		{
			name:           "should return token and nil error",
			email:          "email-1@example.com",
			password:       "inipasswordnya!",
			findByEmailErr: nil,
			validateErr:    nil,
			want:           token,
			wantErr:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur := mocks.NewUserRepository(t)
			tr := mocks.NewTokenHistoryRepository(t)

			uc := NewUserUsecase(ur, tr)

			ur.
				On("FindByEmail", mock.Anything, tc.email).
				Return(&entities.User{
					Model: gorm.Model{
						ID: 1,
					},
					Name:         "user-1",
					Email:        "email-1@example.com",
					PasswordHash: pwHash,
				}, tc.findByEmailErr).
				Once()

			token, err := uc.Login(context.Background(), tc.email, tc.password)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, token)
		})
	}
}

func TestUsecase_GetUserByID(t *testing.T) {
	testCases := []struct {
		name string

		id uint

		findByIDErr error
		want        *model.User
		wantErr     error
	}{
		{
			name:        "should return user and nil error",
			id:          1,
			findByIDErr: nil,
			want: &model.User{
				ID:    1,
				Name:  "user-1",
				Email: "email-1@example.com",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur := mocks.NewUserRepository(t)
			tr := mocks.NewTokenHistoryRepository(t)

			uc := NewUserUsecase(ur, tr)

			ur.
				On("FindByID", mock.Anything, tc.id).
				Return(&entities.User{
					Model: gorm.Model{
						ID: 1,
					},
					Name:  "user-1",
					Email: "email-1@example.com",
				}, tc.findByIDErr).
				Once()

			user, err := uc.GetUserByID(context.Background(), tc.id)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, user)
		})
	}
}

func TestUsecase_RefreshToken(t *testing.T) {
	token := GenerateJWT(nil)
	testCases := []struct {
		name string

		token string

		parseTokenErr error
		findByIDErr   error
		want          string
		wantErr       error
	}{
		{
			name:          "should return token and nil error",
			token:         token,
			parseTokenErr: nil,
			findByIDErr:   nil,
			want:          token,
			wantErr:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ur := mocks.NewUserRepository(t)
			tr := mocks.NewTokenHistoryRepository(t)

			uc := NewUserUsecase(ur, tr)

			tr.
				On("FindByToken", mock.Anything, tc.token).
				Return(nil, nil).
				Once()

			ur.
				On("FindByID", mock.Anything, uint(1)).
				Return(&entities.User{
					Model: gorm.Model{
						ID: 1,
					},
					Name:  "user-1",
					Email: "email-1@example.com",
				}, tc.findByIDErr).
				Maybe()

			token, err := uc.RefreshToken(context.Background(), tc.token)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, token)
		})
	}
}
