package user

import (
	"context"
	"testing"
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/db"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRepository_Create(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string

		user *entities.User

		want    *entities.User
		wantErr error
	}{
		{
			name: "should create new user with id 2",
			user: &entities.User{
				Name:         "user-2",
				Email:        "mail-1@example.com",
				PasswordHash: "inipasswordnya!",
			},
			want: &entities.User{
				Model: gorm.Model{
					ID: 2,
				},
				Name:         "user-2",
				Email:        "mail-1@example.com",
				PasswordHash: "inipasswordnya!",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := db.GetTestDB()
			SeedUser(d, now)
			repo := NewUserRepository(d)

			err := repo.Create(context.Background(), tc.user)

			assert.Equal(t, tc.wantErr, err)
			if tc.want != nil {
				// assert every field of user
				assert.Equal(t, tc.want.ID, tc.user.ID)
				assert.Equal(t, tc.want.Name, tc.user.Name)
				assert.Equal(t, tc.want.Email, tc.user.Email)
				assert.Equal(t, tc.want.PasswordHash, tc.user.PasswordHash)
			}
		})
	}
}

func TestRepository_FindByEmail(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string

		email string

		want    *entities.User
		wantErr error
	}{
		{
			name:  "should return user and nil error",
			email: "email-1@example.com",
			want: &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:         "user-1",
				Email:        "email-1@example.com",
				PasswordHash: "hashed-password-1",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := db.GetTestDB()
			SeedUser(d, now)
			repo := NewUserRepository(d)

			user, err := repo.FindByEmail(context.Background(), tc.email)

			assert.Equal(t, tc.wantErr, err)
			if tc.want != nil {
				// assert every field of user
				assert.Equal(t, tc.want.ID, user.ID)
				assert.Equal(t, tc.want.Name, user.Name)
				assert.Equal(t, tc.want.Email, user.Email)
				assert.Equal(t, tc.want.PasswordHash, user.PasswordHash)
			}
		})
	}
}

func TestRepository_FindByID(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string

		id uint

		want    *entities.User
		wantErr error
	}{
		{
			name: "should return user and nil error",
			id:   1,
			want: &entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:         "user-1",
				Email:        "email-1@example.com",
				PasswordHash: "hashed-password-1",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := db.GetTestDB()
			SeedUser(d, now)
			repo := NewUserRepository(d)

			user, err := repo.FindByID(context.Background(), tc.id)

			assert.Equal(t, tc.wantErr, err)
			if tc.want != nil {
				// assert every field of user
				assert.Equal(t, tc.want.ID, user.ID)
				assert.Equal(t, tc.want.Name, user.Name)
				assert.Equal(t, tc.want.Email, user.Email)
				assert.Equal(t, tc.want.PasswordHash, user.PasswordHash)
			}
		})
	}
}

func SeedUser(d *gorm.DB, now time.Time) {
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
