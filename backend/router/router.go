package router

import (
	"checkpoint/middleware"
	"checkpoint/modules/authentication"
	"checkpoint/modules/project"
	projectRole "checkpoint/modules/project-role"
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

	storageApi := app.Party("/storages")
	storageApi.Use(middleware.AuthMiddleware())
	{
		storageApi.Post("/upload", upload.UploadController)
		storageApi.Get("/get-signed-url", upload.GetSignedURLController)
	}

	projectApi := app.Party("/projects")
	projectApi.Use(middleware.AuthMiddleware())
	{
		projectApi.Post("/", project.CreateProjectController)
		projectApi.Patch("/{id:uuid}", project.UpdateProjectController)
		projectApi.Get("/{id:uuid}", project.GetProjectByIdController)
		projectApi.Delete("/{id:uuid}", project.DeleteProjectByIdController)

		// roles
		projectApi.Post("/roles", middleware.VerifyAuthorizationMiddleware(utils.AuthorizationPermissionData{
			PermissionSubject: "project",
			PermissionAction:  "CREATE",
		}), projectRole.CreateProjectRoleController)
		projectApi.Patch("/roles/{id:uuid}", middleware.VerifyAuthorizationMiddleware(utils.AuthorizationPermissionData{
			PermissionSubject: "project",
			PermissionAction:  "UPDATE",
		}), projectRole.UpdateProjectRoleController)
		projectApi.Get("/roles/{id:uuid}", middleware.VerifyAuthorizationMiddleware(utils.AuthorizationPermissionData{
			PermissionSubject: "project",
			PermissionAction:  "READ",
		}), projectRole.GetProjectRoleByIdController)
		projectApi.Get("/roles", middleware.VerifyAuthorizationMiddleware(utils.AuthorizationPermissionData{
			PermissionSubject: "project",
			PermissionAction:  "READ",
		}), projectRole.GetProjectRolesController)
	}

}
