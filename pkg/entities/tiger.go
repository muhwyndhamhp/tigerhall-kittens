package entities

import (
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"gorm.io/gorm"
)

type Tiger struct {
	gorm.Model
	Name          string      `json:"name"`
	DateOfBirth   time.Time   `json:"date_of_birth"`
	LastSeen      time.Time   `json:"last_seen"`
	LastLatitude  float64     `json:"last_latitude"`
	LastLongitude float64     `json:"last_longitude"`
	Sightings     []*Sighting `json:"sightings"`
}

type TigerUsecase interface {
	CreateTiger(tiger *model.Tiger) (*model.Tiger, error)
	GetTigers(page, pageSize int) ([]*model.Tiger, error)
	GetTigerByID(id uint) (*model.Tiger, error)
}

type TigerRepository interface {
	Create(tiger *Tiger) error
	FindAll(page, pageSize int) ([]Tiger, error)
	FindByID(id uint) (*Tiger, error)
	Update(tiger *Tiger, id uint) error
}
