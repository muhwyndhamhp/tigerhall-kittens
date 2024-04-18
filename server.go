package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/tigerhall-kittens/db"
	"github.com/muhwyndhamhp/tigerhall-kittens/graph"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/sighting"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/tiger"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/modules/user"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/config"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/email"
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
	em := email.NewEmailClient()
	s3 := s3client.NewS3Client()

	queue := make(chan email.SightingEmail)

	userRepo := user.NewUserRepository(d)
	tigerRepo := tiger.NewTigerRepository(d)
	sightingRepo := sighting.NewSightingRepository(d)

	userUsecase := user.NewUserUsecase(userRepo)
	tigerUsecase := tiger.NewTigerUsecase(tigerRepo, sightingRepo, s3)
	sightingUsecase := sighting.NewSightingUsecase(sightingRepo, tigerRepo, userRepo, s3, queue)

	resolver := graph.NewResolver(userUsecase, tigerUsecase, sightingUsecase)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	e.Use(user.AuthMiddleware(userRepo))
	e.GET("/graphiql", echo.WrapHandler(playground.Handler("GraphQL playground", "/query")))
	e.POST("/query", echo.WrapHandler(srv))
	e.GET("/altair", ServeAltair)
	e.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusMovedPermanently, "/altair") })

	go em.QueueConsumer(queue)

	log.Printf("connect to http://localhost:%s/graphiql for GraphiQL playground", port)
	log.Printf("or connect to http://localhost:%s/altair for Altair", port)
	log.Fatal(e.Start(":" + port))
}

func ServeAltair(c echo.Context) error {
	t, err := template.ParseFiles("public/altair.html")
	if err != nil {
		return err
	}

	var o bytes.Buffer
	err = t.ExecuteTemplate(&o, "altair", map[string]interface{}{
		"Endpoint": config.Get(config.BASE_URL) + "/query",
	})
	if err != nil {
		return err
	}

	html := o.String()

	return c.HTML(http.StatusOK, html)
}
