package user

import (
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"gorm.io/gorm"
)

func GenerateJWT(u *entities.User) string {
	if u == nil {
		u = &entities.User{
			Model: gorm.Model{
				ID: 1,
			},
			Name:         "user-1",
			Email:        "email-1@example.com",
			PasswordHash: "inipasswordnya!",
		}
	}
	jwt, _ := u.GenerateToken()

	return jwt
}
