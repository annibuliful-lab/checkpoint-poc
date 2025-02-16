package main

import (
	"checkpoint/auth"
	"checkpoint/db"
	"checkpoint/gql"
	access "checkpoint/gql/directive/access"
	stationApi "checkpoint/gql/directive/station-api"
	uploadmiddleware "checkpoint/gql/upload-middleware"
	"checkpoint/utils/graphql_utils"
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
	err := godotenv.Load()
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
		AllowedHeaders:   []string{"*"},
		// Debug:            true,
	})

	opts := []graphql.SchemaOpt{
		graphql.UseFieldResolvers(),
		graphql.MaxParallelism(20),
		graphql.UseStringDescriptions(),
		graphql.Directives(
			&access.AccessDirective{},
			&stationApi.StationApiDirective{},
		),
	}

	// init graphQL schema
	schema, err := graphql.ParseSchema(string(mergedSchema[:]), gql.GraphqlResolver(), opts...)
	if err != nil {
		panic(err)
	}

	// graphQL handler
	graphQLHandler := corsMiddleware.Handler(
		graphqlws.NewHandlerFunc(
			schema,
			auth.GraphqlContext(uploadmiddleware.Handler(&relay.Handler{Schema: schema})),
			graphqlws.WithContextGenerator(
				graphqlws.ContextGeneratorFunc(auth.WebsocketGraphqlContext),
			),
		),
	)

	http.Handle("/graphql", graphQLHandler)
	http.HandleFunc("/graphiql", http.HandlerFunc(graphql_utils.Graphiql))

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
		db.GetPrimaryClient().Close()
		db.GetRedisClient().Close()

		close(idleConnectionsClosed)
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-idleConnectionsClosed

}
