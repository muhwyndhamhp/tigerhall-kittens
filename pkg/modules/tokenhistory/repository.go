package tokenhistory

import (
	"context"

	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

// Create implements entities.TokenHistoryRepository.
func (r *repo) Create(ctx context.Context, token *entities.TokenHistory) error {
	err := r.db.WithContext(ctx).Create(token).Error
	if err != nil {
		return err
	}

	return nil
}

// FindByToken implements entities.TokenHistoryRepository.
func (r *repo) FindByToken(ctx context.Context, token string) (*entities.TokenHistory, error) {
	var res entities.TokenHistory
	err := r.db.
		WithContext(ctx).
		Where("token = ?", token).
		First(&res).
		Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func NewTokenHistoryRepository(db *gorm.DB) entities.TokenHistoryRepository {
	return &repo{db}
}
