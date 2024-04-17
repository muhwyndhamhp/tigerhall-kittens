package sighting

import (
	"context"

	geo "github.com/kellydunn/golang-geo"
	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/imageproc"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/s3client"
)

type usecase struct {
	repo      entities.SightingRepository
	tigerRepo entities.TigerRepository
	userRepo  entities.UserRepository
	s3        *s3client.S3Client
}

// CreateSighting implements entities.SightingUsecase.
func (u *usecase) CreateSighting(ctx context.Context, sighting *model.NewSighting, userID uint) (*model.Sighting, error) {
	ls, err := u.repo.FindByTigerID(ctx, sighting.TigerID, 1, 1)
	if err != nil {
		return nil, err
	}

	if len(ls) > 0 {
		p0 := geo.NewPoint(ls[0].Latitude, ls[0].Longitude)
		p1 := geo.NewPoint(sighting.Latitude, sighting.Longitude)
		if p1.GreatCircleDistance(p0) <= 5.0 {
			return nil, entities.ErrTigerTooClose
		}
	}

	s := entities.Sighting{
		Date:      sighting.Date,
		Latitude:  sighting.Latitude,
		Longitude: sighting.Longitude,
		TigerID:   sighting.TigerID,
		UserID:    userID,
	}
	if sighting.Image != nil {
		if !imageproc.IsContentTypeValid(sighting.Image.ContentType, sighting.Image.Filename) {
			return nil, entities.ErrInvalidImageType
		}

		r, size, err := imageproc.ResizeImage(sighting.Image.File, sighting.Image.Filename)
		if err != nil {
			return nil, err
		}

		url, err := u.s3.UploadImage(ctx, r, sighting.Image.Filename, sighting.Image.ContentType, int64(size))
		if err != nil {
			return nil, err
		}

		s.ImageURL = url
	}

	err = u.repo.Create(ctx, &s)
	if err != nil {
		return nil, err
	}

	t, err := u.tigerRepo.FindByID(ctx, s.TigerID)
	if err != nil {
		return nil, err
	}

	t.LastSeen = s.Date
	t.LastLatitude = s.Latitude
	t.LastLongitude = s.Longitude

	err = u.tigerRepo.Update(ctx, t, t.ID)
	if err != nil {
		return nil, err
	}

	usr, err := u.userRepo.FindByID(ctx, s.UserID)
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
		ImageURL:  &s.ImageURL,
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
func (u *usecase) GetSightingsByTigerID(ctx context.Context, tigerID uint, page int, pageSize int) ([]*model.Sighting, error) {
	sightings, err := u.repo.FindByTigerID(ctx, tigerID, page, pageSize)
	if err != nil {
		return nil, err
	}

	var result []*model.Sighting
	for _, s := range sightings {
		result = append(result, &model.Sighting{
			ID:        s.ID,
			Date:      s.Date,
			Latitude:  s.Latitude,
			Longitude: s.Longitude,
			TigerID:   s.TigerID,
			UserID:    s.UserID,
			ImageURL:  &s.ImageURL,
		})
	}
	return result, nil
}

func NewSightingUsecase(
	repo entities.SightingRepository,
	tigerRepo entities.TigerRepository,
	userRepo entities.UserRepository,
	s3 *s3client.S3Client,
) entities.SightingUsecase {
	return &usecase{repo, tigerRepo, userRepo, s3}
}
