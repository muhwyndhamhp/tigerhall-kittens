package models

import "gorm.io/gorm"

type Tiger struct {
	gorm.Model
	Name          string      `json:"name"`
	DateOfBirth   string      `json:"dateOfBirth"`
	LastSeen      string      `json:"lastSeen"`
	LastLatitude  float64     `json:"lastLatitude"`
	LastLongitude float64     `json:"lastLongitude"`
	Sightings     []*Sighting `json:"sightings"`
}

type TigerRepository interface {
	Create(tiger *Tiger) error
	FindAll(page, pageSize int) ([]Tiger, error)
	FindByID(id uint) (*Tiger, error)
	Update(tiger *Tiger, id uint) error
}
