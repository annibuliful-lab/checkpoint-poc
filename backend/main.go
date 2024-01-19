package main

import (
	"checkpoint/db"
	"checkpoint/router"
	"context"
	"log"
	"os"
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/kataras/iris/v12"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbClient := db.GetPrimaryClient()

	app := iris.New()

	app.UseRouter(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}))
	app.Use(iris.Compression)

	router.Router(app)

	idleConnsClosed := make(chan struct{})

	iris.RegisterOnInterrupt(func() {
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		dbClient.Close()
		app.Shutdown(ctx)
		close(idleConnsClosed)
	})

	app.Listen(os.Getenv("BACKEND_LISTEN"), iris.WithoutInterruptHandler, iris.WithoutServerError(iris.ErrServerClosed))

	<-idleConnsClosed
}
