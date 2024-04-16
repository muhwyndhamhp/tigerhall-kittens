package models

import "gorm.io/gorm"

type Sighting struct {
	gorm.Model
	Date      string  `json:"date"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	TigerID   uint    `json:"tigerID"`
}

type SightingRepository interface {
	Create(sighting *Sighting) error
	FindByTigerID(tigerID uint) ([]Sighting, error)
}
