package graph

import (
	"testing"
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/db"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/sighting"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/tiger"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/user"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/email"
	s3mock "github.com/muhwyndhamhp/tigerhall-kittens/utils/s3client/mocks"
	"gorm.io/gorm"
)

func Setup(t *testing.T, now time.Time, randomDBErr bool) (*Resolver, *s3mock.S3ClientInterface, chan email.SightingEmail) {
	d := db.GetTestDB()

	SeedDB(d, now, randomDBErr)

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
						PasswordHash: "$2a$10$MGPcG.T8.KzfqkwgPq9TDuiOGLi45guJQ8PQSM.yXMrjeoRs.Wi2C",
					},
				},
			},
		}).Error
		if err != nil {
			panic(err)
		}
	}
}

func GenerateJWT(u *entities.User) string {
	if u == nil {
		u = &entities.User{
			Model: gorm.Model{
				ID: 1,
			},
			Name:         "user-1",
			Email:        "email-1@example.com",
			PasswordHash: "hashed-password-1",
		}
	}

	jwt, _ := u.GenerateToken()

	return jwt
}
