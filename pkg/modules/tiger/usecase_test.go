package tiger

import (
	"context"
	"testing"
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities/mocks"
	s3mocks "github.com/muhwyndhamhp/tigerhall-kittens/utils/s3client/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUsecase_CreateTiger(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string

		createTigerErr    error
		createSightingErr error

		want    *model.Tiger
		wantErr error
	}{
		{
			name:              "should return *model.Tiger and nil error",
			createTigerErr:    nil,
			createSightingErr: nil,
			want: &model.Tiger{
				Name:          "tiger-1",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.550676,
				LastLongitude: 110.828316,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewTigerRepository(t)
			sightingRepo := mocks.NewSightingRepository(t)
			s3 := s3mocks.NewS3ClientInterface(t)

			uc := NewTigerUsecase(repo, sightingRepo, s3)

			repo.
				On("Create", mock.Anything, &entities.Tiger{
					Name:          "tiger-1",
					DateOfBirth:   now,
					LastSeen:      now,
					LastLatitude:  -7.550676,
					LastLongitude: 110.828316,
				}).
				Return(tc.createTigerErr).
				Once()

			sightingRepo.
				On("Create", mock.Anything, mock.Anything).
				Return(tc.createSightingErr).
				Maybe()

			got, err := uc.CreateTiger(context.Background(), &model.NewTiger{
				Name:          "tiger-1",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.550676,
				LastLongitude: 110.828316,
			}, 1)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUsecase_GetTigerByID(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string

		findByIDResp *entities.Tiger
		findByIDErr  error

		want    *model.Tiger
		wantErr error
	}{
		{
			name: "should return *model.Tiger and nil error",
			findByIDResp: &entities.Tiger{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:          "tiger-1",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.550676,
				LastLongitude: 110.828316,
			},
			findByIDErr: nil,
			want: &model.Tiger{
				ID:            1,
				Name:          "tiger-1",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.550676,
				LastLongitude: 110.828316,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewTigerRepository(t)
			sightingRepo := mocks.NewSightingRepository(t)
			s3 := s3mocks.NewS3ClientInterface(t)

			uc := NewTigerUsecase(repo, sightingRepo, s3)

			repo.
				On("FindByID", mock.Anything, uint(1)).
				Return(tc.findByIDResp, tc.findByIDErr).
				Once()

			got, err := uc.GetTigerByID(context.Background(), 1)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUsecase_GetTigers(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name string

		findAllResp  []entities.Tiger
		findAllCount int
		findAllErr   error

		want      []*model.Tiger
		wantCount int
		wantErr   error
	}{
		{
			name: "should return []*model.Tiger and nil error",
			findAllResp: []entities.Tiger{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					Name:          "tiger-1",
					DateOfBirth:   now,
					LastSeen:      now,
					LastLatitude:  -7.550676,
					LastLongitude: 110.828316,
				},
			},
			findAllCount: 1,
			findAllErr:   nil,
			want: []*model.Tiger{
				{
					ID:            1,
					Name:          "tiger-1",
					DateOfBirth:   now,
					LastSeen:      now,
					LastLatitude:  -7.550676,
					LastLongitude: 110.828316,
				},
			},
			wantCount: 1,
			wantErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewTigerRepository(t)
			sightingRepo := mocks.NewSightingRepository(t)
			s3 := s3mocks.NewS3ClientInterface(t)

			uc := NewTigerUsecase(repo, sightingRepo, s3)

			repo.
				On("FindAll", mock.Anything, 1, 10).
				Return(tc.findAllResp, tc.findAllCount, tc.findAllErr).
				Once()

			got, count, err := uc.GetTigers(context.Background(), 1, 10)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}
