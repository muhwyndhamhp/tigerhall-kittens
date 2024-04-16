package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/tigerhall-kittens/db"
	"github.com/muhwyndhamhp/tigerhall-kittens/graph"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/user"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/config"
)

const defaultPort = "8080"

func main() {
	port := config.Get(config.PORT)
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	e := echo.New()

	d := db.GetDB()
	repo := user.NewUserRepository(d)

	e.Use(echo.WrapMiddleware(user.AuthMiddleware(repo)))
	e.GET("/", echo.WrapHandler(playground.Handler("GraphQL playground", "/query")))
	e.POST("/query", echo.WrapHandler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(e.Start(":" + port))
}
