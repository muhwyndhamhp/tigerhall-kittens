package tiger

import (
	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
)

type usecase struct {
	repo         entities.TigerRepository
	sightingRepo entities.SightingRepository
}

// CreateTiger implements entities.TigerUsecase.
func (u *usecase) CreateTiger(tiger *model.Tiger, userID uint) (*model.Tiger, error) {
	t := entities.Tiger{
		Name:          tiger.Name,
		DateOfBirth:   tiger.DateOfBirth,
		LastSeen:      tiger.LastSeen,
		LastLatitude:  tiger.LastLatitude,
		LastLongitude: tiger.LastLongitude,
	}

	err := u.repo.Create(&t)
	if err != nil {
		return nil, err
	}

	sighting := entities.Sighting{
		Date:      tiger.LastSeen,
		Latitude:  tiger.LastLatitude,
		Longitude: tiger.LastLongitude,
		TigerID:   t.ID,
		UserID:    userID,
	}

	err = u.sightingRepo.Create(&sighting)
	if err != nil {
		return nil, err
	}

	return &model.Tiger{
		ID:            t.ID,
		Name:          t.Name,
		DateOfBirth:   t.DateOfBirth,
		LastSeen:      t.LastSeen,
		LastLatitude:  t.LastLatitude,
		LastLongitude: t.LastLongitude,
		Sightings: []*model.Sighting{
			{
				Date:      sighting.Date,
				Latitude:  sighting.Latitude,
				Longitude: sighting.Longitude,
				UserID:    sighting.UserID,
			},
		},
	}, nil
}

// GetTigerByID implements entities.TigerUsecase.
func (u *usecase) GetTigerByID(id uint) (*model.Tiger, error) {
	t, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &model.Tiger{
		ID:            t.ID,
		Name:          t.Name,
		DateOfBirth:   t.DateOfBirth,
		LastSeen:      t.LastSeen,
		LastLatitude:  t.LastLatitude,
		LastLongitude: t.LastLongitude,
		Sightings:     nil,
	}, nil
}

// GetTigers implements entities.TigerUsecase.
func (u *usecase) GetTigers(page int, pageSize int) ([]*model.Tiger, error) {
	tigers, err := u.repo.FindAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	res := make([]*model.Tiger, len(tigers))
	for i, t := range tigers {
		res[i] = &model.Tiger{
			ID:            t.ID,
			Name:          t.Name,
			DateOfBirth:   t.DateOfBirth,
			LastSeen:      t.LastSeen,
			LastLatitude:  t.LastLatitude,
			LastLongitude: t.LastLongitude,
			Sightings:     nil,
		}
	}

	return res, nil
}

func NewTigerUsecase(repo entities.TigerRepository, sightingRepo entities.SightingRepository) entities.TigerUsecase {
	return &usecase{repo, sightingRepo}
}
