package graph

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/db"
	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/sighting"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/tiger"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/user"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/email"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/errs"
	s3mock "github.com/muhwyndhamhp/tigerhall-kittens/utils/s3client/mocks"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/timex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

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

func Setup(t *testing.T, now time.Time, randomDBErr bool) (*Resolver, *s3mock.S3ClientInterface, chan email.SightingEmail) {
	d := db.GetTestDB()

	SeedDB(d, now, randomDBErr)

	timex.SetTestTime(now)

	userRepo := user.NewUserRepository(d)
	tigerRepo := tiger.NewTigerRepository(d)
	sightingRepo := sighting.NewSightingRepository(d)

	mockS3 := s3mock.NewS3ClientInterface(t)

	emailQueue := make(chan email.SightingEmail)

	userUsecase := user.NewUserUsecase(userRepo)
	tigerUsecase := tiger.NewTigerUsecase(tigerRepo, sightingRepo, mockS3)
	sightingUsecase := sighting.NewSightingUsecase(sightingRepo, tigerRepo, userRepo, mockS3, emailQueue)

	r := NewResolver(userUsecase, tigerUsecase, sightingUsecase)

	return r, mockS3, emailQueue
}

func SeedDB(d *gorm.DB, now time.Time, simulateErr bool) {
	if simulateErr {
		err := d.AutoMigrate(&entities.Sighting{}, &entities.User{})
		if err != nil {
			panic(err)
		}

		err = d.Create(&entities.Sighting{
			Date:      now,
			Latitude:  -7.550676,
			Longitude: 110.828316,
			TigerID:   1,
			UserID:    1,
			User: &entities.User{
				Name:         "user-1",
				Email:        "email-1@example.com",
				PasswordHash: "hashed-password-1",
			},
		}).Error
		if err != nil {
			panic(err)
		}
	} else {
		err := d.AutoMigrate(&entities.Tiger{}, &entities.Sighting{}, &entities.User{})
		if err != nil {
			panic(err)
		}
		err = d.Create(&entities.Tiger{
			Name:          "tiger-1",
			DateOfBirth:   now,
			LastSeen:      now,
			LastLatitude:  -7.550676,
			LastLongitude: 110.828316,
			Sightings: []*entities.Sighting{
				{
					Date:      now,
					Latitude:  -7.550676,
					Longitude: 110.828316,
					TigerID:   1,
					UserID:    1,
					User: &entities.User{
						Name:         "user-1",
						Email:        "email-1@example.com",
						PasswordHash: "hashed-password-1",
					},
				},
			},
		}).Error
		if err != nil {
			panic(err)
		}
	}
}
