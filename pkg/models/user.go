package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
}
