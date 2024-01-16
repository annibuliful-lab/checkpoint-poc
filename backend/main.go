package main

import (
	"checkpoint/db"
	"checkpoint/modules/authentication"
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbClient := db.GetClient()

	app := iris.New()

	authApi := app.Party("/auth")
	{
		authApi.Post("/signin", authentication.SignInController)
		authApi.Post("/signout", authentication.SignOutController)
		authApi.Post("/signup", authentication.SignUpController)
	}

	idleConnsClosed := make(chan struct{})

	iris.RegisterOnInterrupt(func() {
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// close all hosts.
		dbClient.Close()
		app.Shutdown(ctx)
		close(idleConnsClosed)
	})

	app.Listen(":8080", iris.WithoutInterruptHandler, iris.WithoutServerError(iris.ErrServerClosed))

	<-idleConnsClosed
}
