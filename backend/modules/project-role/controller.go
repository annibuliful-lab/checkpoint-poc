package projectRole

import (
	"checkpoint/auth"
	"checkpoint/jwt"
	"checkpoint/utils"
	"log"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/samber/lo"
)

func CreateProjectRoleController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)

	var data struct {
		Title         string   `json:"title"`
		PermissionIds []string `json:"permissionIds"`
	}

	err := ctx.ReadJSON(&data)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	payload, _ := jwt.VerifyToken(headers.Token)

	match := auth.VerifyProjectOwner(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(headers.ProjectId),
		AccountId: payload.AccountId,
	})

	if !match {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"message": utils.ForbiddenOperation.Error(),
			"data":    nil,
		})
		return
	}

	projectRole, code, err := CreateProjectRole(CreateProjectRoleData{
		ProjectId:     uuid.MustParse(headers.ProjectId),
		PermissionIds: data.PermissionIds,
		Title:         data.Title,
	})

	if err != nil {
		ctx.StatusCode(code)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	permissions, _, _ := GetProjectRolePermissionsByProjectRoleId(projectRole.ID)

	ctx.StatusCode(201)
	ctx.JSON(iris.Map{
		"message": "created",
		"data": iris.Map{
			"id":          projectRole.ID,
			"title":       projectRole.Title,
			"projectId":   projectRole.ProjectId,
			"createdAt":   projectRole.CreatedAt,
			"updatedAt":   projectRole.UpdatedAt,
			"permissions": permissions,
		},
	})
}

func GetProjectRolesController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)
	payload, _ := jwt.VerifyToken(headers.Token)
	match := auth.VerifyProjectAccount(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(headers.ProjectId),
		AccountId: payload.AccountId,
	})

	if !match {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"message": utils.ForbiddenOperation.Error(),
			"data":    nil,
		})
		return
	}

	limit := ctx.URLParamInt64Default("limit", 20)
	skip := ctx.URLParamInt64Default("skip", 0)
	search := ctx.URLParam("search")

	projectRolesResponse, code, err := GetProjectRoles(GetProjectRolesData{
		ProjectId: uuid.MustParse(headers.ProjectId),
		Search:    search,
		pagination: utils.OffsetPagination{
			Limit: limit,
			Skip:  skip,
		},
	})

	if err != nil {
		log.Println("get-project-role-error", err.Error())
		ctx.StatusCode(code)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	productRoleIds := lo.Map(*projectRolesResponse, func(item ProjectRoleResponse, index int) uuid.UUID {
		return item.ID
	})

	permissions, _, _ := GetProjectRolePermissionsByProjectRoleIds(productRoleIds)

	projectRolesWithPermissions := lo.Map(*projectRolesResponse, func(projectRole ProjectRoleResponse, index int) ProjectRoleResponse {
		rolePermissions := lo.Filter(*permissions, func(item PermissionResponse, index int) bool {
			return item.RoleID == projectRole.ID
		})

		return ProjectRoleResponse{

			ID:          projectRole.ID,
			Title:       projectRole.Title,
			ProjectId:   projectRole.ProjectId,
			CreatedAt:   projectRole.CreatedAt,
			UpdatedAt:   projectRole.UpdatedAt,
			Permissions: rolePermissions,
		}
	})

	ctx.StatusCode(200)
	ctx.JSON(iris.Map{
		"data": projectRolesWithPermissions,
	})

}

func GetProjectRoleByIdController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)
	id, err := uuid.Parse(ctx.Params().Get("id"))

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	payload, _ := jwt.VerifyToken(headers.Token)
	match := auth.VerifyProjectAccount(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(headers.ProjectId),
		AccountId: payload.AccountId,
	})

	if !match {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"message": utils.ForbiddenOperation.Error(),
			"data":    nil,
		})
		return
	}

	projectRole, code, err := GetProjectRoleById(GetProjectRoleByIdData{
		ID:        id,
		ProjectId: uuid.MustParse(headers.ProjectId),
	})

	if err != nil {
		log.Println("get-project-role-error", err.Error())
		ctx.StatusCode(code)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	permissions, _, _ := GetProjectRolePermissionsByProjectRoleId(projectRole.ID)

	ctx.StatusCode(200)
	ctx.JSON(iris.Map{
		"message": "created",
		"data": iris.Map{
			"id":          projectRole.ID,
			"title":       projectRole.Title,
			"projectId":   projectRole.ProjectId,
			"createdAt":   projectRole.CreatedAt,
			"updatedAt":   projectRole.UpdatedAt,
			"permissions": permissions,
		},
	})

}

func UpdateProjectRoleController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)
	id, err := uuid.Parse(ctx.Params().Get("id"))

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	var data struct {
		Title         string   `json:"title"`
		PermissionIds []string `json:"permissionIds"`
	}

	err = ctx.ReadJSON(&data)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	payload, _ := jwt.VerifyToken(headers.Token)

	match := auth.VerifyProjectOwner(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(headers.ProjectId),
		AccountId: payload.AccountId,
	})

	if !match {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"message": utils.ForbiddenOperation.Error(),
			"data":    nil,
		})
		return
	}

	projectRole, code, err := UpdateProjectRole(UpdateProjectRoleData{
		ID:            id,
		PermissionIds: data.PermissionIds,
		Title:         data.Title,
		ProjectId:     uuid.MustParse(headers.ProjectId),
	})

	if err != nil {
		ctx.StatusCode(code)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	permissions, _, _ := GetProjectRolePermissionsByProjectRoleId(projectRole.ID)

	ctx.StatusCode(200)
	ctx.JSON(iris.Map{
		"message": "created",
		"data": iris.Map{
			"id":          projectRole.ID,
			"title":       projectRole.Title,
			"projectId":   projectRole.ProjectId,
			"createdAt":   projectRole.CreatedAt,
			"updatedAt":   projectRole.UpdatedAt,
			"permissions": permissions,
		},
	})
}
