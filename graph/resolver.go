package graph

import "github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userUsecase     entities.UserUsecase
	tigerUsecase    entities.TigerUsecase
	sightingUsecase entities.SightingUsecase
}

func NewResolver(
	userUsecase entities.UserUsecase,
	tigerUsecase entities.TigerUsecase,
	sightingUsecase entities.SightingUsecase,
) *Resolver {
	return &Resolver{
		userUsecase:     userUsecase,
		tigerUsecase:    tigerUsecase,
		sightingUsecase: sightingUsecase,
	}
}
