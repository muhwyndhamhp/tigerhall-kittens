package entities

import (
	"context"
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
	CreateTiger(ctx context.Context, tiger *model.NewTiger, userID uint) (*model.Tiger, error)
	GetTigers(ctx context.Context, page, pageSize int) ([]*model.Tiger, error)
	GetTigerByID(ctx context.Context, id uint) (*model.Tiger, error)
}

type TigerRepository interface {
	Create(ctx context.Context, tiger *Tiger) error
	FindAll(ctx context.Context, page, pageSize int) ([]Tiger, error)
	FindByID(ctx context.Context, id uint) (*Tiger, error)
	Update(ctx context.Context, tiger *Tiger, id uint) error
}
