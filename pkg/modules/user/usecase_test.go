package user

import (
	"context"
	"testing"
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities/mocks"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/timex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUsecase_CreateUser(t *testing.T) {
	now := time.Unix(1713345275, 0) // 17th March 2024 16:14:35
	testCases := []struct {
		name string

		usr *model.NewUser

		createErr error
		want      string
		wantErr   error
	}{
		{
			name: "should return token and nil error",
			usr: &model.NewUser{
				Name:     "user-1",
				Email:    "mail-1@example.com",
				Password: "inipasswordnya!",
			},
			createErr: nil,
			want:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1haWwtMUBleGFtcGxlLmNvbSIsImV4cCI6MTcxMzQzMTY3NSwiaWQiOjAsInVzZXJuYW1lIjoidXNlci0xIn0.Hjcf56Yd3_JDtEI9Ov1VAq-ujwK4s3CIaXiSNWHkPME",
			wantErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)

			uc := NewUserUsecase(repo)

			timex.SetTestTime(now)
			repo.
				On("Create", mock.Anything, mock.Anything).
				Return(tc.createErr).
				Once()

			token, err := uc.CreateUser(context.Background(), tc.usr)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, token)
		})
	}
}

func TestUsecase_Login(t *testing.T) {
	now := time.Unix(1713345275, 0) // 17th March 2024 16:14:35
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
			email:          "mail-1@example.com",
			password:       "inipasswordnya!",
			findByEmailErr: nil,
			validateErr:    nil,
			want:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1haWwtMUBleGFtcGxlLmNvbSIsImV4cCI6MTcxMzQzMTY3NSwiaWQiOjAsInVzZXJuYW1lIjoidXNlci0xIn0.Hjcf56Yd3_JDtEI9Ov1VAq-ujwK4s3CIaXiSNWHkPME",
			wantErr:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)

			uc := NewUserUsecase(repo)

			timex.SetTestTime(now)
			repo.
				On("FindByEmail", mock.Anything, tc.email).
				Return(&entities.User{
					Name:         "user-1",
					Email:        "mail-1@example.com",
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
				Email: "mail-1@example.com",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)

			uc := NewUserUsecase(repo)

			repo.
				On("FindByID", mock.Anything, tc.id).
				Return(&entities.User{
					Model: gorm.Model{
						ID: 1,
					},
					Name:  "user-1",
					Email: "mail-1@example.com",
				}, tc.findByIDErr).
				Once()

			user, err := uc.GetUserByID(context.Background(), tc.id)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, user)
		})
	}
}

func TestUsecase_RefreshToken(t *testing.T) {
	now := time.Unix(1713345275, 0) // 17th March 2024 16:14:35
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
			token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1haWwtMUBleGFtcGxlLmNvbSIsImV4cCI6MTcxMzQzMTY3NSwiaWQiOjAsInVzZXJuYW1lIjoidXNlci0xIn0.Hjcf56Yd3_JDtEI9Ov1VAq-ujwK4s3CIaXiSNWHkPME",
			parseTokenErr: nil,
			findByIDErr:   nil,
			want:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1haWwtMUBleGFtcGxlLmNvbSIsImV4cCI6MTcxMzQzMTY3NSwiaWQiOjAsInVzZXJuYW1lIjoidXNlci0xIn0.Hjcf56Yd3_JDtEI9Ov1VAq-ujwK4s3CIaXiSNWHkPME",
			wantErr:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)

			uc := NewUserUsecase(repo)

			timex.SetTestTime(now)
			repo.
				On("FindByID", mock.Anything, uint(0)).
				Return(&entities.User{
					Name:  "user-1",
					Email: "mail-1@example.com",
				}, tc.findByIDErr).
				Once()

			token, err := uc.RefreshToken(context.Background(), tc.token)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, token)
		})
	}
}
