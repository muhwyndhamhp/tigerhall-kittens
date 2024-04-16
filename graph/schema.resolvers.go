package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"

	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/user"
)

// CreateTiger is the resolver for the createTiger field.
func (r *mutationResolver) CreateTiger(ctx context.Context, input model.NewTiger) (*model.Tiger, error) {
	u, err := user.UserByCtx(ctx)
	if err != nil {
		return nil, err
	}

	t, err := r.tigerUsecase.CreateTiger(&model.Tiger{
		Name:          input.Name,
		DateOfBirth:   input.DateOfBirth,
		LastSeen:      input.LastSeen,
		LastLatitude:  input.LastLatitude,
		LastLongitude: input.LastLongitude,
	}, u.ID)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// CreateSighting is the resolver for the createSighting field.
func (r *mutationResolver) CreateSighting(ctx context.Context, input model.NewSighting) (*model.Sighting, error) {
	u, err := user.UserByCtx(ctx)
	if err != nil {
		return nil, err
	}

	s, err := r.sightingUsecase.CreateSighting(&model.Sighting{
		Date:      input.Date,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		TigerID:   input.TigerID,
		UserID:    u.ID,
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	token, err := r.userUsecase.CreateUser(input.Name, input.Email, input.Password)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (string, error) {
	token, err := r.userUsecase.Login(email, password)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Tigers is the resolver for the tigers field.
func (r *queryResolver) Tigers(ctx context.Context, page int, pageSize int) ([]*model.Tiger, error) {
	tigers, err := r.tigerUsecase.GetTigers(page, pageSize)
	if err != nil {
		return nil, err
	}

	return tigers, nil
}

// SightingByTiger is the resolver for the sightingByTiger field.
func (r *queryResolver) SightingByTiger(ctx context.Context, tigerID uint, page int, pageSize int) ([]*model.Sighting, error) {
	sightings, err := r.sightingUsecase.GetSightingsByTigerID(tigerID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return sightings, nil
}

// Tiger is the resolver for the tiger field.
func (r *sightingResolver) Tiger(ctx context.Context, obj *model.Sighting) (*model.Tiger, error) {
	if obj == nil || obj.TigerID == 0 {
		return nil, nil
	}

	t, err := r.tigerUsecase.GetTigerByID(obj.TigerID)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// User is the resolver for the user field.
func (r *sightingResolver) User(ctx context.Context, obj *model.Sighting) (*model.User, error) {
	if obj == nil || obj.UserID == 0 {
		return nil, nil
	}

	u, err := r.userUsecase.GetUserByID(obj.UserID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Sightings is the resolver for the sightings field.
func (r *tigerResolver) Sightings(ctx context.Context, obj *model.Tiger) ([]*model.Sighting, error) {
	if obj == nil || obj.ID == 0 {
		return nil, nil
	}

	sightings, err := r.sightingUsecase.GetSightingsByTigerID(obj.ID, 0, 0)
	if err != nil {
		return nil, err
	}

	return sightings, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Sighting returns SightingResolver implementation.
func (r *Resolver) Sighting() SightingResolver { return &sightingResolver{r} }

// Tiger returns TigerResolver implementation.
func (r *Resolver) Tiger() TigerResolver { return &tigerResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
	sightingResolver struct{ *Resolver }
	tigerResolver    struct{ *Resolver }
)
