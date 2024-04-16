package user

import (
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/models"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

// FindByID implements models.UserRepository.
func (r *repo) FindByID(id uint) (*models.User, error) {
	var res models.User
	err := r.db.Where("id = ?", id).First(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Create implements models.UserRepository.
func (r *repo) Create(user *models.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

// FindByEmail implements models.UserRepository.
func (r *repo) FindByEmail(email string) (*models.User, error) {
	var res models.User
	err := r.db.Where("email = ?", email).First(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func NewUserRepository(db *gorm.DB) models.UserRepository {
	return &repo{db}
}
