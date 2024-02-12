package main

import (
	"log"

	"checkpoint/auth"
	"checkpoint/db"
	"checkpoint/gql"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/iris-contrib/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/kataras/iris/v12"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbClient := db.GetPrimaryClient()

	app := iris.New()

	app.UseRouter(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH"},
		AllowCredentials: true,
	}))

	app.Use(iris.Compression)
	mergedSchema, err := os.ReadFile("generated.graphql")

	if err != nil {
		log.Fatal("Error loading graphql file")
	}

	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers(), graphql.MaxParallelism(20), graphql.UseStringDescriptions(), graphql.Directives()}
	schema := graphql.MustParseSchema(string(mergedSchema[:]), &gql.Resolver{}, opts...)

	app.Post("/graphql", iris.FromStd(auth.GraphqlContext(&relay.Handler{Schema: schema})))

	idleConnsClosed := make(chan struct{})

	iris.RegisterOnInterrupt(func() {
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		dbClient.Close()
		app.Shutdown(ctx)
		close(idleConnsClosed)
	})

	routes := app.GetRoutes()

	app.Listen(os.Getenv("BACKEND_PORT"), iris.WithoutInterruptHandler, iris.WithoutServerError(iris.ErrServerClosed), func(a *iris.Application) {
		fmt.Println("All routes")
		for _, route := range routes {
			fmt.Println(route)
		}
	})

	<-idleConnsClosed
}
