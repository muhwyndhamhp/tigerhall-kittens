package sighting

import (
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/scopes"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func (r *repo) Create(sighting *entities.Sighting) error {
	err := r.db.Create(sighting).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) FindByTigerID(tigerID uint, page, pageSize int) ([]entities.Sighting, error) {
	var res []entities.Sighting
	err := r.db.
		Scopes(scopes.Paginate(page, pageSize)).
		Where("tiger_id = ?", tigerID).
		Order("date DESC").
		Find(&res).
		Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewSightingRepository(db *gorm.DB) entities.SightingRepository {
	return &repo{db}
}
