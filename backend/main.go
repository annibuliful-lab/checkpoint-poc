package main

import (
	"checkpoint/db"
	"checkpoint/modules/authentication"
	"context"
	"time"

	"github.com/kataras/iris/v12"
)

func main() {
	var dbClient = db.GetClient()

	var app = iris.New()

	var authApi = app.Party("/auth")
	{
		authApi.Post("/signin", authentication.SignIn)
		authApi.Post("/signout", authentication.SignOut)
	}

	var idleConnsClosed = make(chan struct{})

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
