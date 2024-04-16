package entities

import (
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"gorm.io/gorm"
)

type Sighting struct {
	gorm.Model
	Date      time.Time `json:"date"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	TigerID   uint      `json:"tiger_id"`
	UserID    uint      `json:"user_id"`
}

type SightingUsecase interface {
	CreateSighting(sighting *model.Sighting) (*model.Sighting, error)
	GetSightingsByTigerID(tigerID uint, page, pageSize int) ([]*model.Sighting, error)
}

type SightingRepository interface {
	Create(sighting *Sighting) error
	FindByTigerID(tigerID uint, page, pageSize int) ([]Sighting, error)
}
