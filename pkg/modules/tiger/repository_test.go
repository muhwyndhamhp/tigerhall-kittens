package tiger

import (
	"context"
	"testing"
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/db"
	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRepository_Create(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name string

		tiger   *entities.Tiger
		want    *model.Tiger
		wantErr error
	}{
		{
			name: "should create new tiger with id 2",
			tiger: &entities.Tiger{
				Name:          "tiger-2",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.550676,
				LastLongitude: 110.828316,
			},
			want: &model.Tiger{
				ID:            2,
				Name:          "tiger-2",
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
			d := db.GetTestDB()

			SeedDb(d, now)

			r := NewTigerRepository(d)

			err := r.Create(context.Background(), tc.tiger)
			assert.Equal(t, tc.wantErr, err)

			if tc.want != nil {
				// assert every field of tiger
				assert.Equal(t, tc.want.ID, tc.tiger.ID)
				assert.Equal(t, tc.want.Name, tc.tiger.Name)
				assert.Equal(t, tc.want.DateOfBirth, tc.tiger.DateOfBirth)
				assert.Equal(t, tc.want.LastSeen, tc.tiger.LastSeen)
				assert.Equal(t, tc.want.LastLatitude, tc.tiger.LastLatitude)
				assert.Equal(t, tc.want.LastLongitude, tc.tiger.LastLongitude)
			}
		})
	}
}

func TestRepository_FindByID(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name string

		id      uint
		want    *model.Tiger
		wantErr error
	}{
		{
			name: "should retrieve tiger with id 1",
			id:   1,
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
			d := db.GetTestDB()
			SeedDb(d, now)

			r := NewTigerRepository(d)

			got, err := r.FindByID(context.Background(), tc.id)

			assert.Equal(t, tc.wantErr, err)

			if tc.want != nil {
				assert.Equal(t, tc.want.ID, got.ID)
				assert.Equal(t, tc.want.Name, got.Name)
				assert.Equal(t, tc.want.DateOfBirth.Format(time.RFC3339), got.DateOfBirth.Format(time.RFC3339))
				assert.Equal(t, tc.want.LastSeen.Format(time.RFC3339), got.LastSeen.Format(time.RFC3339))
				assert.Equal(t, tc.want.LastLatitude, got.LastLatitude)
				assert.Equal(t, tc.want.LastLongitude, got.LastLongitude)
			}
		})
	}
}

func TestRepository_FindAll(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name string

		page     int
		pageSize int
		want     []*model.Tiger
		wantErr  error
	}{
		{
			name:     "should retrieve all tigers",
			page:     1,
			pageSize: 10,
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
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := db.GetTestDB()
			SeedDb(d, now)

			r := NewTigerRepository(d)

			got, _, err := r.FindAll(context.Background(), tc.page, tc.pageSize)

			assert.Equal(t, tc.wantErr, err)

			for i, want := range tc.want {
				assert.Equal(t, want.ID, got[i].ID)
				assert.Equal(t, want.Name, got[i].Name)
				assert.Equal(t, want.DateOfBirth.Format(time.RFC3339), got[i].DateOfBirth.Format(time.RFC3339))
				assert.Equal(t, want.LastSeen.Format(time.RFC3339), got[i].LastSeen.Format(time.RFC3339))
				assert.Equal(t, want.LastLatitude, got[i].LastLatitude)
				assert.Equal(t, want.LastLongitude, got[i].LastLongitude)
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name string

		tiger   *entities.Tiger
		id      uint
		want    *model.Tiger
		wantErr error
	}{
		{
			name: "should update tiger with id 1",
			tiger: &entities.Tiger{
				Name:          "tiger-1-updated",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.550676,
				LastLongitude: 110.828316,
			},
			id: 1,
			want: &model.Tiger{
				ID:            1,
				Name:          "tiger-1-updated",
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
			d := db.GetTestDB()
			SeedDb(d, now)

			r := NewTigerRepository(d)

			err := r.Update(context.Background(), tc.tiger, tc.id)

			assert.Equal(t, tc.wantErr, err)

			if tc.want != nil {
				got, _ := r.FindByID(context.Background(), tc.id)

				assert.Equal(t, tc.want.ID, got.ID)
				assert.Equal(t, tc.want.Name, got.Name)
				assert.Equal(t, tc.want.DateOfBirth.Format(time.RFC3339), got.DateOfBirth.Format(time.RFC3339))
				assert.Equal(t, tc.want.LastSeen.Format(time.RFC3339), got.LastSeen.Format(time.RFC3339))
				assert.Equal(t, tc.want.LastLatitude, got.LastLatitude)
				assert.Equal(t, tc.want.LastLongitude, got.LastLongitude)
			}
		})
	}
}

func SeedDb(d *gorm.DB, now time.Time) {
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
