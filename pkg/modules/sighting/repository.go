package sighting

import (
	"context"

	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/scopes"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func (r *repo) Create(ctx context.Context, sighting *entities.Sighting) error {
	err := r.db.WithContext(ctx).Create(sighting).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) FindByTigerID(
	ctx context.Context,
	tigerID uint,
	preloads []scopes.Preload,
	page, pageSize int,
) ([]entities.Sighting, int, error) {
	var res []entities.Sighting

	q := r.db.WithContext(ctx).
		Model(&entities.Sighting{}).
		Where("tiger_id = ?", tigerID)

	var count int64
	err := q.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	q = q.Scopes(scopes.Paginate(page, pageSize)).Order("date DESC")

	if len(preloads) > 0 {
		q = q.Scopes(scopes.Preloads(preloads...))
	}

	err = q.Find(&res).Error
	if err != nil {
		return nil, 0, err
	}

	return res, int(count), nil
}

func NewSightingRepository(db *gorm.DB) entities.SightingRepository {
	return &repo{db}
}
