package sighting

import (
	"context"
	"fmt"

	geo "github.com/kellydunn/golang-geo"
	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/email"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/imageproc"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/s3client"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/scopes"
)

type usecase struct {
	repo      entities.SightingRepository
	tigerRepo entities.TigerRepository
	userRepo  entities.UserRepository
	s3        *s3client.S3Client
	ch        chan<- email.SightingEmail
}

// CreateSighting implements entities.SightingUsecase.
func (u *usecase) CreateSighting(ctx context.Context, sighting *model.NewSighting, userID uint) (*model.Sighting, error) {
	ls, _, err := u.repo.FindByTigerID(ctx, sighting.TigerID, []scopes.Preload{}, 1, 1)
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

	go u.queueEmail(t)

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

func (u *usecase) queueEmail(t *entities.Tiger) {
	fmt.Println("Sending email to other users...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sh, _, err := u.repo.FindByTigerID(ctx, t.ID, []scopes.Preload{
		{
			Key:       "User",
			Statement: "",
		},
	}, 1, 1000)
	if err != nil {
		fmt.Println(err)
		return
	}

	sentUsr := map[string]bool{}
	for _, s := range sh {
		if _, ok := sentUsr[s.User.Email]; ok {
			continue
		}

		m := email.SightingEmail{
			DestinationEmail:  s.User.Email,
			TigerName:         t.Name,
			SightingDate:      s.Date.Format("2006-01-02 15:04:05"),
			SightingLatitude:  fmt.Sprintf("%f", s.Latitude),
			SightingLongitude: fmt.Sprintf("%f", s.Longitude),
			ImageURL:          s.ImageURL,
		}
		u.ch <- m

		sentUsr[s.User.Email] = true
	}
}

// GetSightingsByTigerID implements entities.SightingUsecase.
func (u *usecase) GetSightingsByTigerID(ctx context.Context, tigerID uint, page int, pageSize int) ([]*model.Sighting, int, error) {
	sightings, count, err := u.repo.FindByTigerID(ctx, tigerID, []scopes.Preload{}, page, pageSize)
	if err != nil {
		return nil, 0, err
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
	return result, count, nil
}

func NewSightingUsecase(
	repo entities.SightingRepository,
	tigerRepo entities.TigerRepository,
	userRepo entities.UserRepository,
	s3 *s3client.S3Client,
	ch chan<- email.SightingEmail,
) entities.SightingUsecase {
	return &usecase{repo, tigerRepo, userRepo, s3, ch}
}
