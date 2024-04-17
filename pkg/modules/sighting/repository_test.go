package sighting

import (
	"context"
	"testing"
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/db"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/scopes"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRepository_Create(t *testing.T) {
	now := time.Now()

	tc := []struct {
		name string

		sighting *entities.Sighting
		want     *entities.Sighting
		wantErr  error
	}{
		{
			name: "should create new sighting with id 2",
			sighting: &entities.Sighting{
				Date:      now,
				Latitude:  -7.550676,
				Longitude: 110.828316,
				TigerID:   1,
				UserID:    1,
			},
			want: &entities.Sighting{
				Model: gorm.Model{
					ID: 2,
				},
				Date:      now,
				Latitude:  -7.550676,
				Longitude: 110.828316,
				TigerID:   1,
				UserID:    1,
			},
			wantErr: nil,
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			d := db.GetTestDB()
			SeedDB(d, now)

			r := NewSightingRepository(d)

			err := r.Create(context.Background(), c.sighting)

			assert.Equal(t, c.wantErr, err)
			if c.want != nil {
				assert.Equal(t, c.want.ID, c.sighting.ID)
				assert.Equal(t, c.want.Date.Format(time.RFC3339), c.sighting.Date.Format(time.RFC3339))
				assert.Equal(t, c.want.Latitude, c.sighting.Latitude)
				assert.Equal(t, c.want.Longitude, c.sighting.Longitude)
				assert.Equal(t, c.want.TigerID, c.sighting.TigerID)
				assert.Equal(t, c.want.UserID, c.sighting.UserID)
			}
		})
	}
}

func TestRepository_FindByTigerID(t *testing.T) {
	now := time.Now()
	tc := []struct {
		name string

		preloadStmt []scopes.Preload
		want        []entities.Sighting
		wantCount   int
		wantErr     error
	}{
		{
			name: "should return 1 sighting with user preloaded",
			preloadStmt: []scopes.Preload{
				{
					Key: "User",
				},
			},
			want: []entities.Sighting{
				{
					Model: gorm.Model{
						ID: 1,
					},
					Date:      now,
					Latitude:  -7.550676,
					Longitude: 110.828316,
					TigerID:   1,
					UserID:    1,
				},
			},
			wantCount: 1,
			wantErr:   nil,
		},
		{
			name: "should return 1 sighting",
			want: []entities.Sighting{
				{
					Model: gorm.Model{
						ID: 1,
					},
					Date:      now,
					Latitude:  -7.550676,
					Longitude: 110.828316,
					TigerID:   1,
					UserID:    1,
				},
			},
			wantCount: 1,
			wantErr:   nil,
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			d := db.GetTestDB()
			SeedDB(d, now)

			r := NewSightingRepository(d)

			res, count, err := r.FindByTigerID(context.Background(), 1, c.preloadStmt, 1, 10)

			assert.Equal(t, c.wantErr, err)
			assert.Equal(t, c.wantCount, count)
			for i := range c.want {
				assert.Equal(t, c.want[i].ID, res[i].ID)
				assert.Equal(t, c.want[i].Date.Format(time.RFC3339), res[i].Date.Format(time.RFC3339))
				assert.Equal(t, c.want[i].Latitude, res[i].Latitude)
				assert.Equal(t, c.want[i].Longitude, res[i].Longitude)
				assert.Equal(t, c.want[i].TigerID, res[i].TigerID)
				assert.Equal(t, c.want[i].UserID, res[i].UserID)
				if len(c.preloadStmt) > 0 && c.preloadStmt[0].Key == "User" {
					assert.NotNil(t, res[i].User)
				}
			}
		})
	}
}

func SeedDB(d *gorm.DB, now time.Time) {
	err := d.AutoMigrate(&entities.User{}, &entities.Tiger{}, &entities.Sighting{})
	if err != nil {
		panic(err)
	}

	err = d.Create(&entities.Sighting{
		Date:      now,
		Latitude:  -7.550676,
		Longitude: 110.828316,
		TigerID:   1,
		UserID:    1,
	}).Error
	if err != nil {
		panic(err)
	}

	// create user
	err = d.Create(&entities.User{
		Name:         "user",
		Email:        "mail-1@example.com",
		PasswordHash: "random-hash",
	}).Error
	if err != nil {
		panic(err)
	}
}
