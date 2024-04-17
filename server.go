package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/tigerhall-kittens/db"
	"github.com/muhwyndhamhp/tigerhall-kittens/graph"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/sighting"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/tiger"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/user"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/config"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/s3client"
)

const defaultPort = "8080"

func main() {
	port := config.Get(config.PORT)
	if port == "" {
		port = defaultPort
	}

	e := echo.New()

	d := db.GetDB()

	s3 := s3client.NewS3Client()

	userRepo := user.NewUserRepository(d)
	tigerRepo := tiger.NewTigerRepository(d)
	sightingRepo := sighting.NewSightingRepository(d)

	userUsecase := user.NewUserUsecase(userRepo)
	tigerUsecase := tiger.NewTigerUsecase(tigerRepo, sightingRepo)
	sightingUsecase := sighting.NewSightingUsecase(sightingRepo, tigerRepo, userRepo, s3)

	resolver := graph.NewResolver(userUsecase, tigerUsecase, sightingUsecase)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	e.Use(user.AuthMiddleware(userRepo))
	e.GET("/", echo.WrapHandler(playground.Handler("GraphQL playground", "/query")))
	e.POST("/query", echo.WrapHandler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(e.Start(":" + port))
}
