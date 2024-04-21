package sighting

import (
	"context"
	"fmt"

	geo "github.com/kellydunn/golang-geo"
	"github.com/labstack/gommon/log"
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
	s3        s3client.S3ClientInterface
	ch        chan<- email.SightingEmail
}

// CreateSighting implements entities.SightingUsecase.
func (u *usecase) CreateSighting(ctx context.Context, sighting *model.NewSighting, userID uint) (*model.Sighting, error) {
	t, err := u.tigerRepo.FindByID(ctx, sighting.TigerID)
	if err != nil {
		return nil, err
	}

	p0 := geo.NewPoint(t.LastLatitude, t.LastLongitude)
	p1 := geo.NewPoint(sighting.Latitude, sighting.Longitude)
	if p1.GreatCircleDistance(p0) <= 5.0 {
		return nil, entities.ErrTigerTooClose
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

	t.LastSeen = s.Date
	t.LastLatitude = s.Latitude
	t.LastLongitude = s.Longitude

	err = u.tigerRepo.Update(ctx, t, t.ID)
	if err != nil {
		return nil, err
	}

	m := &model.Sighting{
		ID:        s.ID,
		Date:      s.Date,
		Latitude:  s.Latitude,
		Longitude: s.Longitude,
		TigerID:   s.TigerID,
		UserID:    s.UserID,
	}

	if s.ImageURL != "" {
		m.ImageURL = &s.ImageURL
	}

	go u.queueEmail(t, s.ImageURL)
	return m, nil
}

func (u *usecase) queueEmail(t *entities.Tiger, imageURL string) {
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
		log.Error(err)
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
			SightingDate:      t.LastSeen.Format("2006-01-02 15:04:05"),
			SightingLatitude:  fmt.Sprintf("%f", t.LastLatitude),
			SightingLongitude: fmt.Sprintf("%f", t.LastLongitude),
			ImageURL:          imageURL,
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
	s3 s3client.S3ClientInterface,
	ch chan<- email.SightingEmail,
) entities.SightingUsecase {
	return &usecase{repo, tigerRepo, userRepo, s3, ch}
}
