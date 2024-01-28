package router

import (
	"checkpoint/middleware"
	"checkpoint/modules/authentication"
	"checkpoint/modules/project"
	"checkpoint/modules/upload"
	"checkpoint/utils"

	"github.com/kataras/iris/v12"
)

func Router(app *iris.Application) {
	authApi := app.Party("/auth")
	{
		authApi.Post("/signin", authentication.SignInController)
		authApi.Post("/signout", authentication.SignOutController)
		authApi.Post("/signup", authentication.SignUpController)
		authApi.Post("/refresh-token", authentication.RefreshTokenController)
	}

	storageApi := app.Party("/storage")
	storageApi.Use(middleware.AuthMiddleware())
	{
		storageApi.Post("/upload", upload.UploadController)
		storageApi.Get("/get-signed-url", upload.GetSignedURLController)
	}

	projectApi := app.Party("/projects")
	projectApi.Use(middleware.AuthMiddleware())
	{
		projectApi.Post("/", middleware.VerifyAuthorizationMiddleware(utils.AuthorizationPermissionData{
			PermissionSubject: "project",
			PermissionAction:  "CREATE",
		}), project.CreateProjectController)

		projectApi.Patch("/{id:uuid}", middleware.VerifyAuthorizationMiddleware(utils.AuthorizationPermissionData{
			PermissionSubject: "project",
			PermissionAction:  "UPDATE",
		}), project.UpdateProjectController)

		projectApi.Get("/{id:uuid}", project.GetProjectByIdController)
		projectApi.Delete("/{id:uuid}", project.DeleteProjectByIdController)
	}

}
