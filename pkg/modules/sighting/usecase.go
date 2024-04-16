package sighting

import (
	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
)

type usecase struct {
	repo      entities.SightingRepository
	tigerRepo entities.TigerRepository
	userRepo  entities.UserRepository
}

// CreateSighting implements entities.SightingUsecase.
func (u *usecase) CreateSighting(sighting *model.Sighting) (*model.Sighting, error) {
	s := entities.Sighting{
		Date:      sighting.Date,
		Latitude:  sighting.Latitude,
		Longitude: sighting.Longitude,
		TigerID:   sighting.TigerID,
		UserID:    sighting.UserID,
	}

	err := u.repo.Create(&s)
	if err != nil {
		return nil, err
	}

	t, err := u.tigerRepo.FindByID(s.TigerID)
	if err != nil {
		return nil, err
	}

	t.LastSeen = s.Date
	t.LastLatitude = s.Latitude
	t.LastLongitude = s.Longitude

	err = u.tigerRepo.Update(t, t.ID)
	if err != nil {
		return nil, err
	}

	usr, err := u.userRepo.FindByID(s.UserID)
	if err != nil {
		return nil, err
	}

	return &model.Sighting{
		ID:        s.ID,
		Date:      s.Date,
		Latitude:  s.Latitude,
		Longitude: s.Longitude,
		TigerID:   s.TigerID,
		UserID:    s.UserID,
		Tiger: &model.Tiger{
			ID:            t.ID,
			Name:          t.Name,
			LastSeen:      t.LastSeen,
			LastLatitude:  t.LastLatitude,
			LastLongitude: t.LastLongitude,
		},
		User: &model.User{
			ID:    usr.ID,
			Name:  usr.Name,
			Email: usr.Email,
		},
	}, nil
}

// GetSightingsByTigerID implements entities.SightingUsecase.
func (u *usecase) GetSightingsByTigerID(tigerID uint, page int, pageSize int) ([]*model.Sighting, error) {
	panic("unimplemented")
}

func NewSightingUsecase(repo entities.SightingRepository, tigerRepo entities.TigerRepository, userRepo entities.UserRepository) entities.SightingUsecase {
	return &usecase{repo, tigerRepo, userRepo}
}
