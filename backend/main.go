package main

import (
	"checkpoint/auth"
	"checkpoint/gql"
	"checkpoint/gql/directive"
	"context"
	"flag"
	"log"
	"net/http"
	"os/signal"

	"os"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mergedSchema, err := os.ReadFile("generated.graphql")

	if err != nil {
		log.Fatal("Error loading graphql file")
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: true,
		Debug:            true,
	})

	opts := []graphql.SchemaOpt{
		graphql.UseFieldResolvers(),
		graphql.MaxParallelism(20),
		graphql.UseStringDescriptions(),
		graphql.Directives(&directive.AccessDirective{}),
	}

	// init graphQL schema
	s, err := graphql.ParseSchema(string(mergedSchema[:]), gql.GraphqlResolver(), opts...)
	if err != nil {
		panic(err)
	}

	// graphQL handler
	graphQLHandler := corsMiddleware.Handler(graphqlws.NewHandlerFunc(s, auth.GraphqlContext(&relay.Handler{Schema: s})))
	http.Handle("/graphql", graphQLHandler)

	var listenAddress = flag.String("listen", os.Getenv("BACKEND_PORT"), "Listen address.")

	log.Printf("Listening at http://%s", *listenAddress)

	httpServer := http.Server{
		Addr: *listenAddress,
	}

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-idleConnectionsClosed

}
