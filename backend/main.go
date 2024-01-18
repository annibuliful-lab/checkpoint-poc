package main

import (
	router "checkpoint/Router"
	"checkpoint/db"
	"context"
	"log"
	"os"
	"time"

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
