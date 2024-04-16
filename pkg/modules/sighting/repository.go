package sighting

import (
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/models"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func (r *repo) Create(sighting *models.Sighting) error {
	err := r.db.Create(sighting).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) FindByTigerID(tigerID uint) ([]models.Sighting, error) {
	var res []models.Sighting
	err := r.db.Where("tiger_id = ?", tigerID).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewSightingRepository(db *gorm.DB) models.SightingRepository {
	return &repo{db}
}
