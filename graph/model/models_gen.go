// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Mutation struct {
}

type NewSighting struct {
	TigerID   uint      `json:"tigerID"`
	Date      time.Time `json:"date"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

type NewTiger struct {
	Name          string    `json:"name"`
	DateOfBirth   time.Time `json:"dateOfBirth"`
	LastSeen      time.Time `json:"lastSeen"`
	LastLatitude  float64   `json:"lastLatitude"`
	LastLongitude float64   `json:"lastLongitude"`
}

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Query struct {
}

type Sighting struct {
	ID        uint      `json:"id"`
	Date      time.Time `json:"date"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	TigerID   uint      `json:"tigerID"`
	Tiger     *Tiger    `json:"tiger"`
	UserID    uint      `json:"userID"`
	User      *User     `json:"user"`
}

type Tiger struct {
	ID            uint        `json:"id"`
	Name          string      `json:"name"`
	DateOfBirth   time.Time   `json:"dateOfBirth"`
	LastSeen      time.Time   `json:"lastSeen"`
	LastLatitude  float64     `json:"lastLatitude"`
	LastLongitude float64     `json:"lastLongitude"`
	Sightings     []*Sighting `json:"sightings"`
}

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
