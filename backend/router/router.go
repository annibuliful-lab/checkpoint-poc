package router

import (
	"checkpoint/modules/authentication"
	"checkpoint/modules/upload"

	"github.com/kataras/iris/v12"
)

func Router(app *iris.Application) {
	authApi := app.Party("/auth")
	{
		authApi.Post("/signin", authentication.SignInController)
		authApi.Post("/signout", authentication.SignOutController)
		authApi.Post("/signup", authentication.SignUpController)
	}

	storageApi := app.Party("/storage")
	{
		storageApi.Post("/upload", upload.Upload)
		storageApi.Post("/get-signed-url", upload.GetSignedURL)
	}

}
