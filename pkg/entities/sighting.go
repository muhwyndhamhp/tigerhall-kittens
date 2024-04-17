package entities

import (
	"context"
	"errors"
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/scopes"
	"gorm.io/gorm"
)

type Sighting struct {
	gorm.Model
	Date      time.Time `json:"date"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	TigerID   uint      `json:"tiger_id"`
	UserID    uint      `json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID"`
	ImageURL  string    `json:"image_url"`
}

var (
	ErrTigerTooClose    = errors.New("ErrTigerTooClose: tiger is too close, new sightings should be at least more than 5km away from the last sighting")
	ErrInvalidImageType = errors.New("ErrInvalidImageType: invalid image type, only jpeg, jpg, and png are allowed")
)

type SightingUsecase interface {
	CreateSighting(ctx context.Context, sighting *model.NewSighting, userID uint) (*model.Sighting, error)
	GetSightingsByTigerID(ctx context.Context, tigerID uint, page, pageSize int) ([]*model.Sighting, int, error)
}

type SightingRepository interface {
	Create(ctx context.Context, sighting *Sighting) error
	FindByTigerID(ctx context.Context, tigerID uint, preloads []scopes.Preload, page, pageSize int) ([]Sighting, int, error)
}
