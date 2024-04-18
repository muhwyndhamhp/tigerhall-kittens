package graph

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/user"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/email"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestMutation_RefreshToken(t *testing.T) {
	now := time.Now()
	token := GenerateJWT(nil)
	testCases := []struct {
		name string

		token string

		withRandomDBErr bool
		want            string
		wantErr         error
	}{
		{
			name:    "should return token and nil error",
			token:   token,
			want:    token,
			wantErr: nil,
		},
		{
			name:    "should return nil and error given user not found",
			token:   GenerateJWT(&entities.User{Model: gorm.Model{ID: 2}}),
			want:    "",
			wantErr: errs.RespError(entities.ErrUserNotFound),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _, _ := Setup(t, now, tc.withRandomDBErr)

			res, err := r.Mutation().RefreshToken(context.Background(), tc.token)

			assert.Equal(t, tc.want, res)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestMutation_Login(t *testing.T) {
	now := time.Now()
	token := GenerateJWT(nil)
	testCases := []struct {
		name string

		email    string
		password string

		withRandomDBErr bool
		want            string
		wantErr         error
	}{
		{
			name:     "should return token and nil error",
			email:    "email-1@example.com",
			password: "inipasswordnya!",
			want:     token,
			wantErr:  nil,
		},
		{
			name:     "should return nil and error given user not found",
			email:    "email-2@example.com",
			password: "inipasswordnya!",
			want:     "",
			wantErr:  errs.RespError(entities.ErrUserNotFound),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _, _ := Setup(t, now, tc.withRandomDBErr)

			res, err := r.Mutation().Login(context.Background(), tc.email, tc.password)

			fmt.Println(err)
			assert.Equal(t, tc.want, res)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestMutation_CreateUser(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name  string
		input model.NewUser

		withRandomDBErr bool
		want            string
		wantErr         error
	}{
		{
			name: "should return token and nil error",
			input: model.NewUser{
				Name:     "user-2",
				Email:    "email-2@example.com",
				Password: "inipasswordnya!",
			},

			withRandomDBErr: false,
			want: GenerateJWT(&entities.User{
				Model:        gorm.Model{ID: 2},
				Name:         "user-2",
				Email:        "email-2@example.com",
				PasswordHash: "inipasswordnya!",
			}),
			wantErr: nil,
		},
		{
			name: "should return nil and error given user already exists",
			input: model.NewUser{
				Name:     "user-2",
				Email:    "email-1@example.com",
				Password: "inipasswordnya!",
			},

			withRandomDBErr: false,
			want:            "",
			wantErr:         errs.RespError(entities.ErrUserAlreadyExists),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _, _ := Setup(t, now, tc.withRandomDBErr)

			res, err := r.Mutation().CreateUser(context.Background(), tc.input)

			assert.Equal(t, tc.want, res)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestMutation_CreateSighting(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name  string
		input model.NewSighting

		withRandomDBErr bool
		ctx             context.Context
		want            *model.Sighting
		wantErr         error
		wantQueue       *email.SightingEmail
	}{
		{
			name: "should return sighting with id 2, tiger id 1 and user id 1 and nil error",
			input: model.NewSighting{
				TigerID:   1,
				Date:      now,
				Latitude:  -7.250676,
				Longitude: 111.828316,
			},
			ctx: context.WithValue(context.Background(), user.KeyUser, &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:         "user-1",
				Email:        "email-1@example.com",
				PasswordHash: "hashed-password-1",
			}),
			want: &model.Sighting{
				ID:        2,
				UserID:    1,
				TigerID:   1,
				Date:      now,
				Latitude:  -7.250676,
				Longitude: 111.828316,
			},
			wantErr: nil,
			wantQueue: &email.SightingEmail{
				DestinationEmail:  "email-1@example.com",
				TigerName:         "tiger-1",
				SightingDate:      now.Format("2006-01-02 15:04:05"),
				SightingLatitude:  "-7.250676",
				SightingLongitude: "111.828316",
				ImageURL:          "",
			},
		},
		{
			name: "should return nil and error given user not found",
			input: model.NewSighting{
				TigerID:   1,
				Date:      now,
				Latitude:  -7.250676,
				Longitude: 111.828316,
			},
			ctx:     context.WithValue(context.Background(), user.KeyUser, nil),
			want:    nil,
			wantErr: errs.RespError(entities.ErrUserByCtxNotFound),
		},
		{
			name: "should return nil and error given random db error",
			input: model.NewSighting{
				TigerID:   1,
				Date:      now,
				Latitude:  -7.250676,
				Longitude: 111.828316,
			},
			ctx: context.WithValue(context.Background(), user.KeyUser, &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:         "user-1",
				Email:        "email-1@example.com",
				PasswordHash: "hashed-password-1",
			}),
			withRandomDBErr: true,
			want:            nil,
			wantErr:         errs.RespError(errors.New("no such table: tigers")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, mockS3, queue := Setup(t, now, tc.withRandomDBErr)

			mockS3.
				On(
					"UploadImage",
					mock.Anything,
					mock.Anything,
					"filename.jpeg",
					"image/jpeg",
					int64(1000),
				).
				Return("https://example.com/image.jpeg", nil).
				Maybe()

			res, err := r.Mutation().CreateSighting(tc.ctx, tc.input)
			assert.Equal(t, tc.want, res)
			assert.Equal(t, tc.wantErr, err)

			if tc.wantQueue != nil && tc.withRandomDBErr == false {
				assert.Equal(t, *tc.wantQueue, <-queue)
			}
			close(queue)
		})
	}
}

func TestMutation_CreateTiger(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name  string
		input model.NewTiger

		withRandomDBErr bool
		ctx             context.Context
		want            *model.Tiger
		wantErr         error
	}{
		{
			name: "should return tiger with id 2, sighting id 2 and user id 1 and nil error",
			input: model.NewTiger{
				Name:          "tiger-2",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.250676,
				LastLongitude: 111.828316,
			},
			ctx: context.WithValue(context.Background(), user.KeyUser, &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:         "user-1",
				Email:        "email-1@example.com",
				PasswordHash: "hashed-password-1",
			}),
			want: &model.Tiger{
				ID:            2,
				Name:          "tiger-2",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.250676,
				LastLongitude: 111.828316,
			},
			wantErr: nil,
		},
		{
			name: "should return nil and error given user not found",
			input: model.NewTiger{
				Name:          "tiger-2",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.250676,
				LastLongitude: 111.828316,
			},
			ctx:     context.WithValue(context.Background(), user.KeyUser, nil),
			want:    nil,
			wantErr: errs.RespError(entities.ErrUserByCtxNotFound),
		},
		{
			name: "should return nil and error given random db error",
			input: model.NewTiger{
				Name:          "tiger-2",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.250676,
				LastLongitude: 111.828316,
			},
			ctx: context.WithValue(context.Background(), user.KeyUser, &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:         "user-1",
				Email:        "email-1@example.com",
				PasswordHash: "hashed-password-1",
			}),
			withRandomDBErr: true,
			want:            nil,
			wantErr:         errs.RespError(errors.New("no such table: tigers")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, mockS3, _ := Setup(t, now, tc.withRandomDBErr)

			mockS3.
				On(
					"UploadImage",
					mock.Anything,
					mock.Anything,
					"filename.jpeg",
					"image/jpeg",
					int64(1000),
				).
				Return("https://example.com/image.jpeg", nil).
				Maybe()

			res, err := r.Mutation().CreateTiger(tc.ctx, tc.input)

			assert.Equal(t, tc.want, res)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
