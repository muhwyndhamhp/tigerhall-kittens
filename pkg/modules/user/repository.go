package user

import (
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

// FindByID implements entities.UserRepository.
func (r *repo) FindByID(id uint) (*entities.User, error) {
	var res entities.User
	err := r.db.Where("id = ?", id).First(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Create implements entities.UserRepository.
func (r *repo) Create(user *entities.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

// FindByEmail implements entities.UserRepository.
func (r *repo) FindByEmail(email string) (*entities.User, error) {
	var res entities.User
	err := r.db.Where("email = ?", email).First(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func NewUserRepository(db *gorm.DB) entities.UserRepository {
	return &repo{db}
}
