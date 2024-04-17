package tiger

import (
	"context"

	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
)

type usecase struct {
	repo         entities.TigerRepository
	sightingRepo entities.SightingRepository
}

// CreateTiger implements entities.TigerUsecase.
func (u *usecase) CreateTiger(ctx context.Context, tiger *model.NewTiger, userID uint) (*model.Tiger, error) {
	t := entities.Tiger{
		Name:          tiger.Name,
		DateOfBirth:   tiger.DateOfBirth,
		LastSeen:      tiger.LastSeen,
		LastLatitude:  tiger.LastLatitude,
		LastLongitude: tiger.LastLongitude,
	}

	err := u.repo.Create(ctx, &t)
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

	err = u.sightingRepo.Create(ctx, &sighting)
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
	}, nil
}

// GetTigerByID implements entities.TigerUsecase.
func (u *usecase) GetTigerByID(ctx context.Context, id uint) (*model.Tiger, error) {
	t, err := u.repo.FindByID(ctx, id)
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
func (u *usecase) GetTigers(ctx context.Context, page int, pageSize int) ([]*model.Tiger, int, error) {
	tigers, count, err := u.repo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
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

	return res, count, nil
}

func NewTigerUsecase(repo entities.TigerRepository, sightingRepo entities.SightingRepository) entities.TigerUsecase {
	return &usecase{repo, sightingRepo}
}
