package sighting

import (
	"context"
	"testing"
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities/mocks"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/email"
	s3mocks "github.com/muhwyndhamhp/tigerhall-kittens/utils/s3client/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUsecase_CreateSighting(t *testing.T) {
	now := time.Now()
	req := &model.NewSighting{
		TigerID:   101,
		Date:      now,
		Latitude:  -7.550676,
		Longitude: 110.828316,
		// Image:     &graphql.Upload{}, TODO: handle testing for image
	}

	testCases := []struct {
		name         string
		getTigerResp *entities.Tiger
		getTigerErr  error

		createSightingErr error

		updateTigerErr error

		emailFindByTigerResp []entities.Sighting
		emailFindByTigerErr  error

		want    *model.Sighting
		wantErr error
		wantCh  *email.SightingEmail
	}{
		{
			name:              "should return valid model.Sighting given valid input without image and send email",
			createSightingErr: nil,
			getTigerResp: &entities.Tiger{
				Model: gorm.Model{
					ID:        101,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:          "tiger-1",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.250676,
				LastLongitude: 110.828316,
			},
			getTigerErr:    nil,
			updateTigerErr: nil,
			emailFindByTigerResp: []entities.Sighting{
				{
					Date:      now,
					Latitude:  -7.550676,
					Longitude: 110.828316,
					TigerID:   101,
					UserID:    201,
					User: &entities.User{
						Email: "mail-1@example.com",
					},
				},
			},
			emailFindByTigerErr: nil,
			want: &model.Sighting{
				ID:        0,
				Date:      now,
				Latitude:  -7.550676,
				Longitude: 110.828316,
				TigerID:   101,
				UserID:    201,
				ImageURL:  nil,
			},
			wantErr: nil,
			wantCh: &email.SightingEmail{
				DestinationEmail:  "mail-1@example.com",
				TigerName:         "tiger-1",
				SightingDate:      now.Format("2006-01-02 15:04:05"),
				SightingLatitude:  "-7.550676",
				SightingLongitude: "110.828316",
			},
		},
		{
			name:              "should return valid model.Sighting given valid input without image and send no email",
			createSightingErr: nil,
			getTigerResp: &entities.Tiger{
				Model: gorm.Model{
					ID:        101,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:          "tiger-1",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.250676,
				LastLongitude: 110.828316,
			},
			getTigerErr:          nil,
			updateTigerErr:       nil,
			emailFindByTigerResp: []entities.Sighting{},
			emailFindByTigerErr:  nil,
			want: &model.Sighting{
				ID:        0,
				Date:      now,
				Latitude:  -7.550676,
				Longitude: 110.828316,
				TigerID:   101,
				UserID:    201,
				ImageURL:  nil,
			},
			wantErr: nil,
			wantCh:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewSightingRepository(t)
			tigerRepo := mocks.NewTigerRepository(t)
			userRepo := mocks.NewUserRepository(t)

			s3 := s3mocks.NewS3ClientInterface(t)
			ch := make(chan email.SightingEmail)

			usecase := NewSightingUsecase(repo, tigerRepo, userRepo, s3, ch)

			tigerRepo.
				On("FindByID", mock.Anything, req.TigerID).
				Return(tc.getTigerResp, tc.getTigerErr).
				Once()

			repo.
				On("Create", mock.Anything, mock.Anything).
				Return(tc.createSightingErr).
				Maybe()

			tigerRepo.
				On("Update", mock.Anything, mock.Anything, req.TigerID).
				Return(tc.updateTigerErr).
				Maybe()

			repo.
				On("FindByTigerID", mock.Anything, uint(101), mock.Anything, 1, 1000).
				Return(tc.emailFindByTigerResp, len(tc.emailFindByTigerResp), tc.emailFindByTigerErr).
				Maybe()

			go func() {
				for c := range ch {
					if tc.wantCh != nil {
						assert.Equal(t, *tc.wantCh, c)
					}
					close(ch)
				}
			}()

			res, err := usecase.CreateSighting(context.Background(), req, 201)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, res)
		})
	}
}
