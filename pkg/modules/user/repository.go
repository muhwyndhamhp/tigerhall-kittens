package user

import (
	"context"

	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

// FindByID implements entities.UserRepository.
func (r *repo) FindByID(ctx context.Context, id uint) (*entities.User, error) {
	var res entities.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Create implements entities.UserRepository.
func (r *repo) Create(ctx context.Context, user *entities.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

// FindByEmail implements entities.UserRepository.
func (r *repo) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var res entities.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func NewUserRepository(db *gorm.DB) entities.UserRepository {
	return &repo{db}
}
